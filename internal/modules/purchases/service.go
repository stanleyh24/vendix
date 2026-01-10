package purchases

import (
	"context"
	"fmt"
	"time"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/logger"
	"vendix/internal/modules/accounting"
	"vendix/internal/modules/products"

	"github.com/google/uuid"
)

type Service struct {
	repo           *Repository
	productsRepo   *products.Repository
	accountingSvc  *accounting.Service
	accountingRepo *accounting.Repository
	cfg            *config.Config
}

func NewService(db *database.DB, cfg *config.Config) *Service {
	return &Service{
		repo:           NewRepository(db),
		productsRepo:   products.NewRepository(db),
		accountingSvc:  accounting.NewService(db, cfg),
		accountingRepo: accounting.NewRepository(db),
		cfg:            cfg,
	}
}

// Create creates a new purchase order
func (s *Service) Create(ctx context.Context, schema string, req *CreatePurchaseRequest, userID *uuid.UUID) (*PurchaseResponse, error) {
	// Parse supplier ID
	supplierID, err := uuid.Parse(req.SupplierID)
	if err != nil {
		return nil, fmt.Errorf("invalid supplier ID: %w", err)
	}

	// Parse purchase date
	purchaseDate, err := time.Parse("2006-01-02", req.PurchaseDate)
	if err != nil {
		return nil, fmt.Errorf("invalid purchase date format: %w", err)
	}

	// Get next purchase number
	purchaseNumber, err := s.repo.GetNextPurchaseNumber(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to get next purchase number: %w", err)
	}

	// Calculate totals
	var subtotal, taxAmount, total float64
	lines := make([]*PurchaseLine, 0, len(req.Lines))

	for _, lineReq := range req.Lines {
		lineSubtotal := lineReq.Quantity * lineReq.UnitPrice
		lineTaxAmount := lineSubtotal * (lineReq.TaxRate / 100)
		lineTotal := lineSubtotal + lineTaxAmount

		subtotal += lineSubtotal
		taxAmount += lineTaxAmount
		total += lineTotal

		var productID *uuid.UUID
		if lineReq.ProductID != nil {
			pid, err := uuid.Parse(*lineReq.ProductID)
			if err == nil {
				productID = &pid
			}
		}

		line := &PurchaseLine{
			ID:              uuid.New(),
			PurchaseID:      uuid.Nil, // Will be set after purchase creation
			ProductID:       productID,
			Description:     lineReq.Description,
			Quantity:        lineReq.Quantity,
			UnitPrice:       lineReq.UnitPrice,
			TaxRate:         lineReq.TaxRate,
			Subtotal:        lineSubtotal,
			TaxAmount:       lineTaxAmount,
			Total:           lineTotal,
			ReceivedQuantity: 0,
			CreatedAt:       time.Now(),
		}
		lines = append(lines, line)
	}

	// Parse expected delivery date
	var expectedDeliveryDate *time.Time
	if req.ExpectedDeliveryDate != nil && *req.ExpectedDeliveryDate != "" {
		edd, err := time.Parse("2006-01-02", *req.ExpectedDeliveryDate)
		if err == nil {
			expectedDeliveryDate = &edd
		}
	}

	// Determine payment status based on payment method
	paymentStatus := string(PaymentStatusPending)
	if req.PaymentMethod == string(PaymentMethodCash) {
		paymentStatus = string(PaymentStatusPaid)
	}

	// Set payment terms (default: 30 days for credit)
	paymentTerms := 30
	if req.PaymentTerms != nil {
		paymentTerms = *req.PaymentTerms
	}

	// Calculate payment due date for credit purchases
	var paymentDueDate *time.Time
	if req.PaymentMethod == string(PaymentMethodCredit) {
		dueDate := purchaseDate.AddDate(0, 0, paymentTerms)
		paymentDueDate = &dueDate
	}

	// Create purchase
	purchase := &Purchase{
		ID:                  uuid.New(),
		PurchaseNumber:      purchaseNumber,
		SupplierID:          supplierID,
		PurchaseDate:        purchaseDate,
		ExpectedDeliveryDate: expectedDeliveryDate,
		PaymentTerms:        paymentTerms,
		PaymentDueDate:      paymentDueDate,
		PaymentMethod:       req.PaymentMethod,
		PaymentStatus:       paymentStatus,
		Subtotal:            subtotal,
		TaxAmount:           taxAmount,
		Total:               total,
		Currency:            "DOP",
		Reference:           req.Reference,
		Notes:               req.Notes,
		Status:              string(PurchaseStatusDraft),
		CreatedBy:           userID,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Save purchase
	if err := s.repo.Create(ctx, schema, purchase); err != nil {
		return nil, fmt.Errorf("failed to create purchase: %w", err)
	}

	// Save lines
	for _, line := range lines {
		line.PurchaseID = purchase.ID
		if err := s.repo.CreateLine(ctx, schema, line); err != nil {
			logger.Error("Failed to create purchase line", "error", err, "purchase_id", purchase.ID)
			// Continue anyway, but log the error
		}
	}

	// Get supplier name
	supplierName, _ := s.repo.GetSupplierName(ctx, schema, supplierID)

	// Get lines for response
	purchaseLines, _ := s.repo.GetLinesByPurchaseID(ctx, schema, purchase.ID)

	response := &PurchaseResponse{
		Purchase:     *purchase,
		SupplierName: &supplierName,
		Lines:        make([]PurchaseLine, len(purchaseLines)),
	}
	for i, line := range purchaseLines {
		response.Lines[i] = *line
	}

	logger.Info("Purchase created", "purchase_number", purchaseNumber, "supplier_id", supplierID)
	return response, nil
}

// Receive marks a purchase as received and updates inventory
func (s *Service) Receive(ctx context.Context, schema string, purchaseID string, req *ReceivePurchaseRequest) error {
	id, err := uuid.Parse(purchaseID)
	if err != nil {
		return fmt.Errorf("invalid purchase ID: %w", err)
	}

	// Get purchase
	purchase, err := s.repo.GetByID(ctx, schema, id)
	if err != nil {
		return fmt.Errorf("purchase not found: %w", err)
	}

	if purchase.Status == string(PurchaseStatusReceived) {
		return fmt.Errorf("purchase already received")
	}

	if purchase.Status == string(PurchaseStatusCancelled) {
		return fmt.Errorf("cannot receive cancelled purchase")
	}

	// Parse received date (validate format)
	_, err = time.Parse("2006-01-02", req.ReceivedDate)
	if err != nil {
		return fmt.Errorf("invalid received date format: %w", err)
	}

	// Get all lines
	lines, err := s.repo.GetLinesByPurchaseID(ctx, schema, id)
	if err != nil {
		return fmt.Errorf("failed to get purchase lines: %w", err)
	}

	// Create map of line updates
	lineUpdates := make(map[string]float64)
	for _, reqLine := range req.Lines {
		lineUpdates[reqLine.LineID] = reqLine.ReceivedQuantity
	}

	// Update lines and inventory
	for _, line := range lines {
		receivedQty, ok := lineUpdates[line.ID.String()]
		if !ok {
			continue
		}

		// Validate received quantity doesn't exceed ordered
		if receivedQty > line.Quantity {
			return fmt.Errorf("received quantity (%.2f) exceeds ordered quantity (%.2f) for line %s", receivedQty, line.Quantity, line.ID.String())
		}

		// Update line received quantity
		if err := s.repo.UpdateLineReceivedQuantity(ctx, schema, line.ID, receivedQty); err != nil {
			return fmt.Errorf("failed to update line received quantity: %w", err)
		}

		// Update product stock if product is linked
		if line.ProductID != nil {
			if err := s.repo.UpdateProductStock(ctx, schema, *line.ProductID, receivedQty); err != nil {
				logger.Warn("Failed to update product stock", "error", err, "product_id", *line.ProductID)
				// Continue anyway
			} else {
				logger.Info("Product stock updated", "product_id", *line.ProductID, "quantity", receivedQty)
			}
		}
	}

	// Update purchase status
	purchase.Status = string(PurchaseStatusReceived)
	now := time.Now()
	purchase.ReceivedAt = &now
	purchase.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, schema, purchase); err != nil {
		return fmt.Errorf("failed to update purchase status: %w", err)
	}

	// Generate journal entry for received purchase
	if err := s.generatePurchaseJournalEntry(ctx, schema, purchase); err != nil {
		logger.Warn("Failed to generate journal entry for purchase", "error", err, "purchase_id", id)
		// Don't fail the receive operation if journal entry fails
	}

	logger.Info("Purchase received", "purchase_id", id)
	return nil
}

// generatePurchaseJournalEntry creates accounting entries for a purchase
func (s *Service) generatePurchaseJournalEntry(ctx context.Context, schema string, purchase *Purchase) error {
	// Get purchase lines
	lines, err := s.repo.GetLinesByPurchaseID(ctx, schema, purchase.ID)
	if err != nil {
		return fmt.Errorf("failed to get purchase lines: %w", err)
	}

	// Get account mappings
	inventoryAccount, err := s.accountingSvc.GetAccountForTransaction(ctx, schema, accounting.TransactionTypePurchaseInventory, "1131", "Inventario de Mercancías")
	if err != nil {
		return fmt.Errorf("failed to get inventory account: %w", err)
	}

	accountsPayableAccount, err := s.accountingSvc.GetAccountForTransaction(ctx, schema, accounting.TransactionTypePurchaseAP, "2111", "Proveedores")
	if err != nil {
		return fmt.Errorf("failed to get accounts payable account: %w", err)
	}

	// Calculate totals from received lines only
	var totalInventory, totalTax float64
	for _, line := range lines {
		if line.ReceivedQuantity > 0 {
			receivedSubtotal := line.ReceivedQuantity * line.UnitPrice
			receivedTax := receivedSubtotal * (line.TaxRate / 100)
			totalInventory += receivedSubtotal
			totalTax += receivedTax
		}
	}

	if totalInventory == 0 {
		return nil // Nothing received, no journal entry needed
	}

	// Get next entry number from repository
	entryNumber, err := s.accountingRepo.GetNextEntryNumber(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to get next journal entry number: %w", err)
	}

	entry := &accounting.JournalEntry{
		ID:          uuid.New(),
		EntryNumber: entryNumber,
		EntryDate:   purchase.PurchaseDate,
		Description: fmt.Sprintf("Compra de inventario - %s", purchase.PurchaseNumber),
		Reference:   stringPtr(purchase.PurchaseNumber),
		Status:      "posted",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// DEBE: Inventario (1131)
	entry.Lines = append(entry.Lines, accounting.JournalEntryLine{
		ID:             uuid.New(),
		JournalEntryID: entry.ID,
		AccountCode:    inventoryAccount.AccountCode,
		AccountName:    inventoryAccount.AccountName,
		Debit:          totalInventory + totalTax,
		Credit:         0,
		Description:    stringPtr(fmt.Sprintf("Entrada de inventario - %s", purchase.PurchaseNumber)),
		CreatedAt:      time.Now(),
	})

	// HABER: Proveedores (2111) - si es a crédito
	if purchase.PaymentMethod == string(PaymentMethodCredit) {
		entry.Lines = append(entry.Lines, accounting.JournalEntryLine{
			ID:             uuid.New(),
			JournalEntryID: entry.ID,
			AccountCode:    accountsPayableAccount.AccountCode,
			AccountName:    accountsPayableAccount.AccountName,
			Debit:          0,
			Credit:         totalInventory + totalTax,
			Description:    stringPtr(fmt.Sprintf("Cuenta por pagar - %s", purchase.PurchaseNumber)),
			CreatedAt:      time.Now(),
		})
	} else {
		// Si es efectivo, usar la cuenta de caja/banco según el método
		paymentAccount, err := s.accountingSvc.GetAccountForTransaction(ctx, schema, accounting.TransactionTypePaymentCash, "1111", "Caja General")
		if purchase.PaymentMethod == string(PaymentMethodBank) {
			paymentAccount, err = s.accountingSvc.GetAccountForTransaction(ctx, schema, accounting.TransactionTypePaymentBank, "1112", "Banco - Cuenta Corriente")
		}
		if err != nil {
			return fmt.Errorf("failed to get payment account: %w", err)
		}

		entry.Lines = append(entry.Lines, accounting.JournalEntryLine{
			ID:             uuid.New(),
			JournalEntryID: entry.ID,
			AccountCode:    paymentAccount.AccountCode,
			AccountName:    paymentAccount.AccountName,
			Debit:          0,
			Credit:         totalInventory + totalTax,
			Description:    stringPtr(fmt.Sprintf("Pago de compra - %s", purchase.PurchaseNumber)),
			CreatedAt:      time.Now(),
		})
	}

	// Save journal entry
	if err := s.accountingRepo.CreateJournalEntry(ctx, schema, entry); err != nil {
		return fmt.Errorf("failed to create journal entry: %w", err)
	}

	logger.Info("Journal entry created for purchase", "purchase_number", purchase.PurchaseNumber, "entry_number", entryNumber)
	return nil
}

// GetByID gets a purchase by ID
func (s *Service) GetByID(ctx context.Context, schema string, id string) (*PurchaseResponse, error) {
	purchaseID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid purchase ID: %w", err)
	}

	purchase, err := s.repo.GetByID(ctx, schema, purchaseID)
	if err != nil {
		return nil, fmt.Errorf("purchase not found: %w", err)
	}

	lines, err := s.repo.GetLinesByPurchaseID(ctx, schema, purchaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase lines: %w", err)
	}

	supplierName, _ := s.repo.GetSupplierName(ctx, schema, purchase.SupplierID)

	response := &PurchaseResponse{
		Purchase:     *purchase,
		SupplierName: &supplierName,
		Lines:        make([]PurchaseLine, len(lines)),
	}
	for i, line := range lines {
		response.Lines[i] = *line
	}

	return response, nil
}

// List gets all purchases
func (s *Service) List(ctx context.Context, schema string, supplierID *string, status *string, limit, offset int) ([]*PurchaseResponse, error) {
	var supplierUUID *uuid.UUID
	if supplierID != nil {
		uid, err := uuid.Parse(*supplierID)
		if err == nil {
			supplierUUID = &uid
		}
	}

	purchases, err := s.repo.List(ctx, schema, supplierUUID, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list purchases: %w", err)
	}

	responses := make([]*PurchaseResponse, len(purchases))
	for i, purchase := range purchases {
		lines, err := s.repo.GetLinesByPurchaseID(ctx, schema, purchase.ID)
		if err != nil {
			logger.Warn("Failed to get purchase lines", "error", err, "purchase_id", purchase.ID)
			lines = []*PurchaseLine{}
		}

		supplierName, err := s.repo.GetSupplierName(ctx, schema, purchase.SupplierID)
		if err != nil {
			logger.Warn("Failed to get supplier name", "error", err, "supplier_id", purchase.SupplierID)
			supplierName = ""
		}

		var supplierNamePtr *string
		if supplierName != "" {
			supplierNamePtr = &supplierName
		}

		response := &PurchaseResponse{
			Purchase:     *purchase,
			SupplierName: supplierNamePtr,
			Lines:        make([]PurchaseLine, len(lines)),
		}
		for j, line := range lines {
			response.Lines[j] = *line
		}
		responses[i] = response
	}

	return responses, nil
}

// Update updates a purchase
func (s *Service) Update(ctx context.Context, schema string, id string, req *UpdatePurchaseRequest) (*PurchaseResponse, error) {
	purchaseID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid purchase ID: %w", err)
	}

	purchase, err := s.repo.GetByID(ctx, schema, purchaseID)
	if err != nil {
		return nil, fmt.Errorf("purchase not found: %w", err)
	}

	if purchase.Status == string(PurchaseStatusReceived) {
		return nil, fmt.Errorf("cannot update received purchase")
	}

	if purchase.Status == string(PurchaseStatusCancelled) {
		return nil, fmt.Errorf("cannot update cancelled purchase")
	}

	// Update fields
	if req.PurchaseDate != nil {
		purchaseDate, err := time.Parse("2006-01-02", *req.PurchaseDate)
		if err == nil {
			purchase.PurchaseDate = purchaseDate
		}
	}

	if req.ExpectedDeliveryDate != nil {
		if *req.ExpectedDeliveryDate == "" {
			purchase.ExpectedDeliveryDate = nil
		} else {
			edd, err := time.Parse("2006-01-02", *req.ExpectedDeliveryDate)
			if err == nil {
				purchase.ExpectedDeliveryDate = &edd
			}
		}
	}

	if req.PaymentMethod != nil {
		purchase.PaymentMethod = *req.PaymentMethod
	}

	if req.Reference != nil {
		purchase.Reference = req.Reference
	}

	if req.Notes != nil {
		purchase.Notes = req.Notes
	}

	if req.Status != nil {
		purchase.Status = *req.Status
	}

	// Update lines if provided
	if req.Lines != nil && len(req.Lines) > 0 {
		// Delete existing lines
		if err := s.repo.DeleteLinesByPurchaseID(ctx, schema, purchaseID); err != nil {
			return nil, fmt.Errorf("failed to delete existing lines: %w", err)
		}

		// Recalculate totals
		var subtotal, taxAmount, total float64

		for _, lineReq := range req.Lines {
			lineSubtotal := lineReq.Quantity * lineReq.UnitPrice
			lineTaxAmount := lineSubtotal * (lineReq.TaxRate / 100)
			lineTotal := lineSubtotal + lineTaxAmount

			subtotal += lineSubtotal
			taxAmount += lineTaxAmount
			total += lineTotal

			var productID *uuid.UUID
			if lineReq.ProductID != nil {
				pid, err := uuid.Parse(*lineReq.ProductID)
				if err == nil {
					productID = &pid
				}
			}

			line := &PurchaseLine{
				ID:              uuid.New(),
				PurchaseID:      purchaseID,
				ProductID:       productID,
				Description:     lineReq.Description,
				Quantity:        lineReq.Quantity,
				UnitPrice:       lineReq.UnitPrice,
				TaxRate:         lineReq.TaxRate,
				Subtotal:        lineSubtotal,
				TaxAmount:       lineTaxAmount,
				Total:           lineTotal,
				ReceivedQuantity: 0,
				CreatedAt:       time.Now(),
			}

			if err := s.repo.CreateLine(ctx, schema, line); err != nil {
				return nil, fmt.Errorf("failed to create purchase line: %w", err)
			}
		}

		purchase.Subtotal = subtotal
		purchase.TaxAmount = taxAmount
		purchase.Total = total
	}

	purchase.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, schema, purchase); err != nil {
		return nil, fmt.Errorf("failed to update purchase: %w", err)
	}

	return s.GetByID(ctx, schema, id)
}

// Delete deletes a purchase (only if draft)
func (s *Service) Delete(ctx context.Context, schema string, id string) error {
	purchaseID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid purchase ID: %w", err)
	}

	purchase, err := s.repo.GetByID(ctx, schema, purchaseID)
	if err != nil {
		return fmt.Errorf("purchase not found: %w", err)
	}

	if purchase.Status != string(PurchaseStatusDraft) {
		return fmt.Errorf("can only delete draft purchases")
	}

	// Delete lines first
	if err := s.repo.DeleteLinesByPurchaseID(ctx, schema, purchaseID); err != nil {
		return fmt.Errorf("failed to delete purchase lines: %w", err)
	}

	// Delete purchase (would need a delete method in repository)
	// For now, we'll just mark it as cancelled
	purchase.Status = string(PurchaseStatusCancelled)
	purchase.UpdatedAt = time.Now()
	return s.repo.Update(ctx, schema, purchase)
}

// ========== Supplier Payments Service Methods ==========

// CreateSupplierPayment creates a new supplier payment and allocates it to purchases
func (s *Service) CreateSupplierPayment(ctx context.Context, schema string, req *CreateSupplierPaymentRequest, userID *uuid.UUID) (*SupplierPaymentResponse, error) {
	// Parse supplier ID
	supplierID, err := uuid.Parse(req.SupplierID)
	if err != nil {
		return nil, fmt.Errorf("invalid supplier ID: %w", err)
	}

	// Parse payment date
	paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
	if err != nil {
		return nil, fmt.Errorf("invalid payment date format: %w", err)
	}

	// Validate payment method
	if req.PaymentMethod != string(PaymentMethodCash) &&
		req.PaymentMethod != string(PaymentMethodCredit) &&
		req.PaymentMethod != string(PaymentMethodBank) &&
		req.PaymentMethod != string(PaymentMethodCheck) {
		return nil, fmt.Errorf("invalid payment method: %s", req.PaymentMethod)
	}

	// Validate total amount matches allocations
	totalAllocated := 0.0
	for _, alloc := range req.Allocations {
		totalAllocated += alloc.Amount
	}
	
	if totalAllocated != req.Amount {
		return nil, fmt.Errorf("total allocation amount (%.2f) does not match payment amount (%.2f)", totalAllocated, req.Amount)
	}

	// Get next payment number
	paymentNumber, err := s.repo.GetNextSupplierPaymentNumber(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to generate payment number: %w", err)
	}

	// Create payment
	payment := &SupplierPayment{
		ID:            uuid.New(),
		PaymentNumber: paymentNumber,
		SupplierID:    supplierID,
		PaymentMethod: req.PaymentMethod,
		PaymentDate:   paymentDate,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Reference:     req.Reference,
		Notes:         req.Notes,
		Status:        "completed",
		CreatedBy:     userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if payment.Currency == "" {
		payment.Currency = "DOP"
	}

	// Create payment
	if err := s.repo.CreateSupplierPayment(ctx, schema, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Create allocations and update purchase payment status
	for _, allocReq := range req.Allocations {
		purchaseID, err := uuid.Parse(allocReq.PurchaseID)
		if err != nil {
			return nil, fmt.Errorf("invalid purchase ID: %s", allocReq.PurchaseID)
		}

		// Verify purchase exists and belongs to supplier
		purchase, err := s.repo.GetByID(ctx, schema, purchaseID)
		if err != nil {
			return nil, fmt.Errorf("purchase not found: %s", allocReq.PurchaseID)
		}

		if purchase.SupplierID != supplierID {
			return nil, fmt.Errorf("purchase %s does not belong to supplier %s", purchaseID, supplierID)
		}

		// Create allocation
		allocation := &SupplierPaymentAllocation{
			ID:               uuid.New(),
			SupplierPaymentID: payment.ID,
			PurchaseID:       purchaseID,
			Amount:           allocReq.Amount,
			CreatedAt:        time.Now(),
		}

		if err := s.repo.CreateSupplierPaymentAllocation(ctx, schema, allocation); err != nil {
			return nil, fmt.Errorf("failed to create allocation: %w", err)
		}

		// Calculate total allocated to purchase
		totalAllocated, err := s.repo.GetTotalAllocatedToPurchase(ctx, schema, purchaseID)
		if err != nil {
			return nil, fmt.Errorf("failed to get total allocated: %w", err)
		}

		// Update purchase payment status
		var newStatus string
		if totalAllocated >= purchase.Total {
			newStatus = string(PaymentStatusPaid)
		} else if totalAllocated > 0 {
			newStatus = string(PaymentStatusPartial)
		} else {
			newStatus = string(PaymentStatusPending)
		}

		if err := s.repo.UpdatePurchasePaymentStatus(ctx, schema, purchaseID, newStatus); err != nil {
			return nil, fmt.Errorf("failed to update purchase payment status: %w", err)
		}
	}

	// Get allocations for response
	allocations, err := s.repo.GetSupplierPaymentAllocations(ctx, schema, payment.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get allocations: %w", err)
	}

	// Generate journal entry for supplier payment
	if err := s.generateSupplierPaymentJournalEntry(ctx, schema, payment, allocations); err != nil {
		logger.Warn("Failed to generate journal entry for supplier payment", "error", err, "payment_id", payment.ID)
		// Don't fail the payment creation if journal entry fails
	}

	// Get supplier name
	supplierName, err := s.repo.GetSupplierName(ctx, schema, supplierID)
	if err != nil {
		logger.Warn("Failed to get supplier name", "error", err)
	}

	response := &SupplierPaymentResponse{
		SupplierPayment: *payment,
		SupplierName:    &supplierName,
		Allocations:     allocations,
	}

	logger.Info("Supplier payment created", "payment_number", paymentNumber, "supplier_id", supplierID)
	return response, nil
}

// GetSupplierPaymentByID gets a supplier payment by ID
func (s *Service) GetSupplierPaymentByID(ctx context.Context, schema string, id string) (*SupplierPaymentResponse, error) {
	paymentID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid payment ID: %w", err)
	}

	payment, err := s.repo.GetSupplierPaymentByID(ctx, schema, paymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	// Get allocations
	allocations, err := s.repo.GetSupplierPaymentAllocations(ctx, schema, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get allocations: %w", err)
	}

	// Get supplier name
	supplierName, err := s.repo.GetSupplierName(ctx, schema, payment.SupplierID)
	if err != nil {
		logger.Warn("Failed to get supplier name", "error", err)
	}

	response := &SupplierPaymentResponse{
		SupplierPayment: *payment,
		SupplierName:    &supplierName,
		Allocations:     allocations,
	}

	return response, nil
}

// ListSupplierPayments lists supplier payments
func (s *Service) ListSupplierPayments(ctx context.Context, schema string, supplierID *string, limit, offset int) ([]*SupplierPaymentResponse, error) {
	var supplierUUID *uuid.UUID
	if supplierID != nil {
		parsed, err := uuid.Parse(*supplierID)
		if err != nil {
			return nil, fmt.Errorf("invalid supplier ID: %w", err)
		}
		supplierUUID = &parsed
	}

	payments, err := s.repo.ListSupplierPayments(ctx, schema, supplierUUID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}

	responses := make([]*SupplierPaymentResponse, len(payments))
	for i, payment := range payments {
		// Get allocations
		allocations, err := s.repo.GetSupplierPaymentAllocations(ctx, schema, payment.ID)
		if err != nil {
			logger.Warn("Failed to get allocations", "error", err, "payment_id", payment.ID)
			allocations = []*SupplierPaymentAllocation{}
		}

		// Get supplier name
		supplierName, err := s.repo.GetSupplierName(ctx, schema, payment.SupplierID)
		if err != nil {
			logger.Warn("Failed to get supplier name", "error", err)
		}

		responses[i] = &SupplierPaymentResponse{
			SupplierPayment: *payment,
			SupplierName:    &supplierName,
			Allocations:     allocations,
		}
	}

	return responses, nil
}

// GetPurchaseAllocations gets all payment allocations for a purchase
func (s *Service) GetPurchaseAllocations(ctx context.Context, schema string, purchaseID string) ([]*SupplierPaymentAllocation, error) {
	id, err := uuid.Parse(purchaseID)
	if err != nil {
		return nil, fmt.Errorf("invalid purchase ID: %w", err)
	}

	return s.repo.GetPurchaseAllocations(ctx, schema, id)
}

// generateSupplierPaymentJournalEntry creates accounting entries for a supplier payment
func (s *Service) generateSupplierPaymentJournalEntry(ctx context.Context, schema string, payment *SupplierPayment, allocations []*SupplierPaymentAllocation) error {
	// Get account mappings
	accountsPayableAccount, err := s.accountingSvc.GetAccountForTransaction(ctx, schema, accounting.TransactionTypePurchaseAP, "2111", "Proveedores")
	if err != nil {
		return fmt.Errorf("failed to get accounts payable account: %w", err)
	}

	// Get payment account based on payment method
	var paymentAccount *accounting.ChartOfAccounts
	var transactionType string
	
	switch payment.PaymentMethod {
	case string(PaymentMethodCash):
		transactionType = accounting.TransactionTypeSupplierPaymentCash
		paymentAccount, err = s.accountingSvc.GetAccountForTransaction(ctx, schema, transactionType, "1111", "Caja General")
	case string(PaymentMethodBank):
		transactionType = accounting.TransactionTypeSupplierPaymentBank
		paymentAccount, err = s.accountingSvc.GetAccountForTransaction(ctx, schema, transactionType, "1112", "Banco - Cuenta Corriente")
	case string(PaymentMethodCheck):
		// Cheques también van a banco
		transactionType = accounting.TransactionTypeSupplierPaymentBank
		paymentAccount, err = s.accountingSvc.GetAccountForTransaction(ctx, schema, transactionType, "1112", "Banco - Cuenta Corriente")
	default:
		// Por defecto, usar caja
		transactionType = accounting.TransactionTypeSupplierPaymentCash
		paymentAccount, err = s.accountingSvc.GetAccountForTransaction(ctx, schema, transactionType, "1111", "Caja General")
	}
	
	if err != nil {
		return fmt.Errorf("failed to get payment account: %w", err)
	}

	// Get next entry number
	entryNumber, err := s.accountingRepo.GetNextEntryNumber(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to get next journal entry number: %w", err)
	}

	// Build description with purchase numbers
	var purchaseNumbers []string
	for _, alloc := range allocations {
		purchase, err := s.repo.GetByID(ctx, schema, alloc.PurchaseID)
		if err == nil {
			purchaseNumbers = append(purchaseNumbers, purchase.PurchaseNumber)
		}
	}
	
	description := fmt.Sprintf("Pago a proveedor - %s", payment.PaymentNumber)
	if len(purchaseNumbers) > 0 {
		description += fmt.Sprintf(" (Compras: %s)", purchaseNumbers[0])
		if len(purchaseNumbers) > 1 {
			description += fmt.Sprintf(" y %d más", len(purchaseNumbers)-1)
		}
	}

	entry := &accounting.JournalEntry{
		ID:          uuid.New(),
		EntryNumber: entryNumber,
		EntryDate:   payment.PaymentDate,
		Description: description,
		Reference:   stringPtr(payment.PaymentNumber),
		Status:      "posted",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// DEBE: Proveedores (2111) - reduce la cuenta por pagar
	entry.Lines = append(entry.Lines, accounting.JournalEntryLine{
		ID:             uuid.New(),
		JournalEntryID: entry.ID,
		AccountCode:    accountsPayableAccount.AccountCode,
		AccountName:    accountsPayableAccount.AccountName,
		Debit:          payment.Amount,
		Credit:         0,
		Description:    stringPtr(fmt.Sprintf("Pago a proveedor - %s", payment.PaymentNumber)),
		CreatedAt:      time.Now(),
	})

	// HABER: Caja/Banco (1111 o 1112) según método de pago
	entry.Lines = append(entry.Lines, accounting.JournalEntryLine{
		ID:             uuid.New(),
		JournalEntryID: entry.ID,
		AccountCode:    paymentAccount.AccountCode,
		AccountName:    paymentAccount.AccountName,
		Debit:          0,
		Credit:         payment.Amount,
		Description:    stringPtr(fmt.Sprintf("Pago a proveedor - %s", payment.PaymentNumber)),
		CreatedAt:      time.Now(),
	})

	// Save journal entry
	if err := s.accountingRepo.CreateJournalEntry(ctx, schema, entry); err != nil {
		return fmt.Errorf("failed to create journal entry: %w", err)
	}

	logger.Info("Journal entry created for supplier payment", "payment_number", payment.PaymentNumber, "entry_number", entryNumber)
	return nil
}

// Helper function
func stringPtr(s string) *string {
	return &s
}


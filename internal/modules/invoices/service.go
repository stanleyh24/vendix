package invoices

import (
	"context"
	"fmt"
	"time"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/logger"
	"vendix/internal/modules/products"

	"github.com/google/uuid"
)

type Service struct {
	repo         *Repository
	productsRepo *products.Repository
	cfg          *config.Config
}

func NewService(db *database.DB, cfg *config.Config) *Service {
	return &Service{
		repo:         NewRepository(db),
		productsRepo: products.NewRepository(db),
		cfg:          cfg,
	}
}

func (s *Service) Create(ctx context.Context, schema string, req *CreateInvoiceRequest) (*Invoice, error) {
	// Parse dates
	issueDate, err := time.Parse("2006-01-02", req.IssueDate)
	if err != nil {
		return nil, fmt.Errorf("invalid issue date: %w", err)
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return nil, fmt.Errorf("invalid due date: %w", err)
	}

	// Si no se especifica customer_id, usar el cliente genérico
	customerID := req.CustomerID
	isGenericCustomer := false
	if customerID == nil {
		// UUID del cliente genérico (debe existir en cada tenant)
		genericCustomerID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		customerID = &genericCustomerID
		isGenericCustomer = true
		logger.Info("Using generic customer for invoice")
	} else if *customerID == uuid.MustParse("00000000-0000-0000-0000-000000000001") {
		isGenericCustomer = true
	}

	// Generate invoice number
	invoiceNumber, err := s.repo.GetNextInvoiceNumber(ctx, schema, "INV-")
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice number: %w", err)
	}

	// Validar y establecer tipo de NCF
	ncfType := req.NCFType
	if ncfType == "" {
		ncfType = "02" // Por defecto: Consumidor Final
	}
	// Validar que sea un tipo válido
	if ncfType != "01" && ncfType != "02" {
		return nil, fmt.Errorf("invalid ncf_type: must be '01' (Crédito Fiscal) or '02' (Consumidor Final)")
	}
	// Validar que el crédito fiscal solo se use con clientes que no sean genéricos
	if ncfType == "01" && isGenericCustomer {
		return nil, fmt.Errorf("NCF tipo '01' (Crédito Fiscal) requiere un cliente específico con RNC, no se puede usar con cliente genérico")
	}

	// Determinar el status de la factura
	status := "draft" // Por defecto
	if req.Status != nil {
		// Validar que el status sea válido
		if !InvoiceStatus(*req.Status).IsValid() {
			return nil, fmt.Errorf("invalid status: %s. Valid statuses are: draft, pending, paid, sent, cancelled, overdue", *req.Status)
		}
		status = *req.Status
	}

	// Calculate totals
	subtotal := 0.0
	taxAmount := 0.0

	invoice := &Invoice{
		ID:            uuid.New(),
		InvoiceNumber: invoiceNumber,
		NCFType:       ncfType,
		CustomerID:    *customerID,
		IssueDate:     issueDate,
		DueDate:       dueDate,
		Status:        status,
		Currency:      "DOP",
		PaidAmount:    0,
		Notes:         req.Notes,
		Terms:         req.Terms,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Create lines
	invoice.Lines = make([]InvoiceLine, len(req.Lines))
	for i, lineReq := range req.Lines {
		lineSubtotal := lineReq.Quantity * lineReq.UnitPrice
		lineTax := lineSubtotal * lineReq.TaxRate
		lineTotal := lineSubtotal + lineTax

		subtotal += lineSubtotal
		taxAmount += lineTax

		invoice.Lines[i] = InvoiceLine{
			ID:          uuid.New(),
			InvoiceID:   invoice.ID,
			ProductID:   lineReq.ProductID,
			LineNumber:  i + 1,
			Description: lineReq.Description,
			Quantity:    lineReq.Quantity,
			UnitPrice:   lineReq.UnitPrice,
			TaxRate:     lineReq.TaxRate,
			TaxAmount:   lineTax,
			LineTotal:   lineTotal,
			CreatedAt:   time.Now(),
		}
	}

	invoice.Subtotal = subtotal
	invoice.TaxAmount = taxAmount
	invoice.Total = subtotal + taxAmount

	// Crear la factura
	if err := s.repo.Create(ctx, schema, invoice); err != nil {
		logger.Error("Failed to create invoice", "error", err)
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// Reducir stock de productos solo si el tipo de producto es "product" (no servicios)
	for _, line := range invoice.Lines {
		if line.ProductID != nil {
			// Obtener información del producto
			product, err := s.productsRepo.GetByID(ctx, schema, *line.ProductID)
			if err != nil {
				logger.Error("Failed to get product", "error", err, "product_id", line.ProductID)
				continue // No fallar la factura si hay error al obtener el producto
			}

			// Solo reducir stock para productos físicos, no para servicios
			if product.ProductType == "product" {
				// Validar que hay suficiente stock
				if product.StockQuantity < line.Quantity {
					logger.Warn("Insufficient stock", "product", product.Name, "available", product.StockQuantity, "requested", line.Quantity)
					return nil, fmt.Errorf("stock insuficiente para %s: disponible %.2f, solicitado %.2f", product.Name, product.StockQuantity, line.Quantity)
				}

				// Obtener stock actual antes de actualizar
				previousStock := product.StockQuantity

				// Reducir stock (cantidad negativa)
				if err := s.productsRepo.UpdateStock(ctx, schema, *line.ProductID, -line.Quantity); err != nil {
					logger.Error("Failed to update product stock", "error", err, "product_id", line.ProductID)
					return nil, fmt.Errorf("error al actualizar stock del producto: %w", err)
				}

				// Registrar movimiento de inventario
				refType := "invoice"
				newStock := previousStock - line.Quantity
				notes := fmt.Sprintf("Venta registrada en factura %s", invoice.InvoiceNumber)
				if err := s.productsRepo.RecordInventoryMovement(ctx, schema, *line.ProductID, "sale", line.Quantity, previousStock, newStock, &refType, &invoice.ID, &notes, invoice.CreatedBy); err != nil {
					logger.Error("Failed to record inventory movement", "error", err)
					// No fallar la operación si falla el registro de auditoría
				}
			}
		}
	}

	logger.Info("Invoice created", "number", invoice.InvoiceNumber)
	return invoice, nil
}

func (s *Service) GetByID(ctx context.Context, schema string, id string) (*Invoice, error) {
	invoiceID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid invoice ID: %w", err)
	}

	invoice, err := s.repo.GetByID(ctx, schema, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	return invoice, nil
}

func (s *Service) List(ctx context.Context, schema string, status *string) ([]*Invoice, error) {
	invoices, err := s.repo.List(ctx, schema, status)
	if err != nil {
		logger.Error("Failed to list invoices", "error", err)
		return nil, fmt.Errorf("failed to list invoices: %w", err)
	}

	if invoices == nil {
		invoices = []*Invoice{}
	}

	return invoices, nil
}

func (s *Service) Update(ctx context.Context, schema string, id string, req *UpdateInvoiceRequest) (*Invoice, error) {
	invoiceID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid invoice ID: %w", err)
	}

	invoice, err := s.repo.GetByID(ctx, schema, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	if req.Status != nil {
		invoice.Status = *req.Status
	}
	if req.Notes != nil {
		invoice.Notes = req.Notes
	}
	if req.Terms != nil {
		invoice.Terms = req.Terms
	}

	invoice.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, schema, invoice); err != nil {
		logger.Error("Failed to update invoice", "error", err)
		return nil, fmt.Errorf("failed to update invoice: %w", err)
	}

	logger.Info("Invoice updated", "id", id)
	return invoice, nil
}

func (s *Service) SendInvoice(ctx context.Context, schema string, id string) error {
	invoiceID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid invoice ID: %w", err)
	}

	// Update status to sent
	if err := s.repo.UpdateStatus(ctx, schema, invoiceID, "sent"); err != nil {
		return fmt.Errorf("failed to update invoice status: %w", err)
	}

	logger.Info("Invoice sent", "id", id)
	return nil
}

func (s *Service) CancelInvoice(ctx context.Context, schema string, id string) error {
	invoiceID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid invoice ID: %w", err)
	}

	// Update status to cancelled
	if err := s.repo.UpdateStatus(ctx, schema, invoiceID, "cancelled"); err != nil {
		return fmt.Errorf("failed to cancel invoice: %w", err)
	}

	logger.Info("Invoice cancelled", "id", id)
	return nil
}

package returns

import (
	"context"
	"fmt"
	"time"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/logger"
	"vendix/internal/modules/accounting"
	"vendix/internal/modules/creditnotes"
	"vendix/internal/modules/products"
	"vendix/internal/modules/sales"

	"github.com/google/uuid"
)

type Service struct {
	repo           *Repository
	productsRepo   *products.Repository
	salesRepo      *sales.Repository
	accountingSvc  *accounting.Service
	creditNotesSvc *creditnotes.Service
	cfg            *config.Config
}

func NewService(db *database.DB, cfg *config.Config) *Service {
	creditNotesSvc, _ := creditnotes.NewService(db, cfg) // May be nil if storage fails, but we can still proceed
	return &Service{
		repo:           NewRepository(db),
		productsRepo:   products.NewRepository(db),
		salesRepo:      sales.NewRepository(db),
		accountingSvc:  accounting.NewService(db, cfg),
		creditNotesSvc: creditNotesSvc,
		cfg:            cfg,
	}
}

// Create creates a new return
func (s *Service) Create(ctx context.Context, schema string, req *CreateReturnRequest, createdBy *uuid.UUID) (*Return, error) {
	// Validate that either sale_id or invoice_id is provided
	if req.SaleID == nil && req.InvoiceID == nil {
		return nil, fmt.Errorf("either sale_id or invoice_id must be provided")
	}
	if req.SaleID != nil && req.InvoiceID != nil {
		return nil, fmt.Errorf("only one of sale_id or invoice_id can be provided")
	}

	if len(req.Lines) == 0 {
		return nil, fmt.Errorf("at least one line must be provided")
	}

	var returnType ReturnType
	var customerID *uuid.UUID
	var sale *sales.Sale

	// Get original sale or invoice and validate
	if req.SaleID != nil {
		returnType = ReturnTypeSale
		saleObj, err := s.salesRepo.GetByID(ctx, schema, *req.SaleID)
		if err != nil {
			return nil, fmt.Errorf("sale not found: %w", err)
		}
		sale = saleObj

		if sale.Status != string(sales.SaleStatusCompleted) {
			return nil, fmt.Errorf("can only return completed sales")
		}
		customerID = sale.CustomerID
	} else {
		returnType = ReturnTypeInvoice
		// TODO: Add invoice validation when invoice service is available
		return nil, fmt.Errorf("invoice returns not yet implemented")
	}

	// Generate return number
	returnNumber, err := s.repo.GetNextReturnNumber(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to generate return number: %w", err)
	}

	// Validate and build return lines
	subtotal := 0.0
	taxAmount := 0.0
	returnLines := make([]ReturnLine, 0, len(req.Lines))

	for i, lineReq := range req.Lines {
		// Validate that line reference is provided
		if returnType == ReturnTypeSale && lineReq.SaleLineID == nil {
			return nil, fmt.Errorf("sale_line_id is required for line %d", i+1)
		}
		if returnType == ReturnTypeInvoice && lineReq.InvoiceLineID == nil {
			return nil, fmt.Errorf("invoice_line_id is required for line %d", i+1)
		}

		// Get original line information
		var originalLine *SaleLineInfo
		var invoiceLine *InvoiceLineInfo

		if returnType == ReturnTypeSale {
			if lineReq.SaleLineID == nil {
				return nil, fmt.Errorf("sale_line_id is required")
			}
			lineInfo, err := s.repo.GetSaleLineByID(ctx, schema, *lineReq.SaleLineID)
			if err != nil {
				return nil, fmt.Errorf("sale line not found: %w", err)
			}

			// Validate sale line belongs to the sale
			if lineInfo.SaleID != *req.SaleID {
				return nil, fmt.Errorf("sale line does not belong to the specified sale")
			}

			// Validate quantity doesn't exceed available (quantity - returned_quantity)
			availableQty := lineInfo.Quantity - lineInfo.ReturnedQuantity
			if lineReq.Quantity > availableQty {
				return nil, fmt.Errorf("return quantity (%.2f) exceeds available quantity (%.2f) for line", lineReq.Quantity, availableQty)
			}

			originalLine = lineInfo
		} else {
			if lineReq.InvoiceLineID == nil {
				return nil, fmt.Errorf("invoice_line_id is required")
			}
			lineInfo, err := s.repo.GetInvoiceLineByID(ctx, schema, *lineReq.InvoiceLineID)
			if err != nil {
				return nil, fmt.Errorf("invoice line not found: %w", err)
			}

			availableQty := lineInfo.Quantity - lineInfo.ReturnedQuantity
			if lineReq.Quantity > availableQty {
				return nil, fmt.Errorf("return quantity (%.2f) exceeds available quantity (%.2f) for line", lineReq.Quantity, availableQty)
			}

			invoiceLine = lineInfo
		}

		// Use original line prices and calculate totals
		var unitPrice, taxRate float64
		var productID *uuid.UUID
		var description string

		if originalLine != nil {
			unitPrice = originalLine.UnitPrice
			taxRate = originalLine.TaxRate
			productID = originalLine.ProductID
			// Get description from original sale line
			for _, sl := range sale.Lines {
				if sl.ID == *lineReq.SaleLineID {
					description = sl.Description
					break
				}
			}
		} else {
			unitPrice = invoiceLine.UnitPrice
			taxRate = invoiceLine.TaxRate
			productID = invoiceLine.ProductID
			// TODO: Get description from invoice line
			description = fmt.Sprintf("Product return")
		}

		// Calculate line totals
		lineSubtotal := lineReq.Quantity * unitPrice
		lineTax := lineSubtotal * taxRate
		lineTotal := lineSubtotal + lineTax

		subtotal += lineSubtotal
		taxAmount += lineTax

		// Determine restock quantity based on item condition
		restockQty := lineReq.Quantity
		if lineReq.ItemCondition != nil {
			switch *lineReq.ItemCondition {
			case ItemConditionDefective, ItemConditionDamaged:
				// Only restock if restock_quantity is specified
				if lineReq.RestockQuantity != nil {
					restockQty = *lineReq.RestockQuantity
				} else {
					restockQty = 0 // Don't restock defective/damaged items by default
				}
			}
		}

		// Validate product is returnable (if product exists)
		if productID != nil {
			product, err := s.productsRepo.GetByID(ctx, schema, *productID)
			if err == nil {
				// Check if product has is_returnable field (will be nil if column doesn't exist yet)
				// For now, we'll allow all returns, but this can be enhanced
				_ = product
			}
		}

		returnLines = append(returnLines, ReturnLine{
			ID:             uuid.New(),
			ReturnID:       uuid.Nil, // Will be set after return is created
			SaleLineID:     lineReq.SaleLineID,
			InvoiceLineID:  lineReq.InvoiceLineID,
			ProductID:      productID,
			LineNumber:     i + 1,
			Description:    description,
			Quantity:       lineReq.Quantity,
			UnitPrice:      unitPrice,
			TaxRate:        taxRate,
			TaxAmount:      lineTax,
			LineTotal:      lineTotal,
			ItemCondition:  lineReq.ItemCondition,
			RestockQuantity: &restockQty,
			CreatedAt:      time.Now(),
		})
	}

	// Create return
	ret := &Return{
		ID:           uuid.New(),
		ReturnNumber: returnNumber,
		SaleID:       req.SaleID,
		InvoiceID:    req.InvoiceID,
		CustomerID:   customerID,
		ReturnDate:   time.Now(),
		ReturnType:   returnType,
		ReturnReason: req.ReturnReason,
		Notes:        req.Notes,
		Subtotal:     subtotal,
		TaxAmount:    taxAmount,
		Total:        subtotal + taxAmount,
		Status:       ReturnStatusPending,
		CreatedBy:    createdBy,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Lines:        returnLines,
	}

	// Set ReturnID for all lines
	for i := range ret.Lines {
		ret.Lines[i].ReturnID = ret.ID
	}

	// Save to database
	if err := s.repo.Create(ctx, schema, ret); err != nil {
		return nil, fmt.Errorf("failed to create return: %w", err)
	}

	logger.Info("Return created", "return_number", returnNumber, "schema", schema)
	return ret, nil
}

// Approve approves a return
func (s *Service) Approve(ctx context.Context, schema string, id uuid.UUID, approvedBy uuid.UUID) error {
	ret, err := s.repo.GetByID(ctx, schema, id)
	if err != nil {
		return fmt.Errorf("return not found: %w", err)
	}

	if ret.Status != ReturnStatusPending {
		return fmt.Errorf("only pending returns can be approved")
	}

	// Update status
	if err := s.repo.UpdateStatus(ctx, schema, id, ReturnStatusApproved, &approvedBy); err != nil {
		return fmt.Errorf("failed to approve return: %w", err)
	}

	logger.Info("Return approved", "return_id", id, "schema", schema)
	return nil
}

// Reject rejects a return
func (s *Service) Reject(ctx context.Context, schema string, id uuid.UUID, req *RejectReturnRequest, rejectedBy uuid.UUID) error {
	ret, err := s.repo.GetByID(ctx, schema, id)
	if err != nil {
		return fmt.Errorf("return not found: %w", err)
	}

	if ret.Status != ReturnStatusPending {
		return fmt.Errorf("only pending returns can be rejected")
	}

	// Add rejection reason to notes
	notes := req.Reason
	if ret.Notes != nil {
		notes = fmt.Sprintf("%s\nRejected: %s", *ret.Notes, req.Reason)
	}

	// Update status
	// Note: Rejection reason is logged but not stored in notes for now
	// TODO: Add UpdateNotes method to repository if needed
	_ = notes
	if err := s.repo.UpdateStatus(ctx, schema, id, ReturnStatusRejected, &rejectedBy); err != nil {
		return fmt.Errorf("failed to reject return: %w", err)
	}

	logger.Info("Return rejected", "return_id", id, "reason", req.Reason, "schema", schema)
	return nil
}

// Process processes an approved return
func (s *Service) Process(ctx context.Context, schema string, id uuid.UUID, req *ProcessReturnRequest, processedBy *uuid.UUID) error {
	ret, err := s.repo.GetByID(ctx, schema, id)
	if err != nil {
		return fmt.Errorf("return not found: %w", err)
	}

	if ret.Status != ReturnStatusApproved {
		return fmt.Errorf("only approved returns can be processed")
	}

	if !req.RefundMethod.IsValid() {
		return fmt.Errorf("invalid refund method: %s", req.RefundMethod)
	}

	// Determine refund amount
	refundAmount := ret.Total
	if req.RefundAmount != nil {
		if *req.RefundAmount > ret.Total {
			return fmt.Errorf("refund amount cannot exceed return total")
		}
		refundAmount = *req.RefundAmount
	}

	// Restore inventory for each line
	for _, line := range ret.Lines {
		if line.ProductID != nil && line.RestockQuantity != nil && *line.RestockQuantity > 0 {
			product, err := s.productsRepo.GetByID(ctx, schema, *line.ProductID)
			if err != nil {
				logger.Warn("Product not found for return line, skipping inventory restore", "product_id", *line.ProductID)
				continue
			}

			// Only restore stock for physical products
			if product.ProductType == "product" {
				previousStock := product.StockQuantity
				if err := s.productsRepo.UpdateStock(ctx, schema, *line.ProductID, *line.RestockQuantity); err != nil {
					logger.Error("Failed to restore product stock", "error", err, "product_id", *line.ProductID)
					return fmt.Errorf("failed to restore inventory: %w", err)
				}

				// Record inventory movement
				refType := "return"
				newStock := previousStock + *line.RestockQuantity
				notes := fmt.Sprintf("Devolución %s - Restock", ret.ReturnNumber)
				if err := s.productsRepo.RecordInventoryMovement(ctx, schema, *line.ProductID, "return", *line.RestockQuantity, previousStock, newStock, &refType, &ret.ID, &notes, processedBy); err != nil {
					logger.Error("Failed to record inventory movement", "error", err)
					// Don't fail the operation if audit fails
				}
			}
		}

		// Update returned_quantity in original sale/invoice line
		if line.SaleLineID != nil {
			if err := s.repo.UpdateReturnedQuantity(ctx, schema, line.SaleLineID, nil, line.Quantity); err != nil {
				return fmt.Errorf("failed to update sale line returned quantity: %w", err)
			}
		} else if line.InvoiceLineID != nil {
			if err := s.repo.UpdateReturnedQuantity(ctx, schema, nil, line.InvoiceLineID, line.Quantity); err != nil {
				return fmt.Errorf("failed to update invoice line returned quantity: %w", err)
			}
		}
	}

	// Generate credit note if refund method is credit_note and return is from an invoice
	if req.RefundMethod == RefundMethodCreditNote && ret.InvoiceID != nil {
		if s.creditNotesSvc != nil {
			// Prepare return lines info for credit note creation
			returnLines := make([]creditnotes.ReturnLineInfo, len(ret.Lines))
			for i, line := range ret.Lines {
				returnLines[i] = creditnotes.ReturnLineInfo{
					InvoiceLineID: line.InvoiceLineID,
					ProductID:     line.ProductID,
					Description:   line.Description,
					Quantity:      line.Quantity,
					UnitPrice:     line.UnitPrice,
					TaxRate:       line.TaxRate,
					TaxAmount:     line.TaxAmount,
					LineTotal:     line.LineTotal,
				}
			}

			// Create credit note from return
			creditNote, err := s.creditNotesSvc.CreateFromReturn(ctx, schema, ret.ID, *ret.InvoiceID, returnLines, processedBy)
			if err != nil {
				logger.Error("Failed to create credit note from return", "error", err, "return_id", id)
				return fmt.Errorf("failed to create credit note: %w", err)
			}

			logger.Info("Credit note created from return", "credit_note_id", creditNote.ID, "return_id", id)

			// Automatically send credit note to DGII
			if err := s.creditNotesSvc.SendCreditNote(ctx, schema, creditNote.ID); err != nil {
				logger.Error("Failed to send credit note to DGII", "error", err, "credit_note_id", creditNote.ID)
				// Don't fail the return processing if DGII send fails
			}
		} else {
			logger.Warn("Credit notes service not available, skipping credit note generation")
		}
	}

	// Generate accounting journal entry
	if err := s.generateReturnJournalEntry(ctx, schema, ret, req.RefundMethod, processedBy); err != nil {
		logger.Warn("Failed to generate journal entry for return", "error", err, "return_id", id)
		// Don't fail the operation if accounting fails, but log it
	}

	// Update refund information
	refundStatus := "completed"
	if err := s.repo.UpdateRefundInfo(ctx, schema, id, &refundStatus, &req.RefundMethod, refundAmount); err != nil {
		return fmt.Errorf("failed to update refund info: %w", err)
	}

	// Update status to completed
	if err := s.repo.UpdateStatus(ctx, schema, id, ReturnStatusCompleted, processedBy); err != nil {
		return fmt.Errorf("failed to update return status: %w", err)
	}

	logger.Info("Return processed", "return_id", id, "refund_method", req.RefundMethod, "refund_amount", refundAmount, "schema", schema)
	return nil
}

// generateReturnJournalEntry creates accounting entries for a return
func (s *Service) generateReturnJournalEntry(ctx context.Context, schema string, ret *Return, refundMethod RefundMethod, createdBy *uuid.UUID) error {
	// Get account mappings
	salesReturnAccount, err := s.accountingSvc.GetAccountForTransaction(ctx, schema, "sales_return", "4111", "Devoluciones y Descuentos en Ventas")
	if err != nil {
		return fmt.Errorf("failed to get sales return account: %w", err)
	}

	taxReturnAccount, err := s.accountingSvc.GetAccountForTransaction(ctx, schema, "tax_return", "2121", "ITBIS por Pagar")
	if err != nil {
		return fmt.Errorf("failed to get tax return account: %w", err)
	}

	var paymentAccount *accounting.ChartOfAccounts
	switch refundMethod {
	case RefundMethodCash:
		paymentAccount, err = s.accountingSvc.GetAccountForTransaction(ctx, schema, "payment_cash", "1111", "Caja General")
	case RefundMethodCard, RefundMethodTransfer:
		paymentAccount, err = s.accountingSvc.GetAccountForTransaction(ctx, schema, "payment_bank", "1112", "Banco - Cuenta Corriente")
	case RefundMethodCreditNote:
		paymentAccount, err = s.accountingSvc.GetAccountForTransaction(ctx, schema, "sales_credit", "1121", "Clientes")
	default:
		paymentAccount, err = s.accountingSvc.GetAccountForTransaction(ctx, schema, "payment_cash", "1111", "Caja General")
	}
	if err != nil {
		return fmt.Errorf("failed to get payment account: %w", err)
	}

	// Build journal entry lines
	var lines []accounting.CreateJournalEntryLineReq

	// DEBE: Devoluciones y Descuentos en Ventas (reducción de ingresos)
	lines = append(lines, accounting.CreateJournalEntryLineReq{
		AccountCode: salesReturnAccount.AccountCode,
		Debit:       ret.Subtotal,
		Credit:      0,
		Description: stringPtr(fmt.Sprintf("Devolución %s", ret.ReturnNumber)),
	})

	// DEBE: ITBIS por Pagar (reversión de impuesto)
	if ret.TaxAmount > 0 {
		lines = append(lines, accounting.CreateJournalEntryLineReq{
			AccountCode: taxReturnAccount.AccountCode,
			Debit:       ret.TaxAmount,
			Credit:      0,
			Description: stringPtr(fmt.Sprintf("Reversión ITBIS - Devolución %s", ret.ReturnNumber)),
		})
	}

	// HABER: Caja/Banco/Clientes (según método de reembolso)
	lines = append(lines, accounting.CreateJournalEntryLineReq{
		AccountCode: paymentAccount.AccountCode,
		Debit:       0,
		Credit:      ret.Total,
		Description: stringPtr(fmt.Sprintf("Reembolso - Devolución %s", ret.ReturnNumber)),
	})

	// Create journal entry
	entryReq := &accounting.CreateJournalEntryRequest{
		EntryDate:   ret.ReturnDate.Format("2006-01-02"),
		Description: fmt.Sprintf("Devolución %s", ret.ReturnNumber),
		Reference:   stringPtr(ret.ReturnNumber),
		Lines:       lines,
	}

	_, err = s.accountingSvc.CreateJournalEntry(ctx, schema, entryReq)
	if err != nil {
		return fmt.Errorf("failed to create journal entry: %w", err)
	}

	return nil
}

// GetByID retrieves a return by ID
func (s *Service) GetByID(ctx context.Context, schema string, id uuid.UUID) (*Return, error) {
	return s.repo.GetByID(ctx, schema, id)
}

// List retrieves all returns with optional filters
func (s *Service) List(ctx context.Context, schema string, filters *ReturnFilters) ([]*Return, error) {
	return s.repo.List(ctx, schema, filters)
}

// Helper function
func stringPtr(s string) *string {
	return &s
}


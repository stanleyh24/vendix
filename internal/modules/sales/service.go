package sales

import (
	"context"
	"fmt"
	"time"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/logger"
	"vendix/internal/modules/customers"
	"vendix/internal/modules/invoices"
	"vendix/internal/modules/products"

	"github.com/google/uuid"
)

type Service struct {
	repo          *Repository
	productsRepo  *products.Repository
	customersRepo *customers.Repository
	invoiceSvc    *invoices.Service
	cfg           *config.Config
}

func NewService(db *database.DB, cfg *config.Config) *Service {
	invoiceSvc, _ := invoices.NewService(db, cfg)
	return &Service{
		repo:          NewRepository(db),
		productsRepo:  products.NewRepository(db),
		customersRepo: customers.NewRepository(db),
		invoiceSvc:    invoiceSvc,
		cfg:           cfg,
	}
}

func (s *Service) Create(ctx context.Context, schema string, req *CreateSaleRequest) (*Sale, error) {
	// Validar tipo de pago
	if !PaymentType(req.PaymentType).IsValid() {
		return nil, fmt.Errorf("invalid payment type: %s. Valid types are: cash, card, transfer, check", req.PaymentType)
	}

	// Si no se especifica customer_id, usar el cliente genérico
	customerID := req.CustomerID
	if customerID == nil {
		// UUID del cliente genérico (debe existir en cada tenant)
		genericCustomerID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		customerID = &genericCustomerID
		logger.Info("Using generic customer for sale")
	}

	// Generate sale number
	saleNumber, err := s.repo.GetNextSaleNumber(ctx, schema, "SALE-")
	if err != nil {
		return nil, fmt.Errorf("failed to generate sale number: %w", err)
	}

	// Calculate totals
	subtotal := 0.0
	taxAmount := 0.0

	sale := &Sale{
		ID:          uuid.New(),
		SaleNumber:  saleNumber,
		CustomerID:  customerID,
		PaymentType: req.PaymentType,
		Status:      string(SaleStatusCompleted), // Las ventas siempre se crean como completadas
		Notes:       req.Notes,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create lines
	sale.Lines = make([]SaleLine, len(req.Lines))
	for i, lineReq := range req.Lines {
		lineSubtotal := lineReq.Quantity * lineReq.UnitPrice
		lineTax := lineSubtotal * lineReq.TaxRate
		lineTotal := lineSubtotal + lineTax

		subtotal += lineSubtotal
		taxAmount += lineTax

		sale.Lines[i] = SaleLine{
			ID:          uuid.New(),
			SaleID:      sale.ID,
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

	sale.Subtotal = subtotal
	sale.TaxAmount = taxAmount
	sale.Total = subtotal + taxAmount

	// Crear la venta
	if err := s.repo.Create(ctx, schema, sale); err != nil {
		logger.Error("Failed to create sale", "error", err)
		return nil, fmt.Errorf("failed to create sale: %w", err)
	}

	// Reducir stock de productos solo si el tipo de producto es "product" (no servicios)
	for _, line := range sale.Lines {
		if line.ProductID != nil {
			// Obtener información del producto
			product, err := s.productsRepo.GetByID(ctx, schema, *line.ProductID)
			if err != nil {
				logger.Error("Failed to get product", "error", err, "product_id", line.ProductID)
				continue // No fallar la venta si hay error al obtener el producto
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
				refType := "sale"
				newStock := previousStock - line.Quantity
				notes := fmt.Sprintf("Venta registrada en venta %s", sale.SaleNumber)
				if err := s.productsRepo.RecordInventoryMovement(ctx, schema, *line.ProductID, "sale", line.Quantity, previousStock, newStock, &refType, &sale.ID, &notes, sale.CreatedBy); err != nil {
					logger.Error("Failed to record inventory movement", "error", err)
					// No fallar la operación si falla el registro de auditoría
				}
			}
		}
	}

	// Obtener información del cliente para validaciones
	var customer *customers.Customer
	isGenericCustomer := *customerID == uuid.MustParse("00000000-0000-0000-0000-000000000001")
	if !isGenericCustomer {
		var err error
		customer, err = s.customersRepo.GetByID(ctx, schema, *customerID)
		if err != nil {
			return nil, fmt.Errorf("failed to get customer: %w", err)
		}
	}

	// Determinar tipo de NCF (del request o por defecto según cliente)
	ncfType := req.NCFType
	if ncfType == "" {
		// Si no se especifica, usar 02 (Consumidor Final) por defecto
		ncfType = "02"
	}
	
	// Validar que sea un tipo válido
	if ncfType != "01" && ncfType != "02" && ncfType != "15" {
		return nil, fmt.Errorf("invalid ncf_type: must be '01' (Crédito Fiscal), '02' (Consumidor Final), or '15' (Gubernamental)")
	}

	// Validar que el crédito fiscal solo se use con clientes específicos (no genéricos)
	if ncfType == "01" && isGenericCustomer {
		return nil, fmt.Errorf("NCF tipo '01' (Crédito Fiscal) requiere un cliente específico con RNC, no se puede usar con cliente genérico")
	}

	// Validar que el NCF tipo 15 (Gubernamental) solo se use con clientes gubernamentales
	if ncfType == "15" {
		if isGenericCustomer {
			return nil, fmt.Errorf("NCF tipo '15' (Gubernamental) requiere un cliente específico que sea entidad gubernamental")
		}
		if customer == nil || !customer.IsGovernmentEntity {
			return nil, fmt.Errorf("NCF tipo '15' (Gubernamental) solo se puede usar con clientes que sean entidades gubernamentales")
		}
		if customer.TaxID == nil || *customer.TaxID == "" {
			return nil, fmt.Errorf("cliente gubernamental debe tener RNC válido para usar NCF tipo '15'")
		}
	}

	// Crear factura automáticamente con status "paid"
	invoiceReq := &invoices.CreateInvoiceRequest{
		CustomerID: customerID,
		NCFType:    ncfType,
		IssueDate:  time.Now().Format("2006-01-02"),
		DueDate:    time.Now().Format("2006-01-02"), // Mismo día para ventas
		Status:     stringPtr("paid"),               // Las facturas de ventas son pagadas automáticamente
		Notes:      req.Notes,
		Lines:      convertSaleLinesToInvoiceLines(req.Lines),
		WithholdingTaxType: req.WithholdingTaxType,
		WithholdingRate:    req.WithholdingRate,
		WithholdingExempt:  req.WithholdingExempt,
	}

	invoice, err := s.invoiceSvc.Create(ctx, schema, invoiceReq)
	if err != nil {
		logger.Error("Failed to create invoice for sale", "error", err)
		return nil, fmt.Errorf("failed to create invoice for sale: %w", err)
	}

	// Actualizar la venta con el ID de la factura
	sale.InvoiceID = &invoice.ID
	if err := s.repo.Update(ctx, schema, sale); err != nil {
		logger.Error("Failed to update sale with invoice ID", "error", err)
		// No fallar la operación si falla la actualización
	}

	logger.Info("Sale created with invoice", "sale_number", sale.SaleNumber, "invoice_number", invoice.InvoiceNumber)
	return sale, nil
}

func (s *Service) GetByID(ctx context.Context, schema string, id string) (*Sale, error) {
	saleID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid sale ID: %w", err)
	}

	sale, err := s.repo.GetByID(ctx, schema, saleID)
	if err != nil {
		return nil, fmt.Errorf("sale not found: %w", err)
	}

	return sale, nil
}

func (s *Service) List(ctx context.Context, schema string, status *string) ([]*Sale, error) {
	sales, err := s.repo.List(ctx, schema, status)
	if err != nil {
		logger.Error("Failed to list sales", "error", err)
		return nil, fmt.Errorf("failed to list sales: %w", err)
	}

	if sales == nil {
		sales = []*Sale{}
	}

	return sales, nil
}

func (s *Service) Update(ctx context.Context, schema string, id string, req *UpdateSaleRequest) (*Sale, error) {
	saleID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid sale ID: %w", err)
	}

	sale, err := s.repo.GetByID(ctx, schema, saleID)
	if err != nil {
		return nil, fmt.Errorf("sale not found: %w", err)
	}

	if req.Status != nil {
		// Validar que el status sea válido
		if !SaleStatus(*req.Status).IsValid() {
			return nil, fmt.Errorf("invalid status: %s. Valid statuses are: completed, cancelled, refunded", *req.Status)
		}
		sale.Status = *req.Status
	}
	if req.Notes != nil {
		sale.Notes = req.Notes
	}

	sale.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, schema, sale); err != nil {
		logger.Error("Failed to update sale", "error", err)
		return nil, fmt.Errorf("failed to update sale: %w", err)
	}

	logger.Info("Sale updated", "id", id)
	return sale, nil
}

// Helper function to convert sale lines to invoice lines
func convertSaleLinesToInvoiceLines(saleLines []CreateSaleLineReq) []invoices.CreateInvoiceLineReq {
	invoiceLines := make([]invoices.CreateInvoiceLineReq, len(saleLines))
	for i, line := range saleLines {
		invoiceLines[i] = invoices.CreateInvoiceLineReq{
			ProductID:   line.ProductID,
			Description: line.Description,
			Quantity:    line.Quantity,
			UnitPrice:   line.UnitPrice,
			TaxRate:     line.TaxRate,
		}
	}
	return invoiceLines
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// UpdateSaleRequest represents the request to update a sale
type UpdateSaleRequest struct {
	Status *string `json:"status,omitempty"`
	Notes  *string `json:"notes,omitempty"`
}

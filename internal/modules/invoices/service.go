package invoices

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/dgii"
	"vendix/internal/logger"
	"vendix/internal/modules/customers"
	"vendix/internal/modules/ncf"
	"vendix/internal/modules/products"
	"vendix/internal/modules/tenantconfig"
	"vendix/internal/storage"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

type Service struct {
	repo          *Repository
	productsRepo  *products.Repository
	customersRepo *customers.Repository
	storage       *storage.Client
	dgiiService   *dgii.Service
	ncfService    *ncf.Service
	configRepo    *tenantconfig.Repository
	cfg           *config.Config
}

func NewService(db *database.DB, cfg *config.Config) (*Service, error) {
	// Storage client will be initialized per-tenant when needed
	// We'll create it with the default bucket for now, but it will be switched per tenant
	storageClient, err := storage.NewClient(cfg)
	if err != nil {
		logger.Error("Failed to initialize storage client", "error", err)
		// Return service without storage if it fails - allows graceful degradation
		// In production, you might want to fail here
	}

	return &Service{
		repo:          NewRepository(db),
		productsRepo:  products.NewRepository(db),
		customersRepo: customers.NewRepository(db),
		storage:       storageClient,
		dgiiService:   dgii.NewService(cfg),
		ncfService:    ncf.NewService(db),
		configRepo:    tenantconfig.NewRepository(db),
		cfg:           cfg,
	}, nil
}

// getStorageClientForTenant returns a storage client configured for the tenant's bucket
func (s *Service) getStorageClientForTenant(ctx context.Context, schema string) (*storage.Client, error) {
	// Build bucket name for tenant
	bucketName := storage.BuildTenantBucketName(schema)
	
	// Create or get client with tenant bucket
	// If we already have a client, we can switch its bucket
	if s.storage != nil {
		// Create a new client with the tenant's bucket
		tenantClient, err := storage.NewClientWithBucket(s.cfg, bucketName)
		if err != nil {
			return nil, fmt.Errorf("failed to create storage client for tenant: %w", err)
		}
		return tenantClient, nil
	}
	
	return nil, fmt.Errorf("storage client not initialized")
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

	// Obtener información del cliente para validaciones
	var customer *customers.Customer
	if !isGenericCustomer {
		var err error
		customer, err = s.customersRepo.GetByID(ctx, schema, *customerID)
		if err != nil {
			return nil, fmt.Errorf("failed to get customer: %w", err)
		}
	}

	// Validar y establecer tipo de NCF
	ncfType := req.NCFType
	if ncfType == "" {
		ncfType = "02" // Por defecto: Consumidor Final
	}
	// Validar que sea un tipo válido
	if ncfType != "01" && ncfType != "02" && ncfType != "15" {
		return nil, fmt.Errorf("invalid ncf_type: must be '01' (Crédito Fiscal), '02' (Consumidor Final), or '15' (Gubernamental)")
	}
	// Validar que el crédito fiscal solo se use con clientes que no sean genéricos
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

	// Calcular retenciones si aplica
	invoice.WithholdingExempt = false
	if req.WithholdingExempt != nil {
		invoice.WithholdingExempt = *req.WithholdingExempt
	}

	// Calcular retención si el cliente es gubernamental y no está exento
	if !isGenericCustomer && customer != nil && customer.IsGovernmentEntity && !invoice.WithholdingExempt && ncfType == "15" {
		withholdingRate := 0.05 // 5% ISR por defecto
		
		// Usar tasa del request si está especificada
		if req.WithholdingRate != nil && *req.WithholdingRate > 0 {
			withholdingRate = *req.WithholdingRate
		} else if customer.DefaultWithholdingRate != nil && *customer.DefaultWithholdingRate > 0 {
			// Usar tasa por defecto del cliente
			withholdingRate = *customer.DefaultWithholdingRate
		}

		withholdingAmount := invoice.Total * withholdingRate
		invoice.WithholdingTaxAmount = &withholdingAmount
		taxType := "isr"
		if req.WithholdingTaxType != nil && *req.WithholdingTaxType != "" {
			taxType = *req.WithholdingTaxType
		}
		invoice.WithholdingTaxType = &taxType
		invoice.WithholdingRate = &withholdingRate
		invoice.NetAmount = invoice.Total - withholdingAmount
	} else {
		// Si no hay retención, el monto neto es igual al total
		invoice.NetAmount = invoice.Total
	}

	// Generate NCF if enabled (NCF is auto-generated, not provided in request)
	if true {
		// Determine NCF type based on customer
		ncfTypeStr := ncfType
		ncfTypeEnum, err := ncf.GetNCFTypeFromString(ncfTypeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid NCF type: %w", err)
		}

		// Generate NCF
		ncfInfo, err := s.ncfService.GenerateNCF(ctx, schema, ncfTypeEnum)
		if err != nil {
			return nil, fmt.Errorf("failed to generate NCF: %w", err)
		}

		ncfStr := ncfInfo.FullNCF
		invoice.NCF = &ncfStr
		logger.Info("NCF generated for invoice", "ncf", ncfStr, "type", ncfTypeStr)
	}

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

	// Generate and save PDF to storage
	if s.storage != nil {
		if err := s.generateAndSavePDF(ctx, schema, invoice); err != nil {
			logger.Error("Failed to generate and save PDF", "error", err, "invoice_id", invoice.ID)
			// Don't fail invoice creation if PDF generation fails
		}
	} else {
		logger.Warn("Storage client not available, PDF will not be generated", "invoice_id", invoice.ID)
	}

	return invoice, nil
}

// generateAndSavePDF generates a PDF for the invoice and saves it to storage
func (s *Service) generateAndSavePDF(ctx context.Context, schema string, invoice *Invoice) error {
	// Get storage client for tenant bucket
	storageClient, err := s.getStorageClientForTenant(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to get storage client: %w", err)
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Factura "+invoice.InvoiceNumber)
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Cliente: %s", invoice.CustomerName))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Fecha: %s", invoice.IssueDate.Format("2006-01-02")))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Vence: %s", invoice.DueDate.Format("2006-01-02")))
	pdf.Ln(10)

	// Table header
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(100, 8, "Descripción", "B", 0, "L", false, 0, "")
	pdf.CellFormat(20, 8, "Cant.", "B", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, "P. Unit.", "B", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, "Total", "B", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	for _, line := range invoice.Lines {
		pdf.CellFormat(100, 7, line.Description, "", 0, "L", false, 0, "")
		pdf.CellFormat(20, 7, fmt.Sprintf("%.2f", line.Quantity), "", 0, "R", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", line.UnitPrice), "", 0, "R", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", line.LineTotal), "", 1, "R", false, 0, "")
	}

	pdf.Ln(4)
	pdf.Cell(0, 0, "")
	pdf.Ln(2)
	pdf.CellFormat(150, 7, "Subtotal:", "", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", invoice.Subtotal), "", 1, "R", false, 0, "")
	pdf.CellFormat(150, 7, "ITBIS:", "", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", invoice.TaxAmount), "", 1, "R", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(150, 8, "TOTAL:", "", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, fmt.Sprintf("%.2f", invoice.Total), "", 1, "R", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Build object key
	objectKey := storage.BuildInvoicePDFKey(schema, invoice.ID.String(), invoice.InvoiceNumber)

	// Upload to storage
	if err := storageClient.UploadBytes(ctx, buf.Bytes(), objectKey, "application/pdf"); err != nil {
		return fmt.Errorf("failed to upload PDF: %w", err)
	}

	// Generate presigned URL (valid for 7 days)
	pdfURL, err := storageClient.PresignedURL(ctx, objectKey, 7*24*time.Hour)
	if err != nil {
		logger.Error("Failed to generate presigned URL", "error", err)
		// Still update with object key even if presigned URL fails
		pdfURL = fmt.Sprintf("s3://%s/%s", s.cfg.S3Bucket, objectKey)
	}

	// Update invoice with PDF URL
	if err := s.repo.UpdateDocumentURLs(ctx, schema, invoice.ID, &pdfURL, nil); err != nil {
		logger.Error("Failed to update PDF URL", "error", err)
		return fmt.Errorf("failed to update PDF URL: %w", err)
	}

	logger.Info("PDF generated and saved", "invoice_id", invoice.ID, "object_key", objectKey)
	return nil
}

// GetPDFBytes retrieves the PDF bytes for an invoice
func (s *Service) GetPDFBytes(ctx context.Context, schema string, invoiceID uuid.UUID) ([]byte, error) {
	invoice, err := s.repo.GetByID(ctx, schema, invoiceID, []string{})
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	if invoice.PDFURL == nil || *invoice.PDFURL == "" {
		// PDF doesn't exist, generate it
		if err := s.generateAndSavePDF(ctx, schema, invoice); err != nil {
			return nil, fmt.Errorf("failed to generate PDF: %w", err)
		}
		// Reload invoice to get the updated PDF URL
		invoice, err = s.repo.GetByID(ctx, schema, invoiceID, []string{})
		if err != nil {
			return nil, fmt.Errorf("failed to reload invoice: %w", err)
		}
	}

	// Get storage client for tenant bucket
	storageClient, err := s.getStorageClientForTenant(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage client: %w", err)
	}

	// Extract object key from URL (if it's a presigned URL, we need to extract the key)
	objectKey := storage.BuildInvoicePDFKey(schema, invoiceID.String(), invoice.InvoiceNumber)

	// Download from storage
	pdfBytes, err := storageClient.DownloadFile(ctx, objectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to download PDF: %w", err)
	}

	return pdfBytes, nil
}

// GetXMLBytes retrieves the XML bytes for an invoice
func (s *Service) GetXMLBytes(ctx context.Context, schema string, invoiceID uuid.UUID) ([]byte, error) {
	invoice, err := s.repo.GetByID(ctx, schema, invoiceID, []string{})
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	if invoice.XMLURL == nil || *invoice.XMLURL == "" {
		// XML doesn't exist - invoice hasn't been sent to DGII yet
		return nil, fmt.Errorf("XML not available: invoice has not been sent to DGII")
	}

	// Get storage client for tenant bucket
	storageClient, err := s.getStorageClientForTenant(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage client: %w", err)
	}

	// Extract object key from URL (if it's a presigned URL, we need to extract the key)
	objectKey := storage.BuildInvoiceXMLKey(schema, invoiceID.String(), invoice.InvoiceNumber)

	// Download from storage
	xmlBytes, err := storageClient.DownloadFile(ctx, objectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to download XML: %w", err)
	}

	return xmlBytes, nil
}

func (s *Service) GetByID(ctx context.Context, schema string, id string, includes []string) (*Invoice, error) {
	invoiceID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid invoice ID: %w", err)
	}

	invoice, err := s.repo.GetByID(ctx, schema, invoiceID, includes)
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	return invoice, nil
}

func (s *Service) List(ctx context.Context, schema string, status *string, includes []string) ([]*Invoice, error) {
	invoices, err := s.repo.List(ctx, schema, status, includes)
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

	invoice, err := s.repo.GetByID(ctx, schema, invoiceID, []string{})
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

	// Get invoice
	invoice, err := s.repo.GetByID(ctx, schema, invoiceID, []string{})
	if err != nil {
		return fmt.Errorf("invoice not found: %w", err)
	}

	// Validate invoice has NCF
	if invoice.NCF == nil || *invoice.NCF == "" {
		return fmt.Errorf("invoice must have an NCF before sending to DGII")
	}

	// Convert invoice to ElectronicDocument for DGII
	document := s.invoiceToElectronicDocument(ctx, schema, invoice)

	// Process invoice with DGII (sign and send) - get both response and signed document
	processResult, err := s.dgiiService.ProcessInvoiceWithSignedDoc(ctx, document)
	if err != nil {
		logger.Error("Failed to process invoice with DGII", "error", err, "invoice_id", id)
		return fmt.Errorf("failed to send invoice to DGII: %w", err)
	}

	dgiiResponse := processResult.Response
	signedDoc := processResult.SignedDoc

	// Use the signed XML from the adapter if available, otherwise generate our own
	var xmlData []byte
	if signedDoc != nil && signedDoc.SignedXML != "" {
		xmlData = []byte(signedDoc.SignedXML)
	} else {
		// Fallback: generate XML representation
		// Note: In production, the XML should be generated according to DGII's eCF specification
		xmlData = s.generateInvoiceXML(invoice, dgiiResponse)
	}

	// Save XML to storage
	if s.storage != nil && xmlData != nil {
		if err := s.saveInvoiceXML(ctx, schema, invoice, xmlData); err != nil {
			logger.Error("Failed to save XML to storage", "error", err, "invoice_id", id)
			// Don't fail the send operation if storage fails
		}
	} else if xmlData != nil {
		logger.Warn("Storage client not available, XML will not be saved", "invoice_id", id)
	}

	// Update signed_at timestamp if we have signed document
	if signedDoc != nil && signedDoc.SignedAt != "" {
		// Parse signed_at from the signed document
		if signedAt, err := time.Parse(time.RFC3339, signedDoc.SignedAt); err == nil {
			if err := s.repo.UpdateSignedAt(ctx, schema, invoiceID, signedAt); err != nil {
				logger.Error("Failed to update signed_at", "error", err)
			}
		}
	}

	// Update invoice status and DGII information
	if err := s.repo.UpdateDGIIStatus(ctx, schema, invoiceID, "sent", dgiiResponse.TrackingCode, dgiiResponse); err != nil {
		logger.Error("Failed to update invoice DGII status", "error", err)
		return fmt.Errorf("failed to update invoice status: %w", err)
	}

	// Update status to sent
	if err := s.repo.UpdateStatus(ctx, schema, invoiceID, "sent"); err != nil {
		return fmt.Errorf("failed to update invoice status: %w", err)
	}

	logger.Info("Invoice sent to DGII", "id", id, "tracking_code", dgiiResponse.TrackingCode)
	return nil
}

// invoiceToElectronicDocument converts an Invoice to an ElectronicDocument for DGII
func (s *Service) invoiceToElectronicDocument(ctx context.Context, schema string, invoice *Invoice) *dgii.ElectronicDocument {
	// Convert invoice lines to document items
	items := make([]*dgii.DocumentItem, len(invoice.Lines))
	for i, line := range invoice.Lines {
		items[i] = &dgii.DocumentItem{
			LineNumber:  line.LineNumber,
			Description: line.Description,
			Quantity:    line.Quantity,
			UnitPrice:   line.UnitPrice,
			TaxRate:     line.TaxRate,
			TaxAmount:   line.TaxAmount,
			Total:       line.LineTotal,
		}
	}

	// Get tenant config for sender information
	config, err := s.configRepo.Get(ctx, schema)
	sender := &dgii.PartyInfo{
		TaxID:   "00000000000", // Default
		Name:    "Empresa",
		Address: "",
		City:    "",
		Country: "DO",
	}
	if err == nil && config != nil {
		if config.DGIIIRNC != nil && *config.DGIIIRNC != "" {
			sender.TaxID = *config.DGIIIRNC
		}
		if config.CompanyLegalName != "" {
			sender.Name = config.CompanyLegalName
		} else if config.CompanyTradeName != nil && *config.CompanyTradeName != "" {
			sender.Name = *config.CompanyTradeName
		}
		if config.AddressLine1 != nil && *config.AddressLine1 != "" {
			sender.Address = *config.AddressLine1
			if config.AddressLine2 != nil && *config.AddressLine2 != "" {
				sender.Address += ", " + *config.AddressLine2
			}
		}
		if config.City != nil && *config.City != "" {
			sender.City = *config.City
		}
		if config.Country != "" {
			sender.Country = config.Country
		}
	}

	// Receiver info from customer
	receiver := &dgii.PartyInfo{
		TaxID:   "00000000000", // Will be filled from customer if available
		Name:    invoice.CustomerName,
		Country: "DO",
	}

	document := &dgii.ElectronicDocument{
		DocumentType: "Invoice",
		NCF:          *invoice.NCF,
		IssueDate:    invoice.IssueDate.Format("2006-01-02"),
		Sender:       sender,
		Receiver:     receiver,
		Items:        items,
		Totals: &dgii.DocumentTotals{
			Subtotal:  invoice.Subtotal,
			TaxAmount: invoice.TaxAmount,
			Total:     invoice.Total,
			Currency:  invoice.Currency,
		},
		Metadata: map[string]interface{}{
			"invoice_id":      invoice.ID.String(),
			"invoice_number":  invoice.InvoiceNumber,
			"ncf_type":        invoice.NCFType,
			"payment_type":    "cash", // Default, should be inferred from payments
		},
	}

	return document
}

// generateInvoiceXML generates an XML representation of the invoice
// This is a simplified version following DGII eCF format
// In production, this should strictly follow DGII's XML schema specification (XSD)
func (s *Service) generateInvoiceXML(invoice *Invoice, dgiiResponse *dgii.DGIIResponse) []byte {
	// Format date as DD-MM-YYYY (DGII format)
	issueDate := invoice.IssueDate.Format("02-01-2006")
	dueDate := invoice.DueDate.Format("02-01-2006")

	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<eCF>
	<Encabezado>
		<IDVersion>1.0</IDVersion>
		<eNCF>%s</eNCF>
		<TipoCF>%s</TipoCF>
		<FechaEmision>%s</FechaEmision>
		<FechaVencimiento>%s</FechaVencimiento>
		<Emisor>
			<RNCEmisor>13100000000</RNCEmisor>
			<RazonSocialEmisor>Demo Company</RazonSocialEmisor>
		</Emisor>
		<Comprador>
			<RNCComprador>00000000000</RNCComprador>
			<RazonSocialComprador>%s</RazonSocialComprador>
		</Comprador>
	</Encabezado>
	<Detalle>
%s	</Detalle>
	<Total>
		<SubTotal>%.2f</SubTotal>
		<ITBIS>%.2f</ITBIS>
		<Total>%.2f</Total>
		<Moneda>%s</Moneda>
	</Total>
	<DGII>
		<TrackingCode>%s</TrackingCode>
		<Status>%s</Status>
		<ReceivedAt>%s</ReceivedAt>
	</DGII>
</eCF>`,
		*invoice.NCF,
		invoice.NCFType,
		issueDate,
		dueDate,
		invoice.CustomerName,
		s.generateInvoiceLinesXML(invoice.Lines),
		invoice.Subtotal,
		invoice.TaxAmount,
		invoice.Total,
		invoice.Currency,
		dgiiResponse.TrackingCode,
		"sent",
		dgiiResponse.ReceivedAt,
	)

	return []byte(xml)
}

// generateInvoiceLinesXML generates XML for invoice lines
func (s *Service) generateInvoiceLinesXML(lines []InvoiceLine) string {
	var xmlLines string
	for _, line := range lines {
		xmlLines += fmt.Sprintf(`		<Linea>
			<NumeroLinea>%d</NumeroLinea>
			<Descripcion>%s</Descripcion>
			<Cantidad>%.2f</Cantidad>
			<PrecioUnitario>%.2f</PrecioUnitario>
			<TasaITBIS>%.2f</TasaITBIS>
			<ITBISLinea>%.2f</ITBISLinea>
			<MontoLinea>%.2f</MontoLinea>
		</Linea>
`, line.LineNumber, line.Description, line.Quantity, line.UnitPrice, line.TaxRate, line.TaxAmount, line.LineTotal)
	}
	return xmlLines
}

// saveInvoiceXML saves the invoice XML to storage
func (s *Service) saveInvoiceXML(ctx context.Context, schema string, invoice *Invoice, xmlData []byte) error {
	// Get storage client for tenant bucket
	storageClient, err := s.getStorageClientForTenant(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to get storage client: %w", err)
	}

	// Build object key
	objectKey := storage.BuildInvoiceXMLKey(schema, invoice.ID.String(), invoice.InvoiceNumber)

	// Upload to storage
	if err := storageClient.UploadBytes(ctx, xmlData, objectKey, "application/xml"); err != nil {
		return fmt.Errorf("failed to upload XML: %w", err)
	}

	// Generate presigned URL (valid for 7 days)
	xmlURL, err := storageClient.PresignedURL(ctx, objectKey, 7*24*time.Hour)
	if err != nil {
		logger.Error("Failed to generate presigned URL for XML", "error", err)
		// Still update with object key even if presigned URL fails
		xmlURL = fmt.Sprintf("s3://%s/%s", s.cfg.S3Bucket, objectKey)
	}

	// Update invoice with XML URL (preserve existing PDF URL)
	pdfURL := invoice.PDFURL
	if err := s.repo.UpdateDocumentURLs(ctx, schema, invoice.ID, pdfURL, &xmlURL); err != nil {
		logger.Error("Failed to update XML URL", "error", err)
		return fmt.Errorf("failed to update XML URL: %w", err)
	}

	logger.Info("XML saved to storage", "invoice_id", invoice.ID, "object_key", objectKey)
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

// AllocatePayment allocates a payment to one or more invoices
func (s *Service) AllocatePayment(ctx context.Context, schema string, req *CreatePaymentAllocationRequest) error {
	// Validate that payment exists and get its amount
	// Note: We need to import payments module or check payment via repository
	// For now, we'll calculate total allocation amount
	totalAllocation := 0.0
	for _, alloc := range req.Allocations {
		totalAllocation += alloc.Amount
	}

	// Start transaction
	tx, err := s.repo.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Create allocations and update invoices
	for _, allocReq := range req.Allocations {
		// Verify invoice exists and get current balance
		invoice, err := s.repo.GetByID(ctx, schema, allocReq.InvoiceID, []string{})
		if err != nil {
			return fmt.Errorf("invoice not found: %w", err)
		}

		// Calculate current balance
		currentBalance := invoice.Total - invoice.PaidAmount
		if allocReq.Amount > currentBalance {
			return fmt.Errorf("allocation amount (%.2f) exceeds invoice balance (%.2f) for invoice %s", 
				allocReq.Amount, currentBalance, invoice.InvoiceNumber)
		}

		// Create allocation record
		allocation := &PaymentAllocation{
			ID:        uuid.New(),
			PaymentID: req.PaymentID,
			InvoiceID: allocReq.InvoiceID,
			Amount:    allocReq.Amount,
			CreatedAt: time.Now(),
		}

		// Insert allocation
		allocQuery := fmt.Sprintf(`
			INSERT INTO %s.payment_allocations (id, payment_id, invoice_id, amount, created_at)
			VALUES ($1, $2, $3, $4, $5)
		`, schema)
		_, err = tx.ExecContext(ctx, allocQuery, allocation.ID, allocation.PaymentID, allocation.InvoiceID, allocation.Amount, allocation.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to create allocation: %w", err)
		}

		// Update invoice paid_amount
		newPaidAmount := invoice.PaidAmount + allocReq.Amount
		updateQuery := fmt.Sprintf(`
			UPDATE %s.invoices
			SET paid_amount = $1, updated_at = $2
			WHERE id = $3
		`, schema)
		_, err = tx.ExecContext(ctx, updateQuery, newPaidAmount, time.Now(), allocReq.InvoiceID)
		if err != nil {
			return fmt.Errorf("failed to update invoice paid amount: %w", err)
		}

		// Update invoice status if fully paid
		newStatus := invoice.Status
		if newPaidAmount >= invoice.Total {
			newStatus = string(StatusPaid)
		} else if invoice.Status == string(StatusOverdue) {
			// If partially paid, keep as overdue if still past due date
			if time.Now().After(invoice.DueDate) {
				newStatus = string(StatusOverdue)
			} else {
				newStatus = string(StatusPending)
			}
		}

		if newStatus != invoice.Status {
			statusQuery := fmt.Sprintf(`
				UPDATE %s.invoices
				SET status = $1, updated_at = $2
				WHERE id = $3
			`, schema)
			_, err = tx.ExecContext(ctx, statusQuery, newStatus, time.Now(), allocReq.InvoiceID)
			if err != nil {
				return fmt.Errorf("failed to update invoice status: %w", err)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info("Payment allocated to invoices", "payment_id", req.PaymentID, "allocations", len(req.Allocations))
	return nil
}

// GetAccountsReceivableReport generates an accounts receivable aging report
func (s *Service) GetAccountsReceivableReport(ctx context.Context, schema string, asOfDate *string) (*AccountsReceivableReport, error) {
	var reportDate time.Time
	if asOfDate != nil && *asOfDate != "" {
		var err error
		reportDate, err = time.Parse("2006-01-02", *asOfDate)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %w", err)
		}
	} else {
		reportDate = time.Now()
	}

	return s.repo.GetAccountsReceivableReport(ctx, schema, reportDate)
}

// GetInvoiceAllocations gets all payment allocations for an invoice
func (s *Service) GetInvoiceAllocations(ctx context.Context, schema string, invoiceID string) ([]*PaymentAllocation, error) {
	id, err := uuid.Parse(invoiceID)
	if err != nil {
		return nil, fmt.Errorf("invalid invoice ID: %w", err)
	}
	return s.repo.GetAllocationsByInvoice(ctx, schema, id)
}

// GetPaymentAllocations gets all allocations for a payment
func (s *Service) GetPaymentAllocations(ctx context.Context, schema string, paymentID string) ([]*PaymentAllocation, error) {
	id, err := uuid.Parse(paymentID)
	if err != nil {
		return nil, fmt.Errorf("invalid payment ID: %w", err)
	}
	return s.repo.GetAllocationsByPayment(ctx, schema, id)
}

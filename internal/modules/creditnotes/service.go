package creditnotes

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/dgii"
	"vendix/internal/logger"
	"vendix/internal/modules/invoices"
	"vendix/internal/modules/ncf"
	"vendix/internal/modules/tenantconfig"
	"vendix/internal/storage"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

type Service struct {
	repo         *Repository
	invoicesRepo *invoices.Repository
	storage      *storage.Client
	dgiiService  *dgii.Service
	ncfService   *ncf.Service
	configRepo   *tenantconfig.Repository
	cfg          *config.Config
}

func NewService(db *database.DB, cfg *config.Config) (*Service, error) {
	storageClient, err := storage.NewClient(cfg)
	if err != nil {
		logger.Error("Failed to initialize storage client", "error", err)
	}

	return &Service{
		repo:         NewRepository(db),
		invoicesRepo: invoices.NewRepository(db),
		storage:      storageClient,
		dgiiService:  dgii.NewService(cfg),
		ncfService:   ncf.NewService(db),
		configRepo:   tenantconfig.NewRepository(db),
		cfg:          cfg,
	}, nil
}

// CreateFromReturn creates a credit note from a return and original invoice
// This function is called from the returns module when processing a return with credit_note refund method
func (s *Service) CreateFromReturn(ctx context.Context, schema string, returnID uuid.UUID, originalInvoiceID uuid.UUID, returnLines []ReturnLineInfo, createdBy *uuid.UUID) (*CreditNote, error) {
	// Get original invoice
	invoice, err := s.invoicesRepo.GetByID(ctx, schema, originalInvoiceID, []string{})
	if err != nil {
		return nil, fmt.Errorf("original invoice not found: %w", err)
	}

	if invoice.NCF == nil || *invoice.NCF == "" {
		return nil, fmt.Errorf("original invoice must have an NCF to create credit note")
	}

	// Generate credit note number
	creditNoteNumber, err := s.repo.GetNextCreditNoteNumber(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to generate credit note number: %w", err)
	}

	// Generate NCF tipo 04 (Nota de Crédito)
	ncfInfo, err := s.ncfService.GenerateNCF(ctx, schema, ncf.NCFTypeCreditNote)
	if err != nil {
		return nil, fmt.Errorf("failed to generate NCF: %w", err)
	}

	// Calculate totals from return lines
	subtotal := 0.0
	taxAmount := 0.0

	creditNote := &CreditNote{
		ID:                uuid.New(),
		CreditNoteNumber:  creditNoteNumber,
		NCF:               &ncfInfo.FullNCF,
		NCFType:           "04",
		OriginalInvoiceID: &originalInvoiceID,
		ReturnID:          &returnID,
		CustomerID:        invoice.CustomerID,
		IssueDate:         time.Now(),
		Reason:            stringPtr("Devolución de productos"),
		Status:            CreditNoteStatusDraft,
		Currency:          invoice.Currency,
		CreatedBy:         createdBy,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Create lines from return lines (already have correct amounts from return)
	creditNote.Lines = make([]CreditNoteLine, len(returnLines))
	for i, retLine := range returnLines {
		subtotal += retLine.LineTotal - retLine.TaxAmount
		taxAmount += retLine.TaxAmount

		creditNote.Lines[i] = CreditNoteLine{
			ID:                    uuid.New(),
			CreditNoteID:          creditNote.ID,
			OriginalInvoiceLineID: retLine.InvoiceLineID,
			ProductID:             retLine.ProductID,
			LineNumber:            i + 1,
			Description:           retLine.Description,
			Quantity:              retLine.Quantity,
			UnitPrice:             retLine.UnitPrice,
			TaxRate:               retLine.TaxRate,
			TaxAmount:             retLine.TaxAmount,
			LineTotal:             retLine.LineTotal,
			CreatedAt:             time.Now(),
		}
	}

	creditNote.Subtotal = subtotal
	creditNote.TaxAmount = taxAmount
	creditNote.Total = subtotal + taxAmount

	// Save to database
	if err := s.repo.Create(ctx, schema, creditNote); err != nil {
		return nil, fmt.Errorf("failed to create credit note: %w", err)
	}

	logger.Info("Credit note created from return", "credit_note_number", creditNoteNumber, "ncf", ncfInfo.FullNCF, "return_id", returnID)
	return creditNote, nil
}

// ReturnLineInfo represents return line information for creating credit note
type ReturnLineInfo struct {
	InvoiceLineID *uuid.UUID
	ProductID     *uuid.UUID
	Description   string
	Quantity      float64
	UnitPrice     float64
	TaxRate       float64
	TaxAmount     float64
	LineTotal     float64
}

// Create creates a new credit note manually
func (s *Service) Create(ctx context.Context, schema string, req *CreateCreditNoteRequest, createdBy *uuid.UUID) (*CreditNote, error) {
	// Get original invoice
	invoice, err := s.invoicesRepo.GetByID(ctx, schema, req.OriginalInvoiceID, []string{})
	if err != nil {
		return nil, fmt.Errorf("original invoice not found: %w", err)
	}

	if invoice.NCF == nil || *invoice.NCF == "" {
		return nil, fmt.Errorf("original invoice must have an NCF")
	}

	// Parse issue date
	issueDate, err := time.Parse("2006-01-02", req.IssueDate)
	if err != nil {
		return nil, fmt.Errorf("invalid issue date: %w", err)
	}

	// Generate credit note number
	creditNoteNumber, err := s.repo.GetNextCreditNoteNumber(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to generate credit note number: %w", err)
	}

	// Generate NCF tipo 04
	ncfInfo, err := s.ncfService.GenerateNCF(ctx, schema, ncf.NCFTypeCreditNote)
	if err != nil {
		return nil, fmt.Errorf("failed to generate NCF: %w", err)
	}

	// Calculate totals
	subtotal := 0.0
	taxAmount := 0.0

	creditNote := &CreditNote{
		ID:                uuid.New(),
		CreditNoteNumber:  creditNoteNumber,
		NCF:               &ncfInfo.FullNCF,
		NCFType:           "04",
		OriginalInvoiceID: &req.OriginalInvoiceID,
		ReturnID:          req.ReturnID,
		CustomerID:        invoice.CustomerID,
		IssueDate:         issueDate,
		Reason:            req.Reason,
		Status:            CreditNoteStatusDraft,
		Notes:             req.Notes,
		Currency:          invoice.Currency,
		CreatedBy:         createdBy,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Create lines
	creditNote.Lines = make([]CreditNoteLine, len(req.Lines))
	for i, lineReq := range req.Lines {
		lineSubtotal := lineReq.Quantity * lineReq.UnitPrice
		lineTax := lineSubtotal * lineReq.TaxRate
		lineTotal := lineSubtotal + lineTax

		subtotal += lineSubtotal
		taxAmount += lineTax

		creditNote.Lines[i] = CreditNoteLine{
			ID:                    uuid.New(),
			CreditNoteID:          creditNote.ID,
			OriginalInvoiceLineID: lineReq.OriginalInvoiceLineID,
			LineNumber:            i + 1,
			Description:           lineReq.Description,
			Quantity:              lineReq.Quantity,
			UnitPrice:             lineReq.UnitPrice,
			TaxRate:               lineReq.TaxRate,
			TaxAmount:             lineTax,
			LineTotal:             lineTotal,
			CreatedAt:             time.Now(),
		}
	}

	creditNote.Subtotal = subtotal
	creditNote.TaxAmount = taxAmount
	creditNote.Total = subtotal + taxAmount

	// Save to database
	if err := s.repo.Create(ctx, schema, creditNote); err != nil {
		return nil, fmt.Errorf("failed to create credit note: %w", err)
	}

	logger.Info("Credit note created", "credit_note_number", creditNoteNumber, "ncf", ncfInfo.FullNCF)
	return creditNote, nil
}

// SendCreditNote sends a credit note to DGII
func (s *Service) SendCreditNote(ctx context.Context, schema string, id uuid.UUID) error {
	creditNote, err := s.repo.GetByID(ctx, schema, id)
	if err != nil {
		return fmt.Errorf("credit note not found: %w", err)
	}

	if creditNote.NCF == nil || *creditNote.NCF == "" {
		return fmt.Errorf("credit note must have an NCF before sending to DGII")
	}

	// Convert to ElectronicDocument
	document := s.creditNoteToElectronicDocument(ctx, schema, creditNote)

	// Process with DGII
	processResult, err := s.dgiiService.ProcessInvoiceWithSignedDoc(ctx, document)
	if err != nil {
		logger.Error("Failed to process credit note with DGII", "error", err, "credit_note_id", id)
		return fmt.Errorf("failed to send credit note to DGII: %w", err)
	}

	dgiiResponse := processResult.Response
	signedDoc := processResult.SignedDoc

	// Generate XML
	var xmlData []byte
	if signedDoc != nil && signedDoc.SignedXML != "" {
		xmlData = []byte(signedDoc.SignedXML)
	} else {
		xmlData = s.generateCreditNoteXML(ctx, schema, creditNote, dgiiResponse)
	}

	// Save XML to storage
	if s.storage != nil && xmlData != nil {
		if err := s.saveCreditNoteXML(ctx, schema, creditNote, xmlData); err != nil {
			logger.Error("Failed to save XML to storage", "error", err)
		}
	}

	// Update signed_at
	if signedDoc != nil && signedDoc.SignedAt != "" {
		if signedAt, err := time.Parse(time.RFC3339, signedDoc.SignedAt); err == nil {
			if err := s.repo.UpdateSignedAt(ctx, schema, id, signedAt); err != nil {
				logger.Error("Failed to update signed_at", "error", err)
			}
		}
	}

	// Update DGII status and sent_at
	if err := s.repo.UpdateDGIIStatus(ctx, schema, id, "sent", dgiiResponse.TrackingCode, dgiiResponse); err != nil {
		return fmt.Errorf("failed to update DGII status: %w", err)
	}

	if err := s.repo.UpdateSentAt(ctx, schema, id, time.Now()); err != nil {
		logger.Error("Failed to update sent_at", "error", err)
	}

	// Update status to sent
	if err := s.repo.UpdateStatus(ctx, schema, id, CreditNoteStatusSent); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	logger.Info("Credit note sent to DGII", "credit_note_id", id, "ncf", *creditNote.NCF)
	return nil
}

// creditNoteToElectronicDocument converts a CreditNote to ElectronicDocument for DGII
func (s *Service) creditNoteToElectronicDocument(ctx context.Context, schema string, creditNote *CreditNote) *dgii.ElectronicDocument {
	// Convert lines to document items
	items := make([]*dgii.DocumentItem, len(creditNote.Lines))
	for i, line := range creditNote.Lines {
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

	// Get tenant config for sender
	config, err := s.configRepo.Get(ctx, schema)
	sender := &dgii.PartyInfo{
		TaxID:   "00000000000",
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
		}
		// ... similar to invoice
	}

	receiver := &dgii.PartyInfo{
		TaxID:   "00000000000",
		Name:    creditNote.CustomerName,
		Country: "DO",
	}

	// Get original invoice NCF for metadata
	originalNCF := ""
	if creditNote.OriginalInvoiceID != nil {
		if originalInvoice, err := s.invoicesRepo.GetByID(ctx, schema, *creditNote.OriginalInvoiceID, []string{}); err == nil && originalInvoice != nil && originalInvoice.NCF != nil {
			originalNCF = *originalInvoice.NCF
		}
	}

	document := &dgii.ElectronicDocument{
		DocumentType: "CreditNote",
		NCF:          *creditNote.NCF,
		IssueDate:    creditNote.IssueDate.Format("2006-01-02"),
		Sender:       sender,
		Receiver:     receiver,
		Items:        items,
		Totals: &dgii.DocumentTotals{
			Subtotal:  creditNote.Subtotal,
			TaxAmount: creditNote.TaxAmount,
			Total:     creditNote.Total,
			Currency:  creditNote.Currency,
		},
		Metadata: map[string]interface{}{
			"credit_note_id":      creditNote.ID.String(),
			"credit_note_number":  creditNote.CreditNoteNumber,
			"original_invoice_id": func() string {
				if creditNote.OriginalInvoiceID != nil {
					return creditNote.OriginalInvoiceID.String()
				}
				return ""
			}(),
			"original_ncf": originalNCF,
			"ncf_type":     "04",
		},
	}

	return document
}

// generateCreditNoteXML generates XML for credit note
func (s *Service) generateCreditNoteXML(ctx context.Context, schema string, creditNote *CreditNote, dgiiResponse *dgii.DGIIResponse) []byte {
	issueDate := creditNote.IssueDate.Format("02-01-2006")

	// Get original invoice NCF
	originalNCF := ""
	if creditNote.OriginalInvoiceID != nil {
		if invoice, err := s.invoicesRepo.GetByID(ctx, schema, *creditNote.OriginalInvoiceID, []string{}); err == nil && invoice != nil && invoice.NCF != nil {
			originalNCF = *invoice.NCF
		}
	}

	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<eCF>
	<Encabezado>
		<IDVersion>1.0</IDVersion>
		<eNCF>%s</eNCF>
		<TipoCF>04</TipoCF>
		<FechaEmision>%s</FechaEmision>
		<MotivoNotaCredito>%s</MotivoNotaCredito>
		<eNCFModifica>%s</eNCFModifica>
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
		*creditNote.NCF,
		issueDate,
		getReasonOrDefault(creditNote.Reason),
		originalNCF,
		creditNote.CustomerName,
		s.generateCreditNoteLinesXML(creditNote.Lines),
		creditNote.Subtotal,
		creditNote.TaxAmount,
		creditNote.Total,
		creditNote.Currency,
		dgiiResponse.TrackingCode,
		"sent",
		dgiiResponse.ReceivedAt,
	)

	return []byte(xml)
}

// generateCreditNoteLinesXML generates XML for credit note lines
func (s *Service) generateCreditNoteLinesXML(lines []CreditNoteLine) string {
	var xmlLines string
	for _, line := range lines {
		xmlLines += fmt.Sprintf(`		<Linea>
			<NumeroLinea>%d</NumeroLinea>
			<Descripcion>%s</Descripcion>
			<Cantidad>%.2f</Cantidad>
			<PrecioUnitario>%.2f</PrecioUnitario>
			<TasaITBIS>%.2f</TasaITBIS>
			<ITBISLinea>%.2f</ITBISLinea>
			<TotalLinea>%.2f</TotalLinea>
		</Linea>
`, line.LineNumber, line.Description, line.Quantity, line.UnitPrice, line.TaxRate, line.TaxAmount, line.LineTotal)
	}
	return xmlLines
}

// saveCreditNoteXML saves credit note XML to storage
func (s *Service) saveCreditNoteXML(ctx context.Context, schema string, creditNote *CreditNote, xmlData []byte) error {
	storageClient, err := s.getStorageClientForTenant(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to get storage client: %w", err)
	}

	objectKey := storage.BuildCreditNoteXMLKey(schema, creditNote.ID.String(), creditNote.CreditNoteNumber)

	if err := storageClient.UploadBytes(ctx, xmlData, objectKey, "application/xml"); err != nil {
		return fmt.Errorf("failed to upload XML: %w", err)
	}

	xmlURL, err := storageClient.PresignedURL(ctx, objectKey, 7*24*time.Hour)
	if err != nil {
		logger.Error("Failed to generate presigned URL for XML", "error", err)
		xmlURL = fmt.Sprintf("s3://%s/%s", s.cfg.S3Bucket, objectKey)
	}

	if err := s.repo.UpdateXMLURL(ctx, schema, creditNote.ID, xmlURL); err != nil {
		return fmt.Errorf("failed to update XML URL: %w", err)
	}

	return nil
}

// getStorageClientForTenant returns storage client for tenant
func (s *Service) getStorageClientForTenant(ctx context.Context, schema string) (*storage.Client, error) {
	bucketName := storage.BuildTenantBucketName(schema)
	tenantClient, err := storage.NewClientWithBucket(s.cfg, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client for tenant: %w", err)
	}
	return tenantClient, nil
}

// GetByID retrieves a credit note by ID
func (s *Service) GetByID(ctx context.Context, schema string, id uuid.UUID) (*CreditNote, error) {
	return s.repo.GetByID(ctx, schema, id)
}

// List retrieves all credit notes
func (s *Service) List(ctx context.Context, schema string, filters *CreditNoteFilters) ([]*CreditNote, error) {
	return s.repo.List(ctx, schema, filters)
}

// GeneratePDF generates a PDF for a credit note
func (s *Service) GeneratePDF(ctx context.Context, schema string, id uuid.UUID) ([]byte, error) {
	creditNote, err := s.repo.GetByID(ctx, schema, id)
	if err != nil {
		return nil, fmt.Errorf("credit note not found: %w", err)
	}

	// Generate PDF using gofpdf
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "NOTA DE CRÉDITO", "", 0, "C", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, fmt.Sprintf("Número: %s", creditNote.CreditNoteNumber), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	if creditNote.NCF != nil {
		pdf.CellFormat(0, 5, fmt.Sprintf("NCF: %s", *creditNote.NCF), "", 0, "L", false, 0, "")
		pdf.Ln(5)
	}
	pdf.CellFormat(0, 5, fmt.Sprintf("Fecha: %s", creditNote.IssueDate.Format("02/01/2006")), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 5, fmt.Sprintf("Factura Original: %s", creditNote.OriginalInvoiceNumber), "", 0, "L", false, 0, "")
	pdf.Ln(10)

	// Customer info
	pdf.CellFormat(0, 5, fmt.Sprintf("Cliente: %s", creditNote.CustomerName), "", 0, "L", false, 0, "")
	pdf.Ln(5)

	// Lines table
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(100, 6, "Descripción", "1", 0, "L", false, 0, "")
	pdf.CellFormat(25, 6, "Cantidad", "1", 0, "R", false, 0, "")
	pdf.CellFormat(30, 6, "Precio", "1", 0, "R", false, 0, "")
	pdf.CellFormat(35, 6, "Total", "1", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	for _, line := range creditNote.Lines {
		pdf.CellFormat(100, 6, line.Description, "1", 0, "L", false, 0, "")
		pdf.CellFormat(25, 6, fmt.Sprintf("%.2f", line.Quantity), "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 6, fmt.Sprintf("%.2f", line.UnitPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(35, 6, fmt.Sprintf("%.2f", line.LineTotal), "1", 1, "R", false, 0, "")
	}

	// Totals
	pdf.Ln(2)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(155, 6, "Subtotal:", "", 0, "R", false, 0, "")
	pdf.CellFormat(35, 6, fmt.Sprintf("%.2f", creditNote.Subtotal), "", 1, "R", false, 0, "")
	pdf.CellFormat(155, 6, "ITBIS:", "", 0, "R", false, 0, "")
	pdf.CellFormat(35, 6, fmt.Sprintf("%.2f", creditNote.TaxAmount), "", 1, "R", false, 0, "")
	pdf.CellFormat(155, 6, "Total:", "", 0, "R", false, 0, "")
	pdf.CellFormat(35, 6, fmt.Sprintf("%.2f", creditNote.Total), "", 1, "R", false, 0, "")

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func getReasonOrDefault(reason *string) string {
	if reason != nil && *reason != "" {
		return *reason
	}
	return "Devolución de productos"
}

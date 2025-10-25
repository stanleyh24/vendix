package dgii

import (
	"context"
)

// DGIIAdapter defines the interface for DGII integration
type DGIIAdapter interface {
	// SignDocument digitally signs an electronic document
	SignDocument(ctx context.Context, document *ElectronicDocument) (*SignedDocument, error)

	// SendDocument sends a signed document to DGII
	SendDocument(ctx context.Context, signed *SignedDocument) (*DGIIResponse, error)

	// CancelDocument cancels a previously sent document
	CancelDocument(ctx context.Context, documentID string, reason string) error

	// GetDocumentStatus retrieves the status of a document from DGII
	GetDocumentStatus(ctx context.Context, documentID string) (*DocumentStatus, error)
}

// ElectronicDocument represents a document to be sent to DGII
type ElectronicDocument struct {
	DocumentType string                 `json:"document_type"` // Invoice, CreditNote, etc.
	NCF          string                 `json:"ncf"`           // Número de Comprobante Fiscal
	IssueDate    string                 `json:"issue_date"`
	Sender       *PartyInfo             `json:"sender"`
	Receiver     *PartyInfo             `json:"receiver"`
	Items        []*DocumentItem        `json:"items"`
	Totals       *DocumentTotals        `json:"totals"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// SignedDocument represents a digitally signed document
type SignedDocument struct {
	OriginalDocument *ElectronicDocument `json:"original_document"`
	Signature        string              `json:"signature"`
	Certificate      string              `json:"certificate"`
	SignedAt         string              `json:"signed_at"`
	SignedXML        string              `json:"signed_xml"`
}

// DGIIResponse represents the response from DGII
type DGIIResponse struct {
	Success      bool                   `json:"success"`
	TrackingCode string                 `json:"tracking_code"`
	Message      string                 `json:"message"`
	ErrorCode    string                 `json:"error_code,omitempty"`
	ReceivedAt   string                 `json:"received_at"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// DocumentStatus represents the status of a document in DGII
type DocumentStatus struct {
	DocumentID    string `json:"document_id"`
	Status        string `json:"status"`
	TrackingCode  string `json:"tracking_code"`
	Message       string `json:"message"`
	LastUpdatedAt string `json:"last_updated_at"`
}

// PartyInfo represents sender or receiver information
type PartyInfo struct {
	TaxID   string `json:"tax_id"` // RNC or Cédula
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
	Country string `json:"country"`
	Email   string `json:"email,omitempty"`
	Phone   string `json:"phone,omitempty"`
}

// DocumentItem represents a line item in the document
type DocumentItem struct {
	LineNumber  int     `json:"line_number"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TaxRate     float64 `json:"tax_rate"`
	TaxAmount   float64 `json:"tax_amount"`
	Total       float64 `json:"total"`
}

// DocumentTotals represents the totals of a document
type DocumentTotals struct {
	Subtotal  float64 `json:"subtotal"`
	TaxAmount float64 `json:"tax_amount"`
	Total     float64 `json:"total"`
	Currency  string  `json:"currency"`
}

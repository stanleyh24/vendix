package creditnotes

import (
	"time"

	"github.com/google/uuid"
)

// CreditNoteStatus represents the status of a credit note
type CreditNoteStatus string

const (
	CreditNoteStatusDraft     CreditNoteStatus = "draft"     // Borrador
	CreditNoteStatusPending   CreditNoteStatus = "pending"   // Pendiente
	CreditNoteStatusSent      CreditNoteStatus = "sent"      // Enviada a DGII
	CreditNoteStatusCancelled CreditNoteStatus = "cancelled" // Cancelada
)

// CreditNote represents a credit note (Nota de Crédito Fiscal)
type CreditNote struct {
	ID                uuid.UUID        `db:"id" json:"id"`
	CreditNoteNumber  string           `db:"credit_note_number" json:"credit_note_number"`
	NCF               *string          `db:"ncf" json:"ncf,omitempty"`
	NCFType           string           `db:"ncf_type" json:"ncf_type"` // Always "04" for credit notes
	OriginalInvoiceID *uuid.UUID       `db:"original_invoice_id" json:"original_invoice_id,omitempty"`
	ReturnID          *uuid.UUID       `db:"return_id" json:"return_id,omitempty"`
	CustomerID        uuid.UUID        `db:"customer_id" json:"customer_id"`
	IssueDate         time.Time        `db:"issue_date" json:"issue_date"`
	Reason            *string          `db:"reason" json:"reason,omitempty"`
	Status            CreditNoteStatus `db:"status" json:"status"`
	Subtotal          float64          `db:"subtotal" json:"subtotal"`
	TaxAmount         float64          `db:"tax_amount" json:"tax_amount"`
	Total             float64          `db:"total" json:"total"`
	Currency          string           `db:"currency" json:"currency"`
	Notes             *string          `db:"notes" json:"notes,omitempty"`
	DGIIStatus        *string          `db:"dgii_status" json:"dgii_status,omitempty"`
	TrackingCode      *string          `db:"tracking_code" json:"tracking_code,omitempty"`
	SignedAt          *time.Time       `db:"signed_at" json:"signed_at,omitempty"`
	SentAt            *time.Time       `db:"sent_at" json:"sent_at,omitempty"`
	PDFURL            *string          `db:"pdf_url" json:"pdf_url,omitempty"`
	XMLURL            *string          `db:"xml_url" json:"xml_url,omitempty"`
	CreatedBy         *uuid.UUID       `db:"created_by" json:"created_by,omitempty"`
	CreatedAt         time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time        `db:"updated_at" json:"updated_at"`

	// Relations
	Lines                 []CreditNoteLine `json:"lines,omitempty"`
	CustomerName          string           `db:"customer_name" json:"customer_name,omitempty"`
	OriginalInvoiceNumber string           `db:"original_invoice_number" json:"original_invoice_number,omitempty"`
}

// CreditNoteLine represents a line item in a credit note
type CreditNoteLine struct {
	ID                    uuid.UUID  `db:"id" json:"id"`
	CreditNoteID          uuid.UUID  `db:"credit_note_id" json:"credit_note_id"`
	OriginalInvoiceLineID *uuid.UUID `db:"original_invoice_line_id" json:"original_invoice_line_id,omitempty"`
	ProductID             *uuid.UUID `db:"product_id" json:"product_id,omitempty"`
	LineNumber            int        `db:"line_number" json:"line_number"`
	Description           string     `db:"description" json:"description"`
	Quantity              float64    `db:"quantity" json:"quantity"`
	UnitPrice             float64    `db:"unit_price" json:"unit_price"`
	TaxRate               float64    `db:"tax_rate" json:"tax_rate"`
	TaxAmount             float64    `db:"tax_amount" json:"tax_amount"`
	LineTotal             float64    `db:"line_total" json:"line_total"`
	CreatedAt             time.Time  `db:"created_at" json:"created_at"`
}

// CreateCreditNoteRequest represents the request to create a credit note
type CreateCreditNoteRequest struct {
	OriginalInvoiceID uuid.UUID                 `json:"original_invoice_id" validate:"required"`
	ReturnID          *uuid.UUID                `json:"return_id,omitempty"`
	IssueDate         string                    `json:"issue_date" validate:"required"`
	Reason            *string                   `json:"reason,omitempty"`
	Notes             *string                   `json:"notes,omitempty"`
	Lines             []CreateCreditNoteLineReq `json:"lines" validate:"required,min=1"`
}

// CreateCreditNoteLineReq represents a line item in the create request
type CreateCreditNoteLineReq struct {
	OriginalInvoiceLineID *uuid.UUID `json:"original_invoice_line_id,omitempty"`
	Description           string     `json:"description" validate:"required"`
	Quantity              float64    `json:"quantity" validate:"required,gt=0"`
	UnitPrice             float64    `json:"unit_price" validate:"required,gte=0"`
	TaxRate               float64    `json:"tax_rate" validate:"gte=0"`
}

// IsValidStatus checks if the credit note status is valid
func (s CreditNoteStatus) IsValid() bool {
	switch s {
	case CreditNoteStatusDraft, CreditNoteStatusPending, CreditNoteStatusSent, CreditNoteStatusCancelled:
		return true
	default:
		return false
	}
}

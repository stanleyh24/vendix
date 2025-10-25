package invoices

import (
	"time"

	"github.com/google/uuid"
)

// Invoice represents an invoice
type Invoice struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	InvoiceNumber string     `db:"invoice_number" json:"invoice_number"`
	NCF           *string    `db:"ncf" json:"ncf,omitempty"`
	NCFType       string     `db:"ncf_type" json:"ncf_type"` // 01=Crédito Fiscal, 02=Consumidor Final
	CustomerID    uuid.UUID  `db:"customer_id" json:"customer_id"`
	IssueDate     time.Time  `db:"issue_date" json:"issue_date"`
	DueDate       time.Time  `db:"due_date" json:"due_date"`
	Status        string     `db:"status" json:"status"`
	Subtotal      float64    `db:"subtotal" json:"subtotal"`
	TaxAmount     float64    `db:"tax_amount" json:"tax_amount"`
	Total         float64    `db:"total" json:"total"`
	PaidAmount    float64    `db:"paid_amount" json:"paid_amount"`
	Currency      string     `db:"currency" json:"currency"`
	Notes         *string    `db:"notes" json:"notes,omitempty"`
	Terms         *string    `db:"terms" json:"terms,omitempty"`
	DGIIStatus    *string    `db:"dgii_status" json:"dgii_status,omitempty"`
	SignedAt      *time.Time `db:"signed_at" json:"signed_at,omitempty"`
	SentAt        *time.Time `db:"sent_at" json:"sent_at,omitempty"`
	CreatedBy     *uuid.UUID `db:"created_by" json:"created_by,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`

	// Relations
	Lines        []InvoiceLine `json:"lines,omitempty"`
	CustomerName string        `db:"customer_name" json:"customer_name,omitempty"`
}

// InvoiceLine represents a line item in an invoice
type InvoiceLine struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	InvoiceID   uuid.UUID  `db:"invoice_id" json:"invoice_id"`
	ProductID   *uuid.UUID `db:"product_id" json:"product_id,omitempty"`
	LineNumber  int        `db:"line_number" json:"line_number"`
	Description string     `db:"description" json:"description"`
	Quantity    float64    `db:"quantity" json:"quantity"`
	UnitPrice   float64    `db:"unit_price" json:"unit_price"`
	TaxRate     float64    `db:"tax_rate" json:"tax_rate"`
	TaxAmount   float64    `db:"tax_amount" json:"tax_amount"`
	LineTotal   float64    `db:"line_total" json:"line_total"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
}

// InvoiceStatus represents the possible statuses of an invoice
type InvoiceStatus string

const (
	StatusDraft     InvoiceStatus = "draft"     // Borrador - factura en proceso de creación
	StatusPending   InvoiceStatus = "pending"   // Pendiente - factura creada, esperando pago
	StatusPaid      InvoiceStatus = "paid"      // Pagada - factura completamente pagada
	StatusSent      InvoiceStatus = "sent"      // Enviada - factura enviada a DGII
	StatusCancelled InvoiceStatus = "cancelled" // Cancelada - factura cancelada
	StatusOverdue   InvoiceStatus = "overdue"   // Vencida - factura vencida sin pago
)

// IsValidStatus checks if the status is valid
func (s InvoiceStatus) IsValid() bool {
	switch s {
	case StatusDraft, StatusPending, StatusPaid, StatusSent, StatusCancelled, StatusOverdue:
		return true
	default:
		return false
	}
}

// CreateInvoiceRequest represents the request to create an invoice
type CreateInvoiceRequest struct {
	CustomerID *uuid.UUID             `json:"customer_id,omitempty"` // Opcional - si no se envía, usa cliente genérico
	NCFType    string                 `json:"ncf_type"`              // 01=Crédito Fiscal, 02=Consumidor Final (default: 02)
	IssueDate  string                 `json:"issue_date" validate:"required"`
	DueDate    string                 `json:"due_date" validate:"required"`
	Status     *string                `json:"status,omitempty"` // Opcional - si no se envía, usa 'draft' por defecto
	Notes      *string                `json:"notes,omitempty"`
	Terms      *string                `json:"terms,omitempty"`
	Lines      []CreateInvoiceLineReq `json:"lines" validate:"required,min=1"`
}

// CreateInvoiceLineReq represents a line item in the create request
type CreateInvoiceLineReq struct {
	ProductID   *uuid.UUID `json:"product_id,omitempty"`
	Description string     `json:"description" validate:"required"`
	Quantity    float64    `json:"quantity" validate:"required,gt=0"`
	UnitPrice   float64    `json:"unit_price" validate:"required,gte=0"`
	TaxRate     float64    `json:"tax_rate" validate:"gte=0"`
}

// UpdateInvoiceRequest represents the request to update an invoice
type UpdateInvoiceRequest struct {
	Status *string `json:"status,omitempty"`
	Notes  *string `json:"notes,omitempty"`
	Terms  *string `json:"terms,omitempty"`
}

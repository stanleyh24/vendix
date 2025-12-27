package sales

import (
	"time"

	"github.com/google/uuid"
)

// Sale represents a sale transaction
type Sale struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	SaleNumber  string     `db:"sale_number" json:"sale_number"`
	CustomerID  *uuid.UUID `db:"customer_id" json:"customer_id,omitempty"`
	InvoiceID   *uuid.UUID `db:"invoice_id" json:"invoice_id,omitempty"`
	Subtotal    float64    `db:"subtotal" json:"subtotal"`
	TaxAmount   float64    `db:"tax_amount" json:"tax_amount"`
	Total       float64    `db:"total" json:"total"`
	PaymentType string     `db:"payment_type" json:"payment_type"` // cash, card, transfer, etc.
	Status      string     `db:"status" json:"status"`             // completed, cancelled, refunded
	Notes       *string    `db:"notes" json:"notes,omitempty"`
	CreatedBy   *uuid.UUID `db:"created_by" json:"created_by,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`

	// Relations
	Lines        []SaleLine `json:"lines,omitempty"`
	CustomerName string     `db:"customer_name" json:"customer_name,omitempty"`
}

// SaleLine represents a line item in a sale
type SaleLine struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	SaleID      uuid.UUID  `db:"sale_id" json:"sale_id"`
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

// CreateSaleRequest represents the request to create a sale
type CreateSaleRequest struct {
	CustomerID  *uuid.UUID          `json:"customer_id,omitempty"`            // Opcional - si no se envía, usa cliente genérico
	PaymentType string              `json:"payment_type" validate:"required"` // cash, card, transfer, etc.
	NCFType     string              `json:"ncf_type,omitempty"`               // 01=Crédito Fiscal, 02=Consumidor Final (default: 02)
	Notes       *string             `json:"notes,omitempty"`
	Lines       []CreateSaleLineReq `json:"lines" validate:"required,min=1"`
}

// CreateSaleLineReq represents a line item in the create sale request
type CreateSaleLineReq struct {
	ProductID   *uuid.UUID `json:"product_id,omitempty"`
	Description string     `json:"description" validate:"required"`
	Quantity    float64    `json:"quantity" validate:"required,gt=0"`
	UnitPrice   float64    `json:"unit_price" validate:"required,gte=0"`
	TaxRate     float64    `json:"tax_rate" validate:"gte=0"`
}

// SaleStatus represents the possible statuses of a sale
type SaleStatus string

const (
	SaleStatusCompleted SaleStatus = "completed" // Completada - venta finalizada
	SaleStatusCancelled SaleStatus = "cancelled" // Cancelada - venta cancelada
	SaleStatusRefunded  SaleStatus = "refunded"  // Reembolsada - venta reembolsada
)

// IsValidStatus checks if the sale status is valid
func (s SaleStatus) IsValid() bool {
	switch s {
	case SaleStatusCompleted, SaleStatusCancelled, SaleStatusRefunded:
		return true
	default:
		return false
	}
}

// PaymentType represents the possible payment types
type PaymentType string

const (
	PaymentTypeCash     PaymentType = "cash"     // Efectivo
	PaymentTypeCard     PaymentType = "card"     // Tarjeta
	PaymentTypeTransfer PaymentType = "transfer" // Transferencia
	PaymentTypeCheck    PaymentType = "check"    // Cheque
	PaymentTypeMixed    PaymentType = "mixed"    // Mixto (combinación de métodos)
)

// IsValidPaymentType checks if the payment type is valid
func (p PaymentType) IsValid() bool {
	switch p {
	case PaymentTypeCash, PaymentTypeCard, PaymentTypeTransfer, PaymentTypeCheck, PaymentTypeMixed:
		return true
	default:
		return false
	}
}

package returns

import (
	"time"

	"github.com/google/uuid"
)

// ReturnStatus represents the status of a return
type ReturnStatus string

const (
	ReturnStatusPending   ReturnStatus = "pending"   // Pendiente de aprobación
	ReturnStatusApproved  ReturnStatus = "approved"  // Aprobada, pendiente procesamiento
	ReturnStatusRejected  ReturnStatus = "rejected"  // Rechazada
	ReturnStatusCompleted ReturnStatus = "completed" // Completada (inventario restaurado)
	ReturnStatusRefunded  ReturnStatus = "refunded"  // Reembolsada completamente
)

// ReturnType represents the type of return
type ReturnType string

const (
	ReturnTypeSale    ReturnType = "sale"    // Devolución de venta POS
	ReturnTypeInvoice ReturnType = "invoice" // Devolución de factura
)

// RefundMethod represents the refund method
type RefundMethod string

const (
	RefundMethodCash       RefundMethod = "cash"        // Efectivo
	RefundMethodCard       RefundMethod = "card"        // Tarjeta (reversar transacción)
	RefundMethodTransfer   RefundMethod = "transfer"    // Transferencia bancaria
	RefundMethodCreditNote RefundMethod = "credit_note" // Nota de crédito (facturación)
)

// ReturnReason represents the reason for return
type ReturnReason string

const (
	ReturnReasonDefect          ReturnReason = "defect"           // Producto defectuoso
	ReturnReasonWrongItem       ReturnReason = "wrong_item"       // Producto incorrecto
	ReturnReasonCustomerRequest ReturnReason = "customer_request" // Solicitud del cliente
	ReturnReasonDamaged         ReturnReason = "damaged"          // Producto dañado en tránsito
	ReturnReasonNotAsDescribed  ReturnReason = "not_as_described" // No coincide con descripción
)

// ItemCondition represents the condition of returned item
type ItemCondition string

const (
	ItemConditionNew       ItemCondition = "new"       // Nuevo, sin usar
	ItemConditionUsed      ItemCondition = "used"      // Usado pero en buen estado
	ItemConditionDefective ItemCondition = "defective" // Defectuoso
	ItemConditionDamaged   ItemCondition = "damaged"   // Dañado
)

// Return represents a return transaction
type Return struct {
	ID            uuid.UUID   `db:"id" json:"id"`
	ReturnNumber  string      `db:"return_number" json:"return_number"`
	SaleID        *uuid.UUID  `db:"sale_id" json:"sale_id,omitempty"`
	InvoiceID     *uuid.UUID  `db:"invoice_id" json:"invoice_id,omitempty"`
	CustomerID    *uuid.UUID  `db:"customer_id" json:"customer_id,omitempty"`
	ReturnDate    time.Time   `db:"return_date" json:"return_date"`
	ReturnType    ReturnType  `db:"return_type" json:"return_type"`
	ReturnReason  *ReturnReason `db:"return_reason" json:"return_reason,omitempty"`
	Notes         *string      `db:"notes" json:"notes,omitempty"`
	Subtotal      float64      `db:"subtotal" json:"subtotal"`
	TaxAmount     float64      `db:"tax_amount" json:"tax_amount"`
	Total         float64      `db:"total" json:"total"`
	Status        ReturnStatus `db:"status" json:"status"`
	RefundStatus  *string      `db:"refund_status" json:"refund_status,omitempty"`
	RefundMethod  *RefundMethod `db:"refund_method" json:"refund_method,omitempty"`
	RefundAmount  float64      `db:"refund_amount" json:"refund_amount"`
	CreatedBy     *uuid.UUID   `db:"created_by" json:"created_by,omitempty"`
	ApprovedBy    *uuid.UUID   `db:"approved_by" json:"approved_by,omitempty"`
	ApprovedAt    *time.Time   `db:"approved_at" json:"approved_at,omitempty"`
	CreatedAt     time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at" json:"updated_at"`

	// Relations
	Lines         []ReturnLine `json:"lines,omitempty"`
	CustomerName  string       `db:"customer_name" json:"customer_name,omitempty"`
	SaleNumber    *string      `db:"sale_number" json:"sale_number,omitempty"`
	InvoiceNumber *string      `db:"invoice_number" json:"invoice_number,omitempty"`
}

// ReturnLine represents a line item in a return
type ReturnLine struct {
	ID             uuid.UUID     `db:"id" json:"id"`
	ReturnID       uuid.UUID     `db:"return_id" json:"return_id"`
	SaleLineID     *uuid.UUID    `db:"sale_line_id" json:"sale_line_id,omitempty"`
	InvoiceLineID  *uuid.UUID    `db:"invoice_line_id" json:"invoice_line_id,omitempty"`
	ProductID      *uuid.UUID    `db:"product_id" json:"product_id,omitempty"`
	LineNumber     int           `db:"line_number" json:"line_number"`
	Description    string        `db:"description" json:"description"`
	Quantity       float64       `db:"quantity" json:"quantity"`
	UnitPrice      float64       `db:"unit_price" json:"unit_price"`
	TaxRate        float64       `db:"tax_rate" json:"tax_rate"`
	TaxAmount      float64       `db:"tax_amount" json:"tax_amount"`
	LineTotal      float64       `db:"line_total" json:"line_total"`
	ItemCondition  *ItemCondition `db:"item_condition" json:"item_condition,omitempty"`
	RestockQuantity *float64      `db:"restock_quantity" json:"restock_quantity,omitempty"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
}

// CreateReturnRequest represents the request to create a return
type CreateReturnRequest struct {
	SaleID       *uuid.UUID           `json:"sale_id,omitempty"`       // Uno de sale_id o invoice_id debe estar presente
	InvoiceID    *uuid.UUID           `json:"invoice_id,omitempty"`
	ReturnReason *ReturnReason        `json:"return_reason,omitempty"`
	Notes        *string              `json:"notes,omitempty"`
	Lines        []CreateReturnLineReq `json:"lines" validate:"required,min=1"`
}

// CreateReturnLineReq represents a line item in the create return request
type CreateReturnLineReq struct {
	SaleLineID    *uuid.UUID     `json:"sale_line_id,omitempty"`
	InvoiceLineID *uuid.UUID     `json:"invoice_line_id,omitempty"`
	Quantity      float64        `json:"quantity" validate:"required,gt=0"`
	ItemCondition *ItemCondition `json:"item_condition,omitempty"`
	RestockQuantity *float64     `json:"restock_quantity,omitempty"`
}

// ProcessReturnRequest represents the request to process a return
type ProcessReturnRequest struct {
	RefundMethod RefundMethod `json:"refund_method" validate:"required"`
	RefundAmount *float64     `json:"refund_amount,omitempty"` // Opcional, por defecto usa el total de la devolución
}

// RejectReturnRequest represents the request to reject a return
type RejectReturnRequest struct {
	Reason string `json:"reason" validate:"required"`
}

// IsValidStatus checks if the return status is valid
func (s ReturnStatus) IsValid() bool {
	switch s {
	case ReturnStatusPending, ReturnStatusApproved, ReturnStatusRejected, ReturnStatusCompleted, ReturnStatusRefunded:
		return true
	default:
		return false
	}
}

// IsValidReturnType checks if the return type is valid
func (t ReturnType) IsValid() bool {
	switch t {
	case ReturnTypeSale, ReturnTypeInvoice:
		return true
	default:
		return false
	}
}

// IsValidRefundMethod checks if the refund method is valid
func (m RefundMethod) IsValid() bool {
	switch m {
	case RefundMethodCash, RefundMethodCard, RefundMethodTransfer, RefundMethodCreditNote:
		return true
	default:
		return false
	}
}


package purchases

import (
	"time"

	"github.com/google/uuid"
)

// PurchaseStatus represents the status of a purchase order
type PurchaseStatus string

const (
	PurchaseStatusDraft     PurchaseStatus = "draft"
	PurchaseStatusPending   PurchaseStatus = "pending"
	PurchaseStatusReceived  PurchaseStatus = "received"
	PurchaseStatusCancelled PurchaseStatus = "cancelled"
)

// PaymentStatus represents the payment status of a purchase
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPartial   PaymentStatus = "partial"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// PaymentMethod represents how the purchase is paid
type PaymentMethod string

const (
	PaymentMethodCash   PaymentMethod = "cash"
	PaymentMethodCredit PaymentMethod = "credit"
	PaymentMethodBank   PaymentMethod = "bank"
	PaymentMethodCheck  PaymentMethod = "check"
)

// Purchase represents a purchase order
type Purchase struct {
	ID                  uuid.UUID     `db:"id" json:"id"`
	PurchaseNumber      string        `db:"purchase_number" json:"purchase_number"`
	SupplierID          uuid.UUID     `db:"supplier_id" json:"supplier_id"`
	PurchaseDate        time.Time     `db:"purchase_date" json:"purchase_date"`
	ExpectedDeliveryDate *time.Time   `db:"expected_delivery_date" json:"expected_delivery_date,omitempty"`
	PaymentMethod       string        `db:"payment_method" json:"payment_method"`
	PaymentStatus       string        `db:"payment_status" json:"payment_status"`
	PaymentTerms        int           `db:"payment_terms" json:"payment_terms"` // Días de crédito
	PaymentDueDate      *time.Time    `db:"payment_due_date" json:"payment_due_date,omitempty"` // Fecha de vencimiento
	Subtotal            float64       `db:"subtotal" json:"subtotal"`
	TaxAmount           float64       `db:"tax_amount" json:"tax_amount"`
	Total               float64       `db:"total" json:"total"`
	Currency            string        `db:"currency" json:"currency"`
	Reference           *string       `db:"reference" json:"reference,omitempty"`
	Notes               *string       `db:"notes" json:"notes,omitempty"`
	Status              string        `db:"status" json:"status"`
	ReceivedAt          *time.Time    `db:"received_at" json:"received_at,omitempty"`
	CreatedBy           *uuid.UUID    `db:"created_by" json:"created_by,omitempty"`
	CreatedAt           time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time     `db:"updated_at" json:"updated_at"`
}

// PurchaseLine represents a line item in a purchase
type PurchaseLine struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	PurchaseID      uuid.UUID  `db:"purchase_id" json:"purchase_id"`
	ProductID       *uuid.UUID `db:"product_id" json:"product_id,omitempty"`
	Description     string     `db:"description" json:"description"`
	Quantity        float64    `db:"quantity" json:"quantity"`
	UnitPrice       float64    `db:"unit_price" json:"unit_price"`
	TaxRate         float64    `db:"tax_rate" json:"tax_rate"`
	Subtotal        float64    `db:"subtotal" json:"subtotal"`
	TaxAmount       float64    `db:"tax_amount" json:"tax_amount"`
	Total           float64    `db:"total" json:"total"`
	ReceivedQuantity float64   `db:"received_quantity" json:"received_quantity"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
}

// CreatePurchaseRequest DTO
type CreatePurchaseRequest struct {
	SupplierID          string                `json:"supplier_id" validate:"required"`
	PurchaseDate        string                `json:"purchase_date" validate:"required"` // YYYY-MM-DD
	ExpectedDeliveryDate *string              `json:"expected_delivery_date,omitempty"` // YYYY-MM-DD
	PaymentMethod       string                `json:"payment_method" validate:"required"`
	PaymentTerms        *int                  `json:"payment_terms,omitempty"` // Días de crédito (default: 30)
	Reference           *string               `json:"reference,omitempty"`
	Notes               *string               `json:"notes,omitempty"`
	Lines               []CreatePurchaseLineReq `json:"lines" validate:"required,min=1"`
}

// CreatePurchaseLineReq DTO
type CreatePurchaseLineReq struct {
	ProductID   *string  `json:"product_id,omitempty"`
	Description string   `json:"description" validate:"required"`
	Quantity    float64  `json:"quantity" validate:"required,gt=0"`
	UnitPrice   float64  `json:"unit_price" validate:"required,gt=0"`
	TaxRate     float64  `json:"tax_rate"`
}

// UpdatePurchaseRequest DTO
type UpdatePurchaseRequest struct {
	PurchaseDate        *string               `json:"purchase_date,omitempty"`
	ExpectedDeliveryDate *string               `json:"expected_delivery_date,omitempty"`
	PaymentMethod       *string               `json:"payment_method,omitempty"`
	Reference           *string               `json:"reference,omitempty"`
	Notes               *string               `json:"notes,omitempty"`
	Status              *string               `json:"status,omitempty"`
	Lines               []CreatePurchaseLineReq `json:"lines,omitempty"`
}

// ReceivePurchaseRequest DTO for receiving a purchase
type ReceivePurchaseRequest struct {
	ReceivedDate string                    `json:"received_date" validate:"required"` // YYYY-MM-DD
	Lines        []ReceivePurchaseLineReq  `json:"lines" validate:"required,min=1"`
}

// ReceivePurchaseLineReq DTO
type ReceivePurchaseLineReq struct {
	LineID           string  `json:"line_id" validate:"required"`
	ReceivedQuantity float64 `json:"received_quantity" validate:"required,gt=0"`
}

// PurchaseResponse includes supplier information
type PurchaseResponse struct {
	Purchase
	SupplierName *string `json:"supplier_name,omitempty"`
	Lines        []PurchaseLine `json:"lines,omitempty"`
}

// SupplierPayment represents a payment made to a supplier
type SupplierPayment struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	PaymentNumber string     `db:"payment_number" json:"payment_number"`
	SupplierID    uuid.UUID  `db:"supplier_id" json:"supplier_id"`
	PaymentMethod string     `db:"payment_method" json:"payment_method"`
	PaymentDate   time.Time  `db:"payment_date" json:"payment_date"`
	Amount        float64    `db:"amount" json:"amount"`
	Currency      string     `db:"currency" json:"currency"`
	Reference     *string    `db:"reference" json:"reference,omitempty"`
	Notes         *string    `db:"notes" json:"notes,omitempty"`
	Status        string     `db:"status" json:"status"`
	CreatedBy     *uuid.UUID `db:"created_by" json:"created_by,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}

// SupplierPaymentAllocation represents allocation of a payment to a purchase
type SupplierPaymentAllocation struct {
	ID               uuid.UUID `db:"id" json:"id"`
	SupplierPaymentID uuid.UUID `db:"supplier_payment_id" json:"supplier_payment_id"`
	PurchaseID       uuid.UUID `db:"purchase_id" json:"purchase_id"`
	Amount           float64   `db:"amount" json:"amount"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

// CreateSupplierPaymentRequest DTO
type CreateSupplierPaymentRequest struct {
	SupplierID    string                `json:"supplier_id" validate:"required"`
	PaymentMethod string                `json:"payment_method" validate:"required"`
	PaymentDate  string                `json:"payment_date" validate:"required"` // YYYY-MM-DD
	Amount       float64               `json:"amount" validate:"required,gt=0"`
	Currency     string                `json:"currency"`
	Reference    *string               `json:"reference,omitempty"`
	Notes        *string               `json:"notes,omitempty"`
	Allocations  []PaymentAllocationReq `json:"allocations" validate:"required,min=1"` // Asignaciones a compras
}

// PaymentAllocationReq DTO for allocating payment to purchases
type PaymentAllocationReq struct {
	PurchaseID string  `json:"purchase_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
}

// SupplierPaymentResponse includes supplier information
type SupplierPaymentResponse struct {
	SupplierPayment
	SupplierName *string                      `json:"supplier_name,omitempty"`
	Allocations  []*SupplierPaymentAllocation `json:"allocations,omitempty"`
}

// AccountsPayableItem represents an item in the accounts payable aging report
type AccountsPayableItem struct {
	PurchaseID      string    `db:"purchase_id" json:"purchase_id"`
	PurchaseNumber  string    `db:"purchase_number" json:"purchase_number"`
	SupplierID      string    `db:"supplier_id" json:"supplier_id"`
	SupplierName    string    `db:"supplier_name" json:"supplier_name"`
	PurchaseDate    time.Time `db:"purchase_date" json:"purchase_date"`
	ExpectedDeliveryDate *time.Time `db:"expected_delivery_date" json:"expected_delivery_date,omitempty"`
	Total            float64   `db:"total" json:"total"`
	PaidAmount       float64   `db:"paid_amount" json:"paid_amount"`
	Balance          float64   `db:"balance" json:"balance"` // total - paid_amount
	PaymentStatus    string    `db:"payment_status" json:"payment_status"`
	Status           string    `db:"status" json:"status"`
	DaysOverdue      *int      `db:"days_overdue" json:"days_overdue,omitempty"`
	Age0_30          float64   `db:"age_0_30" json:"age_0_30"`   // Balance in 0-30 days
	Age31_60         float64   `db:"age_31_60" json:"age_31_60"` // Balance in 31-60 days
	Age61_90         float64   `db:"age_61_90" json:"age_61_90"` // Balance in 61-90 days
	AgeOver90        float64   `db:"age_over_90" json:"age_over_90"` // Balance over 90 days
}

// AccountsPayableReport represents the complete AP aging report
type AccountsPayableReport struct {
	ReportDate        string                  `json:"report_date"`
	TotalPurchases    int                     `json:"total_purchases"`
	TotalOutstanding  float64                 `json:"total_outstanding"` // Total balance due
	Total0_30         float64                 `json:"total_0_30"`
	Total31_60        float64                 `json:"total_31_60"`
	Total61_90        float64                 `json:"total_61_90"`
	TotalOver90       float64                 `json:"total_over_90"`
	Items             []AccountsPayableItem   `json:"items"`
	SummaryBySupplier map[string]SupplierAPSummary `json:"summary_by_supplier,omitempty"`
}

// SupplierAPSummary represents AP summary grouped by supplier
type SupplierAPSummary struct {
	SupplierID       string  `json:"supplier_id"`
	SupplierName     string  `json:"supplier_name"`
	PurchaseCount    int     `json:"purchase_count"`
	TotalOutstanding float64 `json:"total_outstanding"`
	Total0_30        float64 `json:"total_0_30"`
	Total31_60       float64 `json:"total_31_60"`
	Total61_90       float64 `json:"total_61_90"`
	TotalOver90      float64 `json:"total_over_90"`
}


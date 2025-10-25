package payments

import (
	"time"

	"github.com/google/uuid"
)

// Payment represents a payment or expense
type Payment struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	PaymentNumber   string     `db:"payment_number" json:"payment_number"`
	CustomerID      *uuid.UUID `db:"customer_id" json:"customer_id,omitempty"`
	PaymentMethod   string     `db:"payment_method" json:"payment_method"`
	PaymentDate     time.Time  `db:"payment_date" json:"payment_date"`
	Amount          float64    `db:"amount" json:"amount"`
	Currency        string     `db:"currency" json:"currency"`
	Reference       *string    `db:"reference" json:"reference,omitempty"`
	Notes           *string    `db:"notes" json:"notes,omitempty"`
	StripePaymentID *string    `db:"stripe_payment_id" json:"stripe_payment_id,omitempty"`
	Status          string     `db:"status" json:"status"`
	CreatedBy       *uuid.UUID `db:"created_by" json:"created_by,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`

	// Additional fields for expense categorization
	Category    string `db:"category" json:"category,omitempty"`
	Supplier    string `db:"supplier" json:"supplier,omitempty"`
	Description string `db:"description" json:"description,omitempty"`
}

// CreatePaymentRequest represents the request to create a payment/expense
type CreatePaymentRequest struct {
	CustomerID    *uuid.UUID `json:"customer_id,omitempty"`
	PaymentMethod string     `json:"payment_method" validate:"required"`
	PaymentDate   string     `json:"payment_date" validate:"required"`
	Amount        float64    `json:"amount" validate:"required,gt=0"`
	Currency      string     `json:"currency"`
	Reference     *string    `json:"reference,omitempty"`
	Notes         *string    `json:"notes,omitempty"`

	// For expenses
	Category    string `json:"category,omitempty"`
	Supplier    string `json:"supplier,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdatePaymentRequest represents the request to update a payment
type UpdatePaymentRequest struct {
	Status    *string `json:"status,omitempty"`
	Notes     *string `json:"notes,omitempty"`
	Reference *string `json:"reference,omitempty"`
}

package cashregisters

import (
	"time"

	"github.com/google/uuid"
)

// CashRegister represents a cash register
type CashRegister struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	Name          string     `db:"name" json:"name"`
	Location      *string    `db:"location" json:"location,omitempty"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	InitialBalance float64   `db:"initial_balance" json:"initial_balance"`
	CreatedBy     *uuid.UUID `db:"created_by" json:"created_by,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}

// CashRegisterSession represents a cash register session (opening/closing)
type CashRegisterSession struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	CashRegisterID  uuid.UUID  `db:"cash_register_id" json:"cash_register_id"`
	OpenedBy        uuid.UUID  `db:"opened_by" json:"opened_by"`
	ClosedBy        *uuid.UUID `db:"closed_by" json:"closed_by,omitempty"`
	OpeningBalance  float64    `db:"opening_balance" json:"opening_balance"`
	ExpectedCash    float64    `db:"expected_cash" json:"expected_cash"`
	CountedCash     *float64   `db:"counted_cash" json:"counted_cash,omitempty"`
	ExpectedCard    float64    `db:"expected_card" json:"expected_card"`
	ExpectedTransfer float64    `db:"expected_transfer" json:"expected_transfer"`
	Difference      *float64   `db:"difference" json:"difference,omitempty"`
	Status          string     `db:"status" json:"status"` // open, closed, cancelled
	Notes           *string    `db:"notes" json:"notes,omitempty"`
	JournalEntryID  *uuid.UUID `db:"journal_entry_id" json:"journal_entry_id,omitempty"`
	OpenedAt        time.Time  `db:"opened_at" json:"opened_at"`
	ClosedAt        *time.Time `db:"closed_at" json:"closed_at,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`

	// Relations
	CashRegisterName string    `db:"cash_register_name" json:"cash_register_name,omitempty"`
	OpenedByName     string    `db:"opened_by_name" json:"opened_by_name,omitempty"`
	ClosedByName     *string   `db:"closed_by_name" json:"closed_by_name,omitempty"`
	CountDetails     []CashCountDetail `json:"count_details,omitempty"`
}

// CashCountDetail represents cash count breakdown (bills and coins)
type CashCountDetail struct {
	ID          uuid.UUID `db:"id" json:"id"`
	SessionID   uuid.UUID `db:"session_id" json:"session_id"`
	Denomination float64  `db:"denomination" json:"denomination"`
	Quantity    int       `db:"quantity" json:"quantity"`
	Subtotal    float64   `db:"subtotal" json:"subtotal"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// CreateCashRegisterRequest represents the request to create a cash register
type CreateCashRegisterRequest struct {
	Name          string  `json:"name" validate:"required"`
	Location      *string `json:"location,omitempty"`
	InitialBalance float64 `json:"initial_balance" validate:"gte=0"`
}

// OpenSessionRequest represents the request to open a session
type OpenSessionRequest struct {
	CashRegisterID uuid.UUID `json:"cash_register_id" validate:"required"`
	OpeningBalance float64   `json:"opening_balance" validate:"gte=0"`
	Notes          *string   `json:"notes,omitempty"`
}

// CloseSessionRequest represents the request to close a session
type CloseSessionRequest struct {
	CountedCash float64            `json:"counted_cash" validate:"gte=0"`
	CountDetails []CashCountDetail  `json:"count_details,omitempty"`
	Notes       *string             `json:"notes,omitempty"`
}

// SessionSummary represents a summary of the session
type SessionSummary struct {
	Session           *CashRegisterSession `json:"session"`
	TotalSales        float64              `json:"total_sales"`
	SalesByPaymentType map[string]float64  `json:"sales_by_payment_type"`
	TotalTransactions int                  `json:"total_transactions"`
}


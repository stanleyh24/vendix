package accounting

import (
	"time"

	"github.com/google/uuid"
)

// Account types
const (
	AccountTypeAsset     = "ASSET"
	AccountTypeLiability = "LIABILITY"
	AccountTypeEquity    = "EQUITY"
	AccountTypeRevenue   = "REVENUE"
	AccountTypeExpense   = "EXPENSE"
)

// ChartOfAccounts represents an account in the chart of accounts
type ChartOfAccounts struct {
	ID                uuid.UUID `db:"id" json:"id"`
	AccountCode       string    `db:"account_code" json:"account_code"`
	AccountName       string    `db:"account_name" json:"account_name"`
	AccountType       string    `db:"account_type" json:"account_type"`
	ParentAccountCode *string   `db:"parent_account_code" json:"parent_account_code,omitempty"`
	IsActive          bool      `db:"is_active" json:"is_active"`
	Description       *string   `db:"description" json:"description,omitempty"`
	Balance           float64   `db:"balance" json:"balance"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

// JournalEntry represents a journal entry
type JournalEntry struct {
	ID          uuid.UUID          `db:"id" json:"id"`
	EntryNumber string             `db:"entry_number" json:"entry_number"`
	EntryDate   time.Time          `db:"entry_date" json:"entry_date"`
	Description string             `db:"description" json:"description"`
	Reference   *string            `db:"reference" json:"reference,omitempty"`
	Status      string             `db:"status" json:"status"`
	CreatedBy   *uuid.UUID         `db:"created_by" json:"created_by,omitempty"`
	CreatedAt   time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `db:"updated_at" json:"updated_at"`
	Lines       []JournalEntryLine `json:"lines,omitempty"`
}

// JournalEntryLine represents a line in a journal entry
type JournalEntryLine struct {
	ID             uuid.UUID `db:"id" json:"id"`
	JournalEntryID uuid.UUID `db:"journal_entry_id" json:"journal_entry_id"`
	AccountCode    string    `db:"account_code" json:"account_code"`
	AccountName    string    `db:"account_name" json:"account_name"`
	Debit          float64   `db:"debit" json:"debit"`
	Credit         float64   `db:"credit" json:"credit"`
	Description    *string   `db:"description" json:"description,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

// CreateAccountRequest represents the request to create an account
type CreateAccountRequest struct {
	AccountCode       string  `json:"account_code" validate:"required"`
	AccountName       string  `json:"account_name" validate:"required"`
	AccountType       string  `json:"account_type" validate:"required"`
	ParentAccountCode *string `json:"parent_account_code,omitempty"`
	Description       *string `json:"description,omitempty"`
}

// UpdateAccountRequest represents the request to update an account
type UpdateAccountRequest struct {
	AccountName       *string `json:"account_name,omitempty"`
	AccountType       *string `json:"account_type,omitempty"`
	ParentAccountCode *string `json:"parent_account_code,omitempty"`
	Description       *string `json:"description,omitempty"`
	IsActive          *bool   `json:"is_active,omitempty"`
}

// CreateJournalEntryRequest represents the request to create a journal entry
type CreateJournalEntryRequest struct {
	EntryDate   string                      `json:"entry_date" validate:"required"`
	Description string                      `json:"description" validate:"required"`
	Reference   *string                     `json:"reference,omitempty"`
	Lines       []CreateJournalEntryLineReq `json:"lines" validate:"required,min=2"`
}

// CreateJournalEntryLineReq represents a line in the create request
type CreateJournalEntryLineReq struct {
	AccountCode string  `json:"account_code" validate:"required"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
	Description *string `json:"description,omitempty"`
}

// LedgerEntry represents an entry in the general ledger
type LedgerEntry struct {
	Date           time.Time `db:"date" json:"date"`
	AccountCode    string    `db:"account_code" json:"account_code"`
	AccountName    string    `db:"account_name" json:"account_name"`
	Description    string    `db:"description" json:"description"`
	Reference      string    `db:"reference" json:"reference"`
	Debit          float64   `db:"debit" json:"debit"`
	Credit         float64   `db:"credit" json:"credit"`
	Balance        float64   `db:"balance" json:"balance"`
	JournalEntryID uuid.UUID `db:"journal_entry_id" json:"journal_entry_id"`
}

// TrialBalanceEntry represents an entry in the trial balance
type TrialBalanceEntry struct {
	AccountCode string  `db:"account_code" json:"account_code"`
	AccountName string  `db:"account_name" json:"account_name"`
	AccountType string  `db:"account_type" json:"account_type"`
	Debit       float64 `db:"debit" json:"debit"`
	Credit      float64 `db:"credit" json:"credit"`
}

// FinancialStatementReport represents a financial statement
type FinancialStatementReport struct {
	StatementType string                      `json:"statement_type"`
	PeriodStart   string                      `json:"period_start"`
	PeriodEnd     string                      `json:"period_end"`
	Sections      []FinancialStatementSection `json:"sections"`
	Total         float64                     `json:"total"`
}

// FinancialStatementSection represents a section in a financial statement
type FinancialStatementSection struct {
	Title    string                   `json:"title"`
	Accounts []FinancialStatementItem `json:"accounts"`
	Subtotal float64                  `json:"subtotal"`
}

// FinancialStatementItem represents an item in a financial statement
type FinancialStatementItem struct {
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	Amount      float64 `json:"amount"`
}

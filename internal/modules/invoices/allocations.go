package invoices

import (
	"context"
	"fmt"
	"time"

	"github.com/stanleyh24/vendix/internal/logger"

	"github.com/google/uuid"
)

// PaymentAllocation represents the allocation of a payment to an invoice
type PaymentAllocation struct {
	ID        uuid.UUID `db:"id" json:"id"`
	PaymentID uuid.UUID `db:"payment_id" json:"payment_id"`
	InvoiceID uuid.UUID `db:"invoice_id" json:"invoice_id"`
	Amount    float64   `db:"amount" json:"amount"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// CreatePaymentAllocationRequest represents the request to allocate a payment to invoices
type CreatePaymentAllocationRequest struct {
	PaymentID uuid.UUID                    `json:"payment_id" validate:"required"`
	Allocations []InvoiceAllocationRequest `json:"allocations" validate:"required,min=1"`
}

// InvoiceAllocationRequest represents allocation for a single invoice
type InvoiceAllocationRequest struct {
	InvoiceID uuid.UUID `json:"invoice_id" validate:"required"`
	Amount    float64   `json:"amount" validate:"required,gt=0"`
}

// AccountsReceivableItem represents an item in the accounts receivable report
type AccountsReceivableItem struct {
	InvoiceID      string    `db:"invoice_id" json:"invoice_id"`
	InvoiceNumber  string    `db:"invoice_number" json:"invoice_number"`
	NCF            *string   `db:"ncf" json:"ncf,omitempty"`
	CustomerID     string    `db:"customer_id" json:"customer_id"`
	CustomerName   string    `db:"customer_name" json:"customer_name"`
	IssueDate      time.Time `db:"issue_date" json:"issue_date"`
	DueDate        time.Time `db:"due_date" json:"due_date"`
	Total          float64   `db:"total" json:"total"`
	PaidAmount     float64   `db:"paid_amount" json:"paid_amount"`
	Balance        float64   `db:"balance" json:"balance"` // total - paid_amount
	Status         string    `db:"status" json:"status"`
	DaysOverdue    *int      `db:"days_overdue" json:"days_overdue,omitempty"`
	Age0_30        float64   `db:"age_0_30" json:"age_0_30"`   // Balance in 0-30 days
	Age31_60       float64   `db:"age_31_60" json:"age_31_60"` // Balance in 31-60 days
	Age61_90       float64   `db:"age_61_90" json:"age_61_90"` // Balance in 61-90 days
	AgeOver90      float64   `db:"age_over_90" json:"age_over_90"` // Balance over 90 days
}

// AccountsReceivableReport represents the complete AR aging report
type AccountsReceivableReport struct {
	ReportDate     string                  `json:"report_date"`
	TotalInvoices  int                     `json:"total_invoices"`
	TotalOutstanding float64                `json:"total_outstanding"` // Total balance due
	Total0_30       float64                 `json:"total_0_30"`
	Total31_60      float64                 `json:"total_31_60"`
	Total61_90      float64                 `json:"total_61_90"`
	TotalOver90     float64                 `json:"total_over_90"`
	Items          []AccountsReceivableItem `json:"items"`
	SummaryByCustomer map[string]CustomerARSummary `json:"summary_by_customer,omitempty"`
}

// CustomerARSummary represents AR summary grouped by customer
type CustomerARSummary struct {
	CustomerID     string  `json:"customer_id"`
	CustomerName   string  `json:"customer_name"`
	InvoiceCount   int     `json:"invoice_count"`
	TotalOutstanding float64 `json:"total_outstanding"`
	Total0_30      float64 `json:"total_0_30"`
	Total31_60     float64 `json:"total_31_60"`
	Total61_90     float64 `json:"total_61_90"`
	TotalOver90    float64 `json:"total_over_90"`
}

// CreateAllocation creates a payment allocation to an invoice
func (r *Repository) CreateAllocation(ctx context.Context, schema string, allocation *PaymentAllocation) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.payment_allocations (id, payment_id, invoice_id, amount, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, schema)
	_, err := r.db.ExecContext(ctx, query, allocation.ID, allocation.PaymentID, allocation.InvoiceID, allocation.Amount, allocation.CreatedAt)
	return err
}

// GetAllocationsByPayment gets all allocations for a payment
func (r *Repository) GetAllocationsByPayment(ctx context.Context, schema string, paymentID uuid.UUID) ([]*PaymentAllocation, error) {
	var allocations []*PaymentAllocation
	query := fmt.Sprintf(`
		SELECT id, payment_id, invoice_id, amount, created_at
		FROM %s.payment_allocations
		WHERE payment_id = $1
		ORDER BY created_at ASC
	`, schema)
	err := r.db.SelectContext(ctx, &allocations, query, paymentID)
	return allocations, err
}

// GetAllocationsByInvoice gets all allocations for an invoice
func (r *Repository) GetAllocationsByInvoice(ctx context.Context, schema string, invoiceID uuid.UUID) ([]*PaymentAllocation, error) {
	var allocations []*PaymentAllocation
	query := fmt.Sprintf(`
		SELECT pa.id, pa.payment_id, pa.invoice_id, pa.amount, pa.created_at
		FROM %s.payment_allocations pa
		WHERE pa.invoice_id = $1
		ORDER BY pa.created_at ASC
	`, schema)
	err := r.db.SelectContext(ctx, &allocations, query, invoiceID)
	return allocations, err
}

// UpdateInvoicePaidAmount updates the paid_amount field of an invoice
func (r *Repository) UpdateInvoicePaidAmount(ctx context.Context, schema string, invoiceID uuid.UUID, paidAmount float64) error {
	query := fmt.Sprintf(`
		UPDATE %s.invoices
		SET paid_amount = $1, updated_at = $2
		WHERE id = $3
	`, schema)
	_, err := r.db.ExecContext(ctx, query, paidAmount, time.Now(), invoiceID)
	return err
}

// UpdateInvoiceStatus updates the status of an invoice
func (r *Repository) UpdateInvoiceStatus(ctx context.Context, schema string, invoiceID uuid.UUID, status string) error {
	query := fmt.Sprintf(`
		UPDATE %s.invoices
		SET status = $1, updated_at = $2
		WHERE id = $3
	`, schema)
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), invoiceID)
	return err
}

// GetAccountsReceivableReport generates an accounts receivable aging report
func (r *Repository) GetAccountsReceivableReport(ctx context.Context, schema string, asOfDate time.Time) (*AccountsReceivableReport, error) {
	report := &AccountsReceivableReport{
		ReportDate:        asOfDate.Format("2006-01-02"),
		Items:             []AccountsReceivableItem{},
		SummaryByCustomer: make(map[string]CustomerARSummary),
	}

	query := fmt.Sprintf(`
		SELECT 
			i.id::text as invoice_id,
			i.invoice_number,
			i.ncf,
			i.customer_id::text as customer_id,
			c.name as customer_name,
			i.issue_date,
			i.due_date,
			i.total,
			i.paid_amount,
			(i.total - COALESCE(i.paid_amount, 0)) as balance,
			i.status,
			CASE 
				WHEN i.due_date < $1 AND i.status IN ('pending', 'overdue') THEN 
					($1::date - i.due_date)::int
				ELSE NULL
			END as days_overdue,
			CASE 
				WHEN i.due_date >= $1 - INTERVAL '30 days' AND i.due_date <= $1 THEN 
					(i.total - COALESCE(i.paid_amount, 0))
				WHEN i.due_date > $1 THEN
					(i.total - COALESCE(i.paid_amount, 0))
				ELSE 0
			END as age_0_30,
			CASE 
				WHEN i.due_date >= $1 - INTERVAL '60 days' AND i.due_date < $1 - INTERVAL '30 days' THEN 
					(i.total - COALESCE(i.paid_amount, 0))
				ELSE 0
			END as age_31_60,
			CASE 
				WHEN i.due_date >= $1 - INTERVAL '90 days' AND i.due_date < $1 - INTERVAL '60 days' THEN 
					(i.total - COALESCE(i.paid_amount, 0))
				ELSE 0
			END as age_61_90,
			CASE 
				WHEN i.due_date < $1 - INTERVAL '90 days' THEN 
					(i.total - COALESCE(i.paid_amount, 0))
				ELSE 0
			END as age_over_90
		FROM %s.invoices i
		INNER JOIN %s.customers c ON i.customer_id = c.id
		WHERE i.status IN ('pending', 'overdue')
		  AND (i.total - COALESCE(i.paid_amount, 0)) > 0
		ORDER BY i.due_date ASC, i.invoice_number ASC
	`, schema, schema)

	var items []AccountsReceivableItem
	err := r.db.SelectContext(ctx, &items, query, asOfDate)
	if err != nil {
		logger.Error("Failed to get AR report", "error", err)
		return nil, fmt.Errorf("failed to get AR report: %w", err)
	}

	report.Items = items
	report.TotalInvoices = len(items)

	// Calculate totals and group by customer
	for _, item := range items {
		report.TotalOutstanding += item.Balance
		report.Total0_30 += item.Age0_30
		report.Total31_60 += item.Age31_60
		report.Total61_90 += item.Age61_90
		report.TotalOver90 += item.AgeOver90

		// Group by customer
		customerSummary := report.SummaryByCustomer[item.CustomerID]
		customerSummary.CustomerID = item.CustomerID
		customerSummary.CustomerName = item.CustomerName
		customerSummary.InvoiceCount++
		customerSummary.TotalOutstanding += item.Balance
		customerSummary.Total0_30 += item.Age0_30
		customerSummary.Total31_60 += item.Age31_60
		customerSummary.Total61_90 += item.Age61_90
		customerSummary.TotalOver90 += item.AgeOver90
		report.SummaryByCustomer[item.CustomerID] = customerSummary
	}

	return report, nil
}


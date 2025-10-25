package invoices

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"vendix/internal/database"

	"github.com/google/uuid"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, schema string, invoice *Invoice) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert invoice
	query := fmt.Sprintf(`
		INSERT INTO %s.invoices (id, invoice_number, ncf, ncf_type, customer_id, issue_date, due_date, status, subtotal, tax_amount, total, paid_amount, currency, notes, terms, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`, schema)

	_, err = tx.ExecContext(ctx, query,
		invoice.ID, invoice.InvoiceNumber, invoice.NCF, invoice.NCFType, invoice.CustomerID,
		invoice.IssueDate, invoice.DueDate, invoice.Status, invoice.Subtotal,
		invoice.TaxAmount, invoice.Total, invoice.PaidAmount, invoice.Currency,
		invoice.Notes, invoice.Terms, invoice.CreatedBy, invoice.CreatedAt, invoice.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Insert lines
	lineQuery := fmt.Sprintf(`
		INSERT INTO %s.invoice_lines (id, invoice_id, product_id, line_number, description, quantity, unit_price, tax_rate, tax_amount, line_total, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, schema)

	for _, line := range invoice.Lines {
		_, err = tx.ExecContext(ctx, lineQuery,
			line.ID, line.InvoiceID, line.ProductID, line.LineNumber,
			line.Description, line.Quantity, line.UnitPrice, line.TaxRate,
			line.TaxAmount, line.LineTotal, line.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) GetByID(ctx context.Context, schema string, id uuid.UUID) (*Invoice, error) {
	var invoice Invoice
	query := fmt.Sprintf(`
		SELECT i.id, i.invoice_number, i.ncf, i.ncf_type, i.customer_id, i.issue_date, i.due_date, 
		       i.status, i.subtotal, i.tax_amount, i.total, i.paid_amount, i.currency,
		       i.notes, i.terms, i.dgii_status, i.signed_at, i.sent_at, i.created_by,
		       i.created_at, i.updated_at, c.name as customer_name
		FROM %s.invoices i
		INNER JOIN %s.customers c ON i.customer_id = c.id
		WHERE i.id = $1
	`, schema, schema)

	err := r.db.GetContext(ctx, &invoice, query, id)
	if err != nil {
		return nil, err
	}

	// Get lines
	lines, err := r.getInvoiceLines(ctx, schema, id)
	if err != nil {
		return nil, err
	}

	invoice.Lines = lines
	return &invoice, nil
}

func (r *Repository) getInvoiceLines(ctx context.Context, schema string, invoiceID uuid.UUID) ([]InvoiceLine, error) {
	var lines []InvoiceLine
	query := fmt.Sprintf(`
		SELECT id, invoice_id, product_id, line_number, description, quantity, unit_price, tax_rate, tax_amount, line_total, created_at
		FROM %s.invoice_lines
		WHERE invoice_id = $1
		ORDER BY line_number
	`, schema)

	err := r.db.SelectContext(ctx, &lines, query, invoiceID)
	return lines, err
}

func (r *Repository) List(ctx context.Context, schema string, status *string) ([]*Invoice, error) {
	var invoices []*Invoice
	baseQuery := fmt.Sprintf(`
		SELECT i.id, i.invoice_number, i.ncf, i.ncf_type, i.customer_id, i.issue_date, i.due_date,
		       i.status, i.subtotal, i.tax_amount, i.total, i.paid_amount, i.currency,
		       i.notes, i.terms, i.dgii_status, i.signed_at, i.sent_at, i.created_by,
		       i.created_at, i.updated_at, c.name as customer_name
		FROM %s.invoices i
		INNER JOIN %s.customers c ON i.customer_id = c.id
		WHERE 1=1
	`, schema, schema)

	args := []interface{}{}
	argCount := 1

	if status != nil {
		baseQuery += fmt.Sprintf(" AND i.status = $%d", argCount)
		args = append(args, *status)
		argCount++
	}

	baseQuery += " ORDER BY i.issue_date DESC, i.invoice_number DESC"

	err := r.db.SelectContext(ctx, &invoices, baseQuery, args...)
	if err != nil {
		return nil, err
	}

	// Get lines for each invoice
	for i, invoice := range invoices {
		lines, err := r.getInvoiceLines(ctx, schema, invoice.ID)
		if err != nil {
			return nil, err
		}
		invoices[i].Lines = lines
	}

	return invoices, nil
}

func (r *Repository) Update(ctx context.Context, schema string, invoice *Invoice) error {
	query := fmt.Sprintf(`
		UPDATE %s.invoices
		SET status = $1, notes = $2, terms = $3, updated_at = $4
		WHERE id = $5
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		invoice.Status, invoice.Notes, invoice.Terms, invoice.UpdatedAt, invoice.ID,
	)
	return err
}

func (r *Repository) UpdateStatus(ctx context.Context, schema string, id uuid.UUID, status string) error {
	query := fmt.Sprintf(`
		UPDATE %s.invoices
		SET status = $1, updated_at = $2
		WHERE id = $3
	`, schema)

	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (r *Repository) GetNextInvoiceNumber(ctx context.Context, schema string, prefix string) (string, error) {
	// Use a more robust query that gets the maximum number atomically
	var maxNum sql.NullInt64
	query := fmt.Sprintf(`
		SELECT COALESCE(
			MAX(
				CAST(
					SUBSTRING(invoice_number FROM LENGTH($1) + 1) AS INTEGER
				)
			), 
			0
		) as max_num
		FROM %s.invoices 
		WHERE invoice_number ~ ('^' || $1 || '[0-9]+$')
		  AND LENGTH(invoice_number) = LENGTH($1) + 8
	`, schema)

	err := r.db.GetContext(ctx, &maxNum, query, prefix)
	if err != nil {
		return "", fmt.Errorf("failed to get next invoice number: %w", err)
	}

	// Increment the number
	nextNum := int(maxNum.Int64) + 1

	// Format with padding
	return fmt.Sprintf("%s%08d", prefix, nextNum), nil
}

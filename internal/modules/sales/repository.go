package sales

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

func (r *Repository) Create(ctx context.Context, schema string, sale *Sale) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert sale
	query := fmt.Sprintf(`
		INSERT INTO %s.sales (id, sale_number, customer_id, invoice_id, subtotal, tax_amount, total, payment_type, status, notes, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, schema)

	_, err = tx.ExecContext(ctx, query,
		sale.ID, sale.SaleNumber, sale.CustomerID, sale.InvoiceID, sale.Subtotal,
		sale.TaxAmount, sale.Total, sale.PaymentType, sale.Status, sale.Notes,
		sale.CreatedBy, sale.CreatedAt, sale.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Insert lines
	lineQuery := fmt.Sprintf(`
		INSERT INTO %s.sale_lines (id, sale_id, product_id, line_number, description, quantity, unit_price, tax_rate, tax_amount, line_total, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, schema)

	for _, line := range sale.Lines {
		_, err = tx.ExecContext(ctx, lineQuery,
			line.ID, line.SaleID, line.ProductID, line.LineNumber,
			line.Description, line.Quantity, line.UnitPrice, line.TaxRate,
			line.TaxAmount, line.LineTotal, line.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) GetByID(ctx context.Context, schema string, id uuid.UUID) (*Sale, error) {
	var sale Sale
	query := fmt.Sprintf(`
		SELECT s.id, s.sale_number, s.customer_id, s.invoice_id, s.subtotal, s.tax_amount, s.total, 
		       s.payment_type, s.status, s.notes, s.created_by, s.created_at, s.updated_at,
		       COALESCE(c.name, 'Cliente Genérico') as customer_name
		FROM %s.sales s
		LEFT JOIN %s.customers c ON s.customer_id = c.id
		WHERE s.id = $1
	`, schema, schema)

	err := r.db.GetContext(ctx, &sale, query, id)
	if err != nil {
		return nil, err
	}

	// Get lines
	lines, err := r.getSaleLines(ctx, schema, id)
	if err != nil {
		return nil, err
	}

	sale.Lines = lines
	return &sale, nil
}

func (r *Repository) getSaleLines(ctx context.Context, schema string, saleID uuid.UUID) ([]SaleLine, error) {
	var lines []SaleLine
	query := fmt.Sprintf(`
		SELECT id, sale_id, product_id, line_number, description, quantity, COALESCE(returned_quantity, 0) as returned_quantity, unit_price, tax_rate, tax_amount, line_total, created_at
		FROM %s.sale_lines
		WHERE sale_id = $1
		ORDER BY line_number
	`, schema)

	err := r.db.SelectContext(ctx, &lines, query, saleID)
	return lines, err
}

func (r *Repository) List(ctx context.Context, schema string, status *string) ([]*Sale, error) {
	var sales []*Sale
	baseQuery := fmt.Sprintf(`
		SELECT s.id, s.sale_number, s.customer_id, s.invoice_id, s.subtotal, s.tax_amount, s.total,
		       s.payment_type, s.status, s.notes, s.created_by, s.created_at, s.updated_at,
		       COALESCE(c.name, 'Cliente Genérico') as customer_name
		FROM %s.sales s
		LEFT JOIN %s.customers c ON s.customer_id = c.id
		WHERE 1=1
	`, schema, schema)

	args := []interface{}{}
	argCount := 1

	if status != nil {
		baseQuery += fmt.Sprintf(" AND s.status = $%d", argCount)
		args = append(args, *status)
		argCount++
	}

	baseQuery += " ORDER BY s.created_at DESC, s.sale_number DESC"

	err := r.db.SelectContext(ctx, &sales, baseQuery, args...)
	if err != nil {
		return nil, err
	}

	// Get lines for each sale
	for i, sale := range sales {
		lines, err := r.getSaleLines(ctx, schema, sale.ID)
		if err != nil {
			return nil, err
		}
		sales[i].Lines = lines
	}

	return sales, nil
}

func (r *Repository) Update(ctx context.Context, schema string, sale *Sale) error {
	query := fmt.Sprintf(`
		UPDATE %s.sales
		SET status = $1, notes = $2, updated_at = $3
		WHERE id = $4
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		sale.Status, sale.Notes, sale.UpdatedAt, sale.ID,
	)
	return err
}

func (r *Repository) GetNextSaleNumber(ctx context.Context, schema string, prefix string) (string, error) {
	// Use a more robust query that gets the maximum number atomically
	var maxNum sql.NullInt64
	query := fmt.Sprintf(`
		SELECT COALESCE(
			MAX(
				CAST(
					SUBSTRING(sale_number FROM LENGTH($1) + 1) AS INTEGER
				)
			), 
			0
		) as max_num
		FROM %s.sales 
		WHERE sale_number ~ ('^' || $1 || '[0-9]+$')
		  AND LENGTH(sale_number) = LENGTH($1) + 8
	`, schema)

	err := r.db.GetContext(ctx, &maxNum, query, prefix)
	if err != nil {
		return "", fmt.Errorf("failed to get next sale number: %w", err)
	}

	// Increment the number
	nextNum := int(maxNum.Int64) + 1

	// Format with padding
	return fmt.Sprintf("%s%08d", prefix, nextNum), nil
}

// DailySalesSummary represents a summary of sales for a specific day
type DailySalesSummary struct {
	Date              string  `db:"sale_date" json:"date"`
	CashSubtotal      float64 `db:"cash_subtotal" json:"cash_subtotal"`
	CashTax           float64 `db:"cash_tax" json:"cash_tax"`
	CashTotal         float64 `db:"cash_total" json:"cash_total"`
	CardSubtotal      float64 `db:"card_subtotal" json:"card_subtotal"`
	CardTax           float64 `db:"card_tax" json:"card_tax"`
	CardTotal         float64 `db:"card_total" json:"card_total"`
	TransferSubtotal  float64 `db:"transfer_subtotal" json:"transfer_subtotal"`
	TransferTax       float64 `db:"transfer_tax" json:"transfer_tax"`
	TransferTotal     float64 `db:"transfer_total" json:"transfer_total"`
	CreditSubtotal    float64 `db:"credit_subtotal" json:"credit_subtotal"`
	CreditTax         float64 `db:"credit_tax" json:"credit_tax"`
	CreditTotal       float64 `db:"credit_total" json:"credit_total"`
	TotalSubtotal     float64 `db:"total_subtotal" json:"total_subtotal"`
	TotalTax          float64 `db:"total_tax" json:"total_tax"`
	TotalITBIS        float64 `db:"total_itbis" json:"total_itbis"`
	TotalSelective    float64 `db:"total_selective" json:"total_selective"`
	TotalAmount       float64 `db:"total_amount" json:"total_amount"`
	TotalCost         float64 `db:"total_cost" json:"total_cost"` // Costo total de productos vendidos (COGS)
	TotalTransactions int     `db:"total_transactions" json:"total_transactions"`
}

// GetDailySalesSummary gets a summary of sales for a specific date
func (r *Repository) GetDailySalesSummary(ctx context.Context, schema string, date string) (*DailySalesSummary, error) {
	var summary DailySalesSummary
	query := fmt.Sprintf(`
		SELECT 
			DATE(s.created_at) as sale_date,
			COALESCE(SUM(CASE WHEN s.payment_type = 'cash' THEN s.subtotal ELSE 0 END), 0) as cash_subtotal,
			COALESCE(SUM(CASE WHEN s.payment_type = 'cash' THEN s.tax_amount ELSE 0 END), 0) as cash_tax,
			COALESCE(SUM(CASE WHEN s.payment_type = 'cash' THEN s.total ELSE 0 END), 0) as cash_total,
			COALESCE(SUM(CASE WHEN s.payment_type = 'card' THEN s.subtotal ELSE 0 END), 0) as card_subtotal,
			COALESCE(SUM(CASE WHEN s.payment_type = 'card' THEN s.tax_amount ELSE 0 END), 0) as card_tax,
			COALESCE(SUM(CASE WHEN s.payment_type = 'card' THEN s.total ELSE 0 END), 0) as card_total,
			COALESCE(SUM(CASE WHEN s.payment_type = 'transfer' THEN s.subtotal ELSE 0 END), 0) as transfer_subtotal,
			COALESCE(SUM(CASE WHEN s.payment_type = 'transfer' THEN s.tax_amount ELSE 0 END), 0) as transfer_tax,
			COALESCE(SUM(CASE WHEN s.payment_type = 'transfer' THEN s.total ELSE 0 END), 0) as transfer_total,
			COALESCE(SUM(CASE WHEN s.payment_type = 'check' OR (s.payment_type NOT IN ('cash', 'card', 'transfer') AND (i.status = 'pending' OR i.status = 'overdue')) THEN s.subtotal ELSE 0 END), 0) as credit_subtotal,
			COALESCE(SUM(CASE WHEN s.payment_type = 'check' OR (s.payment_type NOT IN ('cash', 'card', 'transfer') AND (i.status = 'pending' OR i.status = 'overdue')) THEN s.tax_amount ELSE 0 END), 0) as credit_tax,
			COALESCE(SUM(CASE WHEN s.payment_type = 'check' OR (s.payment_type NOT IN ('cash', 'card', 'transfer') AND (i.status = 'pending' OR i.status = 'overdue')) THEN s.total ELSE 0 END), 0) as credit_total,
			COALESCE(SUM(s.subtotal), 0) as total_subtotal,
			COALESCE(SUM(s.tax_amount), 0) as total_tax,
			COALESCE(SUM(s.tax_amount), 0) as total_itbis,
			0.0 as total_selective,
			COALESCE(SUM(s.total), 0) as total_amount,
			COALESCE(SUM(sl.quantity * COALESCE(p.cost, 0)), 0) as total_cost,
			COUNT(DISTINCT s.id) as total_transactions
		FROM %s.sales s
		LEFT JOIN %s.invoices i ON s.invoice_id = i.id
		LEFT JOIN %s.sale_lines sl ON s.id = sl.sale_id
		LEFT JOIN %s.products p ON sl.product_id = p.id
		WHERE DATE(s.created_at) = $1
		  AND s.status = 'completed'
		GROUP BY DATE(s.created_at)
	`, schema, schema, schema, schema)

	err := r.db.GetContext(ctx, &summary, query, date)
	if err != nil {
		// If no sales found, return empty summary with the date
		if err == sql.ErrNoRows {
			summary.Date = date
			return &summary, nil
		}
		return nil, fmt.Errorf("failed to get daily sales summary: %w", err)
	}

	return &summary, nil
}

// GetSalesReportData gets sales data for reporting within a date range
func (r *Repository) GetSalesReportData(ctx context.Context, schema string, startDate, endDate time.Time) ([]SalesReportRow, error) {
	var rows []SalesReportRow
	query := fmt.Sprintf(`
		SELECT 
			s.id::text as sale_id,
			s.sale_number,
			DATE(s.created_at) as sale_date,
			s.customer_id::text as customer_id,
			COALESCE(c.name, 'Cliente Genérico') as customer_name,
			s.payment_type,
			s.subtotal,
			s.tax_amount,
			s.total,
			s.status,
			COALESCE(SUM(sl.quantity * COALESCE(p.cost, 0)), 0) as cost
		FROM %s.sales s
		LEFT JOIN %s.customers c ON s.customer_id = c.id
		LEFT JOIN %s.sale_lines sl ON s.id = sl.sale_id
		LEFT JOIN %s.products p ON sl.product_id = p.id
		WHERE DATE(s.created_at) >= $1
		  AND DATE(s.created_at) <= $2
		  AND s.status = 'completed'
		GROUP BY s.id, s.sale_number, DATE(s.created_at), s.customer_id, c.name, s.payment_type, s.subtotal, s.tax_amount, s.total, s.status
		ORDER BY s.created_at DESC
	`, schema, schema, schema, schema)

	err := r.db.SelectContext(ctx, &rows, query, startDate, endDate)
	return rows, err
}

// SalesReportRow represents a row in the sales report query
type SalesReportRow struct {
	SaleID       string     `db:"sale_id"`
	SaleNumber   string     `db:"sale_number"`
	SaleDate     time.Time  `db:"sale_date"`
	CustomerID   *string    `db:"customer_id"`
	CustomerName string     `db:"customer_name"`
	PaymentType  string     `db:"payment_type"`
	Subtotal     float64    `db:"subtotal"`
	TaxAmount    float64    `db:"tax_amount"`
	Total        float64    `db:"total"`
	Status       string     `db:"status"`
	Cost         float64    `db:"cost"`
}

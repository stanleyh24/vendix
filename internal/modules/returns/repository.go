package returns

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

// Create creates a new return
func (r *Repository) Create(ctx context.Context, schema string, ret *Return) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert return
	query := fmt.Sprintf(`
		INSERT INTO %s.returns (id, return_number, sale_id, invoice_id, customer_id, return_date, return_type, return_reason, notes, 
		                        subtotal, tax_amount, total, status, refund_status, refund_method, refund_amount, 
		                        created_by, approved_by, approved_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	`, schema)

	var returnReason *string
	if ret.ReturnReason != nil {
		reason := string(*ret.ReturnReason)
		returnReason = &reason
	}

	var refundMethod *string
	if ret.RefundMethod != nil {
		method := string(*ret.RefundMethod)
		refundMethod = &method
	}

	_, err = tx.ExecContext(ctx, query,
		ret.ID, ret.ReturnNumber, ret.SaleID, ret.InvoiceID, ret.CustomerID, ret.ReturnDate,
		ret.ReturnType, returnReason, ret.Notes, ret.Subtotal, ret.TaxAmount, ret.Total,
		ret.Status, ret.RefundStatus, refundMethod, ret.RefundAmount,
		ret.CreatedBy, ret.ApprovedBy, ret.ApprovedAt, ret.CreatedAt, ret.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Insert lines
	lineQuery := fmt.Sprintf(`
		INSERT INTO %s.return_lines (id, return_id, sale_line_id, invoice_line_id, product_id, line_number, description, 
		                             quantity, unit_price, tax_rate, tax_amount, line_total, item_condition, restock_quantity, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`, schema)

	for _, line := range ret.Lines {
		var condition *string
		if line.ItemCondition != nil {
			cond := string(*line.ItemCondition)
			condition = &cond
		}

		_, err = tx.ExecContext(ctx, lineQuery,
			line.ID, line.ReturnID, line.SaleLineID, line.InvoiceLineID, line.ProductID, line.LineNumber,
			line.Description, line.Quantity, line.UnitPrice, line.TaxRate, line.TaxAmount, line.LineTotal,
			condition, line.RestockQuantity, line.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetByID retrieves a return by ID
func (r *Repository) GetByID(ctx context.Context, schema string, id uuid.UUID) (*Return, error) {
	var ret Return
	query := fmt.Sprintf(`
		SELECT r.id, r.return_number, r.sale_id, r.invoice_id, r.customer_id, r.return_date, r.return_type, r.return_reason, r.notes,
		       r.subtotal, r.tax_amount, r.total, r.status, r.refund_status, r.refund_method, r.refund_amount,
		       r.created_by, r.approved_by, r.approved_at, r.created_at, r.updated_at,
		       COALESCE(c.name, 'Cliente Genérico') as customer_name,
		       s.sale_number, i.invoice_number
		FROM %s.returns r
		LEFT JOIN %s.customers c ON r.customer_id = c.id
		LEFT JOIN %s.sales s ON r.sale_id = s.id
		LEFT JOIN %s.invoices i ON r.invoice_id = i.id
		WHERE r.id = $1
	`, schema, schema, schema, schema)

	err := r.db.GetContext(ctx, &ret, query, id)
	if err != nil {
		return nil, err
	}

	// Parse enums
	if ret.ReturnReason != nil {
		reason := ReturnReason(*ret.ReturnReason)
		ret.ReturnReason = &reason
	}
	if ret.RefundMethod != nil {
		method := RefundMethod(*ret.RefundMethod)
		ret.RefundMethod = &method
	}

	// Get lines
	lines, err := r.getReturnLines(ctx, schema, id)
	if err != nil {
		return nil, err
	}

	ret.Lines = lines
	return &ret, nil
}

// getReturnLines retrieves all lines for a return
func (r *Repository) getReturnLines(ctx context.Context, schema string, returnID uuid.UUID) ([]ReturnLine, error) {
	var lines []ReturnLine
	query := fmt.Sprintf(`
		SELECT id, return_id, sale_line_id, invoice_line_id, product_id, line_number, description,
		       quantity, unit_price, tax_rate, tax_amount, line_total, item_condition, restock_quantity, created_at
		FROM %s.return_lines
		WHERE return_id = $1
		ORDER BY line_number
	`, schema)

	err := r.db.SelectContext(ctx, &lines, query, returnID)
	if err != nil {
		return nil, err
	}

	// Parse item_condition enum
	for i := range lines {
		if lines[i].ItemCondition != nil {
			cond := ItemCondition(*lines[i].ItemCondition)
			lines[i].ItemCondition = &cond
		}
	}

	return lines, nil
}

// List retrieves all returns with optional filters
func (r *Repository) List(ctx context.Context, schema string, filters *ReturnFilters) ([]*Return, error) {
	query := fmt.Sprintf(`
		SELECT r.id, r.return_number, r.sale_id, r.invoice_id, r.customer_id, r.return_date, r.return_type, r.return_reason, r.notes,
		       r.subtotal, r.tax_amount, r.total, r.status, r.refund_status, r.refund_method, r.refund_amount,
		       r.created_by, r.approved_by, r.approved_at, r.created_at, r.updated_at,
		       COALESCE(c.name, 'Cliente Genérico') as customer_name,
		       s.sale_number, i.invoice_number
		FROM %s.returns r
		LEFT JOIN %s.customers c ON r.customer_id = c.id
		LEFT JOIN %s.sales s ON r.sale_id = s.id
		LEFT JOIN %s.invoices i ON r.invoice_id = i.id
		WHERE 1=1
	`, schema, schema, schema, schema)

	args := []interface{}{}
	argNum := 1

	if filters != nil {
		if filters.Status != nil && *filters.Status != "" {
			query += fmt.Sprintf(" AND r.status = $%d", argNum)
			args = append(args, *filters.Status)
			argNum++
		}
		if filters.SaleID != nil {
			query += fmt.Sprintf(" AND r.sale_id = $%d", argNum)
			args = append(args, *filters.SaleID)
			argNum++
		}
		if filters.InvoiceID != nil {
			query += fmt.Sprintf(" AND r.invoice_id = $%d", argNum)
			args = append(args, *filters.InvoiceID)
			argNum++
		}
		if filters.CustomerID != nil {
			query += fmt.Sprintf(" AND r.customer_id = $%d", argNum)
			args = append(args, *filters.CustomerID)
			argNum++
		}
		if filters.StartDate != nil {
			query += fmt.Sprintf(" AND r.return_date >= $%d", argNum)
			args = append(args, *filters.StartDate)
			argNum++
		}
		if filters.EndDate != nil {
			query += fmt.Sprintf(" AND r.return_date <= $%d", argNum)
			args = append(args, *filters.EndDate)
			argNum++
		}
	}

	query += " ORDER BY r.return_date DESC, r.created_at DESC"

	var returns []*Return
	err := r.db.SelectContext(ctx, &returns, query, args...)
	if err != nil {
		return nil, err
	}

	// Parse enums for each return
	for _, ret := range returns {
		if ret.ReturnReason != nil {
			reason := ReturnReason(*ret.ReturnReason)
			ret.ReturnReason = &reason
		}
		if ret.RefundMethod != nil {
			method := RefundMethod(*ret.RefundMethod)
			ret.RefundMethod = &method
		}
	}

	return returns, nil
}

// ReturnFilters represents filters for listing returns
type ReturnFilters struct {
	Status     *string
	SaleID     *uuid.UUID
	InvoiceID  *uuid.UUID
	CustomerID *uuid.UUID
	StartDate  *time.Time
	EndDate    *time.Time
}

// UpdateStatus updates the status of a return
func (r *Repository) UpdateStatus(ctx context.Context, schema string, id uuid.UUID, status ReturnStatus, approvedBy *uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.returns 
		SET status = $1, approved_by = $2, approved_at = $3, updated_at = $4
		WHERE id = $5
	`, schema)

	var approvedAt *time.Time
	if approvedBy != nil && (status == ReturnStatusApproved || status == ReturnStatusRejected) {
		now := time.Now()
		approvedAt = &now
	}

	_, err := r.db.ExecContext(ctx, query, status, approvedBy, approvedAt, time.Now(), id)
	return err
}

// UpdateRefundInfo updates refund information
func (r *Repository) UpdateRefundInfo(ctx context.Context, schema string, id uuid.UUID, refundStatus *string, refundMethod *RefundMethod, refundAmount float64) error {
	query := fmt.Sprintf(`
		UPDATE %s.returns 
		SET refund_status = $1, refund_method = $2, refund_amount = $3, updated_at = $4
		WHERE id = $5
	`, schema)

	var method *string
	if refundMethod != nil {
		m := string(*refundMethod)
		method = &m
	}

	_, err := r.db.ExecContext(ctx, query, refundStatus, method, refundAmount, time.Now(), id)
	return err
}

// GetNextReturnNumber generates the next return number
func (r *Repository) GetNextReturnNumber(ctx context.Context, schema string) (string, error) {
	query := fmt.Sprintf(`
		SELECT COALESCE(MAX(CAST(SUBSTRING(return_number FROM 'RET-(.*)') AS INTEGER)), 0) + 1
		FROM %s.returns
		WHERE return_number ~ '^RET-[0-9]+$'
	`, schema)

	var nextNum int
	err := r.db.GetContext(ctx, &nextNum, query)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	return fmt.Sprintf("RET-%06d", nextNum), nil
}

// UpdateReturnedQuantity updates the returned_quantity in sale_lines or invoice_lines
func (r *Repository) UpdateReturnedQuantity(ctx context.Context, schema string, saleLineID, invoiceLineID *uuid.UUID, quantity float64) error {
	if saleLineID != nil {
		query := fmt.Sprintf(`
			UPDATE %s.sale_lines 
			SET returned_quantity = returned_quantity + $1, updated_at = $2
			WHERE id = $3
		`, schema)
		_, err := r.db.ExecContext(ctx, query, quantity, time.Now(), *saleLineID)
		return err
	}

	if invoiceLineID != nil {
		query := fmt.Sprintf(`
			UPDATE %s.invoice_lines 
			SET returned_quantity = returned_quantity + $1
			WHERE id = $2
		`, schema)
		_, err := r.db.ExecContext(ctx, query, quantity, *invoiceLineID)
		return err
	}

	return fmt.Errorf("either sale_line_id or invoice_line_id must be provided")
}

// GetSaleLineByID gets a sale line for validation
func (r *Repository) GetSaleLineByID(ctx context.Context, schema string, saleLineID uuid.UUID) (*SaleLineInfo, error) {
	query := fmt.Sprintf(`
		SELECT id, sale_id, product_id, quantity, returned_quantity, unit_price, tax_rate
		FROM %s.sale_lines
		WHERE id = $1
	`, schema)

	var info SaleLineInfo
	err := r.db.GetContext(ctx, &info, query, saleLineID)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// GetInvoiceLineByID gets an invoice line for validation
func (r *Repository) GetInvoiceLineByID(ctx context.Context, schema string, invoiceLineID uuid.UUID) (*InvoiceLineInfo, error) {
	query := fmt.Sprintf(`
		SELECT id, invoice_id, product_id, quantity, returned_quantity, unit_price, tax_rate
		FROM %s.invoice_lines
		WHERE id = $1
	`, schema)

	var info InvoiceLineInfo
	err := r.db.GetContext(ctx, &info, query, invoiceLineID)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// SaleLineInfo represents sale line information for return validation
type SaleLineInfo struct {
	ID             uuid.UUID  `db:"id"`
	SaleID         uuid.UUID  `db:"sale_id"`
	ProductID      *uuid.UUID `db:"product_id"`
	Quantity       float64    `db:"quantity"`
	ReturnedQuantity float64  `db:"returned_quantity"`
	UnitPrice      float64    `db:"unit_price"`
	TaxRate        float64    `db:"tax_rate"`
}

// InvoiceLineInfo represents invoice line information for return validation
type InvoiceLineInfo struct {
	ID             uuid.UUID  `db:"id"`
	InvoiceID      uuid.UUID  `db:"invoice_id"`
	ProductID      *uuid.UUID `db:"product_id"`
	Quantity       float64    `db:"quantity"`
	ReturnedQuantity float64  `db:"returned_quantity"`
	UnitPrice      float64    `db:"unit_price"`
	TaxRate        float64    `db:"tax_rate"`
}

// GetReturnsReportData gets returns data for reporting within a date range
func (r *Repository) GetReturnsReportData(ctx context.Context, schema string, startDate, endDate time.Time) ([]ReturnsReportRow, error) {
	var rows []ReturnsReportRow
	query := fmt.Sprintf(`
		SELECT 
			r.id::text as return_id,
			r.return_number,
			DATE(r.return_date) as return_date,
			r.return_type,
			COALESCE(s.sale_number, i.invoice_number) as original_number,
			r.customer_id::text as customer_id,
			COALESCE(c.name, 'Cliente Genérico') as customer_name,
			r.return_reason,
			r.status,
			r.refund_method,
			r.refund_status,
			r.subtotal,
			r.tax_amount,
			r.total,
			COALESCE(r.refund_amount, 0) as refund_amount
		FROM %s.returns r
		LEFT JOIN %s.customers c ON r.customer_id = c.id
		LEFT JOIN %s.sales s ON r.sale_id = s.id
		LEFT JOIN %s.invoices i ON r.invoice_id = i.id
		WHERE DATE(r.return_date) >= $1
		  AND DATE(r.return_date) <= $2
		ORDER BY r.return_date DESC, r.created_at DESC
	`, schema, schema, schema, schema)

	err := r.db.SelectContext(ctx, &rows, query, startDate, endDate)
	return rows, err
}

// GetReturnsReportProductsData gets product returns data for reporting within a date range
func (r *Repository) GetReturnsReportProductsData(ctx context.Context, schema string, startDate, endDate time.Time) ([]ReturnsReportProductRow, error) {
	var rows []ReturnsReportProductRow
	query := fmt.Sprintf(`
		SELECT 
			rl.product_id::text as product_id,
			COALESCE(p.code, 'N/A') as product_code,
			COALESCE(p.name, rl.description) as product_name,
			SUM(rl.quantity) as quantity,
			COUNT(DISTINCT rl.return_id) as return_count,
			SUM(rl.line_total - rl.tax_amount) as subtotal,
			SUM(rl.tax_amount) as tax,
			SUM(rl.line_total) as total
		FROM %s.return_lines rl
		INNER JOIN %s.returns r ON rl.return_id = r.id
		LEFT JOIN %s.products p ON rl.product_id = p.id
		WHERE DATE(r.return_date) >= $1
		  AND DATE(r.return_date) <= $2
		  AND rl.product_id IS NOT NULL
		GROUP BY rl.product_id, p.code, p.name, rl.description
		ORDER BY quantity DESC, total DESC
		LIMIT 20
	`, schema, schema, schema)

	err := r.db.SelectContext(ctx, &rows, query, startDate, endDate)
	return rows, err
}

// ReturnsReportRow represents a row in the returns report query
type ReturnsReportRow struct {
	ReturnID       string     `db:"return_id"`
	ReturnNumber   string     `db:"return_number"`
	ReturnDate     time.Time  `db:"return_date"`
	ReturnType     string     `db:"return_type"`
	OriginalNumber *string    `db:"original_number"`
	CustomerID     *string    `db:"customer_id"`
	CustomerName   string     `db:"customer_name"`
	ReturnReason   *string    `db:"return_reason"`
	Status         string     `db:"status"`
	RefundMethod   *string    `db:"refund_method"`
	RefundStatus   *string    `db:"refund_status"`
	Subtotal       float64    `db:"subtotal"`
	TaxAmount      float64    `db:"tax_amount"`
	Total          float64    `db:"total"`
	RefundAmount   float64    `db:"refund_amount"`
}

// ReturnsReportProductRow represents a row in the product returns report query
type ReturnsReportProductRow struct {
	ProductID   string  `db:"product_id"`
	ProductCode string  `db:"product_code"`
	ProductName string  `db:"product_name"`
	Quantity    float64 `db:"quantity"`
	ReturnCount int     `db:"return_count"`
	Subtotal    float64 `db:"subtotal"`
	Tax         float64 `db:"tax"`
	Total       float64 `db:"total"`
}


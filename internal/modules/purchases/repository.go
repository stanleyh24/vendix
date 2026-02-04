package purchases

import (
	"context"
	"fmt"
	"time"

	"github.com/stanleyh24/vendix/internal/database"

	"github.com/google/uuid"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new purchase order
func (r *Repository) Create(ctx context.Context, schema string, purchase *Purchase) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.purchases (
			id, purchase_number, supplier_id, purchase_date, expected_delivery_date,
			payment_method, payment_status, payment_terms, payment_due_date,
			subtotal, tax_amount, total, currency,
			reference, notes, status, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`, schema)
	
	_, err := r.db.ExecContext(ctx, query,
		purchase.ID, purchase.PurchaseNumber, purchase.SupplierID, purchase.PurchaseDate,
		purchase.ExpectedDeliveryDate, purchase.PaymentMethod, purchase.PaymentStatus,
		purchase.PaymentTerms, purchase.PaymentDueDate,
		purchase.Subtotal, purchase.TaxAmount, purchase.Total, purchase.Currency,
		purchase.Reference, purchase.Notes, purchase.Status, purchase.CreatedBy,
		purchase.CreatedAt, purchase.UpdatedAt,
	)
	return err
}

// CreateLine creates a purchase line
func (r *Repository) CreateLine(ctx context.Context, schema string, line *PurchaseLine) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.purchase_lines (
			id, purchase_id, product_id, description, quantity, unit_price,
			tax_rate, subtotal, tax_amount, total, received_quantity, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, schema)
	
	_, err := r.db.ExecContext(ctx, query,
		line.ID, line.PurchaseID, line.ProductID, line.Description, line.Quantity,
		line.UnitPrice, line.TaxRate, line.Subtotal, line.TaxAmount, line.Total,
		line.ReceivedQuantity, line.CreatedAt,
	)
	return err
}

// GetByID gets a purchase by ID
func (r *Repository) GetByID(ctx context.Context, schema string, id uuid.UUID) (*Purchase, error) {
	var purchase Purchase
	query := fmt.Sprintf(`
		SELECT id, purchase_number, supplier_id, purchase_date, expected_delivery_date,
			payment_method, payment_status, payment_terms, payment_due_date,
			subtotal, tax_amount, total, currency,
			reference, notes, status, received_at, created_by, created_at, updated_at
		FROM %s.purchases WHERE id = $1
	`, schema)
	
	err := r.db.GetContext(ctx, &purchase, query, id)
	return &purchase, err
}

// GetLinesByPurchaseID gets all lines for a purchase
func (r *Repository) GetLinesByPurchaseID(ctx context.Context, schema string, purchaseID uuid.UUID) ([]*PurchaseLine, error) {
	var lines []*PurchaseLine
	query := fmt.Sprintf(`
		SELECT id, purchase_id, product_id, description, quantity, unit_price,
			tax_rate, subtotal, tax_amount, total, received_quantity, created_at
		FROM %s.purchase_lines WHERE purchase_id = $1 ORDER BY created_at ASC
	`, schema)
	
	err := r.db.SelectContext(ctx, &lines, query, purchaseID)
	return lines, err
}

// List gets all purchases with optional filters
func (r *Repository) List(ctx context.Context, schema string, supplierID *uuid.UUID, status *string, limit, offset int) ([]*Purchase, error) {
	var purchases []*Purchase
	query := fmt.Sprintf(`
		SELECT id, purchase_number, supplier_id, purchase_date, expected_delivery_date,
			payment_method, payment_status, payment_terms, payment_due_date,
			subtotal, tax_amount, total, currency,
			reference, notes, status, received_at, created_by, created_at, updated_at
		FROM %s.purchases WHERE 1=1
	`, schema)
	
	args := []interface{}{}
	argPos := 1
	
	if supplierID != nil {
		query += fmt.Sprintf(" AND supplier_id = $%d", argPos)
		args = append(args, *supplierID)
		argPos++
	}
	
	if status != nil {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, *status)
		argPos++
	}
	
	query += " ORDER BY purchase_date DESC, created_at DESC"
	
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, limit)
		argPos++
		if offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argPos)
			args = append(args, offset)
		}
	}
	
	err := r.db.SelectContext(ctx, &purchases, query, args...)
	return purchases, err
}

// Update updates a purchase
func (r *Repository) Update(ctx context.Context, schema string, purchase *Purchase) error {
	query := fmt.Sprintf(`
		UPDATE %s.purchases SET
			purchase_date = $1, expected_delivery_date = $2, payment_method = $3,
			payment_status = $4, payment_terms = $5, payment_due_date = $6,
			subtotal = $7, tax_amount = $8, total = $9,
			reference = $10, notes = $11, status = $12, received_at = $13,
			updated_at = $14
		WHERE id = $15
	`, schema)
	
	_, err := r.db.ExecContext(ctx, query,
		purchase.PurchaseDate, purchase.ExpectedDeliveryDate, purchase.PaymentMethod,
		purchase.PaymentStatus, purchase.PaymentTerms, purchase.PaymentDueDate,
		purchase.Subtotal, purchase.TaxAmount, purchase.Total,
		purchase.Reference, purchase.Notes, purchase.Status, purchase.ReceivedAt,
		purchase.UpdatedAt, purchase.ID,
	)
	return err
}

// DeleteLinesByPurchaseID deletes all lines for a purchase
func (r *Repository) DeleteLinesByPurchaseID(ctx context.Context, schema string, purchaseID uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s.purchase_lines WHERE purchase_id = $1`, schema)
	_, err := r.db.ExecContext(ctx, query, purchaseID)
	return err
}

// UpdateLineReceivedQuantity updates the received quantity for a line
func (r *Repository) UpdateLineReceivedQuantity(ctx context.Context, schema string, lineID uuid.UUID, receivedQuantity float64) error {
	query := fmt.Sprintf(`
		UPDATE %s.purchase_lines SET received_quantity = $1 WHERE id = $2
	`, schema)
	_, err := r.db.ExecContext(ctx, query, receivedQuantity, lineID)
	return err
}

// GetNextPurchaseNumber gets the next purchase number
func (r *Repository) GetNextPurchaseNumber(ctx context.Context, schema string) (string, error) {
	query := fmt.Sprintf(`
		SELECT COALESCE(MAX(CAST(SUBSTRING(purchase_number FROM '[0-9]+') AS INTEGER)), 0) + 1
		FROM %s.purchases
		WHERE purchase_number ~ '^PUR-[0-9]+$'
	`, schema)
	
	var nextNum int
	err := r.db.GetContext(ctx, &nextNum, query)
	if err != nil {
		return "", err
	}
	
	return fmt.Sprintf("PUR-%06d", nextNum), nil
}

// UpdateProductStock updates product stock quantity
func (r *Repository) UpdateProductStock(ctx context.Context, schema string, productID uuid.UUID, quantity float64) error {
	query := fmt.Sprintf(`
		UPDATE %s.products 
		SET stock_quantity = stock_quantity + $1, updated_at = $2
		WHERE id = $3
	`, schema)
	_, err := r.db.ExecContext(ctx, query, quantity, time.Now(), productID)
	return err
}

// GetSupplierName gets supplier name by ID
func (r *Repository) GetSupplierName(ctx context.Context, schema string, supplierID uuid.UUID) (string, error) {
	query := fmt.Sprintf(`SELECT name FROM %s.suppliers WHERE id = $1`, schema)
	var name string
	err := r.db.GetContext(ctx, &name, query, supplierID)
	return name, err
}

// UpdatePurchaseStatusAndReceivedAt updates purchase status and received_at
func (r *Repository) UpdatePurchaseStatusAndReceivedAt(ctx context.Context, schema string, purchaseID uuid.UUID, status string, receivedAt *time.Time) error {
	query := fmt.Sprintf(`
		UPDATE %s.purchases SET status = $1, received_at = $2, updated_at = $3 WHERE id = $4
	`, schema)
	_, err := r.db.ExecContext(ctx, query, status, receivedAt, time.Now(), purchaseID)
	return err
}

// ========== Supplier Payments Repository Methods ==========

// CreateSupplierPayment creates a new supplier payment
func (r *Repository) CreateSupplierPayment(ctx context.Context, schema string, payment *SupplierPayment) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.supplier_payments (
			id, payment_number, supplier_id, payment_method, payment_date,
			amount, currency, reference, notes, status, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, schema)
	
	_, err := r.db.ExecContext(ctx, query,
		payment.ID, payment.PaymentNumber, payment.SupplierID, payment.PaymentMethod,
		payment.PaymentDate, payment.Amount, payment.Currency, payment.Reference,
		payment.Notes, payment.Status, payment.CreatedBy, payment.CreatedAt, payment.UpdatedAt,
	)
	return err
}

// CreateSupplierPaymentAllocation creates a payment allocation to a purchase
func (r *Repository) CreateSupplierPaymentAllocation(ctx context.Context, schema string, allocation *SupplierPaymentAllocation) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.supplier_payment_allocations (
			id, supplier_payment_id, purchase_id, amount, created_at
		) VALUES ($1, $2, $3, $4, $5)
	`, schema)
	
	_, err := r.db.ExecContext(ctx, query,
		allocation.ID, allocation.SupplierPaymentID, allocation.PurchaseID,
		allocation.Amount, allocation.CreatedAt,
	)
	return err
}

// GetSupplierPaymentByID gets a supplier payment by ID
func (r *Repository) GetSupplierPaymentByID(ctx context.Context, schema string, id uuid.UUID) (*SupplierPayment, error) {
	var payment SupplierPayment
	query := fmt.Sprintf(`
		SELECT id, payment_number, supplier_id, payment_method, payment_date,
			amount, currency, reference, notes, status, created_by, created_at, updated_at
		FROM %s.supplier_payments WHERE id = $1
	`, schema)
	
	err := r.db.GetContext(ctx, &payment, query, id)
	return &payment, err
}

// GetSupplierPaymentAllocations gets all allocations for a supplier payment
func (r *Repository) GetSupplierPaymentAllocations(ctx context.Context, schema string, paymentID uuid.UUID) ([]*SupplierPaymentAllocation, error) {
	var allocations []*SupplierPaymentAllocation
	query := fmt.Sprintf(`
		SELECT id, supplier_payment_id, purchase_id, amount, created_at
		FROM %s.supplier_payment_allocations WHERE supplier_payment_id = $1
		ORDER BY created_at ASC
	`, schema)
	
	err := r.db.SelectContext(ctx, &allocations, query, paymentID)
	return allocations, err
}

// GetPurchaseAllocations gets all payment allocations for a purchase
func (r *Repository) GetPurchaseAllocations(ctx context.Context, schema string, purchaseID uuid.UUID) ([]*SupplierPaymentAllocation, error) {
	var allocations []*SupplierPaymentAllocation
	query := fmt.Sprintf(`
		SELECT id, supplier_payment_id, purchase_id, amount, created_at
		FROM %s.supplier_payment_allocations WHERE purchase_id = $1
		ORDER BY created_at ASC
	`, schema)
	
	err := r.db.SelectContext(ctx, &allocations, query, purchaseID)
	return allocations, err
}

// GetTotalAllocatedToPurchase gets the total amount allocated to a purchase
func (r *Repository) GetTotalAllocatedToPurchase(ctx context.Context, schema string, purchaseID uuid.UUID) (float64, error) {
	query := fmt.Sprintf(`
		SELECT COALESCE(SUM(amount), 0)
		FROM %s.supplier_payment_allocations WHERE purchase_id = $1
	`, schema)
	
	var total float64
	err := r.db.GetContext(ctx, &total, query, purchaseID)
	return total, err
}

// UpdatePurchasePaymentStatus updates the payment status of a purchase
func (r *Repository) UpdatePurchasePaymentStatus(ctx context.Context, schema string, purchaseID uuid.UUID, status string) error {
	query := fmt.Sprintf(`
		UPDATE %s.purchases SET payment_status = $1, updated_at = $2 WHERE id = $3
	`, schema)
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), purchaseID)
	return err
}

// GetNextSupplierPaymentNumber gets the next supplier payment number
func (r *Repository) GetNextSupplierPaymentNumber(ctx context.Context, schema string) (string, error) {
	query := fmt.Sprintf(`
		SELECT COALESCE(MAX(CAST(SUBSTRING(payment_number FROM '[0-9]+') AS INTEGER)), 0) + 1
		FROM %s.supplier_payments
		WHERE payment_number ~ '^SPAY-[0-9]+$'
	`, schema)
	
	var nextNum int
	err := r.db.GetContext(ctx, &nextNum, query)
	if err != nil {
		return "", err
	}
	
	return fmt.Sprintf("SPAY-%06d", nextNum), nil
}

// ListSupplierPayments gets all supplier payments with optional filters
func (r *Repository) ListSupplierPayments(ctx context.Context, schema string, supplierID *uuid.UUID, limit, offset int) ([]*SupplierPayment, error) {
	var payments []*SupplierPayment
	query := fmt.Sprintf(`
		SELECT id, payment_number, supplier_id, payment_method, payment_date,
			amount, currency, reference, notes, status, created_by, created_at, updated_at
		FROM %s.supplier_payments WHERE 1=1
	`, schema)
	
	args := []interface{}{}
	argPos := 1
	
	if supplierID != nil {
		query += fmt.Sprintf(" AND supplier_id = $%d", argPos)
		args = append(args, *supplierID)
		argPos++
	}
	
	query += " ORDER BY payment_date DESC, created_at DESC"
	
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, limit)
		argPos++
		if offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argPos)
			args = append(args, offset)
		}
	}
	
	err := r.db.SelectContext(ctx, &payments, query, args...)
	return payments, err
}

// GetAccountsPayableReport generates an accounts payable aging report
func (r *Repository) GetAccountsPayableReport(ctx context.Context, schema string, asOfDate time.Time) (*AccountsPayableReport, error) {
	report := &AccountsPayableReport{
		ReportDate:        asOfDate.Format("2006-01-02"),
		Items:             []AccountsPayableItem{},
		SummaryBySupplier: make(map[string]SupplierAPSummary),
	}

	query := fmt.Sprintf(`
		WITH purchase_payments AS (
			SELECT 
				purchase_id,
				COALESCE(SUM(amount), 0) as paid_amount
			FROM %s.supplier_payment_allocations
			GROUP BY purchase_id
		)
		SELECT 
			p.id::text as purchase_id,
			p.purchase_number,
			p.supplier_id::text as supplier_id,
			s.name as supplier_name,
			p.purchase_date,
			p.expected_delivery_date,
			p.total,
			COALESCE(pp.paid_amount, 0) as paid_amount,
			(p.total - COALESCE(pp.paid_amount, 0)) as balance,
			p.payment_status,
			p.status,
			CASE 
				WHEN p.purchase_date < $1 AND p.payment_status IN ('pending', 'partial') THEN 
					EXTRACT(DAY FROM ($1 - p.purchase_date))::int
				ELSE NULL
			END as days_overdue,
			CASE 
				WHEN p.purchase_date >= $1 - INTERVAL '30 days' AND p.purchase_date <= $1 THEN 
					(p.total - COALESCE(pp.paid_amount, 0))
				WHEN p.purchase_date > $1 THEN
					(p.total - COALESCE(pp.paid_amount, 0))
				ELSE 0
			END as age_0_30,
			CASE 
				WHEN p.purchase_date >= $1 - INTERVAL '60 days' AND p.purchase_date < $1 - INTERVAL '30 days' THEN 
					(p.total - COALESCE(pp.paid_amount, 0))
				ELSE 0
			END as age_31_60,
			CASE 
				WHEN p.purchase_date >= $1 - INTERVAL '90 days' AND p.purchase_date < $1 - INTERVAL '60 days' THEN 
					(p.total - COALESCE(pp.paid_amount, 0))
				ELSE 0
			END as age_61_90,
			CASE 
				WHEN p.purchase_date < $1 - INTERVAL '90 days' THEN 
					(p.total - COALESCE(pp.paid_amount, 0))
				ELSE 0
			END as age_over_90
		FROM %s.purchases p
		INNER JOIN %s.suppliers s ON p.supplier_id = s.id
		LEFT JOIN purchase_payments pp ON p.id = pp.purchase_id
		WHERE p.payment_status IN ('pending', 'partial')
		  AND (p.total - COALESCE(pp.paid_amount, 0)) > 0
		ORDER BY p.purchase_date ASC, p.purchase_number ASC
	`, schema, schema, schema)

	var items []AccountsPayableItem
	err := r.db.SelectContext(ctx, &items, query, asOfDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get AP report: %w", err)
	}

	report.Items = items
	report.TotalPurchases = len(items)

	// Calculate totals and group by supplier
	for _, item := range items {
		report.TotalOutstanding += item.Balance
		report.Total0_30 += item.Age0_30
		report.Total31_60 += item.Age31_60
		report.Total61_90 += item.Age61_90
		report.TotalOver90 += item.AgeOver90

		// Group by supplier
		supplierSummary := report.SummaryBySupplier[item.SupplierID]
		supplierSummary.SupplierID = item.SupplierID
		supplierSummary.SupplierName = item.SupplierName
		supplierSummary.PurchaseCount++
		supplierSummary.TotalOutstanding += item.Balance
		supplierSummary.Total0_30 += item.Age0_30
		supplierSummary.Total31_60 += item.Age31_60
		supplierSummary.Total61_90 += item.Age61_90
		supplierSummary.TotalOver90 += item.AgeOver90
		report.SummaryBySupplier[item.SupplierID] = supplierSummary
	}

	return report, nil
}


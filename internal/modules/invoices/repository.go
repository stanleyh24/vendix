package invoices

import (
	"context"
	"database/sql"
	"encoding/json"
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
		INSERT INTO %s.invoices (id, invoice_number, ncf, ncf_type, customer_id, issue_date, due_date, status, subtotal, tax_amount, total, paid_amount, currency, notes, terms, pdf_url, xml_url, created_by, created_at, updated_at, withholding_tax_amount, withholding_tax_type, withholding_rate, withholding_exempt, net_amount)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
	`, schema)

	_, err = tx.ExecContext(ctx, query,
		invoice.ID, invoice.InvoiceNumber, invoice.NCF, invoice.NCFType, invoice.CustomerID,
		invoice.IssueDate, invoice.DueDate, invoice.Status, invoice.Subtotal,
		invoice.TaxAmount, invoice.Total, invoice.PaidAmount, invoice.Currency,
		invoice.Notes, invoice.Terms, invoice.PDFURL, invoice.XMLURL, invoice.CreatedBy, invoice.CreatedAt, invoice.UpdatedAt,
		invoice.WithholdingTaxAmount, invoice.WithholdingTaxType, invoice.WithholdingRate, invoice.WithholdingExempt, invoice.NetAmount,
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

func (r *Repository) GetByID(ctx context.Context, schema string, id uuid.UUID, includes []string) (*Invoice, error) {
	var invoice Invoice
	query := fmt.Sprintf(`
		SELECT i.id, i.invoice_number, i.ncf, i.ncf_type, i.customer_id, i.issue_date, i.due_date, 
		       i.status, i.subtotal, i.tax_amount, i.total, i.paid_amount, i.currency,
		       i.notes, i.terms, i.dgii_status, i.signed_at, i.sent_at, i.pdf_url, i.xml_url,
		       i.created_by, i.created_at, i.updated_at, 
		       i.withholding_tax_amount, i.withholding_tax_type, i.withholding_rate, i.withholding_exempt, i.net_amount,
		       COALESCE(c.name, 'Cliente Genérico') as customer_name
		FROM %s.invoices i
		LEFT JOIN %s.customers c ON i.customer_id = c.id
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

	// Load relations if requested
	if contains(includes, "customer") && invoice.CustomerID != uuid.Nil {
		customer, err := r.loadCustomer(ctx, schema, invoice.CustomerID)
		if err == nil && customer != nil {
			invoice.Customer = customer
		}
	}

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

func (r *Repository) List(ctx context.Context, schema string, status *string, includes []string) ([]*Invoice, error) {
	var invoices []*Invoice
	baseQuery := fmt.Sprintf(`
		SELECT i.id, i.invoice_number, i.ncf, i.ncf_type, i.customer_id, i.issue_date, i.due_date,
		       i.status, i.subtotal, i.tax_amount, i.total, i.paid_amount, i.currency,
		       i.notes, i.terms, i.dgii_status, i.signed_at, i.sent_at, i.pdf_url, i.xml_url,
		       i.created_by, i.created_at, i.updated_at,
		       i.withholding_tax_amount, i.withholding_tax_type, i.withholding_rate, i.withholding_exempt, i.net_amount,
		       COALESCE(c.name, 'Cliente Genérico') as customer_name
		FROM %s.invoices i
		LEFT JOIN %s.customers c ON i.customer_id = c.id
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

	// Load relations if requested (batch loading for efficiency)
	if contains(includes, "customer") {
		customerMap, err := r.loadCustomersBatch(ctx, schema, invoices)
		if err == nil {
			for i := range invoices {
				if customer, ok := customerMap[invoices[i].CustomerID]; ok {
					invoices[i].Customer = customer
				}
			}
		}
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

// UpdateDocumentURLs updates the PDF and XML URLs for an invoice
func (r *Repository) UpdateDocumentURLs(ctx context.Context, schema string, invoiceID uuid.UUID, pdfURL, xmlURL *string) error {
	query := fmt.Sprintf(`
		UPDATE %s.invoices
		SET pdf_url = $1, xml_url = $2, updated_at = $3
		WHERE id = $4
	`, schema)

	_, err := r.db.ExecContext(ctx, query, pdfURL, xmlURL, time.Now(), invoiceID)
	return err
}

// UpdateDGIIStatus updates the DGII status and related information for an invoice
func (r *Repository) UpdateDGIIStatus(ctx context.Context, schema string, invoiceID uuid.UUID, status, trackingCode string, dgiiResponse interface{}) error {
	// Serialize dgiiResponse to JSON
	var dgiiResponseJSON []byte
	var err error
	if dgiiResponse != nil {
		dgiiResponseJSON, err = json.Marshal(dgiiResponse)
		if err != nil {
			return fmt.Errorf("failed to marshal dgii response: %w", err)
		}
	}

	query := fmt.Sprintf(`
		UPDATE %s.invoices
		SET dgii_status = $1, dgii_response = $2, sent_at = $3, updated_at = $4
		WHERE id = $5
	`, schema)

	sentAt := time.Now()
	_, err = r.db.ExecContext(ctx, query, status, dgiiResponseJSON, sentAt, time.Now(), invoiceID)
	return err
}

// UpdateSignedAt updates the signed_at timestamp for an invoice
func (r *Repository) UpdateSignedAt(ctx context.Context, schema string, invoiceID uuid.UUID, signedAt time.Time) error {
	query := fmt.Sprintf(`
		UPDATE %s.invoices
		SET signed_at = $1, updated_at = $2
		WHERE id = $3
	`, schema)

	_, err := r.db.ExecContext(ctx, query, signedAt, time.Now(), invoiceID)
	return err
}

// loadCustomer loads a single customer by ID
func (r *Repository) loadCustomer(ctx context.Context, schema string, customerID uuid.UUID) (map[string]interface{}, error) {
	type Customer struct {
		ID           uuid.UUID `db:"id"`
		CustomerType string    `db:"customer_type"`
		TaxID        *string   `db:"tax_id"`
		Name         string    `db:"name"`
		Email        *string   `db:"email"`
		Phone        *string   `db:"phone"`
		Address      *string   `db:"address"`
		City         *string   `db:"city"`
		State        *string   `db:"state"`
		PostalCode   *string   `db:"postal_code"`
		Country      *string   `db:"country"`
		IsActive     bool      `db:"is_active"`
		CreatedAt    time.Time `db:"created_at"`
		UpdatedAt    time.Time `db:"updated_at"`
	}

	var customer Customer
	query := fmt.Sprintf(`
		SELECT id, customer_type, tax_id, name, email, phone, address, city, state, postal_code, country, is_active, created_at, updated_at
		FROM %s.customers
		WHERE id = $1
	`, schema)

	err := r.db.GetContext(ctx, &customer, query, customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return map[string]interface{}{
		"id":            customer.ID.String(),
		"customer_type": customer.CustomerType,
		"tax_id":        customer.TaxID,
		"name":          customer.Name,
		"email":         customer.Email,
		"phone":         customer.Phone,
		"address":       customer.Address,
		"city":          customer.City,
		"state":         customer.State,
		"postal_code":   customer.PostalCode,
		"country":       customer.Country,
		"is_active":     customer.IsActive,
		"created_at":    customer.CreatedAt,
		"updated_at":    customer.UpdatedAt,
	}, nil
}

// loadCustomersBatch loads multiple customers by their IDs (batch loading)
func (r *Repository) loadCustomersBatch(ctx context.Context, schema string, invoices []*Invoice) (map[uuid.UUID]map[string]interface{}, error) {
	if len(invoices) == 0 {
		return make(map[uuid.UUID]map[string]interface{}), nil
	}

	// Collect unique customer IDs
	customerIDs := make(map[uuid.UUID]bool)
	for _, invoice := range invoices {
		if invoice.CustomerID != uuid.Nil {
			customerIDs[invoice.CustomerID] = true
		}
	}

	if len(customerIDs) == 0 {
		return make(map[uuid.UUID]map[string]interface{}), nil
	}

	// Build query with IN clause
	ids := make([]uuid.UUID, 0, len(customerIDs))
	for id := range customerIDs {
		ids = append(ids, id)
	}

	type Customer struct {
		ID           uuid.UUID `db:"id"`
		CustomerType string    `db:"customer_type"`
		TaxID        *string   `db:"tax_id"`
		Name         string    `db:"name"`
		Email        *string   `db:"email"`
		Phone        *string   `db:"phone"`
		Address      *string   `db:"address"`
		City         *string   `db:"city"`
		State        *string   `db:"state"`
		PostalCode   *string   `db:"postal_code"`
		Country      *string   `db:"country"`
		IsActive     bool      `db:"is_active"`
		CreatedAt    time.Time `db:"created_at"`
		UpdatedAt    time.Time `db:"updated_at"`
	}

	var customers []Customer
	query := fmt.Sprintf(`
		SELECT id, customer_type, tax_id, name, email, phone, address, city, state, postal_code, country, is_active, created_at, updated_at
		FROM %s.customers
		WHERE id = ANY($1)
	`, schema)

	err := r.db.SelectContext(ctx, &customers, query, ids)
	if err != nil {
		return nil, err
	}

	// Build map
	customerMap := make(map[uuid.UUID]map[string]interface{})
	for _, customer := range customers {
		customerMap[customer.ID] = map[string]interface{}{
			"id":            customer.ID.String(),
			"customer_type": customer.CustomerType,
			"tax_id":        customer.TaxID,
			"name":          customer.Name,
			"email":         customer.Email,
			"phone":         customer.Phone,
			"address":       customer.Address,
			"city":          customer.City,
			"state":         customer.State,
			"postal_code":   customer.PostalCode,
			"country":       customer.Country,
			"is_active":     customer.IsActive,
			"created_at":    customer.CreatedAt,
			"updated_at":    customer.UpdatedAt,
		}
	}

	return customerMap, nil
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

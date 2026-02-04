package creditnotes

import (
	"context"
	"database/sql"
	"encoding/json"
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

// Create creates a new credit note
func (r *Repository) Create(ctx context.Context, schema string, creditNote *CreditNote) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert credit note
	query := fmt.Sprintf(`
		INSERT INTO %s.credit_notes (id, credit_note_number, ncf, ncf_type, original_invoice_id, return_id, customer_id, 
		                              issue_date, reason, status, subtotal, tax_amount, total, currency, notes,
		                              dgii_status, tracking_code, signed_at, sent_at, pdf_url, xml_url, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
	`, schema)

	_, err = tx.ExecContext(ctx, query,
		creditNote.ID, creditNote.CreditNoteNumber, creditNote.NCF, creditNote.NCFType,
		creditNote.OriginalInvoiceID, creditNote.ReturnID, creditNote.CustomerID,
		creditNote.IssueDate, creditNote.Reason, creditNote.Status,
		creditNote.Subtotal, creditNote.TaxAmount, creditNote.Total, creditNote.Currency,
		creditNote.Notes, creditNote.DGIIStatus, creditNote.TrackingCode,
		creditNote.SignedAt, creditNote.SentAt, creditNote.PDFURL, creditNote.XMLURL,
		creditNote.CreatedBy, creditNote.CreatedAt, creditNote.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Insert lines
	lineQuery := fmt.Sprintf(`
		INSERT INTO %s.credit_note_lines (id, credit_note_id, original_invoice_line_id, product_id, line_number, description,
		                                   quantity, unit_price, tax_rate, tax_amount, line_total, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, schema)

	for _, line := range creditNote.Lines {
		_, err = tx.ExecContext(ctx, lineQuery,
			line.ID, line.CreditNoteID, line.OriginalInvoiceLineID, line.ProductID, line.LineNumber,
			line.Description, line.Quantity, line.UnitPrice, line.TaxRate,
			line.TaxAmount, line.LineTotal, line.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetByID retrieves a credit note by ID
func (r *Repository) GetByID(ctx context.Context, schema string, id uuid.UUID) (*CreditNote, error) {
	var creditNote CreditNote
	query := fmt.Sprintf(`
		SELECT cn.id, cn.credit_note_number, cn.ncf, cn.ncf_type, cn.original_invoice_id, cn.return_id, cn.customer_id,
		       cn.issue_date, cn.reason, cn.status, cn.subtotal, cn.tax_amount, cn.total, cn.currency, cn.notes,
		       cn.dgii_status, cn.tracking_code, cn.signed_at, cn.sent_at, cn.pdf_url, cn.xml_url,
		       cn.created_by, cn.created_at, cn.updated_at,
		       COALESCE(c.name, 'Cliente Genérico') as customer_name,
		       i.invoice_number as original_invoice_number
		FROM %s.credit_notes cn
		LEFT JOIN %s.customers c ON cn.customer_id = c.id
		LEFT JOIN %s.invoices i ON cn.original_invoice_id = i.id
		WHERE cn.id = $1
	`, schema, schema, schema)

	err := r.db.GetContext(ctx, &creditNote, query, id)
	if err != nil {
		return nil, err
	}

	// Get lines
	lines, err := r.getCreditNoteLines(ctx, schema, id)
	if err != nil {
		return nil, err
	}

	creditNote.Lines = lines
	return &creditNote, nil
}

// getCreditNoteLines retrieves all lines for a credit note
func (r *Repository) getCreditNoteLines(ctx context.Context, schema string, creditNoteID uuid.UUID) ([]CreditNoteLine, error) {
	var lines []CreditNoteLine
	query := fmt.Sprintf(`
		SELECT id, credit_note_id, original_invoice_line_id, product_id, line_number, description,
		       quantity, unit_price, tax_rate, tax_amount, line_total, created_at
		FROM %s.credit_note_lines
		WHERE credit_note_id = $1
		ORDER BY line_number
	`, schema)

	err := r.db.SelectContext(ctx, &lines, query, creditNoteID)
	return lines, err
}

// List retrieves all credit notes with optional filters
func (r *Repository) List(ctx context.Context, schema string, filters *CreditNoteFilters) ([]*CreditNote, error) {
	query := fmt.Sprintf(`
		SELECT cn.id, cn.credit_note_number, cn.ncf, cn.ncf_type, cn.original_invoice_id, cn.return_id, cn.customer_id,
		       cn.issue_date, cn.reason, cn.status, cn.subtotal, cn.tax_amount, cn.total, cn.currency, cn.notes,
		       cn.dgii_status, cn.tracking_code, cn.signed_at, cn.sent_at, cn.pdf_url, cn.xml_url,
		       cn.created_by, cn.created_at, cn.updated_at,
		       COALESCE(c.name, 'Cliente Genérico') as customer_name,
		       i.invoice_number as original_invoice_number
		FROM %s.credit_notes cn
		LEFT JOIN %s.customers c ON cn.customer_id = c.id
		LEFT JOIN %s.invoices i ON cn.original_invoice_id = i.id
		WHERE 1=1
	`, schema, schema, schema)

	args := []interface{}{}
	argNum := 1

	if filters != nil {
		if filters.Status != nil && *filters.Status != "" {
			query += fmt.Sprintf(" AND cn.status = $%d", argNum)
			args = append(args, *filters.Status)
			argNum++
		}
		if filters.OriginalInvoiceID != nil {
			query += fmt.Sprintf(" AND cn.original_invoice_id = $%d", argNum)
			args = append(args, *filters.OriginalInvoiceID)
			argNum++
		}
		if filters.ReturnID != nil {
			query += fmt.Sprintf(" AND cn.return_id = $%d", argNum)
			args = append(args, *filters.ReturnID)
			argNum++
		}
		if filters.CustomerID != nil {
			query += fmt.Sprintf(" AND cn.customer_id = $%d", argNum)
			args = append(args, *filters.CustomerID)
			argNum++
		}
		if filters.StartDate != nil {
			query += fmt.Sprintf(" AND cn.issue_date >= $%d", argNum)
			args = append(args, *filters.StartDate)
			argNum++
		}
		if filters.EndDate != nil {
			query += fmt.Sprintf(" AND cn.issue_date <= $%d", argNum)
			args = append(args, *filters.EndDate)
			argNum++
		}
	}

	query += " ORDER BY cn.issue_date DESC, cn.created_at DESC"

	var creditNotes []*CreditNote
	err := r.db.SelectContext(ctx, &creditNotes, query, args...)
	if err != nil {
		// SelectContext nunca devuelve ErrNoRows (solo GetContext/QueryRow lo hace)
		// Pero manejamos el caso por si acaso
		if err == sql.ErrNoRows {
			return []*CreditNote{}, nil
		}
		// Para otros errores, los propagamos
		return nil, err
	}

	// SelectContext siempre devuelve una lista (vacía si no hay resultados)
	// Pero aseguramos que nunca sea nil
	if creditNotes == nil {
		creditNotes = []*CreditNote{}
	}

	// Load lines for each credit note
	for _, cn := range creditNotes {
		lines, err := r.getCreditNoteLines(ctx, schema, cn.ID)
		if err == nil {
			cn.Lines = lines
		}
	}

	return creditNotes, nil
}

// CreditNoteFilters represents filters for listing credit notes
type CreditNoteFilters struct {
	Status            *string
	OriginalInvoiceID *uuid.UUID
	ReturnID          *uuid.UUID
	CustomerID        *uuid.UUID
	StartDate         *time.Time
	EndDate           *time.Time
}

// GetNextCreditNoteNumber generates the next credit note number
func (r *Repository) GetNextCreditNoteNumber(ctx context.Context, schema string) (string, error) {
	query := fmt.Sprintf(`
		SELECT COALESCE(MAX(CAST(SUBSTRING(credit_note_number FROM 'NC-(.*)') AS INTEGER)), 0) + 1
		FROM %s.credit_notes
		WHERE credit_note_number ~ '^NC-[0-9]+$'
	`, schema)

	var nextNum int
	err := r.db.GetContext(ctx, &nextNum, query)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	return fmt.Sprintf("NC-%06d", nextNum), nil
}

// UpdateNCF updates the NCF of a credit note
func (r *Repository) UpdateNCF(ctx context.Context, schema string, id uuid.UUID, ncf string) error {
	query := fmt.Sprintf(`
		UPDATE %s.credit_notes 
		SET ncf = $1, updated_at = $2
		WHERE id = $3
	`, schema)

	_, err := r.db.ExecContext(ctx, query, ncf, time.Now(), id)
	return err
}

// UpdateDGIIStatus updates DGII status and tracking code
func (r *Repository) UpdateDGIIStatus(ctx context.Context, schema string, id uuid.UUID, status, trackingCode string, dgiiResponse interface{}) error {
	// Note: dgiiResponse is kept for future use if we want to store full response
	_, _ = json.Marshal(dgiiResponse)

	query := fmt.Sprintf(`
		UPDATE %s.credit_notes 
		SET dgii_status = $1, tracking_code = $2, updated_at = $3
		WHERE id = $4
	`, schema)

	_, err := r.db.ExecContext(ctx, query, status, trackingCode, time.Now(), id)
	return err
}

// UpdateStatus updates the status of a credit note
func (r *Repository) UpdateStatus(ctx context.Context, schema string, id uuid.UUID, status CreditNoteStatus) error {
	query := fmt.Sprintf(`
		UPDATE %s.credit_notes 
		SET status = $1, updated_at = $2
		WHERE id = $3
	`, schema)

	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

// UpdateSignedAt updates the signed_at timestamp
func (r *Repository) UpdateSignedAt(ctx context.Context, schema string, id uuid.UUID, signedAt time.Time) error {
	query := fmt.Sprintf(`
		UPDATE %s.credit_notes 
		SET signed_at = $1, updated_at = $2
		WHERE id = $3
	`, schema)

	_, err := r.db.ExecContext(ctx, query, signedAt, time.Now(), id)
	return err
}

// UpdateSentAt updates the sent_at timestamp
func (r *Repository) UpdateSentAt(ctx context.Context, schema string, id uuid.UUID, sentAt time.Time) error {
	query := fmt.Sprintf(`
		UPDATE %s.credit_notes 
		SET sent_at = $1, updated_at = $2
		WHERE id = $3
	`, schema)

	_, err := r.db.ExecContext(ctx, query, sentAt, time.Now(), id)
	return err
}

// UpdatePDFURL updates the PDF URL
func (r *Repository) UpdatePDFURL(ctx context.Context, schema string, id uuid.UUID, pdfURL string) error {
	query := fmt.Sprintf(`
		UPDATE %s.credit_notes 
		SET pdf_url = $1, updated_at = $2
		WHERE id = $3
	`, schema)

	_, err := r.db.ExecContext(ctx, query, pdfURL, time.Now(), id)
	return err
}

// UpdateXMLURL updates the XML URL
func (r *Repository) UpdateXMLURL(ctx context.Context, schema string, id uuid.UUID, xmlURL string) error {
	query := fmt.Sprintf(`
		UPDATE %s.credit_notes 
		SET xml_url = $1, updated_at = $2
		WHERE id = $3
	`, schema)

	_, err := r.db.ExecContext(ctx, query, xmlURL, time.Now(), id)
	return err
}

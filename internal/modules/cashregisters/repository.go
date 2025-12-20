package cashregisters

import (
	"context"
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

// Cash Register methods

func (r *Repository) CreateCashRegister(ctx context.Context, schema string, register *CashRegister) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.cash_registers (id, name, location, is_active, initial_balance, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		register.ID, register.Name, register.Location, register.IsActive,
		register.InitialBalance, register.CreatedBy, register.CreatedAt, register.UpdatedAt,
	)
	return err
}

func (r *Repository) GetCashRegisterByID(ctx context.Context, schema string, id uuid.UUID) (*CashRegister, error) {
	var register CashRegister
	query := fmt.Sprintf(`
		SELECT id, name, location, is_active, initial_balance, created_by, created_at, updated_at
		FROM %s.cash_registers
		WHERE id = $1
	`, schema)

	err := r.db.GetContext(ctx, &register, query, id)
	return &register, err
}

func (r *Repository) ListCashRegisters(ctx context.Context, schema string, activeOnly bool) ([]*CashRegister, error) {
	var registers []*CashRegister
	query := fmt.Sprintf(`
		SELECT id, name, location, is_active, initial_balance, created_by, created_at, updated_at
		FROM %s.cash_registers
		WHERE 1=1
	`, schema)

	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY name"

	err := r.db.SelectContext(ctx, &registers, query)
	if registers == nil {
		registers = []*CashRegister{}
	}
	return registers, err
}

// Session methods

func (r *Repository) CreateSession(ctx context.Context, schema string, session *CashRegisterSession) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.cash_register_sessions 
		(id, cash_register_id, opened_by, opening_balance, expected_cash, expected_card, 
		 expected_transfer, status, notes, opened_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		session.ID, session.CashRegisterID, session.OpenedBy, session.OpeningBalance,
		session.ExpectedCash, session.ExpectedCard, session.ExpectedTransfer,
		session.Status, session.Notes, session.OpenedAt, session.CreatedAt, session.UpdatedAt,
	)
	return err
}

func (r *Repository) GetOpenSession(ctx context.Context, schema string, cashRegisterID uuid.UUID) (*CashRegisterSession, error) {
	var session CashRegisterSession
	query := fmt.Sprintf(`
		SELECT id, cash_register_id, opened_by, closed_by, opening_balance,
			   expected_cash, counted_cash, expected_card, expected_transfer,
			   difference, status, notes, journal_entry_id, opened_at, closed_at, 
			   created_at, updated_at
		FROM %s.cash_register_sessions
		WHERE cash_register_id = $1 AND status = 'open'
		ORDER BY opened_at DESC
		LIMIT 1
	`, schema)

	err := r.db.GetContext(ctx, &session, query, cashRegisterID)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *Repository) GetSessionByID(ctx context.Context, schema string, id uuid.UUID) (*CashRegisterSession, error) {
	var session CashRegisterSession
	query := fmt.Sprintf(`
		SELECT s.id, s.cash_register_id, s.opened_by, s.closed_by, s.opening_balance,
			   s.expected_cash, s.counted_cash, s.expected_card, s.expected_transfer,
			   s.difference, s.status, s.notes, s.journal_entry_id, s.opened_at, s.closed_at, 
			   s.created_at, s.updated_at,
			   cr.name as cash_register_name,
			   TRIM(CONCAT(COALESCE(u1.first_name, ''), ' ', COALESCE(u1.last_name, ''))) as opened_by_name,
			   TRIM(CONCAT(COALESCE(u2.first_name, ''), ' ', COALESCE(u2.last_name, ''))) as closed_by_name
		FROM %s.cash_register_sessions s
		LEFT JOIN %s.cash_registers cr ON s.cash_register_id = cr.id
		LEFT JOIN %s.users u1 ON s.opened_by = u1.id
		LEFT JOIN %s.users u2 ON s.closed_by = u2.id
		WHERE s.id = $1
	`, schema, schema, schema, schema)

	err := r.db.GetContext(ctx, &session, query, id)
	if err != nil {
		return nil, err
	}

	// Load count details
	details, err := r.GetCountDetails(ctx, schema, id)
	if err == nil {
		session.CountDetails = details
	}

	return &session, nil
}

func (r *Repository) UpdateSession(ctx context.Context, schema string, session *CashRegisterSession) error {
	query := fmt.Sprintf(`
		UPDATE %s.cash_register_sessions
		SET counted_cash = $1, expected_cash = $2, expected_card = $3, expected_transfer = $4,
			difference = $5, status = $6, closed_by = $7, closed_at = $8, notes = $9, 
			journal_entry_id = $10, updated_at = $11
		WHERE id = $12
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		session.CountedCash, session.ExpectedCash, session.ExpectedCard, session.ExpectedTransfer,
		session.Difference, session.Status, session.ClosedBy, session.ClosedAt,
		session.Notes, session.JournalEntryID, session.UpdatedAt, session.ID,
	)
	return err
}

// GetSalesSummaryByDate gets sales summary by payment type for a date range
func (r *Repository) GetSalesSummaryByDate(ctx context.Context, schema string, start, end time.Time) (map[string]float64, error) {
	query := fmt.Sprintf(`
		SELECT 
			payment_type,
			COUNT(*) as count,
			COALESCE(SUM(total), 0) as total
		FROM %s.sales
		WHERE created_at >= $1 AND created_at <= $2
		  AND status = 'completed'
		GROUP BY payment_type
	`, schema)

	rows, err := r.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summary := map[string]float64{
		"cash":     0,
		"card":     0,
		"transfer": 0,
		"check":    0,
		"mixed":    0,
		"count":    0,
	}

	for rows.Next() {
		var paymentType string
		var count int
		var total float64

		if err := rows.Scan(&paymentType, &count, &total); err != nil {
			return nil, err
		}

		summary[paymentType] = total
		summary["count"] += float64(count)
	}

	return summary, nil
}

// Count details methods

func (r *Repository) CreateCountDetail(ctx context.Context, schema string, detail *CashCountDetail) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.cash_count_details (id, session_id, denomination, quantity, subtotal, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		detail.ID, detail.SessionID, detail.Denomination, detail.Quantity,
		detail.Subtotal, detail.CreatedAt,
	)
	return err
}

func (r *Repository) GetCountDetails(ctx context.Context, schema string, sessionID uuid.UUID) ([]CashCountDetail, error) {
	var details []CashCountDetail
	query := fmt.Sprintf(`
		SELECT id, session_id, denomination, quantity, subtotal, created_at
		FROM %s.cash_count_details
		WHERE session_id = $1
		ORDER BY denomination DESC
	`, schema)

	err := r.db.SelectContext(ctx, &details, query, sessionID)
	return details, err
}

// ListSessions lists sessions with filters
func (r *Repository) ListSessions(ctx context.Context, schema string, cashRegisterID *uuid.UUID, status *string, startDate, endDate *time.Time) ([]*CashRegisterSession, error) {
	var sessions []*CashRegisterSession
	baseQuery := fmt.Sprintf(`
		SELECT s.id, s.cash_register_id, s.opened_by, s.closed_by, s.opening_balance,
			   s.expected_cash, s.counted_cash, s.expected_card, s.expected_transfer,
			   s.difference, s.status, s.notes, s.journal_entry_id, s.opened_at, s.closed_at, 
			   s.created_at, s.updated_at,
			   cr.name as cash_register_name,
			   TRIM(CONCAT(COALESCE(u1.first_name, ''), ' ', COALESCE(u1.last_name, ''))) as opened_by_name,
			   TRIM(CONCAT(COALESCE(u2.first_name, ''), ' ', COALESCE(u2.last_name, ''))) as closed_by_name
		FROM %s.cash_register_sessions s
		LEFT JOIN %s.cash_registers cr ON s.cash_register_id = cr.id
		LEFT JOIN %s.users u1 ON s.opened_by = u1.id
		LEFT JOIN %s.users u2 ON s.closed_by = u2.id
		WHERE 1=1
	`, schema, schema, schema, schema)

	args := []interface{}{}
	argCount := 1

	if cashRegisterID != nil {
		baseQuery += fmt.Sprintf(" AND s.cash_register_id = $%d", argCount)
		args = append(args, *cashRegisterID)
		argCount++
	}

	if status != nil {
		baseQuery += fmt.Sprintf(" AND s.status = $%d", argCount)
		args = append(args, *status)
		argCount++
	}

	if startDate != nil {
		baseQuery += fmt.Sprintf(" AND s.opened_at >= $%d", argCount)
		args = append(args, *startDate)
		argCount++
	}

	if endDate != nil {
		baseQuery += fmt.Sprintf(" AND s.opened_at <= $%d", argCount)
		args = append(args, *endDate)
		argCount++
	}

	baseQuery += " ORDER BY s.opened_at DESC"

	err := r.db.SelectContext(ctx, &sessions, baseQuery, args...)
	if sessions == nil {
		sessions = []*CashRegisterSession{}
	}
	return sessions, err
}


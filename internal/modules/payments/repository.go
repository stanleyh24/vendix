package payments

import (
	"context"
	"fmt"

	"vendix/internal/database"

	"github.com/google/uuid"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, schema string, payment *Payment) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.payments (id, payment_number, customer_id, payment_method, payment_date, amount, currency, reference, notes, status, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		payment.ID, payment.PaymentNumber, payment.CustomerID, payment.PaymentMethod,
		payment.PaymentDate, payment.Amount, payment.Currency, payment.Reference,
		payment.Notes, payment.Status, payment.CreatedBy, payment.CreatedAt, payment.UpdatedAt,
	)

	return err
}

func (r *Repository) GetByID(ctx context.Context, schema string, id uuid.UUID) (*Payment, error) {
	var payment Payment
	query := fmt.Sprintf(`
		SELECT id, payment_number, customer_id, payment_method, payment_date, amount, currency, reference, notes, stripe_payment_id, status, created_by, created_at, updated_at
		FROM %s.payments
		WHERE id = $1
	`, schema)

	err := r.db.GetContext(ctx, &payment, query, id)
	return &payment, err
}

func (r *Repository) List(ctx context.Context, schema string) ([]*Payment, error) {
	var payments []*Payment
	query := fmt.Sprintf(`
		SELECT id, payment_number, customer_id, payment_method, payment_date, amount, currency, reference, notes, stripe_payment_id, status, created_by, created_at, updated_at
		FROM %s.payments
		ORDER BY payment_date DESC, payment_number DESC
	`, schema)

	err := r.db.SelectContext(ctx, &payments, query)
	if err != nil {
		return nil, err
	}

	if payments == nil {
		payments = []*Payment{}
	}

	return payments, nil
}

func (r *Repository) Update(ctx context.Context, schema string, payment *Payment) error {
	query := fmt.Sprintf(`
		UPDATE %s.payments
		SET status = $1, notes = $2, reference = $3, updated_at = $4
		WHERE id = $5
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		payment.Status, payment.Notes, payment.Reference, payment.UpdatedAt, payment.ID,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, schema string, id uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s.payments WHERE id = $1`, schema)
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *Repository) GetNextPaymentNumber(ctx context.Context, schema string, prefix string) (string, error) {
	var lastNumber string
	query := fmt.Sprintf(`
		SELECT payment_number 
		FROM %s.payments 
		WHERE payment_number LIKE $1
		ORDER BY created_at DESC 
		LIMIT 1
	`, schema)

	err := r.db.GetContext(ctx, &lastNumber, query, prefix+"%")
	if err != nil {
		// If no payments exist, start with 1
		return fmt.Sprintf("%s%08d", prefix, 1), nil
	}

	// Parse and increment
	var num int
	fmt.Sscanf(lastNumber, prefix+"%d", &num)
	num++

	return fmt.Sprintf("%s%08d", prefix, num), nil
}

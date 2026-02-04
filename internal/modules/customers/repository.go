package customers

import (
	"context"
	"fmt"

	"github.com/stanleyh24/vendix/internal/database"

	"github.com/google/uuid"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, schema string, customer *Customer) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.customers (id, customer_type, tax_id, name, email, phone, address, city, state, postal_code, country, is_active, is_government_entity, default_withholding_rate, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		customer.ID,
		customer.CustomerType,
		customer.TaxID,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.Address,
		customer.City,
		customer.State,
		customer.PostalCode,
		customer.Country,
		customer.IsActive,
		customer.IsGovernmentEntity,
		customer.DefaultWithholdingRate,
		customer.CreatedAt,
		customer.UpdatedAt,
	)

	return err
}

func (r *Repository) GetByID(ctx context.Context, schema string, id uuid.UUID) (*Customer, error) {
	var customer Customer
	query := fmt.Sprintf(`
		SELECT id, customer_type, tax_id, name, email, phone, address, city, state, postal_code, country, is_active, is_government_entity, default_withholding_rate, created_at, updated_at
		FROM %s.customers
		WHERE id = $1
	`, schema)

	err := r.db.GetContext(ctx, &customer, query, id)
	return &customer, err
}

func (r *Repository) List(ctx context.Context, schema string) ([]*Customer, error) {
	var customers []*Customer
	query := fmt.Sprintf(`
		SELECT id, customer_type, tax_id, name, email, phone, address, city, state, postal_code, country, is_active, is_government_entity, default_withholding_rate, created_at, updated_at
		FROM %s.customers
		ORDER BY name ASC
	`, schema)

	err := r.db.SelectContext(ctx, &customers, query)
	return customers, err
}

func (r *Repository) Update(ctx context.Context, schema string, customer *Customer) error {
	query := fmt.Sprintf(`
		UPDATE %s.customers
		SET name = $1, email = $2, phone = $3, address = $4, city = $5, state = $6, postal_code = $7, is_active = $8, is_government_entity = $9, default_withholding_rate = $10, updated_at = $11
		WHERE id = $12
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.Address,
		customer.City,
		customer.State,
		customer.PostalCode,
		customer.IsActive,
		customer.IsGovernmentEntity,
		customer.DefaultWithholdingRate,
		customer.UpdatedAt,
		customer.ID,
	)

	return err
}

func (r *Repository) Delete(ctx context.Context, schema string, id uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.customers
		WHERE id = $1
	`, schema)

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

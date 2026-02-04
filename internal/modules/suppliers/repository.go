package suppliers

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

func (r *Repository) Create(ctx context.Context, schema string, s *Supplier) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.suppliers (id, name, tax_id, email, phone, address, city, state, postal_code, country, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
	`, schema)
	_, err := r.db.ExecContext(ctx, query, s.ID, s.Name, s.TaxID, s.Email, s.Phone, s.Address, s.City, s.State, s.PostalCode, s.Country, s.IsActive)
	return err
}

func (r *Repository) Update(ctx context.Context, schema string, s *Supplier) error {
	query := fmt.Sprintf(`
		UPDATE %s.suppliers
		SET name = $1, tax_id = $2, email = $3, phone = $4, address = $5, city = $6, state = $7, postal_code = $8, country = $9, is_active = $10, updated_at = NOW()
		WHERE id = $11
	`, schema)
	_, err := r.db.ExecContext(ctx, query, s.Name, s.TaxID, s.Email, s.Phone, s.Address, s.City, s.State, s.PostalCode, s.Country, s.IsActive, s.ID)
	return err
}

func (r *Repository) Delete(ctx context.Context, schema string, id uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s.suppliers WHERE id = $1`, schema)
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *Repository) GetByID(ctx context.Context, schema string, id uuid.UUID) (*Supplier, error) {
	var s Supplier
	query := fmt.Sprintf(`
		SELECT id, name, tax_id, email, phone, address, city, state, postal_code, country, is_active, created_at, updated_at
		FROM %s.suppliers WHERE id = $1
	`, schema)
	if err := r.db.GetContext(ctx, &s, query, id); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) List(ctx context.Context, schema string, onlyActive *bool) ([]*Supplier, error) {
	var list []*Supplier
	base := fmt.Sprintf(`
		SELECT id, name, tax_id, email, phone, address, city, state, postal_code, country, is_active, created_at, updated_at
		FROM %s.suppliers WHERE 1=1
	`, schema)
	args := []interface{}{}
	if onlyActive != nil {
		base += " AND is_active = $1"
		args = append(args, *onlyActive)
	}
	base += " ORDER BY name ASC"
	if err := r.db.SelectContext(ctx, &list, base, args...); err != nil {
		return nil, err
	}
	return list, nil
}

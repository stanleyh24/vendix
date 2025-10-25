package tenants

import (
	"context"
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

func (r *Repository) Create(ctx context.Context, tenant *Tenant) error {
	query := `
		INSERT INTO public.tenants (id, name, slug, schema_name, domain, status, plan_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		tenant.ID,
		tenant.Name,
		tenant.Slug,
		tenant.SchemaName,
		tenant.Domain,
		tenant.Status,
		tenant.PlanID,
		tenant.CreatedAt,
		tenant.UpdatedAt,
	)

	return err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	var tenant Tenant
	query := `
		SELECT id, name, slug, schema_name, domain, status, plan_id, created_at, updated_at
		FROM public.tenants
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &tenant, query, id)
	return &tenant, err
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (*Tenant, error) {
	var tenant Tenant
	query := `
		SELECT id, name, slug, schema_name, domain, status, plan_id, created_at, updated_at
		FROM public.tenants
		WHERE slug = $1 AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &tenant, query, slug)
	return &tenant, err
}

func (r *Repository) List(ctx context.Context) ([]*Tenant, error) {
	var tenants []*Tenant
	query := `
		SELECT id, name, slug, schema_name, domain, status, plan_id, created_at, updated_at
		FROM public.tenants
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &tenants, query)
	return tenants, err
}

func (r *Repository) Update(ctx context.Context, tenant *Tenant) error {
	query := `
		UPDATE public.tenants
		SET name = $1, domain = $2, status = $3, plan_id = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := r.db.ExecContext(ctx, query,
		tenant.Name,
		tenant.Domain,
		tenant.Status,
		tenant.PlanID,
		tenant.UpdatedAt,
		tenant.ID,
	)

	return err
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE public.tenants
		SET deleted_at = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

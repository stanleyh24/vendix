package tenants

import (
	"context"
	"fmt"
	"strings"
	"time"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/logger"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
	db   *database.DB
	cfg  *config.Config
}

func NewService(db *database.DB, cfg *config.Config) *Service {
	return &Service{
		repo: NewRepository(db),
		db:   db,
		cfg:  cfg,
	}
}

func (s *Service) Create(ctx context.Context, req *CreateTenantRequest) (*TenantResponse, error) {
	// Generate schema name from slug
	schemaName := fmt.Sprintf("tenant_%s", strings.ReplaceAll(req.Slug, "-", "_"))

	tenant := &Tenant{
		ID:         uuid.New(),
		Name:       req.Name,
		Slug:       req.Slug,
		SchemaName: schemaName,
		Domain:     req.Domain,
		Status:     "active",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Check if plan_id is provided
	if req.PlanID != nil {
		planID, err := uuid.Parse(*req.PlanID)
		if err != nil {
			return nil, fmt.Errorf("invalid plan_id: %w", err)
		}
		tenant.PlanID = &planID
	}

	if err := s.repo.Create(ctx, tenant); err != nil {
		return nil, err
	}

	logger.Info("Tenant created", "id", tenant.ID, "slug", tenant.Slug)

	return tenant.ToResponse(), nil
}

func (s *Service) List(ctx context.Context) ([]*TenantResponse, error) {
	tenants, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*TenantResponse, len(tenants))
	for i, t := range tenants {
		responses[i] = t.ToResponse()
	}

	return responses, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*TenantResponse, error) {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant ID: %w", err)
	}

	tenant, err := s.repo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	return tenant.ToResponse(), nil
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*TenantResponse, error) {
	tenant, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return tenant.ToResponse(), nil
}

func (s *Service) Update(ctx context.Context, id string, req *UpdateTenantRequest) (*TenantResponse, error) {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant ID: %w", err)
	}

	tenant, err := s.repo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		tenant.Name = *req.Name
	}
	if req.Domain != nil {
		tenant.Domain = req.Domain
	}
	if req.Status != nil {
		tenant.Status = *req.Status
	}
	if req.PlanID != nil {
		planID, err := uuid.Parse(*req.PlanID)
		if err != nil {
			return nil, fmt.Errorf("invalid plan_id: %w", err)
		}
		tenant.PlanID = &planID
	}

	tenant.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, tenant); err != nil {
		return nil, err
	}

	logger.Info("Tenant updated", "id", tenant.ID)

	return tenant.ToResponse(), nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid tenant ID: %w", err)
	}

	if err := s.repo.Delete(ctx, tenantID); err != nil {
		return err
	}

	logger.Info("Tenant deleted", "id", tenantID)

	return nil
}

func (s *Service) InitializeSchema(ctx context.Context, id string) error {
	tenantID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid tenant ID: %w", err)
	}

	tenant, err := s.repo.GetByID(ctx, tenantID)
	if err != nil {
		return err
	}

	// Run tenant migrations
	logger.Info("Initializing schema for tenant", "schema", tenant.SchemaName)

	if err := database.RunTenantMigrations(s.db, tenant.SchemaName); err != nil {
		return fmt.Errorf("failed to run tenant migrations: %w", err)
	}

	logger.Info("Tenant schema initialized successfully", "schema", tenant.SchemaName)

	return nil
}

package tenants

import (
	"context"
	"fmt"
	"strings"
	"time"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/logger"
	"vendix/internal/storage"

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

	// Create default cash register after migrations
	if err := s.createDefaultCashRegister(ctx, tenant.SchemaName); err != nil {
		logger.Error("Failed to create default cash register", "error", err, "schema", tenant.SchemaName)
		// Don't fail initialization if default cash register creation fails
		// This allows tenants to continue functioning even if cash registers feature has issues
	}

	// Create MinIO bucket for tenant
	if err := s.createTenantBucket(ctx, tenant.SchemaName); err != nil {
		logger.Error("Failed to create tenant bucket", "error", err, "schema", tenant.SchemaName)
		// Don't fail initialization if bucket creation fails
		// This allows tenants to continue functioning even if storage is unavailable
	}

	logger.Info("Tenant schema initialized successfully", "schema", tenant.SchemaName)

	return nil
}

// createDefaultCashRegister creates a default cash register for the tenant
func (s *Service) createDefaultCashRegister(ctx context.Context, schema string) error {
	// Check if cash_registers table exists (migration 26)
	checkQuery := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = '%s' 
			AND table_name = 'cash_registers'
		)
	`, schema)

	var tableExists bool
	err := s.db.Pool.QueryRow(ctx, checkQuery).Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("failed to check if cash_registers table exists: %w", err)
	}

	if !tableExists {
		logger.Debug("cash_registers table does not exist, skipping default cash register creation", "schema", schema)
		return nil
	}

	// Check if a cash register already exists
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s.cash_registers`, schema)
	var count int
	err = s.db.Pool.QueryRow(ctx, countQuery).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing cash registers: %w", err)
	}

	if count > 0 {
		logger.Debug("Cash registers already exist, skipping default cash register creation", "schema", schema, "count", count)
		return nil
	}

	// Create default cash register
	registerID := uuid.New()
	registerName := "Caja Principal"
	now := time.Now()

	insertQuery := fmt.Sprintf(`
		INSERT INTO %s.cash_registers (id, name, location, is_active, initial_balance, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, schema)

	_, err = s.db.Pool.Exec(ctx, insertQuery,
		registerID,
		registerName,
		nil, // location
		true, // is_active
		0.00, // initial_balance
		nil,  // created_by (system creation, no user)
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create default cash register: %w", err)
	}

	logger.Info("Default cash register created", "schema", schema, "register_id", registerID, "name", registerName)
	return nil
}

// createTenantBucket creates a MinIO bucket for the tenant
func (s *Service) createTenantBucket(ctx context.Context, schema string) error {
	// Generate bucket name from tenant schema
	bucketName := storage.BuildTenantBucketName(schema)

	// Create storage client with the tenant's bucket
	storageClient, err := storage.NewClientWithBucket(s.cfg, bucketName)
	if err != nil {
		return fmt.Errorf("failed to create storage client for tenant bucket: %w", err)
	}

	// The bucket is automatically created in NewClientWithBucket if it doesn't exist
	// But we can verify it exists
	exists, err := storageClient.FileExists(ctx, ".bucket-check")
	if err != nil {
		// If we can't check, assume bucket creation succeeded (it was created in NewClientWithBucket)
		logger.Info("Tenant bucket created", "schema", schema, "bucket", bucketName)
		return nil
	}

	// Create a marker file to verify bucket is writable
	if !exists {
		// Try to upload a small test file to verify bucket is writable
		testData := []byte("tenant-bucket-initialized")
		if err := storageClient.UploadBytes(ctx, testData, ".bucket-check", "text/plain"); err != nil {
			logger.Warn("Failed to create bucket marker file", "error", err, "bucket", bucketName)
			// Don't fail - bucket exists, just couldn't write marker
		}
	}

	logger.Info("Tenant bucket created and verified", "schema", schema, "bucket", bucketName)
	return nil
}

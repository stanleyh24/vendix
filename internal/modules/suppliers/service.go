package suppliers

import (
	"context"
	"fmt"
	"time"

	"github.com/stanleyh24/vendix/internal/config"
	"github.com/stanleyh24/vendix/internal/database"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
	cfg  *config.Config
}

func NewService(db *database.DB, cfg *config.Config) *Service {
	return &Service{repo: NewRepository(db), cfg: cfg}
}

func (s *Service) Create(ctx context.Context, schema string, req *CreateSupplierRequest) (*Supplier, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	country := "DO"
	if req.Country != nil && *req.Country != "" {
		country = *req.Country
	}
	sup := &Supplier{
		ID:         uuid.New(),
		Name:       req.Name,
		TaxID:      req.TaxID,
		Email:      req.Email,
		Phone:      req.Phone,
		Address:    req.Address,
		City:       req.City,
		State:      req.State,
		PostalCode: req.PostalCode,
		Country:    country,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := s.repo.Create(ctx, schema, sup); err != nil {
		return nil, err
	}
	return sup, nil
}

func (s *Service) Update(ctx context.Context, schema string, id string, req *UpdateSupplierRequest) (*Supplier, error) {
	supplierID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid supplier id: %w", err)
	}
	existing, err := s.repo.GetByID(ctx, schema, supplierID)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.TaxID != nil {
		existing.TaxID = req.TaxID
	}
	if req.Email != nil {
		existing.Email = req.Email
	}
	if req.Phone != nil {
		existing.Phone = req.Phone
	}
	if req.Address != nil {
		existing.Address = req.Address
	}
	if req.City != nil {
		existing.City = req.City
	}
	if req.State != nil {
		existing.State = req.State
	}
	if req.PostalCode != nil {
		existing.PostalCode = req.PostalCode
	}
	if req.Country != nil && *req.Country != "" {
		existing.Country = *req.Country
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	existing.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, schema, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) Delete(ctx context.Context, schema string, id string) error {
	supplierID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid supplier id: %w", err)
	}
	return s.repo.Delete(ctx, schema, supplierID)
}

func (s *Service) GetByID(ctx context.Context, schema string, id string) (*Supplier, error) {
	supplierID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid supplier id: %w", err)
	}
	return s.repo.GetByID(ctx, schema, supplierID)
}

func (s *Service) List(ctx context.Context, schema string, onlyActive *bool) ([]*Supplier, error) {
	return s.repo.List(ctx, schema, onlyActive)
}

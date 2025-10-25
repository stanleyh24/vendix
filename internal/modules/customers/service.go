package customers

import (
	"context"
	"fmt"
	"time"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/logger"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
	cfg  *config.Config
}

func NewService(db *database.DB, cfg *config.Config) *Service {
	return &Service{
		repo: NewRepository(db),
		cfg:  cfg,
	}
}

func (s *Service) Create(ctx context.Context, schema, userID string, req *CreateCustomerRequest) (*CustomerResponse, error) {
	customer := &Customer{
		ID:           uuid.New(),
		CustomerType: req.CustomerType,
		TaxID:        req.TaxID,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Address:      req.Address,
		City:         req.City,
		State:        req.State,
		PostalCode:   req.PostalCode,
		Country:      req.Country,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if customer.Country == nil {
		defaultCountry := "DO"
		customer.Country = &defaultCountry
	}

	if err := s.repo.Create(ctx, schema, customer); err != nil {
		return nil, err
	}

	logger.Info("Customer created", "id", customer.ID, "schema", schema)

	return customer.ToResponse(), nil
}

func (s *Service) List(ctx context.Context, schema string) ([]*CustomerResponse, error) {
	customers, err := s.repo.List(ctx, schema)
	if err != nil {
		return nil, err
	}

	responses := make([]*CustomerResponse, len(customers))
	for i, c := range customers {
		responses[i] = c.ToResponse()
	}

	return responses, nil
}

func (s *Service) GetByID(ctx context.Context, schema, id string) (*CustomerResponse, error) {
	customerID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid customer ID: %w", err)
	}

	customer, err := s.repo.GetByID(ctx, schema, customerID)
	if err != nil {
		return nil, err
	}

	return customer.ToResponse(), nil
}

func (s *Service) Update(ctx context.Context, schema, id string, req *UpdateCustomerRequest) (*CustomerResponse, error) {
	customerID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid customer ID: %w", err)
	}

	customer, err := s.repo.GetByID(ctx, schema, customerID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		customer.Name = *req.Name
	}
	if req.Email != nil {
		customer.Email = req.Email
	}
	if req.Phone != nil {
		customer.Phone = req.Phone
	}
	if req.Address != nil {
		customer.Address = req.Address
	}
	if req.City != nil {
		customer.City = req.City
	}
	if req.State != nil {
		customer.State = req.State
	}
	if req.PostalCode != nil {
		customer.PostalCode = req.PostalCode
	}
	if req.IsActive != nil {
		customer.IsActive = *req.IsActive
	}

	customer.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, schema, customer); err != nil {
		return nil, err
	}

	logger.Info("Customer updated", "id", customer.ID, "schema", schema)

	return customer.ToResponse(), nil
}

func (s *Service) Delete(ctx context.Context, schema, id string) error {
	customerID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid customer ID: %w", err)
	}

	if err := s.repo.Delete(ctx, schema, customerID); err != nil {
		return err
	}

	logger.Info("Customer deleted", "id", customerID, "schema", schema)

	return nil
}

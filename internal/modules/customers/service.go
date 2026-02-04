package customers

import (
	"context"
	"fmt"
	"time"

	"github.com/stanleyh24/vendix/internal/config"
	"github.com/stanleyh24/vendix/internal/database"
	"github.com/stanleyh24/vendix/internal/logger"

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
	// Validar que cliente gubernamental tenga RNC
	isGovernmentEntity := false
	if req.IsGovernmentEntity != nil {
		isGovernmentEntity = *req.IsGovernmentEntity
	}
	
	if isGovernmentEntity {
		if req.TaxID == nil || *req.TaxID == "" {
			return nil, fmt.Errorf("clientes gubernamentales deben tener RNC válido")
		}
		// Validar formato de RNC (9 o 11 dígitos)
		if len(*req.TaxID) != 9 && len(*req.TaxID) != 11 {
			return nil, fmt.Errorf("RNC debe tener 9 o 11 dígitos")
		}
	}

	customer := &Customer{
		ID:                    uuid.New(),
		CustomerType:          req.CustomerType,
		TaxID:                 req.TaxID,
		Name:                  req.Name,
		Email:                 req.Email,
		Phone:                 req.Phone,
		Address:               req.Address,
		City:                  req.City,
		State:                 req.State,
		PostalCode:            req.PostalCode,
		Country:               req.Country,
		IsActive:              true,
		IsGovernmentEntity:    isGovernmentEntity,
		DefaultWithholdingRate: req.DefaultWithholdingRate,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	if customer.Country == nil {
		defaultCountry := "DO"
		customer.Country = &defaultCountry
	}

	if err := s.repo.Create(ctx, schema, customer); err != nil {
		return nil, err
	}

	logger.Info("Customer created", "id", customer.ID, "schema", schema, "is_government_entity", isGovernmentEntity)

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
	if req.IsGovernmentEntity != nil {
		isGovernmentEntity := *req.IsGovernmentEntity
		// Validar que si se marca como gubernamental, tenga RNC
		if isGovernmentEntity {
			if customer.TaxID == nil || *customer.TaxID == "" {
				return nil, fmt.Errorf("no se puede marcar como entidad gubernamental sin RNC válido. Por favor agregue el RNC primero")
			}
			if len(*customer.TaxID) != 9 && len(*customer.TaxID) != 11 {
				return nil, fmt.Errorf("RNC debe tener 9 o 11 dígitos")
			}
		}
		customer.IsGovernmentEntity = isGovernmentEntity
	}
	if req.DefaultWithholdingRate != nil {
		customer.DefaultWithholdingRate = req.DefaultWithholdingRate
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

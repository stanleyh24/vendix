package employees

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

func (s *Service) Create(ctx context.Context, schema, userID string, req *CreateEmployeeRequest) (*EmployeeResponse, error) {
	// Verificar que el código de empleado no exista
	existing, _ := s.repo.GetByCode(ctx, schema, req.EmployeeCode)
	if existing != nil {
		return nil, fmt.Errorf("employee code already exists")
	}

	defaultCountry := "DO"
	if req.Country == nil {
		req.Country = &defaultCountry
	}

	employee := &Employee{
		ID:           uuid.New(),
		EmployeeCode: req.EmployeeCode,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		TaxID:        req.TaxID,
		Address:      req.Address,
		City:         req.City,
		State:        req.State,
		PostalCode:   req.PostalCode,
		Country:      req.Country,
		Department:   req.Department,
		Position:     req.Position,
		HireDate:     req.HireDate,
		Salary:       req.Salary,
		SalaryType:   req.SalaryType,
		IsActive:     true,
		Metadata:     JSONB{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, schema, employee); err != nil {
		return nil, err
	}

	logger.Info("Employee created", "id", employee.ID, "code", employee.EmployeeCode, "schema", schema)

	return employee.ToResponse(), nil
}

func (s *Service) List(ctx context.Context, schema string) ([]*EmployeeResponse, error) {
	employees, err := s.repo.List(ctx, schema)
	if err != nil {
		return nil, err
	}

	responses := make([]*EmployeeResponse, len(employees))
	for i, e := range employees {
		responses[i] = e.ToResponse()
	}

	return responses, nil
}

func (s *Service) ListActive(ctx context.Context, schema string) ([]*EmployeeResponse, error) {
	employees, err := s.repo.ListActive(ctx, schema)
	if err != nil {
		return nil, err
	}

	responses := make([]*EmployeeResponse, len(employees))
	for i, e := range employees {
		responses[i] = e.ToResponse()
	}

	return responses, nil
}

func (s *Service) GetByID(ctx context.Context, schema, id string) (*EmployeeResponse, error) {
	employeeID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid employee ID: %w", err)
	}

	employee, err := s.repo.GetByID(ctx, schema, employeeID)
	if err != nil {
		return nil, err
	}

	return employee.ToResponse(), nil
}

func (s *Service) Update(ctx context.Context, schema, id string, req *UpdateEmployeeRequest) (*EmployeeResponse, error) {
	employeeID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid employee ID: %w", err)
	}

	employee, err := s.repo.GetByID(ctx, schema, employeeID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.FirstName != nil {
		employee.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		employee.LastName = *req.LastName
	}
	if req.Email != nil {
		employee.Email = req.Email
	}
	if req.Phone != nil {
		employee.Phone = req.Phone
	}
	if req.Address != nil {
		employee.Address = req.Address
	}
	if req.City != nil {
		employee.City = req.City
	}
	if req.State != nil {
		employee.State = req.State
	}
	if req.PostalCode != nil {
		employee.PostalCode = req.PostalCode
	}
	if req.Department != nil {
		employee.Department = req.Department
	}
	if req.Position != nil {
		employee.Position = req.Position
	}
	if req.Salary != nil {
		employee.Salary = req.Salary
	}
	if req.SalaryType != nil {
		employee.SalaryType = *req.SalaryType
	}
	if req.IsActive != nil {
		employee.IsActive = *req.IsActive
	}

	employee.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, schema, employee); err != nil {
		return nil, err
	}

	logger.Info("Employee updated", "id", employee.ID, "schema", schema)

	return employee.ToResponse(), nil
}

func (s *Service) Delete(ctx context.Context, schema, id string) error {
	employeeID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid employee ID: %w", err)
	}

	if err := s.repo.Delete(ctx, schema, employeeID); err != nil {
		return err
	}

	logger.Info("Employee deleted", "id", employeeID, "schema", schema)

	return nil
}

package employees

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type JSONB map[string]interface{}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = JSONB{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("employees: cannot scan type %T into JSONB", value)
	}

	if len(bytes) == 0 {
		*j = JSONB{}
		return nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return fmt.Errorf("employees: failed to unmarshal metadata: %w", err)
	}

	*j = JSONB(data)
	return nil
}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return []byte("{}"), nil
	}

	bytes, err := json.Marshal(map[string]interface{}(j))
	if err != nil {
		return nil, fmt.Errorf("employees: failed to marshal metadata: %w", err)
	}

	return bytes, nil
}

type Employee struct {
	ID           uuid.UUID `db:"id" json:"id"`
	EmployeeCode string    `db:"employee_code" json:"employee_code"`
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	Email        *string   `db:"email" json:"email,omitempty"`
	Phone        *string   `db:"phone" json:"phone,omitempty"`
	TaxID        *string   `db:"tax_id" json:"tax_id,omitempty"`
	Address      *string   `db:"address" json:"address,omitempty"`
	City         *string   `db:"city" json:"city,omitempty"`
	State        *string   `db:"state" json:"state,omitempty"`
	PostalCode   *string   `db:"postal_code" json:"postal_code,omitempty"`
	Country      *string   `db:"country" json:"country,omitempty"`
	Department   *string   `db:"department" json:"department,omitempty"`
	Position     *string   `db:"position" json:"position,omitempty"`
	HireDate     *time.Time `db:"hire_date" json:"hire_date,omitempty"`
	Salary       *float64  `db:"salary" json:"salary,omitempty"`
	SalaryType   string    `db:"salary_type" json:"salary_type"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	Metadata     JSONB     `db:"metadata" json:"metadata,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type CreateEmployeeRequest struct {
	EmployeeCode string     `json:"employee_code" validate:"required"`
	FirstName    string     `json:"first_name" validate:"required"`
	LastName     string     `json:"last_name" validate:"required"`
	Email        *string    `json:"email,omitempty"`
	Phone        *string    `json:"phone,omitempty"`
	TaxID        *string    `json:"tax_id,omitempty"`
	Address      *string    `json:"address,omitempty"`
	City         *string    `json:"city,omitempty"`
	State        *string    `json:"state,omitempty"`
	PostalCode   *string    `json:"postal_code,omitempty"`
	Country      *string    `json:"country,omitempty"`
	Department   *string    `json:"department,omitempty"`
	Position     *string    `json:"position,omitempty"`
	HireDate     *time.Time `json:"hire_date,omitempty"`
	Salary       *float64   `json:"salary,omitempty"`
	SalaryType   string     `json:"salary_type" validate:"required,oneof=monthly hourly daily"`
}

type UpdateEmployeeRequest struct {
	FirstName  *string    `json:"first_name,omitempty"`
	LastName   *string    `json:"last_name,omitempty"`
	Email      *string    `json:"email,omitempty"`
	Phone      *string    `json:"phone,omitempty"`
	Address    *string    `json:"address,omitempty"`
	City       *string    `json:"city,omitempty"`
	State      *string    `json:"state,omitempty"`
	PostalCode *string    `json:"postal_code,omitempty"`
	Department *string    `json:"department,omitempty"`
	Position   *string    `json:"position,omitempty"`
	Salary     *float64   `json:"salary,omitempty"`
	SalaryType *string    `json:"salary_type,omitempty"`
	IsActive   *bool      `json:"is_active,omitempty"`
}

type EmployeeResponse struct {
	ID           string     `json:"id"`
	EmployeeCode string     `json:"employee_code"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Email        *string    `json:"email,omitempty"`
	Phone        *string    `json:"phone,omitempty"`
	TaxID        *string    `json:"tax_id,omitempty"`
	Address      *string    `json:"address,omitempty"`
	City         *string    `json:"city,omitempty"`
	State        *string    `json:"state,omitempty"`
	PostalCode   *string    `json:"postal_code,omitempty"`
	Country      *string    `json:"country,omitempty"`
	Department   *string    `json:"department,omitempty"`
	Position     *string    `json:"position,omitempty"`
	HireDate     *time.Time `json:"hire_date,omitempty"`
	Salary       *float64   `json:"salary,omitempty"`
	SalaryType   string     `json:"salary_type"`
	IsActive     bool       `json:"is_active"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (e *Employee) ToResponse() *EmployeeResponse {
	var metadata map[string]interface{}
	if e.Metadata != nil {
		metadata = map[string]interface{}(e.Metadata)
	}

	return &EmployeeResponse{
		ID:           e.ID.String(),
		EmployeeCode: e.EmployeeCode,
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		Email:        e.Email,
		Phone:        e.Phone,
		TaxID:        e.TaxID,
		Address:      e.Address,
		City:         e.City,
		State:        e.State,
		PostalCode:   e.PostalCode,
		Country:      e.Country,
		Department:   e.Department,
		Position:     e.Position,
		HireDate:     e.HireDate,
		Salary:       e.Salary,
		SalaryType:   e.SalaryType,
		IsActive:     e.IsActive,
		Metadata:     metadata,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}


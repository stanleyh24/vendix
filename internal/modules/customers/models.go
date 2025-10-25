package customers

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID           uuid.UUID `db:"id" json:"id"`
	CustomerType string    `db:"customer_type" json:"customer_type"`
	TaxID        *string   `db:"tax_id" json:"tax_id,omitempty"`
	Name         string    `db:"name" json:"name"`
	Email        *string   `db:"email" json:"email,omitempty"`
	Phone        *string   `db:"phone" json:"phone,omitempty"`
	Address      *string   `db:"address" json:"address,omitempty"`
	City         *string   `db:"city" json:"city,omitempty"`
	State        *string   `db:"state" json:"state,omitempty"`
	PostalCode   *string   `db:"postal_code" json:"postal_code,omitempty"`
	Country      *string   `db:"country" json:"country,omitempty"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type CreateCustomerRequest struct {
	CustomerType string  `json:"customer_type" validate:"required"`
	TaxID        *string `json:"tax_id,omitempty"`
	Name         string  `json:"name" validate:"required"`
	Email        *string `json:"email,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	Address      *string `json:"address,omitempty"`
	City         *string `json:"city,omitempty"`
	State        *string `json:"state,omitempty"`
	PostalCode   *string `json:"postal_code,omitempty"`
	Country      *string `json:"country,omitempty"`
}

type UpdateCustomerRequest struct {
	Name       *string `json:"name,omitempty"`
	Email      *string `json:"email,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	Address    *string `json:"address,omitempty"`
	City       *string `json:"city,omitempty"`
	State      *string `json:"state,omitempty"`
	PostalCode *string `json:"postal_code,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

type CustomerResponse struct {
	ID           string    `json:"id"`
	CustomerType string    `json:"customer_type"`
	TaxID        *string   `json:"tax_id,omitempty"`
	Name         string    `json:"name"`
	Email        *string   `json:"email,omitempty"`
	Phone        *string   `json:"phone,omitempty"`
	Address      *string   `json:"address,omitempty"`
	City         *string   `json:"city,omitempty"`
	State        *string   `json:"state,omitempty"`
	PostalCode   *string   `json:"postal_code,omitempty"`
	Country      *string   `json:"country,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (c *Customer) ToResponse() *CustomerResponse {
	return &CustomerResponse{
		ID:           c.ID.String(),
		CustomerType: c.CustomerType,
		TaxID:        c.TaxID,
		Name:         c.Name,
		Email:        c.Email,
		Phone:        c.Phone,
		Address:      c.Address,
		City:         c.City,
		State:        c.State,
		PostalCode:   c.PostalCode,
		Country:      c.Country,
		IsActive:     c.IsActive,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

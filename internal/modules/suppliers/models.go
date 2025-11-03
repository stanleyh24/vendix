package suppliers

import (
	"time"

	"github.com/google/uuid"
)

// Supplier represents a vendor/supplier entity
type Supplier struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	TaxID      *string   `db:"tax_id" json:"tax_id,omitempty"`
	Email      *string   `db:"email" json:"email,omitempty"`
	Phone      *string   `db:"phone" json:"phone,omitempty"`
	Address    *string   `db:"address" json:"address,omitempty"`
	City       *string   `db:"city" json:"city,omitempty"`
	State      *string   `db:"state" json:"state,omitempty"`
	PostalCode *string   `db:"postal_code" json:"postal_code,omitempty"`
	Country    string    `db:"country" json:"country"`
	IsActive   bool      `db:"is_active" json:"is_active"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// CreateSupplierRequest DTO
type CreateSupplierRequest struct {
	Name       string  `json:"name"`
	TaxID      *string `json:"tax_id"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	Address    *string `json:"address"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	PostalCode *string `json:"postal_code"`
	Country    *string `json:"country"`
}

// UpdateSupplierRequest DTO
type UpdateSupplierRequest struct {
	Name       *string `json:"name"`
	TaxID      *string `json:"tax_id"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	Address    *string `json:"address"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	PostalCode *string `json:"postal_code"`
	Country    *string `json:"country"`
	IsActive   *bool   `json:"is_active"`
}

package products

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	Code          string     `db:"code" json:"code"`
	Name          string     `db:"name" json:"name"`
	Description   *string    `db:"description" json:"description,omitempty"`
	ProductType   string     `db:"product_type" json:"product_type"`
	Unit          string     `db:"unit" json:"unit"`
	Price         float64    `db:"price" json:"price"`
	Cost          *float64   `db:"cost" json:"cost,omitempty"`
	TaxRate       float64    `db:"tax_rate" json:"tax_rate"`
	StockQuantity float64    `db:"stock_quantity" json:"stock_quantity"`
	SupplierID    *uuid.UUID `db:"supplier_id" json:"supplier_id,omitempty"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}

type CreateProductRequest struct {
	Code          string     `json:"code" validate:"required"`
	Name          string     `json:"name" validate:"required"`
	Description   *string    `json:"description,omitempty"`
	ProductType   string     `json:"product_type" validate:"required"`
	Unit          string     `json:"unit" validate:"required"`
	Price         float64    `json:"price" validate:"required"`
	Cost          *float64   `json:"cost,omitempty"`
	TaxRate       float64    `json:"tax_rate"`
	StockQuantity float64    `json:"stock_quantity"`
	SupplierID    *uuid.UUID `json:"supplier_id,omitempty"`
}

type UpdateProductRequest struct {
	Name          *string    `json:"name,omitempty"`
	Description   *string    `json:"description,omitempty"`
	Unit          *string    `json:"unit,omitempty"`
	Price         *float64   `json:"price,omitempty"`
	Cost          *float64   `json:"cost,omitempty"`
	TaxRate       *float64   `json:"tax_rate,omitempty"`
	StockQuantity *float64   `json:"stock_quantity,omitempty"`
	SupplierID    *uuid.UUID `json:"supplier_id,omitempty"`
	IsActive      *bool      `json:"is_active,omitempty"`
}

type ProductResponse struct {
	ID            string    `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Description   *string   `json:"description,omitempty"`
	ProductType   string    `json:"product_type"`
	Unit          string    `json:"unit"`
	Price         float64   `json:"price"`
	Cost          *float64  `json:"cost,omitempty"`
	TaxRate       float64   `json:"tax_rate"`
	StockQuantity float64   `json:"stock_quantity"`
	SupplierID    *string   `json:"supplier_id,omitempty"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (p *Product) ToResponse() *ProductResponse {
	return &ProductResponse{
		ID:            p.ID.String(),
		Code:          p.Code,
		Name:          p.Name,
		Description:   p.Description,
		ProductType:   p.ProductType,
		Unit:          p.Unit,
		Price:         p.Price,
		Cost:          p.Cost,
		TaxRate:       p.TaxRate,
		StockQuantity: p.StockQuantity,
		SupplierID: func() *string {
			if p.SupplierID == nil {
				return nil
			}
			s := p.SupplierID.String()
			return &s
		}(),
		IsActive:  p.IsActive,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

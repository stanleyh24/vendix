package products

import (
	"context"
	"fmt"
	"time"

	"vendix/internal/config"
	"vendix/internal/database"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
	cfg  *config.Config
}

func NewService(db *database.DB, cfg *config.Config) *Service {
	return &Service{repo: NewRepository(db), cfg: cfg}
}

func (s *Service) Create(ctx context.Context, schema string, req *CreateProductRequest) (*ProductResponse, error) {
	product := &Product{
		ID:            uuid.New(),
		Code:          req.Code,
		Name:          req.Name,
		Description:   req.Description,
		ProductType:   req.ProductType,
		Unit:          req.Unit,
		Price:         req.Price,
		Cost:          req.Cost,
		TaxRate:       req.TaxRate,
		StockQuantity: req.StockQuantity,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := s.repo.Create(ctx, schema, product); err != nil {
		return nil, err
	}
	return product.ToResponse(), nil
}

func (s *Service) List(ctx context.Context, schema string) ([]*ProductResponse, error) {
	products, err := s.repo.List(ctx, schema)
	if err != nil {
		return nil, err
	}
	responses := make([]*ProductResponse, len(products))
	for i, p := range products {
		responses[i] = p.ToResponse()
	}
	return responses, nil
}

func (s *Service) GetByID(ctx context.Context, schema, id string) (*ProductResponse, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}
	product, err := s.repo.GetByID(ctx, schema, productID)
	if err != nil {
		return nil, err
	}
	return product.ToResponse(), nil
}

func (s *Service) Update(ctx context.Context, schema, id string, req *UpdateProductRequest) (*ProductResponse, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}
	product, err := s.repo.GetByID(ctx, schema, productID)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = req.Description
	}
	if req.Unit != nil {
		product.Unit = *req.Unit
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Cost != nil {
		product.Cost = req.Cost
	}
	if req.TaxRate != nil {
		product.TaxRate = *req.TaxRate
	}
	if req.StockQuantity != nil {
		product.StockQuantity = *req.StockQuantity
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	product.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, schema, product); err != nil {
		return nil, err
	}
	return product.ToResponse(), nil
}

func (s *Service) Delete(ctx context.Context, schema, id string) error {
	productID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}
	return s.repo.Delete(ctx, schema, productID)
}

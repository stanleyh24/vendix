package products

import (
	"context"
	"fmt"
	"time"

	"vendix/internal/database"

	"github.com/google/uuid"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, schema string, product *Product) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.products (id, code, name, description, product_type, unit, price, cost, tax_rate, stock_quantity, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, schema)
	_, err := r.db.ExecContext(ctx, query, product.ID, product.Code, product.Name, product.Description, product.ProductType, product.Unit, product.Price, product.Cost, product.TaxRate, product.StockQuantity, product.IsActive, product.CreatedAt, product.UpdatedAt)
	return err
}

func (r *Repository) GetByID(ctx context.Context, schema string, id uuid.UUID) (*Product, error) {
	var product Product
	query := fmt.Sprintf(`SELECT id, code, name, description, product_type, unit, price, cost, tax_rate, stock_quantity, is_active, created_at, updated_at FROM %s.products WHERE id = $1`, schema)
	err := r.db.GetContext(ctx, &product, query, id)
	return &product, err
}

func (r *Repository) List(ctx context.Context, schema string) ([]*Product, error) {
	var products []*Product
	query := fmt.Sprintf(`SELECT id, code, name, description, product_type, unit, price, cost, tax_rate, stock_quantity, is_active, created_at, updated_at FROM %s.products ORDER BY name ASC`, schema)
	err := r.db.SelectContext(ctx, &products, query)
	return products, err
}

func (r *Repository) Update(ctx context.Context, schema string, product *Product) error {
	query := fmt.Sprintf(`UPDATE %s.products SET name = $1, description = $2, unit = $3, price = $4, cost = $5, tax_rate = $6, stock_quantity = $7, is_active = $8, updated_at = $9 WHERE id = $10`, schema)
	_, err := r.db.ExecContext(ctx, query, product.Name, product.Description, product.Unit, product.Price, product.Cost, product.TaxRate, product.StockQuantity, product.IsActive, product.UpdatedAt, product.ID)
	return err
}

func (r *Repository) Delete(ctx context.Context, schema string, id uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s.products WHERE id = $1`, schema)
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// UpdateStock actualiza la cantidad en stock de un producto
func (r *Repository) UpdateStock(ctx context.Context, schema string, productID uuid.UUID, quantity float64) error {
	query := fmt.Sprintf(`UPDATE %s.products SET stock_quantity = stock_quantity + $1, updated_at = $2 WHERE id = $3`, schema)
	_, err := r.db.ExecContext(ctx, query, quantity, time.Now(), productID)
	return err
}

// GetStockQuantity obtiene la cantidad actual en stock de un producto
func (r *Repository) GetStockQuantity(ctx context.Context, schema string, productID uuid.UUID) (float64, error) {
	var quantity float64
	query := fmt.Sprintf(`SELECT stock_quantity FROM %s.products WHERE id = $1`, schema)
	err := r.db.GetContext(ctx, &quantity, query, productID)
	return quantity, err
}

// RecordInventoryMovement registra un movimiento de inventario para auditoría
func (r *Repository) RecordInventoryMovement(ctx context.Context, schema string, productID uuid.UUID, movementType string, quantity, previousStock, newStock float64, referenceType *string, referenceID *uuid.UUID, notes *string, createdBy *uuid.UUID) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.inventory_movements (id, product_id, movement_type, quantity, previous_stock, new_stock, reference_type, reference_id, notes, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, schema)
	_, err := r.db.ExecContext(ctx, query, uuid.New(), productID, movementType, quantity, previousStock, newStock, referenceType, referenceID, notes, createdBy, time.Now())
	return err
}

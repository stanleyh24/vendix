package sales

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"

	"github.com/stanleyh24/vendix/internal/config"
	"github.com/stanleyh24/vendix/internal/database"
	"github.com/stanleyh24/vendix/internal/testutil"
)

func mustPtr[T any](v T) *T { return &v }

// seedProduct inserts a simple product record in the tenant schema for testing.
func seedProduct(t *testing.T, db *database.DB, schema string, name string, price, taxRate, stock float64) uuid.UUID {
	t.Helper()
	id := uuid.New()
	q := `
		INSERT INTO %s.products (id, code, name, product_type, unit, price, tax_rate, is_active, stock_quantity, created_at, updated_at)
		VALUES ($1, $2, $3, 'product', 'unit', $4, $5, true, $6, NOW(), NOW())
	`
	query := fmt.Sprintf(q, schema)
	if _, err := db.ExecContext(context.Background(), query, id, id.String()[0:8], name, price, taxRate, stock); err != nil {
		t.Fatalf("failed to seed product: %v", err)
	}
	return id
}

func TestService_Create_SimpleSale(t *testing.T) {
	cfg := config.Load()
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		t.Skipf("database not available: %v", err)
	}

	defer db.Close()

	schema := testutil.CreateIsolatedSchema(t, db)
	defer testutil.DropSchema(t, db, schema)
	testutil.MigrateTenant(t, db, schema)

	// Seed one product with stock
	productID := seedProduct(t, db, schema, "Producto Prueba", 100.00, 0.18, 50)

	svc := NewService(db, cfg)

	notes := "Venta de prueba"
	req := &CreateSaleRequest{
		CustomerID:  nil, // fuerza cliente genérico
		PaymentType: string(PaymentTypeCash),
		Notes:       &notes,
		Lines: []CreateSaleLineReq{
			{
				ProductID:   &productID,
				Description: "Producto Prueba",
				Quantity:    2,
				UnitPrice:   100.0,
				TaxRate:     0.18,
			},
		},
	}

	sale, err := svc.Create(context.Background(), schema, req)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if sale.ID == uuid.Nil {
		t.Fatalf("expected sale ID to be set")
	}
	if sale.Total <= 0 {
		t.Fatalf("expected total > 0, got %v", sale.Total)
	}
	if sale.InvoiceID == nil {
		t.Fatalf("expected invoice to be created and linked")
	}
	if sale.Status != string(SaleStatusCompleted) {
		t.Fatalf("expected status completed, got %s", sale.Status)
	}
}

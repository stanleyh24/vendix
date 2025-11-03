package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"

	"vendix/internal/config"
	"vendix/internal/database"
)

// NewTestDB loads config from env and returns a connected DB for tests.
func NewTestDB(t *testing.T) *database.DB {
	t.Helper()
	cfg := config.Load()
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		t.Fatalf("failed to connect test database: %v", err)
	}
	return db
}

// CreateIsolatedSchema creates a random schema name for test isolation.
func CreateIsolatedSchema(t *testing.T, db *database.DB) string {
	t.Helper()
	schema := fmt.Sprintf("test_%s", uuid.New().String()[0:8])
	if err := db.CreateSchema(context.Background(), schema); err != nil {
		t.Fatalf("failed to create test schema %s: %v", schema, err)
	}
	return schema
}

// MigrateTenant runs tenant migrations on the provided schema.
func MigrateTenant(t *testing.T, db *database.DB, schema string) {
	t.Helper()
	if err := database.RunTenantMigrations(db, schema); err != nil {
		t.Fatalf("failed to run tenant migrations on %s: %v", schema, err)
	}
}

// DropSchema drops the schema after tests.
func DropSchema(t *testing.T, db *database.DB, schema string) {
	t.Helper()
	if err := db.DropSchema(context.Background(), schema); err != nil {
		t.Fatalf("failed to drop test schema %s: %v", schema, err)
	}
}

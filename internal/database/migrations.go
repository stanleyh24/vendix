package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/stanleyh24/vendix/internal/logger"
)

// Migration represents a database migration
type Migration struct {
	Version int
	Name    string
	SQL     string
}

// RunPublicMigrations runs migrations for the public schema
func RunPublicMigrations(db *DB) error {
	ctx := context.Background()

	// Create migrations tracking table
	if err := createMigrationsTable(ctx, db, "public"); err != nil {
		return err
	}

	// Get public migrations
	migrations := getPublicMigrations()

	// Run migrations
	for _, migration := range migrations {
		applied, err := isMigrationApplied(ctx, db, "public", migration.Version)
		if err != nil {
			return err
		}

		if applied {
			logger.Debug("Migration already applied",
				"schema", "public",
				"version", migration.Version,
				"name", migration.Name,
			)
			continue
		}

		logger.Info("Running migration",
			"schema", "public",
			"version", migration.Version,
			"name", migration.Name,
		)

		if err := runMigration(ctx, db, "public", migration); err != nil {
			return fmt.Errorf("failed to run migration %d: %w", migration.Version, err)
		}

		logger.Info("Migration completed",
			"schema", "public",
			"version", migration.Version,
		)
	}

	return nil
}

// RunTenantMigrations runs migrations for a tenant schema
func RunTenantMigrations(db *DB, tenantSchema string) error {
	ctx := context.Background()

	// Create schema if it doesn't exist
	if err := db.CreateSchema(ctx, tenantSchema); err != nil {
		return err
	}

	// Create migrations tracking table
	if err := createMigrationsTable(ctx, db, tenantSchema); err != nil {
		return err
	}

	// Get tenant migrations
	migrations := getTenantMigrations()

	// Run migrations
	for _, migration := range migrations {
		applied, err := isMigrationApplied(ctx, db, tenantSchema, migration.Version)
		if err != nil {
			return err
		}

		if applied {
			logger.Debug("Migration already applied",
				"schema", tenantSchema,
				"version", migration.Version,
				"name", migration.Name,
			)
			continue
		}

		logger.Info("Running migration",
			"schema", tenantSchema,
			"version", migration.Version,
			"name", migration.Name,
		)

		if err := runMigration(ctx, db, tenantSchema, migration); err != nil {
			return fmt.Errorf("failed to run migration %d: %w", migration.Version, err)
		}

		logger.Info("Migration completed",
			"schema", tenantSchema,
			"version", migration.Version,
		)
	}

	return nil
}

func createMigrationsTable(ctx context.Context, db *DB, schema string) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.schema_migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`, schema)

	_, err := db.Pool.Exec(ctx, query)
	return err
}

func isMigrationApplied(ctx context.Context, db *DB, schema string, version int) (bool, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s.schema_migrations WHERE version = $1", schema)
	err := db.Pool.QueryRow(ctx, query, version).Scan(&count)
	return count > 0, err
}

func runMigration(ctx context.Context, db *DB, schema string, migration Migration) error {
	// Start transaction
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Set search path
	if _, err := tx.Exec(ctx, fmt.Sprintf("SET search_path TO %s, public", schema)); err != nil {
		return err
	}

	// Run migration SQL
	if _, err := tx.Exec(ctx, migration.SQL); err != nil {
		return err
	}

	// Record migration
	recordQuery := fmt.Sprintf(`
		INSERT INTO %s.schema_migrations (version, name)
		VALUES ($1, $2)
	`, schema)
	if _, err := tx.Exec(ctx, recordQuery, migration.Version, migration.Name); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// LoadMigrationsFromDir loads migrations from a directory
func LoadMigrationsFromDir(dir string) ([]Migration, error) {
	var migrations []Migration

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Parse version from filename (e.g., 001_create_users.sql)
		var version int
		var name string
		_, err := fmt.Sscanf(file.Name(), "%d_%s", &version, &name)
		if err != nil {
			continue
		}

		// Read SQL content
		content, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    strings.TrimSuffix(name, ".sql"),
			SQL:     string(content),
		})
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

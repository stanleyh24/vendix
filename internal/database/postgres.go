package database

import (
	"context"
	"fmt"
	"time"

	"github.com/stanleyh24/vendix/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
	Pool *pgxpool.Pool
}

func NewPostgresDB(cfg *config.Config) (*DB, error) {
	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	// Create sqlx connection
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create pgxpool connection for advanced features
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.DBMaxOpenConns)
	poolConfig.MinConns = int32(cfg.DBMaxIdleConns)
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &DB{
		DB:   db,
		Pool: pool,
	}, nil
}

func (db *DB) Close() error {
	if db.Pool != nil {
		db.Pool.Close()
	}
	if db.DB != nil {
		return db.DB.Close()
	}
	return nil
}

// SetSearchPath sets the schema search path for a connection
func (db *DB) SetSearchPath(ctx context.Context, schema string) error {
	query := fmt.Sprintf("SET search_path TO %s, public", schema)
	_, err := db.Pool.Exec(ctx, query)
	return err
}

// CreateSchema creates a new schema if it doesn't exist
func (db *DB) CreateSchema(ctx context.Context, schema string) error {
	query := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)
	_, err := db.Pool.Exec(ctx, query)
	return err
}

// DropSchema drops a schema and all its objects
func (db *DB) DropSchema(ctx context.Context, schema string) error {
	query := fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema)
	_, err := db.Pool.Exec(ctx, query)
	return err
}

// SchemaExists checks if a schema exists
func (db *DB) SchemaExists(ctx context.Context, schema string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM information_schema.schemata WHERE schema_name = $1
		)
	`
	err := db.Pool.QueryRow(ctx, query, schema).Scan(&exists)
	return exists, err
}

package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/adorufus/imgupper/config"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Database wraps the SQL DB connection
type Database struct {
	*sql.DB
}

// New creates a new database connection
func New(cfg config.DatabaseConfig) (*Database, error) {
	db, err := sql.Open(cfg.Driver, cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MaxIdle)
	db.SetConnMaxLifetime(cfg.Timeout)

	// Check connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{DB: db}, nil
}

// PingContext checks database connection with context
func (d *Database) PingContext(ctx context.Context) error {
	return d.DB.PingContext(ctx)
}

// internal/repository/health.go
package repository

import (
	"context"

	"github.com/adorufus/imgupper/pkg/database"
)

// HealthRepository defines the health repository interface
type HealthRepository interface {
	CheckConnection(ctx context.Context) error
}

// healthRepository implements HealthRepository
type healthRepository struct {
	db *database.Database
}

// NewHealthRepository creates a new HealthRepository
func NewHealthRepository(db *database.Database) HealthRepository {
	return &healthRepository{
		db: db,
	}
}

// CheckConnection checks database connection
func (r *healthRepository) CheckConnection(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

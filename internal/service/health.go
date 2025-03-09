// internal/service/health.go
package service

import (
	"context"
)

// HealthService defines the health service interface
type HealthService interface {
	CheckDatabase(ctx context.Context) error
}

// healthService implements HealthService
type healthService struct {
	deps Deps
}

// NewHealthService creates a new HealthService
func NewHealthService(deps Deps) HealthService {
	return &healthService{
		deps: deps,
	}
}

// CheckDatabase checks database connection
func (s *healthService) CheckDatabase(ctx context.Context) error {
	return s.deps.Repos.Health.CheckConnection(ctx)
}

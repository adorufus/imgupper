package service

import (
	"time"

	"github.com/adorufus/imgupper/internal/repository"
	"github.com/adorufus/imgupper/pkg/logger"
)

// Deps contains dependencies for services
type Deps struct {
	Repos  *repository.Repositories
	Logger logger.Logger
}

// Services contains all application services
type Services struct {
	User   UserService
	Health HealthService
	Auth   AuthService
	Cr2    Cr2Service
}

// NewServices creates a new Services instance
func NewServices(deps Deps, jwtSecret string) *Services {
	tokenDuration := 24 * time.Hour

	return &Services{
		User:   NewUserService(deps),
		Health: NewHealthService(deps),
		Cr2:    NewCr2Srvice(deps),
		Auth:   NewAuthService(deps, jwtSecret, tokenDuration),
	}
}

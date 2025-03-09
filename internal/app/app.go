package app

import (
	"github.com/adorufus/imgupper/config"
	"github.com/adorufus/imgupper/internal/handler"
	"github.com/adorufus/imgupper/internal/repository"
	"github.com/adorufus/imgupper/internal/service"
	"github.com/adorufus/imgupper/pkg/cloudflare"
	"github.com/adorufus/imgupper/pkg/database"
	"github.com/adorufus/imgupper/pkg/logger"
	"github.com/adorufus/imgupper/pkg/middleware"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gorilla/mux"
)

type App struct {
	Router     *mux.Router
	Config     *config.Config
	Logger     logger.Logger
	DB         *database.Database
	Cr2        *s3.Client
	Handlers   *handler.Handlers
	Services   *service.Services
	Repository *repository.Repositories
}

// NewApp creates a new application with all dependencies
func NewApp(cfg *config.Config) (*App, error) {
	// Initialize logger
	log, err := logger.New(cfg.Logger)
	if err != nil {
		return nil, err
	}

	// Initialize database
	db, err := database.New(cfg.Database)
	if err != nil {
		return nil, err
	}

	cr2, err := cloudflare.NewR2Client(cfg.Cloudflare)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	repos := repository.NewRepositories(db, cr2)

	// Initialize services with repositories
	services := service.NewServices(service.Deps{
		Repos:  repos,
		Logger: log,
	}, cfg.JWT.Secret)

	// Configure JWT middleware
	jwtConfig := middleware.JWTConfig{
		Secret:         cfg.JWT.Secret,
		ExpirationTime: cfg.JWT.ExpirationTime,
	}

	// Initialize handlers with services
	handlers := handler.NewHandlers(handler.Deps{
		Services:  services,
		Logger:    log,
		JWTConfig: jwtConfig,
	})

	// Initialize router with handlers
	router := mux.NewRouter()
	handlers.RegisterRoutes(router)

	return &App{
		Router:     router,
		Config:     cfg,
		Logger:     log,
		DB:         db,
		Cr2:        cr2,
		Handlers:   handlers,
		Services:   services,
		Repository: repos,
	}, nil
}

func (a *App) Close() error {
	return a.DB.Close()
}

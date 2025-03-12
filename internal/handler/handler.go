package handler

import (
	"github.com/adorufus/imgupper/internal/service"
	"github.com/adorufus/imgupper/pkg/logger"
	"github.com/adorufus/imgupper/pkg/middleware"
	"github.com/gorilla/mux"
)

// Deps contains dependencies for handlers
type Deps struct {
	Services  *service.Services
	Logger    logger.Logger
	JWTConfig middleware.JWTConfig
}

// Handlers contains all HTTP handlers
type Handlers struct {
	deps   Deps
	user   *UserHandler
	health *HealthHandler
	auth   *AuthHandler
	cr2    *Cr2Handler
}

// NewHandlers creates a new Handlers instance
func NewHandlers(deps Deps) *Handlers {
	return &Handlers{
		deps:   deps,
		user:   NewUserHandler(deps),
		health: NewHealthHandler(deps),
		auth:   NewAuthHandler(deps),
		cr2:    NewCr2Handler(deps),
	}
}

// RegisterRoutes registers all routes to the router
func (h *Handlers) RegisterRoutes(router *mux.Router) {
	// API v1 routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Health check
	api.HandleFunc("/health", h.health.Check).Methods("GET")

	// Auth routes
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", h.auth.Register).Methods("POST")
	auth.HandleFunc("/login", h.auth.Login).Methods("POST")

	// User routes - protected with JWT middleware
	users := api.PathPrefix("/users").Subrouter()
	users.Use(middleware.JWTAuth(h.deps.JWTConfig))
	users.HandleFunc("", h.user.Create).Methods("POST")
	users.HandleFunc("", h.user.GetAll).Methods("GET")
	users.HandleFunc("/{id}", h.user.GetByID).Methods("GET")
	users.HandleFunc("/{id}", h.user.Update).Methods("PUT")
	users.HandleFunc("/{id}", h.user.Delete).Methods("DELETE")

	object := api.PathPrefix("/object").Subrouter()
	object.Use(middleware.JWTAuth(h.deps.JWTConfig))
	object.HandleFunc("/upload", h.cr2.ObjectUpload).Methods("POST")
	object.HandleFunc("/mine", h.cr2.ObjectFetchByUserId).Methods("GET")
}

package handler

import (
	"net/http"

	"github.com/adorufus/imgupper/pkg/httputil"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	deps Deps
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(deps Deps) *HealthHandler {
	return &HealthHandler{
		deps: deps,
	}
}

// Check handles health check endpoint
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	if err := h.deps.Services.Health.CheckDatabase(r.Context()); err != nil {
		h.deps.Logger.Error("Health check failed", "error", err)
		httputil.JSONResponse(w, map[string]string{
			"status":  "error",
			"message": "Database connection failed",
		}, http.StatusServiceUnavailable)
		return
	}

	httputil.JSONResponse(w, map[string]string{
		"status":  "ok",
		"message": "Service is healthy",
	}, http.StatusOK)
}

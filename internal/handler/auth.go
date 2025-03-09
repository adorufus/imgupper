package handler

import (
	"encoding/json"
	"net/http"

	"github.com/adorufus/imgupper/internal/model"
	"github.com/adorufus/imgupper/pkg/httputil"
)

// AuthHandler handles auth-related requests
type AuthHandler struct {
	deps Deps
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(deps Deps) *AuthHandler {
	return &AuthHandler{
		deps: deps,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := h.deps.Services.Auth.Register(r.Context(), req)
	if err != nil {
		h.deps.Logger.Error("Registration failed", "error", err)
		httputil.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	httputil.JSONResponse(w, resp, http.StatusCreated)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := h.deps.Services.Auth.Login(r.Context(), req)
	if err != nil {
		h.deps.Logger.Error("Login failed", "error", err)
		httputil.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	httputil.JSONResponse(w, resp, http.StatusOK)
}

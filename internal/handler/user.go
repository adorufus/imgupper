package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/adorufus/imgupper/internal/model"
	"github.com/adorufus/imgupper/pkg/httputil"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	deps Deps
}

func NewUserHandler(deps Deps) *UserHandler {
	return &UserHandler{
		deps: deps,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		httputil.ErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createdUser, err := h.deps.Services.User.Create(r.Context(), user)

	if err != nil {
		h.deps.Logger.Error("Failed to create user", "error", err)
		httputil.ErrorResponse(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	httputil.JSONResponse(w, createdUser, http.StatusCreated)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		httputil.ErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.deps.Services.User.GetByID(r.Context(), id)
	if err != nil {
		h.deps.Logger.Error("Failed to get user", "error", err, "id", id)
		httputil.ErrorResponse(w, "User not found", http.StatusNotFound)
		return
	}

	httputil.JSONResponse(w, user, http.StatusOK)
}

// GetAll gets all users
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.deps.Services.User.GetAll(r.Context())
	if err != nil {
		h.deps.Logger.Error("Failed to get users", "error", err)
		httputil.ErrorResponse(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	httputil.JSONResponse(w, users, http.StatusOK)
}

// Update updates a user
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		httputil.ErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		httputil.ErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user.ID = id
	updatedUser, err := h.deps.Services.User.Update(r.Context(), user)
	if err != nil {
		h.deps.Logger.Error("Failed to update user", "error", err, "id", id)
		httputil.ErrorResponse(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	httputil.JSONResponse(w, updatedUser, http.StatusOK)
}

// Delete deletes a user
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		httputil.ErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.deps.Services.User.Delete(r.Context(), id); err != nil {
		h.deps.Logger.Error("Failed to delete user", "error", err, "id", id)
		httputil.ErrorResponse(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	httputil.JSONResponse(w, map[string]string{"message": "User deleted successfully"}, http.StatusOK)
}

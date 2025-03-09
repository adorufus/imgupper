package handler

import (
	"net/http"
	"strconv"

	"github.com/adorufus/imgupper/internal/model"
	"github.com/adorufus/imgupper/pkg/httputil"
)

type Cr2Handler struct {
	deps Deps
}

func NewCr2Handler(deps Deps) *Cr2Handler {
	return &Cr2Handler{
		deps: deps,
	}
}

func (h *Cr2Handler) ObjectUpload(w http.ResponseWriter, r *http.Request) {
	var req model.CR2UploadRequest

	if userIDStr := r.FormValue("user_id"); userIDStr != "" {
		userId, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			httputil.ErrorResponse(w, "Invalid user_id format", http.StatusBadRequest)
			return
		}

		req.UserID = userId
	} else {
		httputil.ErrorResponse(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(300 << 20) // 10MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get file from form
	object, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer object.Close()

	response, err := h.deps.Services.Cr2.ObjectUpload(r.Context(), req, object, handler)
	if err != nil {
		h.deps.Logger.Error("Unable to upload file", "error", err)
		httputil.ErrorResponse(w, "Unable to upload file "+err.Error(), http.StatusBadRequest)
		return
	}

	httputil.JSONResponse(w, response, http.StatusCreated)
}

func (h *Cr2Handler) ObjectFetchById(w http.ResponseWriter, r *http.Request) {

	response, err := h.deps.Services.Cr2.ObjectFetchById(r.Context(), 0)
	if err != nil {
		h.deps.Logger.Error("Unable to fetch object", "error", err)
		httputil.ErrorResponse(w, "Unable to fetch object: "+err.Error(), http.StatusNotFound)
		return
	}

	httputil.JSONResponse(w, response, http.StatusOK)
}

func (h *Cr2Handler) ObjectFetchByUserId(w http.ResponseWriter, r *http.Request) {

	response, err := h.deps.Services.Cr2.ObjectFetchByUserId(r.Context())
	if err != nil {
		h.deps.Logger.Error("Unable to fetch objects from user", "error", err)
		httputil.ErrorResponse(w, "Unable to fetch objects from user: "+err.Error(), http.StatusNotFound)
		return
	}

	httputil.JSONResponse(w, response, http.StatusOK)
}

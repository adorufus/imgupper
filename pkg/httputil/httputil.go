package httputil

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse writes an error response
func ErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	JSONResponse(w, map[string]string{"error": message}, statusCode)
}

// JSONResponse writes a JSON response
func JSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// If encoding fails, write a simple error message
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Failed to encode response"}`))
		}
	}
}

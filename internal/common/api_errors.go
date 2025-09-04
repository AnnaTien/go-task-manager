package common

import (
	"encoding/json"
	"net/http"
)

// APIErrorResponse defines a consistent structure for API error messages.
type APIErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// WriteError sends a structured JSON error response.
func WriteError(w http.ResponseWriter, status int, message string, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := APIErrorResponse{
		Status:  status,
		Message: message,
		Details: details,
	}
	json.NewEncoder(w).Encode(response)
}

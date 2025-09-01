package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sergicanet9/scv-go-tools/v3/observability"
)

// ErrorResponse sets the HTTP status code, and writes a plain text message to the client.
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprint(w, message)
}

// ErrorResponse sets the HTTP status code and writes a JSON-encoded payload to the client.
func SuccessResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		observability.Logger().Printf("failed to encode success response: %v", err)
	}
}

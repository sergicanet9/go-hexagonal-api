package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	wrappersv4 "github.com/sergicanet9/go-hexagonal-api/scvv4/wrappers"
	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
)

// SuccessResponse sets the HTTP status code and writes a JSON-encoded payload to the client.
func SuccessResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	httpResponse(w, statusCode, payload)
}

// ErrorResponse maps the HTTP status code, and writes a JSON-encoded error message to the client.
func ErrorResponse(w http.ResponseWriter, err error) {
	var statusCode int
	switch {
	case errors.Is(err, wrappers.ValidationErr):
		statusCode = http.StatusBadRequest
	case errors.Is(err, wrappers.NonExistentErr):
		statusCode = http.StatusNotFound
	case errors.Is(err, wrappers.UnauthorizedErr):
		statusCode = http.StatusUnauthorized
	case errors.Is(err, wrappersv4.UnauthenticatedErr):
		statusCode = http.StatusForbidden
	default:
		statusCode = http.StatusInternalServerError
	}

	payload := map[string]string{"error": err.Error()}
	httpResponse(w, statusCode, payload)
}

func httpResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(payload)
	if err != nil {
		response, _ = json.Marshal(map[string]string{"error": fmt.Sprintf("failed to marshal the response: %v", err)})
		statusCode = http.StatusInternalServerError
	}
	w.WriteHeader(statusCode)
	w.Write(response)
}

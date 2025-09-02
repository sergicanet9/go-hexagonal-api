package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorResponse sets the HTTP status code and writes a JSON-encoded payload to the client.
func SuccessResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	responseJSON(w, statusCode, payload)
}

// ErrorResponse sets the HTTP status code, and writes a JSON-encoded error message to the client.
func ErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	payload := map[string]string{"error": err.Error()}
	responseJSON(w, statusCode, payload)
}

func responseJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(payload)
	if err != nil {
		response, _ = json.Marshal(map[string]string{"error": fmt.Sprintf("failed to marshal the response: %v", err)})
		statusCode = http.StatusInternalServerError
	}
	w.WriteHeader(statusCode)
	w.Write(response)
}

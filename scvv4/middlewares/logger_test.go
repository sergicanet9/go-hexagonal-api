package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sergicanet9/go-hexagonal-api/scvv4/utils"
	"github.com/stretchr/testify/assert"
)

// TestLogger_HandlerOk checks that the middleware preserves the status code and response when the handler returns a successful response
func TestLogger_HandlerOk(t *testing.T) {
	// Arrange
	var url = "http://testing"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	expectedResponse := map[string]string{"response": "test-response"}

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.SuccessResponse(w, http.StatusOK, expectedResponse)
	})
	handlerToTest := Logger()(handlerFunc)

	// Act
	handlerToTest.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusOK, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response: %s", err)
	}
	assert.Equal(t, expectedResponse, response)
}

// TestLogger_HandlerError checks that the middleware preserves the status code and response when the handler returns an error response
func TestLogger_HandlerError(t *testing.T) {
	// Arrange
	var url = "http://testing"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	expectedResponse := map[string]string{"error": "test-error"}

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("test-error"))
	})
	handlerToTest := Logger()(handlerFunc)

	// Act
	handlerToTest.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response: %s", err)
	}
	assert.Equal(t, expectedResponse, response)
}

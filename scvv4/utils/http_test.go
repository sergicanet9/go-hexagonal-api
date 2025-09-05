package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSuccessResponse_Ok checks that SuccessResponse returns the expected response when everything goes as expected
func TestSuccessResponse_Ok(t *testing.T) {
	// Arrange
	var url = "http://testing"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	expectedResponse := map[string]string{"response": "test-response"}

	handlerToTest := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SuccessResponse(w, http.StatusOK, expectedResponse)
	})

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

// TestErrorResponse_Ok checks that ErrorRespnse returns the expected error response when everything goes as expected
func TestErrorResponse_Ok(t *testing.T) {
	// Arrange
	var url = "http://testing"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	expectedResponse := map[string]string{"error": "test-error-response"}

	handlerToTest := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ErrorResponse(w, errors.New("test-error-response"))
	})

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

// TestHTTPResponse_PayloadNotMarshalled checks that HTTPResponse returns the expected response when the response cannot be marshalled
func TestHTTPResponse_PayloadNotMarshalled(t *testing.T) {
	// Arrange
	var url = "http://testing"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	notMarshableResponse := map[string]interface{}{"response": make(chan int)}
	expectedResponse := map[string]string(map[string]string{"error": "failed to marshal the response: json: unsupported type: chan int"})

	handlerToTest := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpResponse(w, http.StatusOK, notMarshableResponse)
	})

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

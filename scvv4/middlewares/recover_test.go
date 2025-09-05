package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRecover_NoPanic checks that the middleware does not return an error when no panic happens in the handler
func TestRecover_NoPanic(t *testing.T) {
	// Arrange
	var url = "http://testing"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handlerToTest := Recover(handlerFunc)

	// Act
	handlerToTest.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusOK, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
}

// TestRecover_Panic checks that the middleware returns an error when there is a panic in the handler
func TestRecover_Panic(t *testing.T) {
	// Arrange
	var url = "http://testing"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	expectedResponse := map[string]string{"error": "recovered from panic during HTTP call GET , Panic: test panic"}

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})
	handlerToTest := Recover(handlerFunc)

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

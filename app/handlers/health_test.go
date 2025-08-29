package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sergicanet9/go-hexagonal-api/config"
)

// TestHealthCheck_Ok checks that healthCheck handler does not return an error when everything goes as expected
func TestHealthCheck_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	cfg := config.Config{}
	healthHandler := NewHealthHandler(cfg)
	SetHealthRoutes(r, healthHandler)

	rr := httptest.NewRecorder()
	url := "http://testing/health"
	req := httptest.NewRequest(http.MethodGet, url, nil)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusOK, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
}

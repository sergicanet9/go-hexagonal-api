package integration

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHealthCheck_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)

		// Act
		url := fmt.Sprintf("%s:%d/go-hexagonal-api/%s/health", cfg.Address, cfg.Port, database)

		req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", contentType)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		defer resp.Body.Close()

		// Assert
		if want, got := http.StatusOK, resp.StatusCode; want != got {
			t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
		}
	})
}
package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestJWT checks that the JWT middleware correctly handles all expected scenarios
func TestJWT(t *testing.T) {
	cases := []struct {
		name             string
		jwtToken         string
		requiredClaims   []string
		expectedCode     int
		expectedResponse map[string]string
	}{
		{
			name:             "Valid token and claims",
			jwtToken:         "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek",
			expectedCode:     http.StatusOK,
			expectedResponse: nil,
		},
		{
			name:             "Missing token",
			jwtToken:         "",
			expectedCode:     http.StatusUnauthorized,
			expectedResponse: map[string]string{"error": "authorization token is not provided"},
		},
		{
			name:             "Malformed token",
			jwtToken:         "123",
			expectedCode:     http.StatusUnauthorized,
			expectedResponse: map[string]string{"error": "invalid token format, should be Bearer + {token}"},
		},
		{
			name:             "Invalid signin method",
			jwtToken:         "Bearer eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.e30.",
			expectedCode:     http.StatusUnauthorized,
			expectedResponse: map[string]string{"error": "invalid token: signin method not valid"},
		},
		{
			name:             "Invalid secret",
			jwtToken:         "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.df4nTfuWWdndfrlIxF0iWUrrcANrM4bzKdbYa9VeAj8",
			expectedCode:     http.StatusUnauthorized,
			expectedResponse: map[string]string{"error": "invalid token: signature is invalid"},
		},
		{
			name:             "Missing required claim",
			jwtToken:         "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek",
			requiredClaims:   []string{"test-claim"},
			expectedCode:     http.StatusForbidden,
			expectedResponse: map[string]string{"error": "insufficient permissions: required claim 'test-claim' not found"},
		},
	}

	secret := "test-secret"
	url := "http://testing"

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, url, nil)
			if tt.jwtToken != "" {
				req.Header.Add("Authorization", tt.jwtToken)
			}

			handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			handlerToTest := JWT(secret, tt.requiredClaims...)(handlerFunc)

			handlerToTest.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedResponse != nil {
				var response map[string]string
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("unexpected error parsing the response: %s", err)
				}
				assert.Equal(t, tt.expectedResponse, response)
			}
		})
	}
}

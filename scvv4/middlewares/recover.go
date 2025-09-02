package middlewares

import (
	"fmt"
	"net/http"

	"github.com/sergicanet9/go-hexagonal-api/scvv4/utils"
)

// Recover is an HTTP middleware that recovers from panics, logs the error and writes to the response body of the incomming call
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Errorf("recovered from panic during HTTP call %s %s, Panic: %v", r.Method, r.URL.Path, err)
				utils.ErrorResponse(w, http.StatusInternalServerError, message)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

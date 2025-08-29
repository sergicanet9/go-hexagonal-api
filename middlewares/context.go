package middlewares

import (
	"context"
	"net/http"
)

// SetRequestContext sets the request context to the received context
func SetRequestContext(ctx context.Context) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

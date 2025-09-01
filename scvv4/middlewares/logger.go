package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/observability"
)

type responseWriterWrap struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseWriterWrap) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriterWrap) Write(data []byte) (int, error) {
	rw.body.Write(data)
	return rw.ResponseWriter.Write(data)
}

// Logger is configurable HTTP middleware that logs details of the incomming call
func Logger(skippedPaths ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			for _, prefix := range skippedPaths {
				if strings.HasPrefix(r.URL.Path, prefix) {
					next.ServeHTTP(w, r)
					return
				}
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				observability.Logger().Printf("failed to read request body, skiping body loging: %v", err)
			}

			r.Body = io.NopCloser(bytes.NewBuffer(body))

			rw := &responseWriterWrap{w, http.StatusOK, &bytes.Buffer{}}

			next.ServeHTTP(rw, r)

			latency := time.Since(start)

			observability.Logger().Printf(
				"HTTP Call: %s %s - Req Body: %s - Status: %d - Latency: %s - Resp Body: %s",
				r.Method, r.URL.Path, body, rw.statusCode, latency, rw.body.String(),
			)
		})
	}
}

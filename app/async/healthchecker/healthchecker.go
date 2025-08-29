package healthchecker

import (
	"context"
	"net/http"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/observability"
)

const contentType = "application/json"

func Run(ctx context.Context, cancel context.CancelFunc, url string, interval time.Duration) {
	defer cancel()
	defer func() {
		if rec := recover(); rec != nil {
			observability.Logger().Print("recovered panic in async process: %w", rec)
		}
	}()

	for ctx.Err() == nil {
		<-time.After(interval)

		start := time.Now()

		req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
		if err != nil {
			observability.Logger().Printf("async process failure, error: %s", err)
			continue
		}

		req.Header.Set("Content-Type", contentType)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			observability.Logger().Printf("async process failure, error: %s", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			observability.Logger().Printf("async process failure, error: expected status code 200 but got %d", resp.StatusCode)
			continue
		}

		elapsed := time.Since(start)
		observability.Logger().Printf("Health Check complete, time elapsed: %s", elapsed)
	}
}

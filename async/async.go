package async

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sergicanet9/go-hexagonal-api/config"
)

const contentType = "application/json"

type async struct {
	config config.Config
}

func New(cfg config.Config) *async {
	return &async{
		config: cfg,
	}
}

func (a async) Run(ctx context.Context, cancel context.CancelFunc) func() error {
	return func() error {
		go healthCheck(ctx, cancel, fmt.Sprintf("http://:%d/health", a.config.Port), a.config.Async.Interval.Duration)

		for ctx.Err() == nil {
			<-time.After(1 * time.Second)
		}
		return fmt.Errorf("async process stopped")
	}
}

func healthCheck(ctx context.Context, cancel context.CancelFunc, url string, interval time.Duration) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Print("recovered panic in async process: %w", rec)
		}
	}()
	defer cancel()

	for ctx.Err() == nil {
		<-time.After(interval)

		start := time.Now()

		req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
		if err != nil {
			log.Printf("async process failure, error: %s", err)
			continue
		}

		req.Header.Set("Content-Type", contentType)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("async process failure, error: %s", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("async process failure, error: expected status code 200 but got %d", resp.StatusCode)
			continue
		}

		elapsed := time.Since(start)
		log.Printf("Health Check complete, time elapsed: %s", elapsed)
	}
}

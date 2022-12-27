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
		go healthCheck(ctx, cancel, a.config.Address, a.config.Port, a.config.Async.Interval.Duration)

		for ctx.Err() == nil {
			<-time.After(1 * time.Second)
		}
		return fmt.Errorf("async process stopped")
	}
}

func healthCheck(ctx context.Context, cancel context.CancelFunc, address string, port int, interval time.Duration) {
	defer cancel()
	defer func() {
		if rec := recover(); rec != nil {
			log.Print("recovered panic in async process: %w", rec)
		}
	}()

	for {
		<-time.After(interval)

		start := time.Now()

		url := fmt.Sprintf("%s:%d/api/health", address, port)

		req, err := http.NewRequest(http.MethodPost, url, http.NoBody)
		if err != nil {
			log.Printf("async process failure, error: %s", err)
			continue
		}

		req.Header.Set("Content-Type", contentType)

		_, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("async process failure, error: %s", err)
			continue
		}

		elapsed := time.Since(start)
		log.Printf("Health Check complete, time elapsed: %s", elapsed)
	}
}

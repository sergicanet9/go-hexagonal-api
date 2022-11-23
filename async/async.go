package async

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/sergicanet9/go-hexagonal-api/config"
)

const contentType = "application/json"

type async struct {
	config  config.Config
	address string
}

func New(cfg config.Config, address string) *async {
	return &async{
		config:  cfg,
		address: address,
	}
}

func (a async) Run(ctx context.Context) {
	var g multierror.Group
	g.Go(healthCheck(ctx, a.address, a.config.Port, a.config.Async.Interval.Duration))

	if err := g.Wait().ErrorOrNil(); err != nil {
		log.Printf("async process stopped, error: %s", err)
	}
}

func healthCheck(ctx context.Context, address string, port int, interval time.Duration) func() error {
	return func() error {
		for ctx.Err() == nil {
			<-time.After(interval)

			start := time.Now()

			url := fmt.Sprintf("%s:%d/api/health", address, port)

			req, err := http.NewRequest(http.MethodPost, url, http.NoBody)
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

			defer resp.Body.Close()

			elapsed := time.Since(start)
			log.Printf("Health Check complete, time elapsed: %s", elapsed)
		}
		return ctx.Err()
	}
}

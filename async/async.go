package async

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/sergicanet9/go-mongo-restapi/config"
	"go.mongodb.org/mongo-driver/mongo"
)

const contentType = "application/json"

type Async struct {
	config config.Config
	db     *mongo.Database
}

func NewAsync(cfg config.Config, db *mongo.Database) *Async {
	return &Async{
		config: cfg,
		db:     db,
	}
}

func (a Async) Run(ctx context.Context) {
	var g multierror.Group
	g.Go(healthCheck(ctx, a.config.Address, a.config.Port, a.config.Async.Interval.Duration))

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

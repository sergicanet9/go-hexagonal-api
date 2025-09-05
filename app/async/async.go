package async

import (
	"context"
	"fmt"

	"github.com/sergicanet9/go-hexagonal-api/app/async/healthchecker"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/scv-go-tools/v4/observability"
)

type async struct {
	config config.Config
}

func New(cfg config.Config) async {
	return async{
		config: cfg,
	}
}

func (a async) Run(ctx context.Context, cancel context.CancelFunc) func() error {
	return func() error {
		go healthchecker.RunHTTP(ctx, cancel, fmt.Sprintf("http://:%d/health", a.config.HTTPPort), a.config.Async.Interval.Duration)
		go healthchecker.RunGRPC(ctx, cancel, fmt.Sprintf(":%d", a.config.GRPCPort), a.config.Async.Interval.Duration)

		<-ctx.Done()
		observability.Logger().Printf("Async process stopped")
		return nil
	}
}

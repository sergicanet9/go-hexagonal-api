package main

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/jessevdk/go-flags"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sergicanet9/go-hexagonal-api/app/api"
	"github.com/sergicanet9/go-hexagonal-api/app/async"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/scvv4/observability"
)

func main() {
	var opts struct {
		Version     string `long:"ver" description:"Version" required:"true"`
		Environment string `long:"env" description:"Environment" choice:"local" choice:"prod" required:"true"`
		HTTPPort    int    `long:"hport" description:"Running HTTP port" required:"true"`
		GRPCPort    int    `long:"gport" description:"Running gRPC port" required:"true"`
		Database    string `long:"db" description:"The database adapter to use" choice:"mongo" choice:"postgres" required:"true"`
		DSN         string `long:"dsn" description:"DSN of the selected database" required:"true"`
		NewRelicKey string `long:"nrkey" description:"New Relic Key" required:"false"`
	}

	args, err := flags.Parse(&opts)
	if err != nil {
		observability.Logger().Fatal(fmt.Errorf("provided flags not valid: %s, %w", args, err))
	}

	cfg, err := config.ReadConfig(opts.Version, opts.Environment, opts.HTTPPort, opts.GRPCPort, opts.Database, opts.DSN, opts.NewRelicKey, "config")
	if err != nil {
		observability.Logger().Fatal(fmt.Errorf("cannot parse config file for env %s: %w", opts.Environment, err))
	}

	var newrelicApp *newrelic.Application
	if cfg.NewRelicKey != "" {
		appName := fmt.Sprintf("go-hexagonal-api-%s-%s", cfg.Database, cfg.Environment)
		newrelicApp, err = observability.SetupNewRelic(appName, cfg.NewRelicKey)
		if err != nil {
			observability.Logger().Fatalf("could not set up new relic: %s", err)
		}
	}

	observability.Logger().Printf("Version: %s", cfg.Version)
	observability.Logger().Printf("Environment: %s", cfg.Environment)
	observability.Logger().Printf("Database: %s", cfg.Database)

	var g multierror.Group
	ctx, cancel := context.WithCancel(context.Background())
	grpcServerReady := make(chan struct{})

	a := api.New(ctx, cfg, newrelicApp)
	g.Go(a.RunGRPC(ctx, cancel, grpcServerReady))
	g.Go(a.RunHTTP(ctx, cancel, grpcServerReady))

	if cfg.Async.Run {
		async := async.New(cfg)
		g.Go(async.Run(ctx, cancel))
	}

	<-ctx.Done()
	observability.Logger().Printf("context canceled, the application will terminate...")

	done := make(chan struct{})
	go func() {
		g.Wait()
		close(done)
	}()

	select {
	case <-done:
		observability.Logger().Printf("application terminated gracefully")
	case <-time.After(10 * time.Second):
		observability.Logger().Fatalf("some processes did not terminate gracefully, application termination forced")
	}
}

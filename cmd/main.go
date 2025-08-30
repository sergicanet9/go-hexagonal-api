package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/jessevdk/go-flags"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sergicanet9/go-hexagonal-api/app/api"
	"github.com/sergicanet9/go-hexagonal-api/app/async"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/scv-go-tools/v3/observability"
)

// @title Go Hexagonal API
// @description Powered by scv-go-tools - https://github.com/sergicanet9/scv-go-tools

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	var opts struct {
		Version     string `long:"ver" description:"Version" required:"true"`
		Environment string `long:"env" description:"Environment" choice:"local" choice:"prod" required:"true"`
		Port        int    `long:"port" description:"Running port" required:"true"`
		Database    string `long:"db" description:"The database adapter to use" choice:"mongo" choice:"postgres" required:"true"`
		DSN         string `long:"dsn" description:"DSN of the selected database" required:"true"`
		NewRelicKey string `long:"nrkey" description:"New Relic Key" required:"false"`
	}

	args, err := flags.Parse(&opts)
	if err != nil {
		observability.Logger().Fatal(fmt.Errorf("provided flags not valid: %s, %w", args, err))
	}

	cfg, err := config.ReadConfig(opts.Version, opts.Environment, opts.Port, opts.Database, opts.DSN, opts.NewRelicKey, "config")
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

	var g multierror.Group
	ctx, cancel := context.WithCancel(context.Background())

	a := api.New(ctx, cfg, newrelicApp)
	// g.Go(a.Run(ctx, cancel))
	g.Go(a.RunGRPC(ctx, cancel))
	g.Go(a.RunGRPCGateway(ctx, cancel))

	if cfg.Async.Run {
		async := async.New(cfg)
		g.Go(async.Run(ctx, cancel))
	}

	if err := g.Wait().ErrorOrNil(); err != nil {
		observability.Logger().Fatal(err)
	}
}

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/jessevdk/go-flags"
	"github.com/sergicanet9/go-hexagonal-api/app/api"
	"github.com/sergicanet9/go-hexagonal-api/async"
	"github.com/sergicanet9/go-hexagonal-api/config"
)

// @title Go Hexagonal API
// @description Powered by scv-go-tools - https://github.com/sergicanet9/scv-go-tools

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	var opts struct {
		Version     string `long:"ver" description:"Version" required:"true"`
		Environment string `long:"env" description:"Environment" choice:"local" choice:"dev" required:"true"`
		Port        int    `long:"port" description:"Running port" required:"true"`
		Database    string `long:"db" description:"The database adapter to use" choice:"mongo" choice:"postgres" required:"true"`
		DSN         string `long:"dsn" description:"DSN of the selected database" required:"true"`
	}

	args, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(fmt.Errorf("provided flags not valid: %s, %w", args, err))
	}

	cfg, err := config.ReadConfig(opts.Version, opts.Environment, opts.Port, opts.Database, opts.DSN)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot parse config file for env %s: %w", opts.Environment, err))
	}

	var g multierror.Group
	ctx, cancel := context.WithCancel(context.Background())

	a := api.New(ctx, cfg)
	g.Go(a.Run(ctx, cancel))

	if cfg.Async.Run {
		async := async.New(cfg)
		g.Go(async.Run(ctx, cancel))
	}

	if err := g.Wait().ErrorOrNil(); err != nil {
		log.Fatal(err)
	}
}

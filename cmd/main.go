package main

import (
	"context"
	"flag"
	"fmt"
	"log"

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
	defaultPath := "."
	defaultPort := 8080
	defaultV, defaultEnv, defaultDB := "debug", "local", "mongo"
	versionF := flag.String("v", defaultV, "version")
	environmentF := flag.String("env", defaultEnv, "environment")
	portF := flag.Int("p", defaultPort, "port")
	databaseF := flag.String("db", defaultDB, "database")
	flag.Parse()

	cfg, err := config.ReadConfig(*versionF, *environmentF, *portF, *databaseF, defaultPath)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot parse config file in path %s for env %s: %w", defaultPath, *environmentF, err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a, addr := api.New(ctx, cfg)
	if cfg.Async.Run {
		async := async.New(cfg, addr)
		go async.Run(ctx)
	}

	a.Run(ctx)
}

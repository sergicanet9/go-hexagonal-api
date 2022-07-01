package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/sergicanet9/go-hexagonal-api/api"
	"github.com/sergicanet9/go-hexagonal-api/config"
)

// @title Go Hexagonal API
// @description Powered by scv-go-framework - https://github.com/sergicanet9/scv-go-framework

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	a := api.API{}
	a.Initialize(ctx, cfg)
	a.Run()
}

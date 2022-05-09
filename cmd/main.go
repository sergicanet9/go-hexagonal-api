package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/sergicanet9/go-mongo-restapi/api"
	"github.com/sergicanet9/go-mongo-restapi/config"
)

// @title Go Mongo RestAPI
// @description Powered by scv-go-framework - https://github.com/sergicanet9/scv-go-framework

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defaultPath := "."
	defaultEnv, defaultV := "local", "debug"
	environmentF := flag.String("env", defaultEnv, "environment")
	versionF := flag.String("v", defaultV, "version")
	flag.Parse()

	cfg, err := config.ReadConfig(*environmentF, *versionF, defaultPath)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot parse config file in path %s for env %s: %w", defaultPath, *environmentF, err))
	}

	a := api.API{}
	a.Initialize(ctx, cfg)
	a.Run()
}

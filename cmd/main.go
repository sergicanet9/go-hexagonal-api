package main

import (
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
	defaultPath := "."
	defaultEnv := "local"
	envF := flag.String("env", defaultEnv, "environment")
	flag.Parse()

	cfg, err := config.ReadConfig(*envF, defaultPath)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot parse config file in path %s for env %s: %w", defaultPath, *envF, err))
	}

	a := api.API{}
	a.Initialize(cfg)
	a.Run()
}

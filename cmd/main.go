package main

import "github.com/scanet9/go-mongo-restapi/api"

// @title Go Mongo RestAPI
// @version 1.0
// @description Powered by scv-go-framework - https://github.com/scanet9/scv-go-framework

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization (Format: Bearer {yourToken})

func main() {
	a := api.API{}
	a.Initialize()
	a.Run()
}

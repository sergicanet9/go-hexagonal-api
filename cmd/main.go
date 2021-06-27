package main

import "github.com/scanet9/go-mongo-restapi/api"

func main() {
	a := api.API{}
	a.Initialize()
	a.Run()
}

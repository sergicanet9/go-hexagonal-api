package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scanet9/go-mongo-restapi/api/handlers"
	"github.com/scanet9/go-mongo-restapi/business"
	"github.com/scanet9/go-mongo-restapi/config"
	"github.com/scanet9/scv-go-framework/infrastructure"
)

// API struct
type API struct {
	Router *mux.Router
}

// Initialize API
func (a *API) Initialize() {
	router := mux.NewRouter()
	a.Router = router

	db := infrastructure.ConnectDB(config.DbName, config.DbConnectionString)
	userService := business.NewUserService(db)
	handlers.SetUserRoutes(a.Router, userService)
}

// Run API
func (a *API) Run() {
	log.Printf("Listening on port %d", config.APIPort)
	log.Printf("Open http://localhost:%d in the browser", config.APIPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.APIPort), a.Router))
}

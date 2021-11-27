package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scanet9/go-mongo-restapi/api/handlers"
	"github.com/scanet9/go-mongo-restapi/business"
	"github.com/scanet9/go-mongo-restapi/config"
	_ "github.com/scanet9/go-mongo-restapi/docs" // docs is generated by Swag CLI, needs to be imported.
	infrastructure "github.com/scanet9/scv-go-framework/v2/infrastructure/mongo"
	httpSwagger "github.com/swaggo/http-swagger"
)

// API struct
type API struct {
	config config.Config
	router *mux.Router
}

// Initialize API
func (a *API) Initialize(cfg config.Config) {
	a.config = cfg

	router := mux.NewRouter()
	a.router = router

	a.router.PathPrefix("/swagger").Handler(
		httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("http://%s:%d/swagger/doc.json", a.config.Address, a.config.Port))),
	)

	db, err := infrastructure.ConnectMongoDB(a.config.DbName, a.config.DbConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	userService := business.NewUserService(a.config, db)
	handlers.SetUserRoutes(a.config, a.router, userService)
}

// Run API
func (a *API) Run() {
	log.Printf("Environment: %s", a.config.Env)
	log.Printf("Listening on port %d", a.config.Port)
	log.Printf("Open http://%s:%d/swagger/index.html in the browser", a.config.Address, a.config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", a.config.Port), a.router))
}

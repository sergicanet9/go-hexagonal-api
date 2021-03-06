package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/sergicanet9/go-hexagonal-api/app/docs" // docs is generated by Swag CLI, needs to be imported.
	"github.com/sergicanet9/go-hexagonal-api/app/handlers"
	"github.com/sergicanet9/go-hexagonal-api/async"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	"github.com/sergicanet9/go-hexagonal-api/core/services"
	"github.com/sergicanet9/go-hexagonal-api/infrastructure/mongo"
	"github.com/sergicanet9/go-hexagonal-api/infrastructure/postgres"
	"github.com/sergicanet9/scv-go-tools/v3/infrastructure"
	httpSwagger "github.com/swaggo/http-swagger"
)

// API struct
type API struct {
	config  config.Config
	address string
	router  *mux.Router
}

// Initialize API
func (a *API) Initialize(ctx context.Context, cfg config.Config) {
	a.config = cfg

	router := mux.NewRouter()
	a.router = router

	var userRepo ports.UserRepository
	switch a.config.Database {
	case "mongo":
		db, err := infrastructure.ConnectMongoDB(ctx, a.config.MongoDBName, a.config.MongoConnectionString)
		if err != nil {
			log.Fatal(err)
		}
		userRepo = mongo.NewUserRepository(db)
		a.address = a.config.MongoAddress
	case "postgres":
		db, err := infrastructure.ConnectPostgresDB(a.config.PostgresConnectionString)
		if err != nil {
			log.Fatal(err)
		}
		userRepo = postgres.NewUserRepository(db)
		a.address = a.config.PostgresAddress
	default:
		log.Fatalf("database flag %s not valid", a.config.Database)
	}

	userService := services.NewUserService(a.config, userRepo)

	handlers.SetHealthRoutes(ctx, a.config, a.router)
	handlers.SetUserRoutes(ctx, a.config, a.router, userService)

	a.router.PathPrefix("/swagger").Handler(
		httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("%s:%d/swagger/doc.json", a.address, a.config.Port))),
	)

	if a.config.Async.Run {
		async := async.NewAsync(a.config, a.address)
		go async.Run(ctx)
	}
}

// Run API
func (a *API) Run() {
	log.Printf("Version: %s", a.config.Version)
	log.Printf("Environment: %s", a.config.Environment)
	log.Printf("Database: %s", a.config.Database)
	log.Printf("Listening on port %d", a.config.Port)
	log.Printf("Open %s:%d/swagger/index.html in the browser", a.address, a.config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", a.config.Port), a.router))
}

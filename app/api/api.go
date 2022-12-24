package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/gorilla/mux"
	_ "github.com/sergicanet9/go-hexagonal-api/app/docs" // docs is generated by Swag CLI, needs to be imported.
	"github.com/sergicanet9/go-hexagonal-api/app/handlers"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	"github.com/sergicanet9/go-hexagonal-api/core/services"
	"github.com/sergicanet9/go-hexagonal-api/db/mongo"
	"github.com/sergicanet9/go-hexagonal-api/db/postgres"
	"github.com/sergicanet9/scv-go-tools/v3/infrastructure"
	httpSwagger "github.com/swaggo/http-swagger"
)

type api struct {
	config   config.Config
	address  string
	services svs
}

type svs struct {
	user ports.UserService
}

// New creates a new API
func New(ctx context.Context, cfg config.Config) (a api, addr string) {
	a.config = cfg

	var userRepo ports.UserRepository
	switch a.config.Database {
	case "mongo":
		db, err := infrastructure.ConnectMongoDB(ctx, a.config.DSN)
		if err != nil {
			log.Fatal(err)
		}

		userRepo = mongo.NewUserRepository(db)
		a.address = a.config.MongoAddress
	case "postgres":
		db, err := infrastructure.ConnectPostgresDB(ctx, a.config.DSN)
		if err != nil {
			log.Fatal(err)
		}

		_, filePath, _, _ := runtime.Caller(0)
		migrationsDir := filepath.Join(filePath, "../../..", cfg.PostgresMigrationsDir)
		err = infrastructure.MigratePostgresDB(db, migrationsDir)
		if err != nil {
			log.Fatal(err)
		}

		userRepo = postgres.NewUserRepository(db)
		a.address = a.config.PostgresAddress
	default:
		log.Fatalf("database flag %s not valid", a.config.Database)
	}

	a.services.user = services.NewUserService(a.config, userRepo)
	return a, a.address
}

// Run API
func (a *api) Run(ctx context.Context) {
	router := mux.NewRouter()

	handlers.SetHealthRoutes(ctx, a.config, router)
	handlers.SetUserRoutes(ctx, a.config, router, a.services.user)

	router.PathPrefix("/swagger").HandlerFunc(httpSwagger.WrapHandler)

	log.Printf("Version: %s", a.config.Version)
	log.Printf("Environment: %s", a.config.Environment)
	log.Printf("Database: %s", a.config.Database)
	log.Printf("Listening on port %d", a.config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", a.config.Port), router))
}

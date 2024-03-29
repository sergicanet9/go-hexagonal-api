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
	"github.com/sergicanet9/go-hexagonal-api/infrastructure/mongo"
	"github.com/sergicanet9/go-hexagonal-api/infrastructure/postgres"
	"github.com/sergicanet9/scv-go-tools/v3/infrastructure"
	httpSwagger "github.com/swaggo/http-swagger"
)

type api struct {
	config   config.Config
	services svs
}

type svs struct {
	user ports.UserService
}

// New creates a new API
func New(ctx context.Context, cfg config.Config) (a api) {
	a.config = cfg

	var userRepo ports.UserRepository
	switch a.config.Database {
	case "mongo":
		db, err := infrastructure.ConnectMongoDB(ctx, a.config.DSN)
		if err != nil {
			log.Fatal(err)
		}

		userRepo, err = mongo.NewUserRepository(ctx, db)
		if err != nil {
			log.Fatal(err)
		}
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
	default:
		log.Fatalf("database flag %s not valid", a.config.Database)
	}

	a.services.user = services.NewUserService(a.config, userRepo)
	return a
}

// Run API
func (a *api) Run(ctx context.Context, cancel context.CancelFunc) func() error {
	return func() error {
		defer cancel()

		router := mux.NewRouter()

		handlers.SetHealthRoutes(ctx, a.config, router)
		handlers.SetUserRoutes(ctx, a.config, router, a.services.user)
		router.PathPrefix("/swagger").HandlerFunc(httpSwagger.WrapHandler)

		log.Printf("Version: %s", a.config.Version)
		log.Printf("Environment: %s", a.config.Environment)
		log.Printf("Database: %s", a.config.Database)
		log.Printf("Listening on port %d", a.config.Port)

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", a.config.Port),
			Handler: router,
		}
		go shutdown(ctx, server)
		return server.ListenAndServe()
	}
}

func shutdown(ctx context.Context, server *http.Server) {
	<-ctx.Done()
	log.Printf("Shutting down API gracefully...")
	server.Shutdown(ctx)
}

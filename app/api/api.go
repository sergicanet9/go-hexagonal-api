package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/fullstorydev/grpcui/standalone"
	"github.com/gorilla/mux"
	grpcRuntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	"github.com/newrelic/go-agent/v3/newrelic"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/sergicanet9/go-hexagonal-api/app/handlers"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	"github.com/sergicanet9/go-hexagonal-api/core/services"
	"github.com/sergicanet9/go-hexagonal-api/infrastructure/mongo"
	"github.com/sergicanet9/go-hexagonal-api/infrastructure/postgres"
	"github.com/sergicanet9/go-hexagonal-api/proto/gen/v1/go/pb"
	"github.com/sergicanet9/scv-go-tools/v4/api/interceptors"
	"github.com/sergicanet9/scv-go-tools/v4/api/middlewares"
	"github.com/sergicanet9/scv-go-tools/v4/infrastructure"
	"github.com/sergicanet9/scv-go-tools/v4/observability"
)

type api struct {
	config      config.Config
	services    svs
	newrelicApp *newrelic.Application
}

type svs struct {
	user ports.UserService
}

// New creates a new API
func New(ctx context.Context, cfg config.Config, nrApp *newrelic.Application) (a api) {
	a.config = cfg
	a.newrelicApp = nrApp

	var userRepo ports.UserRepository
	switch a.config.Database {
	case "mongo":
		db, err := infrastructure.ConnectMongoDB(ctx, a.config.DSN)
		if err != nil {
			observability.Logger().Fatal(err)
		}

		userRepo, err = mongo.NewUserRepository(ctx, db)
		if err != nil {
			observability.Logger().Fatal(err)
		}
	case "postgres":
		db, err := infrastructure.ConnectPostgresDB(ctx, a.config.DSN)
		if err != nil {
			observability.Logger().Fatal(err)
		}

		_, filePath, _, _ := runtime.Caller(0)
		migrationsDir := filepath.Join(filePath, "../../..", cfg.PostgresMigrationsDir)
		err = infrastructure.MigratePostgresDB(db, migrationsDir)
		if err != nil {
			observability.Logger().Fatal(err)
		}

		userRepo = postgres.NewUserRepository(db)
	default:
		observability.Logger().Fatalf("database flag %s not valid", a.config.Database)
	}

	a.services.user = services.NewUserService(a.config, userRepo)
	return a
}

func (a *api) RunGRPC(ctx context.Context, cancel context.CancelFunc, grpcServerReady chan struct{}) func() error {
	return func() error {
		defer cancel()

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.config.GRPCPort))
		if err != nil {
			log.Fatalf("failed to listen on gRPC port: %s", err)
		}

		healthHander := handlers.NewHealthHandler(ctx, a.config)
		userHandler := handlers.NewUserHandler(ctx, a.config, a.services.user)

		methodPolicies := []interceptors.MethodPolicy{}
		methodPolicies = append(methodPolicies, userHandler.JWTMethodPolicies()...)

		server := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				interceptors.UnaryLogger(),
				interceptors.UnaryRecover(),
				nrgrpc.UnaryServerInterceptor(a.newrelicApp),
				interceptors.UnaryJWT(a.config.JWTSecret, methodPolicies),
			),
			grpc.ChainStreamInterceptor(
				interceptors.StreamLogger(),
				interceptors.StreamRecover(),
				nrgrpc.StreamServerInterceptor(a.newrelicApp),
				interceptors.StreamJWT(a.config.JWTSecret, methodPolicies),
			),
		)

		pb.RegisterHealthServiceServer(server, healthHander)
		pb.RegisterUserServiceServer(server, userHandler)

		reflection.Register(server)

		close(grpcServerReady)

		go shutdownGRPC(ctx, server)

		observability.Logger().Printf("Server listening on gRPC port %d", a.config.GRPCPort)
		return server.Serve(lis)
	}
}

func shutdownGRPC(ctx context.Context, server *grpc.Server) {
	<-ctx.Done()
	observability.Logger().Printf("Shutting down gRPC server gracefully...")
	server.GracefulStop()
}

func (a *api) RunHTTP(ctx context.Context, cancel context.CancelFunc, grpcServerReady chan struct{}) func() error {
	return func() error {
		defer cancel()

		<-grpcServerReady

		grpcServerAddr := fmt.Sprintf(":%d", a.config.GRPCPort)

		gmux := grpcRuntime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		err := pb.RegisterHealthServiceHandlerFromEndpoint(ctx, gmux, grpcServerAddr, opts)
		if err != nil {
			observability.Logger().Fatalf("failed to register health handler gateway: %s", err)
		}

		err = pb.RegisterUserServiceHandlerFromEndpoint(ctx, gmux, grpcServerAddr, opts)
		if err != nil {
			observability.Logger().Fatalf("failed to register user handler gateway: %s", err)

		}

		httpRouter := mux.NewRouter()
		httpRouter.Use(middlewares.Logger("/swagger", "/docs.swagger.json", "/grpcui"))
		httpRouter.Use(middlewares.Recover)

		conn, err := grpc.NewClient(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			observability.Logger().Fatalf("failed to connect to gRPC server: %s", err)
		}
		defer conn.Close()

		grpcuiHandler, err := standalone.HandlerViaReflection(ctx, conn, grpcServerAddr)
		if err != nil {
			observability.Logger().Fatalf("error creating grpcui handler: %s", err)
		}
		httpRouter.PathPrefix("/grpcui/").Handler(http.StripPrefix("/grpcui", grpcuiHandler))

		swaggerHandler := httpSwagger.Handler(httpSwagger.URL("/docs.swagger.json"))
		httpRouter.PathPrefix("/docs.swagger.json").Handler(http.FileServer(http.Dir("proto/gen/openapi")))
		httpRouter.PathPrefix("/swagger").Handler(swaggerHandler)

		httpRouter.PathPrefix("/").Handler(gmux)

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", a.config.HTTPPort),
			Handler: httpRouter,
		}
		go shutdownHTTP(ctx, server)

		observability.Logger().Printf("gRPC-Gateway server listening on HTTP port %d, proxying to gRPC port %d", a.config.HTTPPort, a.config.GRPCPort)
		return http.ListenAndServe(fmt.Sprintf(":%d", a.config.HTTPPort), httpRouter)
	}
}

func shutdownHTTP(ctx context.Context, server *http.Server) {
	<-ctx.Done()
	observability.Logger().Printf("Shutting down HTTP server gracefully...")
	server.Shutdown(ctx)
}

package test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/sergicanet9/go-hexagonal-api/api"
	"github.com/sergicanet9/go-hexagonal-api/config"
	infraMongo "github.com/sergicanet9/scv-go-framework/v2/infrastructure/mongo"
	infraPostgres "github.com/sergicanet9/scv-go-framework/v2/infrastructure/postgres"
)

const (
	contentType           = "application/json"
	mongoInternalPort     = "27017/tcp"
	mongoDBName           = "test-db"
	mongoConnectionEnv    = "mongoConnection"
	postgresUser          = "postgres"
	postgresPassword      = "test"
	postgresDBName        = "test-db"
	postgresInternalPort  = "5432/tcp"
	postgresConnectionEnv = "postgresConnection"
	jwtSecret             = "eaeBbXUxks"
	nonExpiryToken        = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXV0aG9yaXplZCI6dHJ1ZX0.cCKM32os5ROKxeE3IiDWoOyRew9T8puzPUKurPhrDug"
)

// TestMain does the setup before running the tests and the teardown afterwards
func TestMain(m *testing.M) {
	// Uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	mongoResource := setupMongo(pool)
	postgresResource := setupPostgres(pool)

	// Runs the tests
	code := m.Run()

	// When it´s done, kill and remove the containers
	if err = pool.Purge(mongoResource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}
	os.Unsetenv(mongoConnectionEnv)

	if err = pool.Purge(postgresResource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}
	os.Unsetenv(postgresConnectionEnv)

	os.Exit(code)
}

func Databases(t *testing.T, f func(*testing.T, string), databases ...string) {
	t.Helper()

	// if no databases specified, test is going to run on both
	if len(databases) == 0 {
		databases = []string{"mongo", "postgres"}
	}

	for _, db := range databases {
		t.Run(db, func(t *testing.T) {
			f(t, db)
		})
	}
}

// New starts a testing instance of the API and returns its config
func New(t *testing.T, database string) config.Config {
	t.Helper()

	cfg, err := testConfig(database)
	if err != nil {
		t.Fatal(err)
	}

	a := api.API{}
	a.Initialize(context.Background(), cfg)
	go func() {
		a.Run()
	}()

	return cfg
}

func setupMongo(pool *dockertest.Pool) *dockertest.Resource {
	// Pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "3.0",
		Env: []string{
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}
	connectionString := fmt.Sprintf("mongodb://localhost:%s", resource.GetPort(mongoInternalPort))
	os.Setenv(mongoConnectionEnv, connectionString)

	// Exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err := infraMongo.ConnectMongoDB(context.Background(), mongoDBName, connectionString)
		if err != nil {
			return err
		}

		return db.Client().Ping(context.Background(), nil)
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	return resource
}

func setupPostgres(pool *dockertest.Pool) *dockertest.Resource {
	// Pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.6",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", postgresUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", postgresPassword),
			fmt.Sprintf("POSTGRES_DB=%s", postgresDBName),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}
	connectionString := fmt.Sprintf("host=localhost port=%s dbname=%s user=%s password=%s sslmode=disable", resource.GetPort(postgresInternalPort), postgresDBName, postgresUser, postgresPassword)
	os.Setenv(postgresConnectionEnv, connectionString)

	// Exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	var db *sql.DB
	if err := pool.Retry(func() error {
		db, err = infraPostgres.ConnectPostgresDB(connectionString)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	// Migrates the database
	workingDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	migrationsDirectory := filepath.Join(filepath.Dir(workingDirectory), "adapters/repositories/postgres/migrations")

	goose.SetTableName("public.goose_db_version")
	if err := goose.Up(db, migrationsDirectory); err != nil {
		log.Fatalf("could not run the database migrations: %s", err)
	}

	return resource
}

func testConfig(database string) (c config.Config, err error) {
	c.Version = "Integration tests"
	c.Environment = "Integration tests"
	port, err := freePort()
	if err != nil {
		return c, err
	}
	c.Port = port
	c.Database = database

	c.Address = "http://localhost"
	c.MongoConnectionString = os.Getenv(mongoConnectionEnv)
	c.MongoDBName = mongoDBName
	c.PostgresConnectionString = os.Getenv(postgresConnectionEnv)
	c.JWTSecret = jwtSecret
	c.Timeout = config.Duration{Duration: 5 * time.Second}

	return c, nil
}

func freePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, errors.WithStack(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}

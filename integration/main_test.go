package integration

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sergicanet9/go-hexagonal-api/app/api"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/scv-go-tools/v3/api/utils"
	"github.com/sergicanet9/scv-go-tools/v3/test"
)

const (
	contentType           = "application/json"
	mongoContainerPort    = "27017/tcp"
	mongoDBName           = "test-db"
	mongoDSNEnv           = "mongoDSN"
	postgresUser          = "postgres"
	postgresPassword      = "test"
	postgresDBName        = "test-db"
	postgresContainerPort = "5432/tcp"
	postgresDSNEnv        = "postgresDSN"
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
	os.Unsetenv(mongoDSNEnv)

	if err = pool.Purge(postgresResource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}
	os.Unsetenv(postgresDSNEnv)

	os.Exit(code)
}

func setupMongo(pool *dockertest.Pool) *dockertest.Resource {
	// Pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
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
	dsn := fmt.Sprintf("mongodb://localhost:%s/%s", resource.GetPort(mongoContainerPort), mongoDBName)
	os.Setenv(mongoDSNEnv, dsn)

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
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", postgresUser, postgresPassword, resource.GetPort(postgresContainerPort), postgresDBName)
	os.Setenv(postgresDSNEnv, dsn)

	return resource
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

	cfg, err := testConfig(t, database)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	a := api.New(ctx, cfg)
	run := a.Run(ctx, cancel)
	go run()

	<-time.After(100 * time.Millisecond) // waiting time for letting the API start completely
	return cfg
}

func testConfig(t *testing.T, database string) (c config.Config, err error) {
	c.Version = "Integration tests"
	c.Environment = "Integration tests"
	c.Port = test.FreePort(t)
	c.Database = database
	switch database {
	case "mongo":
		c.DSN = os.Getenv(mongoDSNEnv)
	case "postgres":
		c.DSN = os.Getenv(postgresDSNEnv)
	default:
		return config.Config{}, fmt.Errorf("database flag %s not valid", database)
	}

	c.Address = "http://localhost"
	c.PostgresMigrationsDir = "db/postgres/migrations"
	c.JWTSecret = jwtSecret
	c.Timeout = utils.Duration{Duration: 5 * time.Second}

	return c, nil
}

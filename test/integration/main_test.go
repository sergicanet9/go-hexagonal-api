package integration

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sergicanet9/go-hexagonal-api/app/api"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/scv-go-tools/v3/api/utils"
	"github.com/sergicanet9/scv-go-tools/v3/infrastructure"
	"github.com/sergicanet9/scv-go-tools/v3/testutils"
)

const (
	contentType           = "application/json"
	mongoDBName           = "test-db"
	mongoUser             = "mongo"
	mongoPassword         = "test"
	mongoContainerPort    = "27017/tcp"
	mongoDSNEnv           = "mongoDSN"
	postgresDBName        = "test-db"
	postgresUser          = "postgres"
	postgresPassword      = "test"
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
		log.Panicf("could not connect to docker: %s", err)
	}

	mongoResource := setupMongo(pool)
	postgresResource := setupPostgres(pool)

	// Runs the tests
	code := m.Run()

	// When it´s done, kill and remove the containers
	if err = pool.Purge(mongoResource); err != nil {
		log.Panicf("could not purge resource: %s", err)
	}
	os.Unsetenv(mongoDSNEnv)

	if err = pool.Purge(postgresResource); err != nil {
		log.Panicf("could not purge resource: %s", err)
	}
	os.Unsetenv(postgresDSNEnv)

	os.Exit(code)
}

func setupMongo(pool *dockertest.Pool) *dockertest.Resource {
	// creates filekey
	_, filePath, _, _ := runtime.Caller(0)
	fileKey, err := os.CreateTemp(path.Dir(filePath), "")
	if err != nil {
		log.Panic(err)
	}
	defer os.Remove(fileKey.Name())

	bytes := []byte(`secret123`)
	err = os.WriteFile(fileKey.Name(), bytes, 0644)
	if err != nil {
		log.Panic(err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
		Env: []string{
			fmt.Sprintf("MONGO_INITDB_DATABASE=%s", mongoDBName),
			fmt.Sprintf("MONGO_INITDB_ROOT_USERNAME=%s", mongoUser),
			fmt.Sprintf("MONGO_INITDB_ROOT_PASSWORD=%s", mongoPassword),
		},
		Mounts: []string{fmt.Sprintf("%s:/auth/file.key", fileKey.Name())},
		Cmd:    []string{"--keyFile", "/auth/file.key", "--replSet", "rs0", "--bind_ip_all"},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Panicf("could not start resource: %s", err)
	}

	exitCode, err := resource.Exec([]string{"/bin/sh", "-c", "chown mongodb:mongodb /auth/file.key"}, dockertest.ExecOptions{})
	if err != nil {
		log.Panicf("failure executing command in the resource: %s", err)
	}
	if exitCode != 0 {
		log.Panicf("failure executing command in the resource, exit code was %d", exitCode)
	}

	dsn := fmt.Sprintf("mongodb://%s:%s@localhost:%s/%s?authSource=admin&connect=direct", mongoUser, mongoPassword, resource.GetPort(mongoContainerPort), mongoDBName)
	os.Setenv(mongoDSNEnv, dsn)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 10 * time.Second
	err = pool.Retry(func() error {
		_, err = infrastructure.ConnectMongoDB(context.Background(), dsn)
		return err
	})
	if err != nil {
		log.Panicf("Could not connect to docker: %s", err)
	}

	exitCode, err = resource.Exec([]string{"/bin/sh", "-c", fmt.Sprintf("echo 'rs.initiate().ok' | mongosh -u %s -p %s --quiet", mongoUser, mongoPassword)}, dockertest.ExecOptions{})
	if err != nil {
		log.Panicf("failure executing command in the resource: %s", err)
	}
	if exitCode != 0 {
		log.Panicf("failure executing command in the resource, exit code was %d", exitCode)
	}

	return resource
}

func setupPostgres(pool *dockertest.Pool) *dockertest.Resource {
	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.6",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", postgresUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", postgresPassword),
			fmt.Sprintf("POSTGRES_DB=%s", postgresDBName),
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Panicf("could not start resource: %s", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", postgresUser, postgresPassword, resource.GetPort(postgresContainerPort), postgresDBName)
	os.Setenv(postgresDSNEnv, dsn)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 10 * time.Second
	err = pool.Retry(func() error {
		_, err = infrastructure.ConnectPostgresDB(context.Background(), dsn)
		return err
	})
	if err != nil {
		log.Panicf("Could not connect to docker: %s", err)
	}

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
	nrApp, _ := newrelic.NewApplication(
		newrelic.ConfigAppName("go-new-relic-test"),
		newrelic.ConfigLicense("0123456789012345678901234567890123456789"),
	)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	a := api.New(ctx, cfg, nrApp)
	run := a.Run(ctx, cancel)
	go run()

	<-time.After(100 * time.Millisecond) // waiting time for letting the API start completely
	return cfg
}

func testConfig(t *testing.T, database string) (c config.Config, err error) {
	c.Version = "Integration tests"
	c.Environment = "Integration tests"
	c.HTTPPort = testutils.FreePort(t)
	c.Database = database
	switch database {
	case "mongo":
		c.DSN = os.Getenv(mongoDSNEnv)
	case "postgres":
		c.DSN = os.Getenv(postgresDSNEnv)
	default:
		return config.Config{}, fmt.Errorf("database flag %s not valid", database)
	}

	c.PostgresMigrationsDir = "infrastructure/postgres/migrations"
	c.JWTSecret = jwtSecret
	c.Timeout = utils.Duration{Duration: 30 * time.Second}

	return c, nil
}

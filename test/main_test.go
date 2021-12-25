package test

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/pkg/errors"
	"github.com/sergicanet9/go-mongo-restapi/api"
	"github.com/sergicanet9/go-mongo-restapi/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	contentType       = "application/json"
	mongoInternalPort = "27017/tcp"
	mongoPortEnv      = "mongoPort"
)

// TestMain does the setup before running the tests and the teardown afterwards
func TestMain(m *testing.M) {
	// Uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	// Pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mongo", "3.0", nil)
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}
	os.Setenv(mongoPortEnv, resource.GetPort(mongoInternalPort))

	// Exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://localhost:%s", os.Getenv(mongoPortEnv))))
		if err != nil {
			return err
		}

		return client.Ping(context.Background(), nil)
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	// Runs the tests
	code := m.Run()

	// When it´s done, kill and remove the container
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Unsetenv(mongoPortEnv)
	os.Exit(code)
}

// New starts a testing instance of the API and returns its config
func New(t *testing.T) config.Config {
	t.Helper()

	c, err := testConfig()
	if err != nil {
		t.Fatal(err)
	}

	a := api.API{}
	a.Initialize(c)
	go func() {
		a.Run()
	}()

	return c
}

func testConfig() (c config.Config, err error) {
	c.Env = "Integration tests"

	port, err := freePort()
	if err != nil {
		return c, err
	}
	c.Port = port
	c.Address = "http://localhost"
	c.DBConnectionString = fmt.Sprintf("mongodb://localhost:%s", os.Getenv(mongoPortEnv))
	c.DBName = "test-db"
	c.JWTSecret = "eaeBbXUxks"

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

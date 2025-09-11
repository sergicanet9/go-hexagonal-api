package v1

import (
	"context"
	"testing"

	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/stretchr/testify/assert"
)

// TestHealthCheck_WhenEnvironmentIsNotLocal_DSNFiltered checks that healthCheck handler does return the expected response filtering the DSN
func TestHealthCheck_WhenEnvironmentIsNotLocal_DSNIsFiltered(t *testing.T) {
	// Arrange
	cfg := config.Config{
		Version:     "test-version",
		Environment: "test-environment",
		Database:    "test-database",
		HTTPPort:    1,
		GRPCPort:    2,
		DSN:         "test-dsn",
	}
	healthHandler := NewHealthHandler(context.Background(), cfg)
	expectedDSN := "***FILTERED***"

	// Act
	resp, err := healthHandler.HealthCheck(context.Background(), nil)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, cfg.Version, resp.Version)
	assert.Equal(t, cfg.Environment, resp.Environment)
	assert.Equal(t, cfg.Database, resp.Database)
	assert.Equal(t, cfg.HTTPPort, int(resp.HttpPort))
	assert.Equal(t, cfg.GRPCPort, int(resp.GrpcPort))
	assert.Equal(t, expectedDSN, resp.Dsn)
}

// TestHealthCheck_WhenEnvironmentIsLocal_DSNIsNotFiltered checks that healthCheck handler does return the expected response without filtering the DSN
func TestHealthCheck_WhenEnvironmentIsLocal_DSNIsNotFiltered(t *testing.T) {
	// Arrange
	cfg := config.Config{
		Version:     "test-version",
		Environment: "local",
		Database:    "test-database",
		HTTPPort:    1,
		GRPCPort:    2,
		DSN:         "test-dsn",
	}
	healthHandler := NewHealthHandler(context.Background(), cfg)

	// Act
	resp, err := healthHandler.HealthCheck(context.Background(), nil)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, cfg.Version, resp.Version)
	assert.Equal(t, cfg.Environment, resp.Environment)
	assert.Equal(t, cfg.Database, resp.Database)
	assert.Equal(t, cfg.HTTPPort, int(resp.HttpPort))
	assert.Equal(t, cfg.GRPCPort, int(resp.GrpcPort))
	assert.Equal(t, cfg.DSN, resp.Dsn)
}

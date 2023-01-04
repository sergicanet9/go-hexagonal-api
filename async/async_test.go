package async

import (
	"context"
	"testing"
	"time"

	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/stretchr/testify/assert"
)

// TestNew_Ok checks that New creates a new async struct with the expected values
func TestNew_Ok(t *testing.T) {
	// Arrange
	expectedConfig := config.Config{}
	// Act
	async := New(expectedConfig)

	// Assert
	assert.Equal(t, expectedConfig, async.config)
}

// TestRun_ContextCancelled checks that Run finishes and returns the expected error when the context gets cancelled.
func TestRun_ContextCancelled(t *testing.T) {
	// Arrange
	async := &async{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	expectedError := "async process stopped"

	// Act
	errFunc := async.Run(ctx, cancel)

	// Assert
	assert.Equal(t, expectedError, errFunc().Error())
}

// TestHealthCheck_SuccessfulURLThenContextCancelled checks that healthCheck finishes and returns the expected error when there is a timeout while the given URL was returning success.
func TestHealthCheck_SuccessfulURLThenContextCancelled(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	successURL := "http://www.google.com"
	expectedError := context.DeadlineExceeded.Error()

	// Act
	healthCheck(ctx, cancel, successURL, time.Second)

	// Assert
	assert.Equal(t, expectedError, ctx.Err().Error())
}

// TestHealthCheck_NotFoundURLThenContextCancelled checks that healthCheck finishes and returns the expected error when there is a timeout while the given URL was returning not found.
func TestHealthCheck_NotFoundURLThenContextCancelled(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	notFoundURL := "http://testing.com/random-url-to-be-not-found"
	expectedError := context.DeadlineExceeded.Error()

	// Act
	healthCheck(ctx, cancel, notFoundURL, time.Second)

	// Assert
	assert.Equal(t, expectedError, ctx.Err().Error())
}

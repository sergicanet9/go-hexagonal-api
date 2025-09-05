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

// TestRun_ContextCancelled checks that Run finishes when the context gets cancelled
func TestRun_ContextCancelled(t *testing.T) {
	// Arrange
	async := &async{
		config: config.Config{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)

	// Act
	errFunc := async.Run(ctx, cancel)

	// Assert
	assert.Nil(t, errFunc())
}

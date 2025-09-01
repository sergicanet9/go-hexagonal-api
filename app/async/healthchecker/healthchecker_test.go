package healthchecker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestRunHTTP_ContextCancelled checks that the HTTP healthchecker runs until the context gets cancelled
func TestRunHTTP_ContextCancelled(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	url := "http://www.google.com"
	expectedError := context.DeadlineExceeded.Error()

	// Act
	RunHTTP(ctx, cancel, url, time.Second)

	// Assert
	assert.Equal(t, expectedError, ctx.Err().Error())
}

// TestRunGRPC_ContextCancelled checks that the gRPC healthchecker runs until the context gets cancelled
func TestRunGRPC_ContextCancelled(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	target := "dns:///www.google.com:443"
	expectedError := context.DeadlineExceeded.Error()

	// Act
	RunGRPC(ctx, cancel, target, time.Second)

	// Assert
	assert.Equal(t, expectedError, ctx.Err().Error())
}

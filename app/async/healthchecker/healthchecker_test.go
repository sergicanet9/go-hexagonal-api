package healthchecker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestRun_SuccessfulURLThenContextCancelled checks that healthCheck finishes and returns the expected error when the context gets cancelled
// while the given URL was returning success
func TestRun_SuccessfulURLThenContextCancelled(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	successURL := "http://www.google.com"
	expectedError := context.DeadlineExceeded.Error()

	// Act
	Run(ctx, cancel, successURL, time.Second)

	// Assert
	assert.Equal(t, expectedError, ctx.Err().Error())
}

// TestRun_InvalidURLThenContextCancelled checks that healthCheck finishes and returns the expected error when the context gets cancelled
// while the given URL was invalid
func TestRun_InvalidURLThenContextCancelled(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	notFoundURL := "http://€@$testing.com"
	expectedError := context.DeadlineExceeded.Error()

	// Act
	Run(ctx, cancel, notFoundURL, time.Second)

	// Assert
	assert.Equal(t, expectedError, ctx.Err().Error())
}

// TestRun_NotFoundURLThenContextCancelled checks that healthCheck finishes and returns the expected error when the context gets cancelled
// while the given URL was returning not found
func TestRun_NotFoundURLThenContextCancelled(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	notFoundURL := "http://testing.com/random-url-to-be-not-found"
	expectedError := context.DeadlineExceeded.Error()

	// Act
	Run(ctx, cancel, notFoundURL, time.Second)

	// Assert
	assert.Equal(t, expectedError, ctx.Err().Error())
}

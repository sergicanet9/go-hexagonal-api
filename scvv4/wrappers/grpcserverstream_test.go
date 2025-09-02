package wrappers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestContext_Ok that the GRPCServerStream correctly returns the context it was initialized with
func TestContext_Ok(t *testing.T) {
	// Arrange
	expectedCtx := context.WithValue(context.Background(), "key", "value")
	ss := NewGRPCServerStream(expectedCtx)

	// Act
	actualCtx := ss.Context()

	// Assert
	assert.Equal(t, expectedCtx, actualCtx)
	assert.Equal(t, "value", actualCtx.Value("key"))
}

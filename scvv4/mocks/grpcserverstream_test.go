package mocks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestContext_Ok that the MockGRPCServerStream correctly returns the context it was initialized with
func TestContext_Ok(t *testing.T) {
	// Arrange
	expectedCtx := context.WithValue(context.Background(), "key", "value")
	ss := &MockGRPCServerStream{
		Ctx: expectedCtx,
	}

	// Act
	actualCtx := ss.Context()

	// Assert
	assert.Equal(t, expectedCtx, actualCtx)
	assert.Equal(t, "value", actualCtx.Value("key"))
}

package observability

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLogger_ReturnsSingleton checks that the same logger is always returned
func TestLogger_ReturnsSingleton(t *testing.T) {
	// Act
	logger1 := Logger()
	logger2 := Logger()

	// Assert
	assert.NotNil(t, logger1)
	assert.NotNil(t, logger2)
	assert.Equal(t, logger1, logger2)
}

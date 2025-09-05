package observability

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSetupNewRelic_InvalidKey checks that an error is returned when an invalid New Relic key is provided
func TestSetupNewRelic_InvalidKey(t *testing.T) {
	// Arrange
	appName := "test-app"
	invalidKey := "invalid-key"

	// Act
	_, err := SetupNewRelic(appName, invalidKey)

	// Assert
	assert.NotEmpty(t, err)
}

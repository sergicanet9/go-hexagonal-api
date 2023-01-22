package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsValid_Ok checks that String returns the expected response when everything goes as expected
func TestString_Ok(t *testing.T) {
	// Arrange
	var adminClaim UserClaim
	expectedClaim := "admin"
	adminClaim = 0

	// Act
	claim := adminClaim.String()

	// Assert
	assert.Equal(t, expectedClaim, claim)
}

// TestIsValid_Ok checks that IsValid returns the expected response when everything goes as expected
func TestIsValid_Ok(t *testing.T) {
	// Arrange
	var adminClaim UserClaim

	// Act
	isValid := adminClaim.IsValid()

	// Assert
	assert.True(t, isValid)
}

// TestGetUserClaims_Ok checks that GetUserClaims returns the expected response when everything goes as expected
func TestGetUserClaims_Ok(t *testing.T) {
	// Arrange
	expectedClaims := map[int]string{0: "admin"}

	// Act
	resp := GetUserClaims()

	// Assert
	assert.Equal(t, expectedClaims, resp)
}

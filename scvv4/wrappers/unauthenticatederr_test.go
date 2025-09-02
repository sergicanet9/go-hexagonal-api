package wrappers

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewUnauthenticatedErr_Ok checks that NewUnauthenticatedErr returns the expected error type when receives an error
func TestNewUnauthenticatedErr_Ok(t *testing.T) {
	// Arrange
	err := fmt.Errorf("test error")

	// Act
	gotErr := NewUnauthenticatedErr(err)

	// Assert
	assert.NotEmpty(t, gotErr)
	assert.IsType(t, UnauthenticatedErr, gotErr)
}

// TestNewUnauthenticatedErr_NilErr checks that NewUnauthenticatedErr returns nil when receives a nil error
func TestNewUnauthenticatedErr_NilErr(t *testing.T) {
	// Arrange
	var err error

	// Act
	gotErr := NewUnauthenticatedErr(err)

	// Assert
	assert.Nil(t, gotErr)
}

// TestUnauthenticatedErrError_Ok checks that Error returns the expected error message of the receiver
func TestUnauthenticatedErrError_Ok(t *testing.T) {
	// Arrange
	expectedMsg := "test error"
	err := NewUnauthenticatedErr(errors.New(expectedMsg))

	// Act
	gotMsg := err.Error()

	// Assert
	assert.Equal(t, expectedMsg, gotMsg)
}

// TestUnauthenticatedErrIs_True checks that Is returns true when the receiver is an unauthenticatedError
func TestUnauthenticatedErrIs_True(t *testing.T) {
	// Arrange
	err := NewUnauthenticatedErr(fmt.Errorf("test"))

	// Act
	isUnauthenticatedErr := errors.Is(err, UnauthenticatedErr)

	// Assert
	assert.True(t, isUnauthenticatedErr)
}

// TestUnauthenticatedErrIs_False checks that Is returns false when the receiver is not an unauthenticatedError
func TestUnauthenticatedErrIs_False(t *testing.T) {
	// Arrange
	err := fmt.Errorf("test")

	// Act
	isUnauthenticatedErr := errors.Is(err, UnauthenticatedErr)

	// Assert
	assert.False(t, isUnauthenticatedErr)
}

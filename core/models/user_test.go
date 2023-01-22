package models

import (
	"testing"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
	"github.com/stretchr/testify/assert"
)

// TestValidateCreateUserReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateCreateUserReq_Ok(t *testing.T) {
	// Arrange
	req := CreateUserReq{
		Email:        "test@test.com",
		PasswordHash: "test",
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateCreateUserReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateCreateUserReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := CreateUserReq{}
	expectedError := "email cannot be empty | password cannot be empty"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestValidateLoginUserReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateLoginUserReq_Ok(t *testing.T) {
	// Arrange
	req := LoginUserReq{
		Email:    "test@test.com",
		Password: "test",
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateLoginUserReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateLoginUserReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := LoginUserReq{}
	expectedError := "email cannot be empty | password cannot be empty"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

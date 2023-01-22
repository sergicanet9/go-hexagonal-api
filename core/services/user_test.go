package services

import (
	"context"
	"testing"

	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/entities"
	"github.com/sergicanet9/go-hexagonal-api/core/models"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	"github.com/sergicanet9/go-hexagonal-api/test/mocks"
	"github.com/sergicanet9/scv-go-tools/v3/testutils"
	"github.com/stretchr/testify/assert"
)

// TestNewUserService_Ok checks that NewUserService creates a new userService struct
func TestNewUserService_Ok(t *testing.T) {
	// Arrange
	cfg := config.Config{}
	userRepositoryMock := mocks.NewUserRepository(t)

	// Act
	service := NewUserService(cfg, userRepositoryMock)

	// Assert
	assert.NotEmpty(t, service)
}

// TestLogin_Ok checks that Login returns the expected response when a valid request is provided
func TestLogin_Ok(t *testing.T) {
	// Arrange
	req := models.LoginUserReq{
		Email:    "test@test.com",
		Password: "test",
	}

	filter := map[string]interface{}{"email": req.Email}
	var result []interface{}
	expectedUser := entities.User{
		Email:        req.Email,
		PasswordHash: "$2a$10$NexA3QvmeUMPME6GVhFaX.C4A.y2VIPBwRNrV0c2DncjCAWSBnINK",
	}
	result = append(result, &expectedUser)

	var nilPointer *int
	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.Get), context.Background(), filter, nilPointer, nilPointer).Return(result, nil).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	resp, err := service.Login(context.Background(), req)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, models.UserResp(expectedUser), resp.User)
}

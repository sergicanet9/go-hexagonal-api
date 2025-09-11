package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/entities"
	"github.com/sergicanet9/go-hexagonal-api/core/models"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	"github.com/sergicanet9/go-hexagonal-api/test/mocks"
	"github.com/sergicanet9/scv-go-tools/v4/testutils"
	"github.com/sergicanet9/scv-go-tools/v4/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

// TestLogin_Ok checks that Login returns the expected response when a valid request is received
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
		ClaimIDs:     []int32{0},
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
	assert.Equal(t, models.GetUserResp(expectedUser), resp.User)
}

// TestLogin_NotFound checks that Login returns an error when the user is not found
func TestLogin_NotFound(t *testing.T) {
	// Arrange
	req := models.LoginUserReq{
		Email:    "test@test.com",
		Password: "test",
	}

	filter := map[string]interface{}{"email": req.Email}
	expectedError := fmt.Sprintf("email %s not found", req.Email)

	var nilPointer *int
	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.Get), context.Background(), filter, nilPointer, nilPointer).Return(nil, wrappers.NonExistentErr).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	_, err := service.Login(context.Background(), req)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.NonExistentErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestLogin_InvalidRequest checks that Login returns an error when the received request is not valid
func TestLogin_InvalidRequest(t *testing.T) {
	// Arrange
	req := models.LoginUserReq{
		Email:    "",
		Password: "",
	}

	expectedError := "email cannot be empty | password cannot be empty"

	service := &userService{
		config:     config.Config{},
		repository: nil,
	}

	// Act
	_, err := service.Login(context.Background(), req)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestLogin_IncorrectPassword checks that Login returns an error when the received password does not match with the one returned from the repository
func TestLogin_IncorrectPassword(t *testing.T) {
	// Arrange
	req := models.LoginUserReq{
		Email:    "test@test.com",
		Password: "incorrect-password",
	}

	filter := map[string]interface{}{"email": req.Email}
	var result []interface{}
	expectedUser := entities.User{
		Email:        req.Email,
		PasswordHash: "$2a$10$NexA3QvmeUMPME6GVhFaX.C4A.y2VIPBwRNrV0c2DncjCAWSBnINK",
	}
	result = append(result, &expectedUser)

	expectedError := "password incorrect"

	var nilPointer *int
	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.Get), context.Background(), filter, nilPointer, nilPointer).Return(result, nil).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	_, err := service.Login(context.Background(), req)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestLogin_InvalidClaims checks that Login returns an error when the claims returned from the repository are not valid
func TestLogin_InvalidClaims(t *testing.T) {
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
		ClaimIDs:     []int32{3},
	}
	result = append(result, &expectedUser)

	expectedError := "claim 3 is not valid"

	var nilPointer *int
	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.Get), context.Background(), filter, nilPointer, nilPointer).Return(result, nil).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	_, err := service.Login(context.Background(), req)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestCreate_Ok checks that Create returns the expected response when a valid request is received
func TestCreate_Ok(t *testing.T) {
	// Arrange
	req := models.CreateUserReq{
		Email:    "test@test.com",
		Password: "test",
	}

	expectedResponse := models.CreateUserResp{
		ID: "new-id",
	}

	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.Create), mock.Anything, mock.AnythingOfType("entities.User")).Return(expectedResponse.ID, nil).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	resp, err := service.Create(context.Background(), req)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, resp)
}

// TestCreate_CreateError checks that Create returns an error when the Create function from the repository fails
func TestCreate_CreateError(t *testing.T) {
	// Arrange
	req := models.CreateUserReq{
		Email:    "test@test.com",
		Password: "test",
	}

	expectedError := "repository-error"

	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.Create), mock.Anything, mock.AnythingOfType("entities.User")).Return("", errors.New(expectedError)).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	_, err := service.Create(context.Background(), req)

	// Assert
	assert.NotEmpty(t, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestCreate_InvalidRequest checks that Create returns an error when the received request is not valid
func TestCreate_InvalidRequest(t *testing.T) {
	// Arrange
	req := models.CreateUserReq{
		Email:    "",
		Password: "",
	}

	expectedError := "email cannot be empty | password cannot be empty"

	service := &userService{
		config:     config.Config{},
		repository: nil,
	}

	// Act
	_, err := service.Create(context.Background(), req)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestCreate_InvalidClaims checks that Create returns an error when the received claims in the request are not valid
func TestCreate_InvalidClaims(t *testing.T) {
	// Arrange
	req := models.CreateUserReq{
		Email:    "test@test.com",
		Password: "test",
		ClaimIDs: []int32{3},
	}

	expectedError := "claim 3 is not valid"

	service := &userService{
		config:     config.Config{},
		repository: nil,
	}

	// Act
	_, err := service.Create(context.Background(), req)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateMany_Ok checks that CreateMany does not return an error when a valid request is received
func TestCreateMany_Ok(t *testing.T) {
	// Arrange
	req := []models.CreateUserReq{
		{
			Email:    "test@test.com",
			Password: "test",
		},
	}

	expectedResponse := models.CreateManyUserResp{
		IDs: []string{"new-id"},
	}

	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.CreateMany), mock.Anything, mock.AnythingOfType("[]interface {}")).Return(expectedResponse.IDs, nil).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	resp, err := service.CreateMany(context.Background(), req)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, resp)
}

// TestCreateMany_CreateManyError checks that CreateMany returns an error when the CreateMany function from the repository fails
func TestCreateMany_CreateManyError(t *testing.T) {
	// Arrange
	req := []models.CreateUserReq{
		{
			Email:    "test@test.com",
			Password: "test",
		},
	}

	expectedError := "repository-error"

	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.CreateMany), mock.Anything, mock.AnythingOfType("[]interface {}")).Return([]string{}, errors.New(expectedError)).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	_, err := service.CreateMany(context.Background(), req)

	// Assert
	assert.NotEmpty(t, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateMany_InvalidRequest checks that CreateMany returns an error when one of the users in the received request is not valid
func TestCreateMany_InvalidRequest(t *testing.T) {
	// Arrange
	req := []models.CreateUserReq{
		{
			Email:    "",
			Password: "",
		},
	}

	expectedError := "email cannot be empty | password cannot be empty"

	service := &userService{
		config:     config.Config{},
		repository: nil,
	}

	// Act
	_, err := service.CreateMany(context.Background(), req)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateMany_InvalidClaims checks that CreateMany returns an error when the received claims of one of the users in the request are not valid
func TestCreateMany_InvalidClaims(t *testing.T) {
	// Arrange
	req := []models.CreateUserReq{
		{
			Email:    "test@test.com",
			Password: "test",
			ClaimIDs: []int32{3},
		},
	}

	expectedError := "claim 3 is not valid"

	service := &userService{
		config:     config.Config{},
		repository: nil,
	}

	// Act
	_, err := service.CreateMany(context.Background(), req)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestGetAll_Ok checks that GetAll returns the expected response when everything goes as expected
func TestGetAll_Ok(t *testing.T) {
	// Arrange
	var result []interface{}
	expectedUser := entities.User{
		Email:        "test@test.com",
		PasswordHash: "$2a$10$NexA3QvmeUMPME6GVhFaX.C4A.y2VIPBwRNrV0c2DncjCAWSBnINK",
		ClaimIDs:     []int32{0},
	}
	result = append(result, &expectedUser)

	var nilPointer *int
	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.Get), mock.Anything, map[string]interface{}{}, nilPointer, nilPointer).Return(result, nil).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	resp, err := service.GetAll(context.Background())

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, models.GetUserResp(expectedUser), resp[0])
}

// TestGetAll_NoResourcesFound checks that GetAll does not return an error when the repository does not return an user
func TestGetAll_NoResourcesFound(t *testing.T) {
	// Arrange
	var nilPointer *int
	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.Get), mock.Anything, map[string]interface{}{}, nilPointer, nilPointer).Return(nil, wrappers.NonExistentErr).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	resp, err := service.GetAll(context.Background())

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 0, len(resp))
}

// TestGetByID_Ok checks that GetByID returns the expected response when a valid ID is received
func TestGetByID_Ok(t *testing.T) {
	// Arrange
	expectedUser := entities.User{
		ID:           "test-id",
		Email:        "test@test.com",
		PasswordHash: "$2a$10$NexA3QvmeUMPME6GVhFaX.C4A.y2VIPBwRNrV0c2DncjCAWSBnINK",
		ClaimIDs:     []int32{0},
	}

	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.GetByID), context.Background(), expectedUser.ID).Return(&expectedUser, nil).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	resp, err := service.GetByID(context.Background(), expectedUser.ID)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, models.GetUserResp(expectedUser), resp)
}

// TestGetByID_Ok checks that GetByID returns tan error when the provided ID does not exist
func TestGetByID_NotFound(t *testing.T) {
	// Arrange
	nonExistentID := "non-existent-id"
	expectedError := fmt.Sprintf("ID %s not found", nonExistentID)

	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.GetByID), context.Background(), nonExistentID).Return(nil, wrappers.NonExistentErr).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	_, err := service.GetByID(context.Background(), nonExistentID)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.NonExistentErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestUpdate_Ok checks that Update does not return an error when everything goes as expected

func TestUpdate_Ok(t *testing.T) {
	// Arrange
	testParam := "test"
	testEmail := "test@test.com"
	testClaims := []int32{0}
	id := "test-id"

	req := models.UpdateUserReq{
		Name:        &testParam,
		Surnames:    &testParam,
		Email:       &testEmail,
		NewPassword: &testParam,
		OldPassword: &testParam,
		ClaimIDs:    &testClaims,
	}

	existingUser := entities.User{
		PasswordHash: "$2a$10$NexA3QvmeUMPME6GVhFaX.C4A.y2VIPBwRNrV0c2DncjCAWSBnINK",
	}

	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.GetByID), context.Background(), id).Return(&existingUser, nil).Once()
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.Update), context.Background(), id, mock.AnythingOfType("entities.User")).Return(nil).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	err := service.Update(context.Background(), id, req)

	// Assert
	assert.Nil(t, err)
}

// TestUpdate_NotFound checks that Update returns an error when the provided ID does not exist
func TestUpdate_NotFound(t *testing.T) {
	// Arrange
	nonExistentID := "non-existent-id"
	expectedError := fmt.Sprintf("ID %s not found", nonExistentID)

	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.GetByID), context.Background(), nonExistentID).Return(nil, wrappers.NonExistentErr).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	err := service.Update(context.Background(), nonExistentID, models.UpdateUserReq{})

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.NonExistentErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestUpdate_IncorrectPassword checks that Update returns an error when the provided old password is not correct
func TestUpdate_IncorrectPassword(t *testing.T) {
	// Arrange
	newPassword := "new-password"
	incorrectOldPassword := "incorrect-password"
	id := "test-id"

	req := models.UpdateUserReq{
		NewPassword: &newPassword,
		OldPassword: &incorrectOldPassword,
	}

	existingUser := entities.User{
		PasswordHash: "$2a$10$NexA3QvmeUMPME6GVhFaX.C4A.y2VIPBwRNrV0c2DncjCAWSBnINK",
	}

	expectedError := "password incorrect"

	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.GetByID), context.Background(), id).Return(&existingUser, nil).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	err := service.Update(context.Background(), id, req)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestUpdate_InvalidClaims checks that Update returns an error when the provided new claims are not valid
func TestUpdate_InvalidClaims(t *testing.T) {
	// Arrange
	invalidClaims := []int32{3}
	id := "test-id"

	req := models.UpdateUserReq{
		ClaimIDs: &invalidClaims,
	}

	expectedError := "claim 3 is not valid"

	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.GetByID), context.Background(), id).Return(&entities.User{}, nil).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	err := service.Update(context.Background(), id, req)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestDelete_Ok checks that Delete does not return an error when everything goes as expected
func TestDelete_Ok(t *testing.T) {
	// Arrange
	testID := "test-id"
	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.Delete), context.Background(), testID).Return(nil).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	err := service.Delete(context.Background(), testID)

	// Assert
	assert.Nil(t, err)
}

// TestDelete_NotFound checks that Delete returns an error when the provided ID does not exist
func TestDelete_NotFound(t *testing.T) {
	// Arrange
	nonExistentID := "non-existent-id"
	expectedError := fmt.Sprintf("ID %s not found", nonExistentID)

	userRepositoryMock := mocks.NewUserRepository(t)
	userRepositoryMock.On(testutils.FunctionName(t, ports.UserRepository.Delete), context.Background(), nonExistentID).Return(wrappers.NonExistentErr).Once()

	service := &userService{
		config:     config.Config{},
		repository: userRepositoryMock,
	}

	// Act
	err := service.Delete(context.Background(), nonExistentID)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.NonExistentErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestGetUserClaims_Ok checks that GetUserClaims returns the expected response when everything goes as expected
func TestGetUserClaims_Ok(t *testing.T) {
	// Arrange
	service := &userService{
		config:     config.Config{},
		repository: nil,
	}

	expectedClaims := map[int]string{0: "admin"}

	// Act
	resp := service.GetUserClaims(context.Background())

	// Assert
	assert.Equal(t, expectedClaims, resp)
}

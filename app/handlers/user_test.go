package handlers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/models"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	"github.com/sergicanet9/go-hexagonal-api/proto/gen/go/pb"
	"github.com/sergicanet9/go-hexagonal-api/test/mocks"
	"github.com/sergicanet9/scv-go-tools/v3/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TestLoginUser_Ok checks that the Login handler returns the expected response on a valid request
func TestLoginUser_Ok(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	expectedUser := models.GetUserResp{
		ID:        "test-id",
		Name:      "Test",
		Surnames:  "User",
		Email:     "test@test.com",
		ClaimIDs:  []int32{0},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	expectedResp := models.LoginUserResp{
		User:  expectedUser,
		Token: "test-token",
	}

	userService.On(testutils.FunctionName(t, ports.UserService.Login), mock.Anything, mock.AnythingOfType("models.LoginUserReq")).Return(expectedResp, nil).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	req := &pb.LoginUserRequest{
		Email:    "test@test.com",
		Password: "test",
	}

	// Act
	resp, err := handler.Login(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "test-token", resp.Token)
	assert.Equal(t, expectedUser.ID, resp.User.Id)
	assert.Equal(t, expectedUser.Name, resp.User.Name)
	assert.Equal(t, expectedUser.Surnames, resp.User.Surnames)
	assert.Equal(t, expectedUser.Email, resp.User.Email)
	assert.Equal(t, expectedUser.ClaimIDs, resp.User.ClaimIds)
	assert.True(t, expectedUser.CreatedAt.Equal(resp.User.CreatedAt.AsTime()))
	assert.True(t, expectedUser.UpdatedAt.Equal(resp.User.UpdatedAt.AsTime()))
}

// TestLoginUser_ServiceError checks that the Login handler returns a gRPC error when the service fails
func TestLoginUser_ServiceError(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	expectedError := "service-error"
	userService.On(testutils.FunctionName(t, ports.UserService.Login), mock.Anything, mock.AnythingOfType("models.LoginUserReq")).Return(models.LoginUserResp{}, errors.New(expectedError)).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	req := &pb.LoginUserRequest{
		Email:    "test@test.com",
		Password: "test",
	}

	// Act
	_, err := handler.Login(context.Background(), req)

	// Assert
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, expectedError, st.Message())
}

// TestCreateUser_Ok checks that the Create handler returns the expected response on a valid request
func TestCreateUser_Ok(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	expectedResp := models.CreateUserResp{
		ID: "new-id",
	}
	userService.On(testutils.FunctionName(t, ports.UserService.Create), mock.Anything, mock.AnythingOfType("models.CreateUserReq")).Return(expectedResp, nil).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	req := &pb.CreateUserRequest{
		Email:    "test@test.com",
		Password: "test",
	}

	// Act
	resp, err := handler.Create(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResp.ID, resp.Id)
}

// TestCreateUser_ServiceError checks that the Create handler returns a gRPC error when the service fails
func TestCreateUser_ServiceError(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	expectedError := "service-error"
	userService.On(testutils.FunctionName(t, ports.UserService.Create), mock.Anything, mock.AnythingOfType("models.CreateUserReq")).Return(models.CreateUserResp{}, errors.New(expectedError)).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	req := &pb.CreateUserRequest{
		Email:    "test@test.com",
		Password: "test",
	}

	// Act
	_, err := handler.Create(context.Background(), req)

	// Assert
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, expectedError, st.Message())
}

// TestCreateManyUsers_Ok checks that the CreateMany handler returns the expected response on a valid request
func TestCreateManyUsers_Ok(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	expectedResp := models.CreateManyUserResp{
		IDs: []string{"id-1", "id-2"},
	}
	userService.On(testutils.FunctionName(t, ports.UserService.CreateMany), mock.Anything, mock.AnythingOfType("[]models.CreateUserReq")).Return(expectedResp, nil).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	req := &pb.CreateManyUsersRequest{
		Users: []*pb.CreateUserRequest{
			{Email: "test1@test.com", Password: "test"},
			{Email: "test2@test.com", Password: "test"},
		},
	}

	// Act
	resp, err := handler.CreateMany(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResp.IDs, resp.Ids)
}

// TestCreateManyUsers_ServiceError checks that the CreateMany handler returns a gRPC error when the service fails
func TestCreateManyUsers_ServiceError(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	expectedError := "service-error"
	userService.On(testutils.FunctionName(t, ports.UserService.CreateMany), mock.Anything, mock.AnythingOfType("[]models.CreateUserReq")).Return(models.CreateManyUserResp{}, errors.New(expectedError)).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	req := &pb.CreateManyUsersRequest{
		Users: []*pb.CreateUserRequest{
			{Email: "test1@test.com", Password: "test"},
		},
	}

	// Act
	_, err := handler.CreateMany(context.Background(), req)

	// Assert
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, expectedError, st.Message())
}

// TestGetAllUsers_Ok checks that the GetAll handler returns the expected response
func TestGetAllUsers_Ok(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	expectedUsers := []models.GetUserResp{
		{
			ID:        "test-id",
			Email:     "test@test.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	userService.On(testutils.FunctionName(t, ports.UserService.GetAll), mock.Anything).Return(expectedUsers, nil).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	// Act
	resp, err := handler.GetAll(context.Background(), &emptypb.Empty{})

	// Assert
	assert.NoError(t, err)
	assert.Len(t, resp.Users, 1)
	assert.Equal(t, expectedUsers[0].ID, resp.Users[0].Id)
	assert.Equal(t, expectedUsers[0].Name, resp.Users[0].Name)
	assert.Equal(t, expectedUsers[0].Surnames, resp.Users[0].Surnames)
	assert.Equal(t, expectedUsers[0].Email, resp.Users[0].Email)
	assert.Equal(t, expectedUsers[0].ClaimIDs, resp.Users[0].ClaimIds)
	assert.True(t, expectedUsers[0].CreatedAt.Equal(resp.Users[0].CreatedAt.AsTime()))
	assert.True(t, expectedUsers[0].UpdatedAt.Equal(resp.Users[0].UpdatedAt.AsTime()))
}

// TestGetAllUsers_ServiceError checks that the GetAll handler returns a gRPC error when the service fails
func TestGetAllUsers_ServiceError(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	expectedError := "service-error"
	userService.On(testutils.FunctionName(t, ports.UserService.GetAll), mock.Anything).Return([]models.GetUserResp{}, errors.New(expectedError)).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	// Act
	_, err := handler.GetAll(context.Background(), &emptypb.Empty{})

	// Assert
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, expectedError, st.Message())
}

// TestGetUserByEmail_Ok checks that the GetByEmail handler returns the expected response on a valid request
func TestGetUserByEmail_Ok(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	testEmail := "test@test.com"
	expectedUser := models.GetUserResp{
		ID:        "test-id",
		Email:     testEmail,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userService.On(testutils.FunctionName(t, ports.UserService.GetByEmail), mock.Anything, testEmail).Return(expectedUser, nil).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	req := &pb.GetUserByEmailRequest{Email: testEmail}

	// Act
	resp, err := handler.GetByEmail(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, resp.Id)
	assert.Equal(t, expectedUser.Name, resp.Name)
	assert.Equal(t, expectedUser.Surnames, resp.Surnames)
	assert.Equal(t, expectedUser.Email, resp.Email)
	assert.Equal(t, expectedUser.ClaimIDs, resp.ClaimIds)
	assert.True(t, expectedUser.CreatedAt.Equal(resp.CreatedAt.AsTime()))
	assert.True(t, expectedUser.UpdatedAt.Equal(resp.UpdatedAt.AsTime()))
}

// TestGetUserByEmail_ServiceError checks that the GetByEmail handler returns a gRPC error when the service fails
func TestGetUserByEmail_ServiceError(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	expectedError := "service-error"
	testEmail := "test@test.com"
	userService.On(testutils.FunctionName(t, ports.UserService.GetByEmail), mock.Anything, testEmail).Return(models.GetUserResp{}, errors.New(expectedError)).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	req := &pb.GetUserByEmailRequest{Email: testEmail}

	// Act
	_, err := handler.GetByEmail(context.Background(), req)

	// Assert
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, expectedError, st.Message())
}

// TestGetUserByID_Ok checks that the GetByID handler returns the expected response on a valid request
func TestGetUserByID_Ok(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	testID := "test-id"
	expectedUser := models.GetUserResp{
		ID:        testID,
		Email:     "test@test.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userService.On(testutils.FunctionName(t, ports.UserService.GetByID), mock.Anything, testID).Return(expectedUser, nil).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	req := &pb.GetUserByIDRequest{Id: testID}

	// Act
	resp, err := handler.GetByID(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, resp.Id)
	assert.Equal(t, expectedUser.Name, resp.Name)
	assert.Equal(t, expectedUser.Surnames, resp.Surnames)
	assert.Equal(t, expectedUser.Email, resp.Email)
	assert.Equal(t, expectedUser.ClaimIDs, resp.ClaimIds)
	assert.True(t, expectedUser.CreatedAt.Equal(resp.CreatedAt.AsTime()))
	assert.True(t, expectedUser.UpdatedAt.Equal(resp.UpdatedAt.AsTime()))
}

// TestGetUserByID_ServiceError checks that the GetByID handler returns a gRPC error when the service fails
func TestGetUserByID_ServiceError(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	expectedError := "service-error"
	testID := "test-id"
	userService.On(testutils.FunctionName(t, ports.UserService.GetByID), mock.Anything, testID).Return(models.GetUserResp{}, errors.New(expectedError)).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	req := &pb.GetUserByIDRequest{Id: testID}

	// Act
	_, err := handler.GetByID(context.Background(), req)

	// Assert
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, expectedError, st.Message())
}

// TestUpdateUser_Ok checks that the Update handler returns no error on a valid request
func TestUpdateUser_Ok(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	testID := "test-id"
	userService.On(testutils.FunctionName(t, ports.UserService.Update), mock.Anything, testID, mock.AnythingOfType("models.UpdateUserReq")).Return(nil).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	req := &pb.UpdateUserRequest{Id: testID}

	// Act
	_, err := handler.Update(context.Background(), req)

	// Assert
	assert.NoError(t, err)
}

// TestUpdateUser_ServiceError checks that the Update handler returns a gRPC error when the service fails
func TestUpdateUser_ServiceError(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	expectedError := "service-error"
	testID := "test-id"
	userService.On(testutils.FunctionName(t, ports.UserService.Update), mock.Anything, testID, mock.AnythingOfType("models.UpdateUserReq")).Return(errors.New(expectedError)).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	req := &pb.UpdateUserRequest{Id: testID}

	// Act
	_, err := handler.Update(context.Background(), req)

	// Assert
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, expectedError, st.Message())
}

// TestGetUserClaims_Ok checks that the GetClaims handler returns the expected response.
func TestGetUserClaims_Ok(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	expectedClaims := map[int]string{
		1: "test-claim-1",
		2: "test-claim-2",
	}
	userService.On(testutils.FunctionName(t, ports.UserService.GetUserClaims), mock.Anything).Return(expectedClaims).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	// Act
	resp, err := handler.GetClaims(context.Background(), &emptypb.Empty{})

	// Assert
	assert.NoError(t, err)
	assert.Len(t, resp.Claims, 2)

	// Verify the claims
	claimsMap := make(map[int32]string)
	for _, claim := range resp.Claims {
		claimsMap[claim.Id] = claim.Value
	}
	assert.Equal(t, expectedClaims[1], claimsMap[1])
	assert.Equal(t, expectedClaims[2], claimsMap[2])
}

// TestDeleteUser_Ok checks that the Delete handler returns no error on a valid request
func TestDeleteUser_Ok(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	testID := "test-id"
	userService.On(testutils.FunctionName(t, ports.UserService.Delete), mock.Anything, testID).Return(nil).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	req := &pb.DeleteUserRequest{Id: testID}

	// Act
	_, err := handler.Delete(context.Background(), req)

	// Assert
	assert.NoError(t, err)
}

// TestDeleteUser_ServiceError checks that the Delete handler returns a gRPC error when the service fails
func TestDeleteUser_ServiceError(t *testing.T) {
	// Arrange
	userService := mocks.NewUserService(t)
	expectedError := "service-error"
	testID := "test-id"
	userService.On(testutils.FunctionName(t, ports.UserService.Delete), mock.Anything, testID).Return(errors.New(expectedError)).Once()

	cfg := config.Config{}
	handler := NewUserHandler(context.Background(), cfg, userService)

	req := &pb.DeleteUserRequest{Id: testID}

	// Act
	_, err := handler.Delete(context.Background(), req)

	// Assert
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, expectedError, st.Message())
}

package handlersold

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/models"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	"github.com/sergicanet9/go-hexagonal-api/test/mocks"
	"github.com/sergicanet9/scv-go-tools/v3/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestLoginUser_Ok checks that LoginUser handler returns the expected response when a valid request is received
func TestLoginUser_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	expectedResponse := models.LoginUserResp{
		User:  models.UserResp{},
		Token: "test-token",
	}
	userService.On(testutils.FunctionName(t, ports.UserService.Login), mock.Anything, mock.AnythingOfType("models.LoginUserReq")).Return(expectedResponse, nil).Once()

	cfg := config.Config{}
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/users/login"
	body := models.LoginUserReq{
		Email:    "test@test.com",
		Password: "test",
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(b))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusOK, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response models.LoginUserResp
	if err = json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedResponse, response)
}

// TestLoginUser_InvalidRequest checks that LoginUser handler returns an error when the received request is not valid
func TestLoginUser_InvalidRequest(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	expectedError := map[string]string(map[string]string{"error": "invalid character 'i' looking for beginning of value"})

	cfg := config.Config{}
	userHandler := NewUserHandler(context.Background(), cfg, nil)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/users/login"
	invalidBody := []byte(`{"Email":invalid-type}`)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(invalidBody))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedError, response)
}

// TestLoginUser_LoginError checks that LoginUser handler returns an error when the Login function from the service fails
func TestLoginUser_LoginError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	expectedError := "service-error"
	userService.On(testutils.FunctionName(t, ports.UserService.Login), mock.Anything, mock.AnythingOfType("models.LoginUserReq")).Return(models.LoginUserResp{}, errors.New(expectedError)).Once()

	cfg := config.Config{}
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/users/login"
	body := models.LoginUserReq{
		Email:    "test@test.com",
		Password: "test",
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(b))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, map[string]string(map[string]string{"error": expectedError}), response)
}

// TestCreateUser_Ok checks that CreateUser handler returns the expected response when a valid request is received
func TestCreateUser_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	expectedResponse := models.CreationResp{
		InsertedID: "new-id",
	}
	userService.On(testutils.FunctionName(t, ports.UserService.Create), mock.Anything, mock.AnythingOfType("models.CreateUserReq")).Return(expectedResponse, nil).Once()

	cfg := config.Config{}
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/users"
	body := models.CreateUserReq{
		Email:        "test@test.com",
		PasswordHash: "test",
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(b))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusCreated, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response models.CreationResp
	if err = json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedResponse, response)
}

// TestCreateUser_InvalidRequest checks that CreateUser handler returns an error when the received request is not valid
func TestCreateUser_InvalidRequest(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	expectedError := map[string]string(map[string]string{"error": "invalid character 'i' looking for beginning of value"})

	cfg := config.Config{}
	userHandler := NewUserHandler(context.Background(), cfg, nil)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/users"
	invalidBody := []byte(`{"Email":invalid-type}`)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(invalidBody))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedError, response)
}

// TestCreateUser_CreateError checks that CreateUser handler returns an error when the Create function from the service fails
func TestCreateUser_CreateError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	expectedError := "service-error"
	userService.On(testutils.FunctionName(t, ports.UserService.Create), mock.Anything, mock.AnythingOfType("models.CreateUserReq")).Return(models.CreationResp{}, errors.New(expectedError)).Once()

	cfg := config.Config{}
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/users"
	body := models.CreateUserReq{
		Email:        "test@test.com",
		PasswordHash: "test",
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(b))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, map[string]string(map[string]string{"error": expectedError}), response)
}

// TestCreateManyUsers_Ok checks that CreateManyUsers handler returns the expected response when a valid request is received
func TestCreateManyUsers_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	expectedResponse := models.MultiCreationResp{
		InsertedIDs: []string{"new-id"},
	}
	userService.On(testutils.FunctionName(t, ports.UserService.CreateMany), mock.Anything, mock.AnythingOfType("[]models.CreateUserReq")).Return(expectedResponse, nil).Once()

	cfg := config.Config{}
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/users/many"
	body := []models.CreateUserReq{
		{
			Email:        "test@test.com",
			PasswordHash: "test",
		},
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(b))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusCreated, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response models.MultiCreationResp
	if err = json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedResponse, response)
}

// TestCreateManyUsers_InvalidRequest checks that CreateManyUsers handler returns an error when the received request is not valid
func TestCreateManyUsers_InvalidRequest(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	expectedError := map[string]string(map[string]string{"error": "invalid character 'i' looking for beginning of value"})

	cfg := config.Config{}
	userHandler := NewUserHandler(context.Background(), cfg, nil)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/users/many"
	invalidBody := []byte(`[{"Email":invalid-type}]`)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(invalidBody))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedError, response)
}

// TestCreateManyUsers_CreateManyError checks that CreateManyUsers handler returns an error when the CreateMany function from the service fails
func TestCreateManyUsers_CreateManyError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	expectedError := "service-error"
	userService.On(testutils.FunctionName(t, ports.UserService.CreateMany), mock.Anything, mock.AnythingOfType("[]models.CreateUserReq")).Return(models.MultiCreationResp{}, errors.New(expectedError)).Once()

	cfg := config.Config{}
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/users/many"
	body := []models.CreateUserReq{
		{
			Email:        "test@test.com",
			PasswordHash: "test",
		},
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(b))

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, map[string]string(map[string]string{"error": expectedError}), response)
}

// TestGetAllUsers_Ok checks that GetAllUsers handler returns the expected response when everything goes as expected
func TestGetAllUsers_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	expectedResponse := []models.UserResp{
		{
			Email: "test@test.com",
		},
	}
	userService.On(testutils.FunctionName(t, ports.UserService.GetAll), mock.Anything).Return(expectedResponse, nil).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/users"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusOK, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response []models.UserResp
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedResponse, response)
}

// TestGetAllUsers_GetAllError checks that GetAllUsers handler returns an error when the GetAll function from the service fails
func TestGetAllUsers_GetAllError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	expectedError := "service-error"
	userService.On(testutils.FunctionName(t, ports.UserService.GetAll), mock.Anything).Return([]models.UserResp{}, errors.New(expectedError)).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/users"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, map[string]string(map[string]string{"error": expectedError}), response)
}

// TestGetUserByEmail_Ok checks that GetUserByEmail handler returns the expected response when everything goes as expected
func TestGetUserByEmail_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	expectedResponse := models.UserResp{
		Email: "test@test.com",
	}
	userService.On(testutils.FunctionName(t, ports.UserService.GetByEmail), mock.Anything, expectedResponse.Email).Return(expectedResponse, nil).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := fmt.Sprintf("http://testing/v1/users/email/%s", expectedResponse.Email)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusOK, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response models.UserResp
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedResponse, response)
}

// TestGetUserByEmail_GetByEmailError checks that GetUserByEmail handler returns an error when the GetByEmail function from the service fails
func TestGetUserByEmail_GetByEmailError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	expectedError := "service-error"
	testEmail := "test@test.com"
	userService.On(testutils.FunctionName(t, ports.UserService.GetByEmail), mock.Anything, testEmail).Return(models.UserResp{}, errors.New(expectedError)).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := fmt.Sprintf("http://testing/v1/users/email/%s", testEmail)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, map[string]string(map[string]string{"error": expectedError}), response)
}

// TestGetUserByID_Ok checks that GetUserByID handler returns the expected response when everything goes as expected
func TestGetUserByID_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	expectedResponse := models.UserResp{
		ID: "test-id",
	}
	userService.On(testutils.FunctionName(t, ports.UserService.GetByID), mock.Anything, expectedResponse.ID).Return(expectedResponse, nil).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := fmt.Sprintf("http://testing/v1/users/%s", expectedResponse.ID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusOK, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response models.UserResp
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedResponse, response)
}

// TestGetUserByID_GetUserByIDError checks that GetUserByID handler returns an error when the GetByID function from the service fails
func TestGetUserByID_GetByIDError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	expectedError := "service-error"
	testID := "test-id"
	userService.On(testutils.FunctionName(t, ports.UserService.GetByID), mock.Anything, testID).Return(models.UserResp{}, errors.New(expectedError)).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := fmt.Sprintf("http://testing/v1/users/%s", testID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, map[string]string(map[string]string{"error": expectedError}), response)
}

// TestUpdateUser_Ok checks that UpdateUser handler returns the expected response when a valid request is received
func TestUpdateUser_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	testID := "test-id"
	userService.On(testutils.FunctionName(t, ports.UserService.Update), mock.Anything, testID, mock.AnythingOfType("models.UpdateUserReq")).Return(nil).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := fmt.Sprintf("http://testing/v1/users/%s", testID)
	testEmail := "test@test.com"
	body := models.UpdateUserReq{
		Email: &testEmail,
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(b))
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusOK, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
}

// TesUpdateUser_InvalidRequest checks that UpdateUser handler returns an error when the received request is not valid
func TestUpdateUser_InvalidRequest(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	expectedError := map[string]string(map[string]string{"error": "invalid character 'i' looking for beginning of value"})

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	userHandler := NewUserHandler(context.Background(), cfg, nil)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	testID := "test-id"
	url := fmt.Sprintf("http://testing/v1/users/%s", testID)
	invalidBody := []byte(`{"Email":invalid-type}`)
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(invalidBody))
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, expectedError, response)
}

// TestUpdateUser_UpdateError checks that UpdateUser handler returns an error when the Update function from the service fails
func TestUpdateUser_UpdateError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	expectedError := "service-error"
	testID := "test-id"
	userService.On(testutils.FunctionName(t, ports.UserService.Update), mock.Anything, testID, mock.AnythingOfType("models.UpdateUserReq")).Return(errors.New(expectedError)).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := fmt.Sprintf("http://testing/v1/users/%s", testID)
	testEmail := "test@test.com"
	body := models.UpdateUserReq{
		Email: &testEmail,
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(b))
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, map[string]string(map[string]string{"error": expectedError}), response)
}

// TestDeleteUser_Ok checks that DeleteUser handler returns the expected response when everything goes as expected
func TestDeleteUser_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	testID := "test-id"
	userService.On(testutils.FunctionName(t, ports.UserService.Delete), mock.Anything, testID).Return(nil).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := fmt.Sprintf("http://testing/v1/users/%s", testID)
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	headerName := "Authorization"
	jwtOkAdmin := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXV0aG9yaXplZCI6dHJ1ZX0.ojgkTERSbf1y8rVLNNF-u70hp_GmJOcsmdB5PLcdews"
	req.Header.Add(headerName, jwtOkAdmin)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusOK, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
}

// TestDeleteUser_DeleteError checks that DeleteUser handler returns an error when the Delete function from the service fails
func TestDeleteUser_DeleteError(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	expectedError := "service-error"
	testID := "test-id"
	userService.On(testutils.FunctionName(t, ports.UserService.Delete), mock.Anything, testID).Return(errors.New(expectedError)).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := fmt.Sprintf("http://testing/v1/users/%s", testID)
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	headerName := "Authorization"
	jwtOkAdmin := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXV0aG9yaXplZCI6dHJ1ZX0.ojgkTERSbf1y8rVLNNF-u70hp_GmJOcsmdB5PLcdews"
	req.Header.Add(headerName, jwtOkAdmin)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusInternalServerError, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", req.URL, err)
	}
	assert.Equal(t, map[string]string(map[string]string{"error": expectedError}), response)
}

// TestGetUserClaims_Ok checks that GetUserClaims handler returns the expected response when everything goes as expected
func TestGetUserClaims_Ok(t *testing.T) {
	// Arrange
	r := mux.NewRouter()

	userService := mocks.NewUserService(t)
	expectedResponse := map[int]string{
		0: "test-claim",
	}
	userService.On(testutils.FunctionName(t, ports.UserService.GetUserClaims), mock.Anything).Return(expectedResponse).Once()

	cfg := config.Config{}
	cfg.JWTSecret = "test-secret"
	userHandler := NewUserHandler(context.Background(), cfg, userService)
	SetUserRoutes(r, userHandler)

	rr := httptest.NewRecorder()
	url := "http://testing/v1/claims"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	headerName := "Authorization"
	jwtOk := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek"
	req.Header.Add(headerName, jwtOk)

	// Act
	r.ServeHTTP(rr, req)

	// Assert
	if want, got := http.StatusOK, rr.Code; want != got {
		t.Fatalf("unexpected http status code: want=%d but got=%d", want, got)
	}
}

package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"testing"

	"github.com/sergicanet9/go-mongo-restapi/models/entities"
	"github.com/sergicanet9/go-mongo-restapi/models/requests"
	"github.com/sergicanet9/go-mongo-restapi/models/responses"
	infrastructure "github.com/sergicanet9/scv-go-framework/v2/infrastructure/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_Login_Ok(t *testing.T) {
	// Arrange
	cfg := New(t)
	testUser := getNewTestUser()
	testUser.Email = "testlogin@test.com"
	err := insertUser(&testUser, cfg.DBName, cfg.DBConnectionString)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	body := requests.User{
		Email:        "testlogin@test.com",
		PasswordHash: "test",
	}
	b, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("%s:%d/api/users/login", cfg.Address, cfg.Port)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	// Assert
	if want, got := http.StatusOK, resp.StatusCode; want != got {
		t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
	}
	var response responses.Login
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
	}
	assert.Equal(t, testUser.ID, response.User.ID)
	assert.NotEmpty(t, response.Token)
}

func Test_Create_Created(t *testing.T) {
	// Arrange
	cfg := New(t)
	testUser := getNewTestUser()

	// Act
	body := requests.User(testUser)
	b, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("%s:%d/api/users", cfg.Address, cfg.Port)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	// Assert
	if want, got := http.StatusCreated, resp.StatusCode; want != got {
		t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
	}
	var response responses.Creation
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
	}
	assert.NotEmpty(t, response.InsertedID)
	createdUser, err := findUser(response.InsertedID, cfg.DBName, cfg.DBConnectionString)
	if err != nil {
		t.Fatalf("unexpected error while finding the created user: %s", err)
	}
	assert.Equal(t, testUser.Name, createdUser.Name)
	assert.Equal(t, testUser.Surnames, createdUser.Surnames)
	assert.Equal(t, testUser.Email, createdUser.Email)
}

func Test_GetAllUsers_Ok(t *testing.T) {
	// Arrange
	cfg := New(t)

	// Act
	url := fmt.Sprintf("%s:%d/api/users", cfg.Address, cfg.Port)

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", nonExpiryToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	// Assert
	if want, got := http.StatusOK, resp.StatusCode; want != got {
		t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
	}
	var response []responses.User
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
	}
	assert.NotNil(t, response)
}

func Test_GetByEmail_Ok(t *testing.T) {
	// Arrange
	cfg := New(t)
	testUser := getNewTestUser()
	err := insertUser(&testUser, cfg.DBName, cfg.DBConnectionString)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	url := fmt.Sprintf("%s:%d/api/users/email/%s", cfg.Address, cfg.Port, testUser.Email)

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", nonExpiryToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	// Assert
	if want, got := http.StatusOK, resp.StatusCode; want != got {
		t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
	}
	var response responses.User
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
	}
	assert.Equal(t, testUser.Email, response.Email)
}

func Test_GetByID_Ok(t *testing.T) {
	// Arrange
	cfg := New(t)
	testUser := getNewTestUser()
	err := insertUser(&testUser, cfg.DBName, cfg.DBConnectionString)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	url := fmt.Sprintf("%s:%d/api/users/%s", cfg.Address, cfg.Port, testUser.ID.Hex())

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", nonExpiryToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	// Assert
	if want, got := http.StatusOK, resp.StatusCode; want != got {
		t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
	}
	var response responses.User
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
	}
	assert.Equal(t, testUser.ID, response.ID)
}

func Test_Update_Ok(t *testing.T) {
	// Arrange
	cfg := New(t)
	testUser := getNewTestUser()
	err := insertUser(&testUser, cfg.DBName, cfg.DBConnectionString)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	testUser.Name = "modified"
	testUser.Surnames = "modified"
	body := requests.User(testUser)
	b, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("%s:%d/api/users/%s", cfg.Address, cfg.Port, testUser.ID.Hex())

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", nonExpiryToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	// Assert
	if want, got := http.StatusOK, resp.StatusCode; want != got {
		t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
	}
	updatedUser, err := findUser(testUser.ID, cfg.DBName, cfg.DBConnectionString)
	if err != nil {
		t.Fatalf("unexpected error while finding the created user: %s", err)
	}
	assert.Equal(t, testUser.Name, updatedUser.Name)
	assert.Equal(t, testUser.Surnames, updatedUser.Surnames)
	assert.Equal(t, testUser.Email, updatedUser.Email)
}

func Test_Delete_Ok(t *testing.T) {
	// Arrange
	cfg := New(t)
	testUser := getNewTestUser()
	err := insertUser(&testUser, cfg.DBName, cfg.DBConnectionString)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	body := requests.User(testUser)
	b, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("%s:%d/api/users/%s", cfg.Address, cfg.Port, testUser.ID.Hex())

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", nonExpiryToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	// Assert
	if want, got := http.StatusOK, resp.StatusCode; want != got {
		t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
	}
	_, err = findUser(testUser.ID, cfg.DBName, cfg.DBConnectionString)
	assert.Equal(t, mongo.ErrNoDocuments, err)
}

func Test_GetClaims_Ok(t *testing.T) {
	// Arrange
	cfg := New(t)

	// Act
	url := fmt.Sprintf("%s:%d/api/claims", cfg.Address, cfg.Port)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", nonExpiryToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	// Assert
	if want, got := http.StatusOK, resp.StatusCode; want != got {
		t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
	}
	var response map[int]string
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
	}
	assert.NotEmpty(t, response)
}

func getNewTestUser() entities.User {
	return entities.User{
		ID:           primitive.NewObjectID(),
		Name:         "test",
		Surnames:     "test",
		Email:        fmt.Sprintf("test%d@test.com", rand.Int()),
		PasswordHash: "$2a$10$Q71DDcyvQhzt2K1EbRp1cOh4ToUh9de9ETsixwXGOVeRorTh8tjN2", // test hashed
	}
}

func insertUser(u *entities.User, dbname, connection string) error {
	db, err := infrastructure.ConnectMongoDB(dbname, connection)
	if err != nil {
		return err
	}

	_, err = db.Collection(entities.CollectionNameUser).InsertOne(context.Background(), u)
	return err
}

func findUser(ID primitive.ObjectID, dbname, connection string) (entities.User, error) {
	db, err := infrastructure.ConnectMongoDB(dbname, connection)
	if err != nil {
		return entities.User{}, err
	}

	var u entities.User
	err = db.Collection(entities.CollectionNameUser).FindOne(context.Background(), bson.M{"_id": ID}).Decode(&u)
	return u, err
}

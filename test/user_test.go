package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/domain"
	"github.com/sergicanet9/go-hexagonal-api/core/dto/requests"
	"github.com/sergicanet9/go-hexagonal-api/core/dto/responses"
	"github.com/sergicanet9/scv-go-tools/v3/infrastructure"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_LoginUser_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser := getNewTestUser()
		testUser.Email = "testlogin@test.com"
		err := insertUser(&testUser, cfg)
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
		var response responses.LoginUser
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
		}
		assert.Equal(t, testUser.ID, response.User.ID)
		assert.NotEmpty(t, response.Token)
	})
}

func Test_CreateUser_Created(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
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
		createdUser, err := findUser(response.InsertedID, cfg)
		if err != nil {
			t.Fatalf("unexpected error while finding the created user: %s", err)
		}
		assert.Equal(t, testUser.Name, createdUser.Name)
		assert.Equal(t, testUser.Surnames, createdUser.Surnames)
		assert.Equal(t, testUser.Email, createdUser.Email)
	})
}

func Test_GetAllUsers_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)

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
	})
}

func Test_GetUserByEmail_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser := getNewTestUser()
		err := insertUser(&testUser, cfg)
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
	})
}

func Test_GetUserByID_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser := getNewTestUser()
		err := insertUser(&testUser, cfg)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		url := fmt.Sprintf("%s:%d/api/users/%s", cfg.Address, cfg.Port, testUser.ID)

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
	})
}

func Test_UpdateUser_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser := getNewTestUser()
		err := insertUser(&testUser, cfg)
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

		url := fmt.Sprintf("%s:%d/api/users/%s", cfg.Address, cfg.Port, testUser.ID)

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
		updatedUser, err := findUser(testUser.ID, cfg)
		if err != nil {
			t.Fatalf("unexpected error while finding the created user: %s", err)
		}
		assert.Equal(t, testUser.Name, updatedUser.Name)
		assert.Equal(t, testUser.Surnames, updatedUser.Surnames)
		assert.Equal(t, testUser.Email, updatedUser.Email)
	})
}

func Test_DeleteUser_Ok(t *testing.T) {
	notFoundError := map[string]error{"mongo": mongo.ErrNoDocuments, "postgres": sql.ErrNoRows}
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser := getNewTestUser()
		err := insertUser(&testUser, cfg)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		body := requests.User(testUser)
		b, err := json.Marshal(body)
		if err != nil {
			log.Fatal(err)
		}

		url := fmt.Sprintf("%s:%d/api/users/%s", cfg.Address, cfg.Port, testUser.ID)

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
		_, err = findUser(testUser.ID, cfg)
		assert.Equal(t, notFoundError[database], err)
	})
}

func Test_GetUserClaims_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)

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
	})
}

func getNewTestUser() domain.User {
	return domain.User{
		Name:         "test",
		Surnames:     "test",
		Email:        fmt.Sprintf("test%d@test.com", rand.Int()),
		PasswordHash: "$2a$10$Q71DDcyvQhzt2K1EbRp1cOh4ToUh9de9ETsixwXGOVeRorTh8tjN2", // test hashed
	}
}

func insertUser(u *domain.User, cfg config.Config) error {
	switch cfg.Database {
	case "mongo":
		db, err := infrastructure.ConnectMongoDB(context.Background(), cfg.MongoDBName, cfg.MongoConnectionString)
		if err != nil {
			return err
		}

		result, err := db.Collection(domain.EntityNameUser).InsertOne(context.Background(), u)
		u.ID = result.InsertedID.(primitive.ObjectID).Hex()
		return err

	case "postgres":
		db, err := infrastructure.ConnectPostgresDB(cfg.PostgresConnectionString)
		if err != nil {
			return err
		}

		q := `
		INSERT INTO users (name, surnames, email, password_hash, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, name, surnames, email, password_hash, created_at, updated_at;
		`

		row := db.QueryRowContext(
			context.Background(), q, u.Name, u.Surnames, u.Email, u.PasswordHash, time.Now().UTC(), time.Now().UTC(),
		)

		err = row.Scan(&u.ID, &u.Name, &u.Surnames, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
		return err

	default:
		return fmt.Errorf("database flag %s not valid", cfg.Database)
	}
}

func findUser(ID string, cfg config.Config) (domain.User, error) {
	switch cfg.Database {
	case "mongo":
		db, err := infrastructure.ConnectMongoDB(context.Background(), cfg.MongoDBName, cfg.MongoConnectionString)
		if err != nil {
			return domain.User{}, err
		}

		objectID, err := primitive.ObjectIDFromHex(ID)
		if err != nil {
			return domain.User{}, err
		}

		var u domain.User
		err = db.Collection(domain.EntityNameUser).FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&u)
		return u, err

	case "postgres":
		db, err := infrastructure.ConnectPostgresDB(cfg.PostgresConnectionString)
		if err != nil {
			return domain.User{}, err
		}

		q := `
		SELECT id, name, surnames, email, password_hash, created_at, updated_at
			FROM users WHERE id = $1;
		`

		row := db.QueryRowContext(context.Background(), q, ID)

		var u domain.User
		err = row.Scan(&u.ID, &u.Name, &u.Surnames, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
		return u, err

	default:
		return domain.User{}, fmt.Errorf("database flag %s not valid", cfg.Database)
	}
}

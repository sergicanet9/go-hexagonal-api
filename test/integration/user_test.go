package integration

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/entities"
	"github.com/sergicanet9/go-hexagonal-api/proto/gen/go/pb"
	"github.com/sergicanet9/scv-go-tools/v3/infrastructure"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/encoding/protojson"
)

// TestLoginUser_Ok checks that Login endpoint returns the expected response when everything goes as expected
func TestLoginUser_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser, password := getNewTestUser()
		err := insertUser(&testUser, cfg)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		body := &pb.LoginUserRequest{
			Email:    testUser.Email,
			Password: password,
		}
		b, err := protojson.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("http://:%d/v1/users/login", cfg.HTTPPort)

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

		var response pb.LoginUserResponse
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("unexpected error reading the response while calling %s: %s", resp.Request.URL, err)
		}
		if err := protojson.Unmarshal(bodyBytes, &response); err != nil {
			t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
		}

		assert.Equal(t, testUser.ID, response.User.Id)
		assert.NotEmpty(t, response.Token)
	})
}

// TestCreateUser checks that CreateUser endpoint returns the expected response when everything goes as expected
func TestCreateUser_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser, password := getNewTestUser()

		// Act
		body := mapUserToCreateUserReq(testUser, password)
		b, err := protojson.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("http://:%d/v1/users", cfg.HTTPPort)

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

		var response pb.CreateUserResponse
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("unexpected error reading the response while calling %s: %s", resp.Request.URL, err)
		}
		if err := protojson.Unmarshal(bodyBytes, &response); err != nil {
			t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
		}

		assert.NotEmpty(t, response.Id)
		createdUser, err := findUser(response.Id, cfg)
		if err != nil {
			t.Fatalf("unexpected error while finding the created user: %s", err)
		}
		assert.Equal(t, testUser.Name, createdUser.Name)
		assert.Equal(t, testUser.Surnames, createdUser.Surnames)
		assert.Equal(t, testUser.Email, createdUser.Email)
		assert.NotNil(t, createdUser.PasswordHash)
		assert.Equal(t, testUser.ClaimIDs, createdUser.ClaimIDs)
		assert.NotNil(t, testUser.CreatedAt)
		assert.NotNil(t, testUser.UpdatedAt)
	})
}

// TestCreateManyUsers_Ok checks that CreateManyUsers endpoint returns the expected response when everything goes as expected
func TestCreateManyUsers_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		user1, password1 := getNewTestUser()
		user2, password2 := getNewTestUser()
		users := []entities.User{user1, user2}

		// Act
		body := &pb.CreateManyUsersRequest{
			Users: []*pb.CreateUserRequest{
				mapUserToCreateUserReq(user1, password1),
				mapUserToCreateUserReq(user2, password2),
			},
		}

		b, err := protojson.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("http://:%d/v1/users/many", cfg.HTTPPort)

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
		var response pb.CreateManyUsersResponse
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("unexpected error reading the response while calling %s: %s", resp.Request.URL, err)
		}
		if err := protojson.Unmarshal(bodyBytes, &response); err != nil {
			t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
		}

		assert.Equal(t, 2, len(response.Ids))
		for i, id := range response.Ids {
			createdUser, err := findUser(id, cfg)
			if err != nil {
				t.Fatalf("unexpected error while finding the created user: %s", err)
			}
			assert.Equal(t, users[i].Name, createdUser.Name)
			assert.Equal(t, users[i].Surnames, createdUser.Surnames)
			assert.Equal(t, users[i].Email, createdUser.Email)
			assert.NotNil(t, users[i].PasswordHash)
			assert.Equal(t, users[i].ClaimIDs, createdUser.ClaimIDs)
			assert.NotNil(t, users[i].CreatedAt)
			assert.NotNil(t, users[i].UpdatedAt)
		}
	})
}

// TestGetAllUsers_Ok checks that GetAllUsers endpoint returns the expected response when everything goes as expected
func TestGetAllUsers_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser, _ := getNewTestUser()
		err := insertUser(&testUser, cfg)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		url := fmt.Sprintf("http://:%d/v1/users", cfg.HTTPPort)

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

		var response pb.GetAllUsersResponse
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("unexpected error reading the response while calling %s: %s", resp.Request.URL, err)
		}
		if err := protojson.Unmarshal(bodyBytes, &response); err != nil {
			t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
		}

		assert.NotEmpty(t, response.Users)
	})
}

// TestGetUserByEmail_Ok checks that GetUserByEmail endpoint returns the expected response when everything goes as expected
func TestGetUserByEmail_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser, _ := getNewTestUser()
		err := insertUser(&testUser, cfg)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		url := fmt.Sprintf("http://:%d/v1/users/email/%s", cfg.HTTPPort, testUser.Email)

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

		var response pb.GetUserResponse
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("unexpected error reading the response while calling %s: %s", resp.Request.URL, err)
		}
		if err := protojson.Unmarshal(bodyBytes, &response); err != nil {
			t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
		}

		assert.Equal(t, testUser.Name, response.Name)
		assert.Equal(t, testUser.Surnames, response.Surnames)
		assert.Equal(t, testUser.Email, response.Email)
		assert.Equal(t, testUser.ClaimIDs, response.ClaimIds)
	})
}

// TestGetUserByID_Ok checks that GetUserByID endpoint returns the expected response when everything goes as expected
func TestGetUserByID_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser, _ := getNewTestUser()
		err := insertUser(&testUser, cfg)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		url := fmt.Sprintf("http://:%d/v1/users/%s", cfg.HTTPPort, testUser.ID)

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

		var response pb.GetUserResponse
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("unexpected error reading the response while calling %s: %s", resp.Request.URL, err)
		}
		if err := protojson.Unmarshal(bodyBytes, &response); err != nil {
			t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
		}

		assert.Equal(t, testUser.Name, response.Name)
		assert.Equal(t, testUser.Surnames, response.Surnames)
		assert.Equal(t, testUser.Email, response.Email)
		assert.Equal(t, testUser.ClaimIDs, response.ClaimIds)
	})
}

// TestUpdateUser_Ok checks that UpdateUser endpoint returns the expected response when everything goes as expected
func TestUpdateUser_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser, _ := getNewTestUser()

		err := insertUser(&testUser, cfg)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		testUser.Name = "modified"
		testUser.Surnames = "modified"
		testUser.ClaimIDs = []int32{0}
		body := &pb.UpdateUserRequest{
			Name:     &testUser.Name,
			Surnames: &testUser.Surnames,
			Claims:   &pb.ClaimIds{Ids: testUser.ClaimIDs},
		}
		b, err := protojson.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("http://:%d/v1/users/%s", cfg.HTTPPort, testUser.ID)

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
		assert.NotNil(t, updatedUser.PasswordHash)
		assert.Equal(t, testUser.ClaimIDs, updatedUser.ClaimIDs)
		assert.NotEqual(t, testUser.CreatedAt, testUser.UpdatedAt)
	})
}

// TestDeleteUser_Ok checks that DeleteUser endpoint returns the expected response when everything goes as expected
func TestDeleteUser_Ok(t *testing.T) {
	notFoundError := map[string]error{"mongo": mongo.ErrNoDocuments, "postgres": sql.ErrNoRows}
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser, _ := getNewTestUser()
		err := insertUser(&testUser, cfg)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		url := fmt.Sprintf("http://:%d/v1/users/%s", cfg.HTTPPort, testUser.ID)

		req, err := http.NewRequest(http.MethodDelete, url, nil)
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

// TestGetUserClaims_Ok checks that GetUserClaims endpoint returns the expected response when everything goes as expected
func TestGetUserClaims_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)

		// Act
		url := fmt.Sprintf("http://:%d/v1/claims", cfg.HTTPPort)

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

		var response pb.GetClaimsResponse
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("unexpected error reading the response while calling %s: %s", resp.Request.URL, err)
		}
		if err := protojson.Unmarshal(bodyBytes, &response); err != nil {
			t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
		}

		assert.NotEmpty(t, response.Claims)
	})
}

// HELP FUNCTIONS
func getNewTestUser() (u entities.User, pwd string) {
	u = entities.User{
		Name:         "test",
		Surnames:     "test",
		Email:        fmt.Sprintf("test%d@test.com", rand.Int()),
		PasswordHash: "$2a$10$Cr1oVDUOUoCT3ZSbLanruO5oIdu9YIqoXWFD7iaR8uKvWjgIoSnqa",
		ClaimIDs:     nil,
	}
	pwd = "test"
	return
}

func insertUser(u *entities.User, cfg config.Config) error {
	now := time.Now().UTC()

	switch cfg.Database {
	case "mongo":
		db, err := infrastructure.ConnectMongoDB(context.Background(), cfg.DSN)
		if err != nil {
			return err
		}
		u.CreatedAt = now
		u.UpdatedAt = now
		result, err := db.Collection(entities.EntityNameUser).InsertOne(context.Background(), u)
		u.ID = result.InsertedID.(primitive.ObjectID).Hex()
		return err

	case "postgres":
		db, err := infrastructure.ConnectPostgresDB(context.Background(), cfg.DSN)
		if err != nil {
			return err
		}

		q := `
		INSERT INTO users (name, surnames, email, password_hash, claim_ids, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, name, surnames, email, password_hash, claim_ids, created_at, updated_at;
		`

		row := db.QueryRowContext(
			context.Background(), q, u.Name, u.Surnames, u.Email, u.PasswordHash, pq.Array(u.ClaimIDs), now, now,
		)

		err = row.Scan(&u.ID, &u.Name, &u.Surnames, &u.Email, &u.PasswordHash, pq.Array(&u.ClaimIDs), &u.CreatedAt, &u.UpdatedAt)
		return err

	default:
		return fmt.Errorf("database flag %s not valid", cfg.Database)
	}
}

func findUser(ID string, cfg config.Config) (entities.User, error) {
	switch cfg.Database {
	case "mongo":
		db, err := infrastructure.ConnectMongoDB(context.Background(), cfg.DSN)
		if err != nil {
			return entities.User{}, err
		}

		objectID, err := primitive.ObjectIDFromHex(ID)
		if err != nil {
			return entities.User{}, err
		}

		var u entities.User
		err = db.Collection(entities.EntityNameUser).FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&u)
		return u, err

	case "postgres":
		db, err := infrastructure.ConnectPostgresDB(context.Background(), cfg.DSN)
		if err != nil {
			return entities.User{}, err
		}

		q := `
		SELECT id, name, surnames, email, password_hash, claim_ids, created_at, updated_at
			FROM users WHERE id = $1;
		`

		row := db.QueryRowContext(context.Background(), q, ID)

		var u entities.User
		err = row.Scan(&u.ID, &u.Name, &u.Surnames, &u.Email, &u.PasswordHash, pq.Array(&u.ClaimIDs), &u.CreatedAt, &u.UpdatedAt)
		return u, err

	default:
		return entities.User{}, fmt.Errorf("database flag %s not valid", cfg.Database)
	}
}

func mapUserToCreateUserReq(user entities.User, password string) *pb.CreateUserRequest {
	return &pb.CreateUserRequest{
		Name:     user.Name,
		Surnames: user.Surnames,
		Email:    user.Email,
		Password: password,
		ClaimIds: user.ClaimIDs,
	}

}

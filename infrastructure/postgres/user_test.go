package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/sergicanet9/go-hexagonal-api/core/entities"
	"github.com/sergicanet9/scv-go-tools/v3/infrastructure"
	"github.com/sergicanet9/scv-go-tools/v3/mocks"
	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
	"github.com/stretchr/testify/assert"
)

// TestNewUserRepository_Ok checks that NewUserRepository creates a new userRepository struct
func TestNewUserRepository_Ok(t *testing.T) {
	// Arrange
	_, db := mocks.NewSqlDB(t)
	defer db.Close()

	// Act
	repo := NewUserRepository(db)

	// Assert
	assert.NotEmpty(t, repo)
}

// TestCreate_Ok checks that Create returns the expected response when a valid entity is received
func TestCreate_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newUser := entities.User{}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectQuery("INSERT INTO users").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	// Act
	id, err := repo.Create(context.Background(), newUser)

	// Assert
	assert.Equal(t, expectedID, id)
	assert.Nil(t, err)
}

// TestCreate_InsertError checks that Create returns an error when the insert statement fails
func TestCreate_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newUser := entities.User{}
	expectedError := "insert error"
	mock.ExpectQuery("INSERT INTO users").WillReturnError(errors.New(expectedError))

	// Act
	_, err := repo.Create(context.Background(), newUser)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGet_Ok checks that Get returns the expected response when a valid filter is received
func TestGet_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedUser := entities.User{
		ID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	filter := map[string]interface{}{"email": "test-email", "name": "test-name"}
	skip := 1
	take := 1
	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surnames", "email", "password_hash", "claims", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Surnames, expectedUser.Email, expectedUser.PasswordHash, pq.Array(expectedUser.Claims), expectedUser.CreatedAt, expectedUser.UpdatedAt))

	// Act
	result, err := repo.Get(context.Background(), filter, &skip, &take)

	// Assert
	assert.Nil(t, err)
	assert.True(t, len(result) == 1)

	entity := *(result[0].(*entities.User))
	assert.Equal(t, expectedUser, entity)
}

// TestGet_SelectError checks that Get returns an error when the select query fails
func TestGet_SelectError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnError(errors.New(expectedError))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGet_NoResourcesFound checks that Get returns an error when no resources are found
func TestGet_NoResourcesFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surnames", "email", "password_hash", "claims", "created_at", "updated_at"}))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestGetByID_Ok checks that GetByID returns the expected response when the received ID has a valid format
func TestGetByID_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedUser := entities.User{
		ID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surnames", "email", "password_hash", "claims", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Surnames, expectedUser.Email, expectedUser.PasswordHash, pq.Array(expectedUser.Claims), expectedUser.CreatedAt, expectedUser.UpdatedAt))

	// Act
	result, err := repo.GetByID(context.Background(), expectedUser.ID)

	// Assert
	assert.Nil(t, err)

	entity := *(result.(*entities.User))
	assert.Equal(t, expectedUser, entity)
}

// TestGetByID_SelectError checks that GetByID returns an error when the select query fails
func TestGetByID_SelectError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnError(errors.New(expectedError))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetByID_ResourceNotFound checks that GetByID returns an error when the resource is not found
func TestGetByID_ResourceNotFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surnames", "email", "password_hash", "claims", "created_at", "updated_at"}))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestUpdate_Ok checks that Update does not return an error when the received ID has a valid format
func TestUpdate_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Update(context.Background(), "", entities.User{})

	// Assert
	assert.Nil(t, err)
}

// TestUpdate_UpdateError checks that Update returns an error when the update statement fails
func TestUpdate_UpdateError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "update error"
	mock.ExpectExec("UPDATE users").WillReturnError(errors.New(expectedError))

	// Act
	err := repo.Update(context.Background(), "", entities.User{})

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestUpdate_NotUpdatedError checks that Update returns an error when the update statement does not update any document
func TestUpdate_NotUpdatedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Update(context.Background(), "", entities.User{})

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestDelete_Ok checks that Delete does not return an error when the received ID has a valid format
func TestDelete_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM users").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Nil(t, err)
}

// TestDelte_DeleteError checks that Delete returns an error when the delete statement fails
func TestDelte_DeleteError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "delete error"
	mock.ExpectExec("DELETE FROM users").WillReturnError(errors.New(expectedError))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestDelete_NotDeletedError checks that Delete returns an error when the delete statement does not delete any document
func TestDelete_NotDeletedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM users").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestCreateMany_Ok checks that CreateMany does not return an error when everything goes as expected
func TestCreateMany_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newUsers := []interface{}{entities.User{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO users").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
	mock.ExpectCommit()

	// Act
	ids, err := repo.CreateMany(context.Background(), newUsers)

	// Assert
	assert.True(t, len(ids) == 1)
	assert.Equal(t, expectedID, ids[0])
	assert.Nil(t, err)
}

// TestCreateMany_InsertError checks that CreateMany returns an error when the insert statement fails
func TestCreateMany_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newUsers := []interface{}{entities.User{}}
	expectedError := "insert error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO users").WillReturnError(errors.New(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newUsers)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateMany_BeginError checks that CreateMany returns an error when the begin statement fails
func TestCreateMany_BeginError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newUsers := []interface{}{entities.User{}}
	expectedError := "begin error"
	mock.ExpectBegin().WillReturnError(errors.New(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newUsers)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateMany_CommitError checks that CreateMany returns an error when the commit statement fails
func TestCreateMany_CommitError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newUsers := []interface{}{entities.User{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	expectedError := "commit error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO users").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
	mock.ExpectCommit().WillReturnError(errors.New(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newUsers)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

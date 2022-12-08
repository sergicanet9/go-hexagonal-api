package postgres

import (
	"context"
	"testing"

	"github.com/sergicanet9/go-hexagonal-api/core/entities"
	"github.com/sergicanet9/scv-go-tools/v3/infrastructure"
	"github.com/sergicanet9/scv-go-tools/v3/mocks"
	"github.com/stretchr/testify/assert"
)

// TestNewUserRepository_Ok checks that TestNewUserRepository creates a new userRepository struct
func TestNewUserRepository_Ok(t *testing.T) {
	// Arrange
	_, db := mocks.NewSqlDB(t)
	defer db.Close()

	// Act
	repo := NewUserRepository(db)

	// Assert
	assert.NotEmpty(t, repo)
}

// TODO...
func TestCreate_Ok(t *testing.T) {
	// Arrange
	_, db := mocks.NewSqlDB(t)
	defer db.Close()
	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newUser := entities.User{}

	// Act
	_, _ = repo.Create(context.Background(), newUser)

	// Assert
	// assert.IsType(t, newUser.ID, id)
	// assert.Equal(t, nil, err)
}

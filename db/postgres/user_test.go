package postgres

import (
	"testing"

	"github.com/sergicanet9/scv-go-tools/v3/mocks"
	"github.com/stretchr/testify/assert"
)

// TestNewUserRepository_Ok checks that TestNewUserRepository creates a new userRepository struct
func TestNewUserRepository_Ok(t *testing.T) {
	_, db := mocks.NewSqlDB(t)
	defer db.Close()

	// Act
	repo := NewUserRepository(db)

	// Assert
	assert.NotEmpty(t, repo)

}

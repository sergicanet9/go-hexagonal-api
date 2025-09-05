package mongo

import (
	"context"
	"testing"

	"github.com/sergicanet9/go-hexagonal-api/core/entities"
	"github.com/sergicanet9/scv-go-tools/v4/infrastructure"
	"github.com/sergicanet9/scv-go-tools/v4/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// TestNewUserRepository_Ok checks that NewUserRepository creates a new userRepository struct
func TestNewUserRepository_Ok(t *testing.T) {
	mt := mocks.NewMongoDB(t)

	mt.Run("", func(mt *mtest.T) {
		// Arrange
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		// Act
		repo, err := NewUserRepository(context.Background(), mt.DB)

		// Assert
		assert.NotEmpty(t, repo)
		assert.Nil(t, err)
	})
}

// TestCreateMany_Ok checks that CreateMany does not return an error when everything goes as expected
func TestCreateMany_Ok(t *testing.T) {
	mt := mocks.NewMongoDB(t)

	mt.Run("", func(mt *mtest.T) {
		// Arrange
		repo := userRepository{
			infrastructure.MongoRepository{
				DB:         mt.DB,
				Collection: mt.DB.Collection(entities.EntityNameUser),
				Target:     entities.User{},
			},
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		entityToAdd := entities.User{}
		newEntities := []interface{}{entityToAdd}

		// Act
		ids, err := repo.CreateMany(context.Background(), newEntities)

		// Assert
		assert.True(t, len(ids) == 1)
		assert.IsType(t, entityToAdd.ID, ids[0])
		assert.Nil(t, err)
	})
}

// TestCreateMany_CreateError checks that CreateMany returns an error when Create fails
func TestInsertMany_CreateError(t *testing.T) {
	mt := mocks.NewMongoDB(t)

	mt.Run("", func(mt *mtest.T) {
		// Arrange
		repo := userRepository{
			infrastructure.MongoRepository{
				DB:         mt.DB,
				Collection: mt.DB.Collection(entities.EntityNameUser),
				Target:     entities.User{},
			},
		}

		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		newEntities := []interface{}{entities.User{}}

		// Act
		_, err := repo.CreateMany(context.Background(), newEntities)

		// Assert
		assert.NotEmpty(t, err)
	})
}

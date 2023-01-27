package mongo

import (
	"context"
	"testing"

	"github.com/sergicanet9/go-hexagonal-api/core/entities"
	"github.com/sergicanet9/scv-go-tools/v3/infrastructure"
	"github.com/sergicanet9/scv-go-tools/v3/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// TestNewUserRepository_Ok checks that NewUserRepository creates a new userRepository struct
func TestNewUserRepository_Ok(t *testing.T) {
	mt := mocks.NewMongoDB(t)
	defer mt.Close()

	mt.Run("", func(mt *mtest.T) {
		// Act
		repo, err := NewUserRepository(context.Background(), mt.DB)

		// Assert
		assert.NotEmpty(t, repo)
		assert.Nil(t, err)
	})
}

// TestInsertMany_Ok checks that InsertMany does not return an error when everything goes as expected
func TestInsertMany_Ok(t *testing.T) {
	mt := mocks.NewMongoDB(t)
	defer mt.Close()

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

		newEntities := []interface{}{entities.User{}}

		// Act
		err := repo.InsertMany(context.Background(), newEntities)

		// Assert
		assert.Nil(t, err)
	})
}

// TestInsertMany_CreateError checks that InsertMany returns an error when Create fails
func TestInsertMany_CreateError(t *testing.T) {
	mt := mocks.NewMongoDB(t)
	defer mt.Close()

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
		err := repo.InsertMany(context.Background(), newEntities)

		// Assert
		assert.NotEmpty(t, err)
	})
}

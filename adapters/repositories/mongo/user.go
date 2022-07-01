package mongo

import (
	"context"

	"github.com/sergicanet9/go-mongo-restapi/core/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository struct of an user repository for mongo
type UserRepository struct {
	MongoRepository
}

// NewUserRepository creates a user repository for mongo
func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{
		MongoRepository{
			collection,
			domain.User{},
		},
	}
}

func (r *UserRepository) Test(ctx context.Context) error {
	return nil
}

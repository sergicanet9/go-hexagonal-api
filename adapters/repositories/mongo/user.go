package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository struct of an user repository for mongo
type UserRepository struct {
	MongoRepository
}

// NewUserRepository creates a user repository for mongo
func NewUserRepository(collection *mongo.Collection, target interface{}) *UserRepository {
	return &UserRepository{
		MongoRepository{
			collection,
			target,
		},
	}
}

func (r *UserRepository) Test(ctx context.Context) error {
	return nil
}

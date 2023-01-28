package mongo

import (
	"context"

	"github.com/sergicanet9/go-hexagonal-api/core/entities"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	"github.com/sergicanet9/scv-go-tools/v3/infrastructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// userRepository adapter of an user repository for mongo.
type userRepository struct {
	infrastructure.MongoRepository
}

// NewUserRepository creates a user repository for mongo
func NewUserRepository(ctx context.Context, db *mongo.Database) (ports.UserRepository, error) {
	r := &userRepository{
		infrastructure.MongoRepository{
			DB:         db,
			Collection: db.Collection(entities.EntityNameUser),
			Target:     entities.User{},
		},
	}

	_, err := r.Collection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	return r, err
}

func (r *userRepository) CreateMany(ctx context.Context, users []interface{}) ([]string, error) {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := r.DB.Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	var result []string
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		for _, entity := range users {
			id, err := r.Create(sessionContext, entity)
			if err != nil {
				return nil, err
			}
			result = append(result, id)
		}
		return nil, nil
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)
	return result, err
}

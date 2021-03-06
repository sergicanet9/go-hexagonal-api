package mongo

import (
	"context"

	"github.com/sergicanet9/go-hexagonal-api/core/domain"
	"github.com/sergicanet9/scv-go-tools/v3/infrastructure"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// UserRepository adapter of an user repository for mongo.
type UserRepository struct {
	infrastructure.MongoRepository
}

// NewUserRepository creates a user repository for mongo
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		infrastructure.MongoRepository{
			DB:         db,
			Collection: db.Collection(domain.EntityNameUser),
			Target:     domain.User{},
		},
	}
}

func (r *UserRepository) InsertMany(ctx context.Context, entities []interface{}) error {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := r.DB.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		for _, entity := range entities {
			_, err = r.Create(sessionContext, entity)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	_, err = session.WithTransaction(context.Background(), callback, txnOpts)
	return err
}

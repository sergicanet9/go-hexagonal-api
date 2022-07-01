package mongo

import (
	"context"

	"github.com/sergicanet9/go-mongo-restapi/core/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// UserRepository struct of an user repository for mongo
type UserRepository struct {
	MongoRepository
}

// NewUserRepository creates a user repository for mongo
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		MongoRepository{
			db,
			db.Collection(domain.EntityNameUser),
			domain.User{},
		},
	}
}

func (r *UserRepository) InsertMany(ctx context.Context, entities []interface{}) error {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := r.db.Client().StartSession()
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

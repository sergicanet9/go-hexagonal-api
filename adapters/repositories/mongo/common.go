package mongo

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository struct of a mongo repository
type MongoRepository struct {
	collection *mongo.Collection
	target     interface{}
}

// Create creates an entity in the repository's collection
func (r *MongoRepository) Create(ctx context.Context, entity interface{}) (string, error) {
	result, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Get gets the documents mathing the filter in the repository's collection
func (r *MongoRepository) Get(ctx context.Context, filter map[string]interface{}, skip, take *int) ([]interface{}, error) {
	var result []interface{}

	var skip64, take64 int64
	if skip != nil {
		skip64 = int64(*skip)
	}
	if take != nil {
		take64 = int64(*take)
	}
	cur, err := r.collection.Find(ctx, filter, &options.FindOptions{Skip: &skip64, Limit: &take64})
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		entry := reflect.New(reflect.TypeOf(r.target)).Interface()
		if err := cur.Decode(entry); err != nil {
			return nil, err
		}
		result = append(result, entry)
	}

	return result, nil
}

// GetByID get the document with the specified ID in the repository's collection
func (r *MongoRepository) GetByID(ctx context.Context, ID string) (interface{}, error) {
	result := r.target

	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": _id}
	err = r.collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Update updates the document with the specified ID in the repository's collection
func (r *MongoRepository) Update(ctx context.Context, ID string, entity interface{}, upsert bool) error {
	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": entity}
	opts := options.Update().SetUpsert(upsert)
	result, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if result.ModifiedCount < 1 && result.UpsertedCount < 1 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// Delete deletes the document with the specified ID in the repository's collection
func (r *MongoRepository) Delete(ctx context.Context, ID string) error {
	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": _id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount < 1 {
		return mongo.ErrNoDocuments
	}
	return nil
}

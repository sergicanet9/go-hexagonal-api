package responses

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Creation response struct
type Creation struct {
	InsertedID primitive.ObjectID `json:"insertedId" bson:"insertedId"`
}

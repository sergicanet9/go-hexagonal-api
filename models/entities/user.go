package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CollectionNameUser contains the name of the mongodb collection for the entity
const CollectionNameUser = "users"

// User struct
type User struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Surnames     string             `json:"surnames" bson:"surnames"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"passwordHash" bson:"passwordHash"`
}

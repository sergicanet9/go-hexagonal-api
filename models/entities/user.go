package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CollectionNameUser contains the name of the mongodb collection for the entity
const CollectionNameUser = "users"

// User struct
type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name"`
	Surnames     string             `bson:"surnames"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"passwordHash"`
}

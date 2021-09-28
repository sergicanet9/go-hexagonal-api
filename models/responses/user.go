package responses

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User response struct
type User struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	Surnames     string             `json:"surnames"`
	Email        string             `json:"email"`
	PasswordHash string             `json:"-"`
}

// Login response struct
type Login struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// Creation response struct
type Creation struct {
	InsertedID primitive.ObjectID `json:"insertedId"`
}

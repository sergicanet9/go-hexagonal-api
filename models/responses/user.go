package responses

import (
	"time"

	"github.com/sergicanet9/go-mongo-restapi/models/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User response struct
type User struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	Surnames     string             `json:"surnames"`
	Email        string             `json:"email"`
	PasswordHash string             `json:"-"`
	Claims       []entities.Claim   `bson:"claims"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
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

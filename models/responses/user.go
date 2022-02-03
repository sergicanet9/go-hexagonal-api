package responses

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User response struct
type User struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	Surnames     string             `json:"surnames"`
	Email        string             `json:"email"`
	PasswordHash string             `json:"-"`
	Claims       []int              `json:"claims"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
}

// LoginUser response struct
type LoginUser struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

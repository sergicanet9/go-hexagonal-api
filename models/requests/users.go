package requests

import "go.mongodb.org/mongo-driver/bson/primitive"

// User request struct
type User struct {
	ID           primitive.ObjectID `json:"-"`
	Name         string             `json:"name"`
	Surnames     string             `json:"surnames"`
	Email        string             `json:"email"`
	PasswordHash string             `json:"password"`
}

// Login request struct
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Update struct {
	ID          *primitive.ObjectID `json:"-"`
	Name        *string             `json:"name"`
	Surnames    *string             `json:"surnames"`
	Email       *string             `json:"email"`
	OldPassword *string             `json:"old_password"`
	NewPassword *string             `json:"new_password"`
}

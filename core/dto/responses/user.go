package responses

import (
	"time"
)

// User response struct
type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Surnames     string    `json:"surnames"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Claims       []int64   `json:"claims"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// LoginUser response struct
type LoginUser struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

package requests

import (
	"time"
)

// User request struct
type User struct {
	ID           string    `json:"-"`
	Name         string    `json:"name"`
	Surnames     string    `json:"surnames"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password"`
	Claims       []int64   `json:"claims"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

// LoginUser request struct
type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUser struct {
	ID          string     `json:"-"`
	Name        *string    `json:"name"`
	Surnames    *string    `json:"surnames"`
	Email       *string    `json:"email"`
	OldPassword *string    `json:"old_password"`
	NewPassword *string    `json:"new_password"`
	Claims      *[]int64   `json:"claims"`
	CreatedAt   *time.Time `json:"-"`
	UpdatedAt   *time.Time `json:"-"`
}

package models

import (
	"time"
)

// UserReq user request struct
type UserReq struct {
	ID           string    `json:"-"`
	Name         string    `json:"name"`
	Surnames     string    `json:"surnames"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password"`
	Claims       []int64   `json:"claims"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

// UserResp user response struct
type UserResp struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Surnames     string    `json:"surnames"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Claims       []int64   `json:"claims"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// LoginUserReq login user request struct
type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUserResp login user response struct
type LoginUserResp struct {
	User  UserResp `json:"user"`
	Token string   `json:"token"`
}

// UpdateUserReq update user request struct
type UpdateUserReq struct {
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

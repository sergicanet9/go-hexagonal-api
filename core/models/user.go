package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
)

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

// CreateUserReq user request struct
type CreateUserReq struct {
	ID           string    `json:"-"`
	Name         string    `json:"name"`
	Surnames     string    `json:"surnames"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password"`
	Claims       []int64   `json:"claims"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

// Validate checks that a given CreateUserReq is valid
func (req CreateUserReq) Validate() error {
	var msgs []string

	if req.Email == "" {
		msgs = append(msgs, "email cannot be empty")
	}
	if req.PasswordHash == "" {
		msgs = append(msgs, "password cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}

// LoginUserReq login user request struct
type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate checks that a given LoginUserReq is valid
func (req LoginUserReq) Validate() error {
	var msgs []string

	if req.Email == "" {
		msgs = append(msgs, "email cannot be empty")
	}
	if req.Password == "" {
		msgs = append(msgs, "password cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
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

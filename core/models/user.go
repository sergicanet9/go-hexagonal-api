package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
)

// LoginUserReq login user request struct
type LoginUserReq struct {
	Email    string
	Password string
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
		return wrappers.NewValidationErr(fmt.Errorf("%s", strings.Join(msgs, " | ")))
	}

	return nil
}

// LoginUserResp login user response struct
type LoginUserResp struct {
	User  GetUserResp
	Token string
}

// CreateUserReq create user request struct
type CreateUserReq struct {
	Name     string
	Surnames string
	Email    string
	Password string
	ClaimIDs []int32
}

// Validate checks that a given CreateUserReq is valid
func (req CreateUserReq) Validate() error {
	var msgs []string

	if req.Email == "" {
		msgs = append(msgs, "email cannot be empty")
	}
	if req.Password == "" {
		msgs = append(msgs, "password cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf("%s", strings.Join(msgs, " | ")))
	}

	return nil
}

// CreateUserResp create user respponse struct
type CreateUserResp struct {
	InsertedID string
}

// CreateManyUserResp create many user response struct
type CreateManyUserResp struct {
	InsertedIDs []string
}

// UpdateUserReq update user request struct
type UpdateUserReq struct {
	Name        *string
	Surnames    *string
	Email       *string
	OldPassword *string
	NewPassword *string
	ClaimIDs    *[]int32
}

// GetUserResp user response struct
type GetUserResp struct {
	ID           string
	Name         string
	Surnames     string
	Email        string
	PasswordHash string
	ClaimIDs     []int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

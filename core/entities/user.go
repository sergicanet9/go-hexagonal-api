package entities

import (
	"time"
)

// EntityNameUser contains the name of the entity
const EntityNameUser = "users"

// Claim type
type Claim int

const (
	admin = iota
)

func (claim Claim) String() string {
	return [...]string{"admin"}[claim]
}

func (claim Claim) IsValid() bool {
	return claim >= admin && claim <= admin
}

func GetClaims() map[int]string {
	claims := make(map[int]string)
	for i, claim := range [...]string{"admin"} {
		claims[i] = claim
	}
	return claims
}

// User struct
type User struct {
	ID           string    `bson:"_id,omitempty"`
	Name         string    `bson:"name"`
	Surnames     string    `bson:"surnames"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password_hash"`
	Claims       []int64   `bson:"claims"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

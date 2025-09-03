package entities

import (
	"time"
)

// EntityNameUser contains the name of the entity
const EntityNameUser = "users"

// UserClaim type
type UserClaim int

const (
	admin = iota
)

func (claim UserClaim) String() string {
	return [...]string{"admin"}[claim]
}

func (claim UserClaim) IsValid() bool {
	return claim >= admin && claim <= admin
}

func GetUserClaims() map[int]string {
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
	ClaimIDs     []int32   `bson:"claim_ids"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

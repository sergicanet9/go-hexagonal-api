package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CollectionNameUser contains the name of the mongodb collection for the entity
const CollectionNameUser = "users"

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
	ID           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name"`
	Surnames     string             `bson:"surnames"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"passwordHash"`
	Claims       []Claim            `bson:"claims"`
	CreatedAt    time.Time          `bson:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"`
}

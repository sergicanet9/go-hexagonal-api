package responses

import "github.com/scanet9/go-mongo-restapi/models/entities"

// Login response struct
type Login struct {
	User  entities.User `json:"user" bson:"user"`
	Token string        `json:"token" bson:"token"`
}

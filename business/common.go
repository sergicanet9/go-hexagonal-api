package business

import (
	infrastructure "github.com/scanet9/scv-go-framework/v2/infrastructure/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

//Service struct
type Service struct {
	db   *mongo.Database
	repo infrastructure.MongoRepository
}

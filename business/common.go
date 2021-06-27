package business

import (
	"github.com/scanet9/scv-go-framework/infrastructure"
	"go.mongodb.org/mongo-driver/mongo"
)

//Service struct
type Service struct {
	db   *mongo.Database
	repo infrastructure.MongoRepository
}

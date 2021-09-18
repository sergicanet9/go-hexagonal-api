package business

import (
	"github.com/scanet9/go-mongo-restapi/config"
	infrastructure "github.com/scanet9/scv-go-framework/v2/infrastructure/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

//Service struct
type Service struct {
	config config.Config
	db     *mongo.Database
	repo   infrastructure.MongoRepository
}

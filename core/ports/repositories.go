package ports

import (
	"context"

	"github.com/sergicanet9/scv-go-tools/v3/repository"
)

type UserRepository interface {
	repository.Repository
	InsertMany(ctx context.Context, entities []interface{}) error
}

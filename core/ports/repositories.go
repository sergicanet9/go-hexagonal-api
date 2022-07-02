package ports

import (
	"context"

	"github.com/sergicanet9/scv-go-tools/v3/ports"
)

type UserRepository interface {
	ports.Repository
	InsertMany(ctx context.Context, entities []interface{}) error
}

package ports

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, entity interface{}) (string, error)
	Get(ctx context.Context, filter map[string]interface{}, skip, take *int) ([]interface{}, error)
	GetByID(ctx context.Context, ID string) (interface{}, error)
	Update(ctx context.Context, ID string, entity interface{}, upsert bool) error
	Delete(ctx context.Context, ID string) error
}

type UserRepository interface {
	Repository
	Test(ctx context.Context) error
}

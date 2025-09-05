package ports

import (
	"context"

	"github.com/sergicanet9/go-hexagonal-api/core/models"
	"github.com/sergicanet9/scv-go-tools/v3/repository"
)

// UserRepositoy interface
type UserRepository interface {
	repository.Repository
	CreateMany(ctx context.Context, entities []interface{}) ([]string, error)
}

// UserService interface
type UserService interface {
	Login(ctx context.Context, credentials models.LoginUserReq) (models.LoginUserResp, error)
	Create(ctx context.Context, user models.CreateUserReq) (models.CreateUserResp, error)
	CreateMany(ctx context.Context, users []models.CreateUserReq) (models.CreateManyUserResp, error)
	GetAll(ctx context.Context) ([]models.GetUserResp, error)
	GetByEmail(ctx context.Context, email string) (models.GetUserResp, error)
	GetByID(ctx context.Context, ID string) (models.GetUserResp, error)
	Update(ctx context.Context, ID string, user models.UpdateUserReq) error
	Delete(ctx context.Context, ID string) error
	GetUserClaims(ctx context.Context) map[int]string
}

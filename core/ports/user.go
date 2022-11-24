package ports

import (
	"context"

	"github.com/sergicanet9/go-hexagonal-api/core/models"
	"github.com/sergicanet9/scv-go-tools/v3/repository"
)

// UserRepositoy interface
type UserRepository interface {
	repository.Repository
	InsertMany(ctx context.Context, entities []interface{}) error
}

// UserService interface
type UserService interface {
	Login(ctx context.Context, credentials models.LoginUserReq) (models.LoginUserResp, error)
	Create(ctx context.Context, u models.CreateUserReq) (models.CreationResp, error)
	GetAll(ctx context.Context) ([]models.UserResp, error)
	GetByEmail(ctx context.Context, email string) (models.UserResp, error)
	GetByID(ctx context.Context, ID string) (models.UserResp, error)
	Update(ctx context.Context, ID string, u models.UpdateUserReq) error
	Delete(ctx context.Context, ID string) error
	GetClaims(ctx context.Context) map[int]string
	AtomicTransationProof(ctx context.Context) error
}

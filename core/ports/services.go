package ports

import (
	"context"

	"github.com/sergicanet9/go-mongo-restapi/core/dto/requests"
	"github.com/sergicanet9/go-mongo-restapi/core/dto/responses"
)

type UserService interface {
	Login(ctx context.Context, credentials requests.LoginUser) (responses.LoginUser, error)
	Create(ctx context.Context, u requests.User) (responses.Creation, error)
	GetAll(ctx context.Context) ([]responses.User, error)
	GetByEmail(ctx context.Context, email string) (responses.User, error)
	GetByID(ctx context.Context, ID string) (responses.User, error)
	Update(ctx context.Context, ID string, u requests.UpdateUser) error
	Delete(ctx context.Context, ID string) error
	GetClaims(ctx context.Context) (map[int]string, error)
	AtomicTransationProof(ctx context.Context) error
}

package handlers

import (
	"context"
	"errors"

	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	"github.com/sergicanet9/go-hexagonal-api/proto/gen/go/pb"
	"github.com/sergicanet9/go-hexagonal-api/scvv4/interceptors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userHandler struct {
	ctx context.Context
	cfg config.Config
	svc ports.UserService
	pb.UnimplementedUserServiceServer
}

// NewUserHandler creates a new user handler
func NewUserHandler(ctx context.Context, cfg config.Config, svc ports.UserService) *userHandler {
	return &userHandler{
		ctx: ctx,
		cfg: cfg,
		svc: svc,
	}
}

func (u *userHandler) JWTMethodPolicies() []interceptors.MethodPolicy {
	return []interceptors.MethodPolicy{
		{
			MethodName:     "/user.UserService/GetAll",
			RequiredClaims: nil,
		},
	}
}

func (u *userHandler) GetAll(_ *emptypb.Empty, stream pb.UserService_GetAllServer) error {
	return errors.New("not implementesssd")
}

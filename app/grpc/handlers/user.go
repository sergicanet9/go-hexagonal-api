package handlers

import (
	"context"
	"errors"

	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	pb "github.com/sergicanet9/go-hexagonal-api/proto/protogen"
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

func (h *userHandler) GetAll(_ context.Context, _ *emptypb.Empty) (*pb.UserService_GetAllUsersServer, error) {
	return nil, errors.New("not implemented")
}

package handlers

import (
	"context"
	"errors"

	"github.com/sergicanet9/go-hexagonal-api/app/mappers"
	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	"github.com/sergicanet9/go-hexagonal-api/proto/gen/go/pb"
	"github.com/sergicanet9/go-hexagonal-api/scvv4/interceptors"
	"github.com/sergicanet9/go-hexagonal-api/scvv4/utils"
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

// JWTMethodPolicies defines custom JWT method policies
func (u *userHandler) JWTMethodPolicies() []interceptors.MethodPolicy {
	srv := "user.UserService"
	methods := []struct {
		name   string
		claims []string
	}{
		{"GetAll", nil},
		{"GetByEmail", nil},
		{"GetByID", nil},
		{"Update", nil},
		{"GetClaims", nil},
		{"Delete", []string{"admin"}},
	}

	var policies []interceptors.MethodPolicy
	for _, m := range methods {
		policies = append(policies, interceptors.MethodPolicy{
			MethodName:     "/" + srv + "/" + m.name,
			RequiredClaims: m.claims,
		})
	}
	return policies
}

func (u *userHandler) Login(_ context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	ctx, cancel := context.WithTimeout(u.ctx, u.cfg.Timeout.Duration)
	defer cancel()

	loginReq := mappers.LoginUserReqToModel(req)
	resp, err := u.svc.Login(ctx, loginReq)
	if err != nil {
		return nil, utils.ToGRPC(err)
	}

	return mappers.LoginUserRespToPB(resp), nil
}

func (u *userHandler) Create(_ context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(u.ctx, u.cfg.Timeout.Duration)
	defer cancel()

	createReq := mappers.CreateUserReqToModel(req)
	resp, err := u.svc.Create(ctx, createReq)
	if err != nil {
		return nil, utils.ToGRPC(err)
	}

	return mappers.CreateUserRespToPB(resp), nil
}

func (u *userHandler) CreateMany(_ context.Context, req *pb.CreateManyUsersRequest) (*pb.CreateManyUsersResponse, error) {
	return nil, errors.New("not implemented")
}

func (u *userHandler) GetAll(_ context.Context, _ *emptypb.Empty) (*pb.GetAllUsersResponse, error) {
	return nil, errors.New("not implemented")
}

func (u *userHandler) GetByEmail(_ context.Context, req *pb.GetUserByEmailRequest) (*pb.UserResponse, error) {
	return nil, errors.New("not implemented")
}

func (u *userHandler) GetByID(_ context.Context, req *pb.GetUserByIDRequest) (*pb.UserResponse, error) {
	return nil, errors.New("not implemented")
}

func (u *userHandler) Update(_ context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	return nil, errors.New("not implemented")
}

func (u *userHandler) GetClaims(_ context.Context, _ *emptypb.Empty) (*pb.GetClaimsResponse, error) {
	return nil, errors.New("not implemented")
}

func (u *userHandler) Delete(_ context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	return nil, errors.New("not implemented")
}

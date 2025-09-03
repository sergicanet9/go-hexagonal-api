package handlers

import (
	"context"

	"github.com/sergicanet9/go-hexagonal-api/config"
	"github.com/sergicanet9/go-hexagonal-api/core/models"
	"github.com/sergicanet9/go-hexagonal-api/core/ports"
	"github.com/sergicanet9/go-hexagonal-api/proto/gen/go/pb"
	"github.com/sergicanet9/go-hexagonal-api/scvv4/interceptors"
	"github.com/sergicanet9/go-hexagonal-api/scvv4/utils"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	loginReq := models.LoginUserReq{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := u.svc.Login(ctx, loginReq)
	if err != nil {
		return nil, utils.ToGRPC(err)
	}

	loginResp := &pb.LoginUserResponse{
		User: &pb.GetUserResponse{
			Id:        resp.User.ID,
			Name:      resp.User.Name,
			Surnames:  resp.User.Surnames,
			Email:     resp.User.Email,
			ClaimIds:  resp.User.ClaimIDs,
			CreatedAt: timestamppb.New(resp.User.CreatedAt),
			UpdatedAt: timestamppb.New(resp.User.UpdatedAt),
		},
		Token: resp.Token,
	}
	return loginResp, nil
}

func (u *userHandler) Create(_ context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(u.ctx, u.cfg.Timeout.Duration)
	defer cancel()

	createReq := models.CreateUserReq{
		Name:     req.Name,
		Surnames: req.Surnames,
		Email:    req.Email,
		Password: req.Password,
		ClaimIDs: req.ClaimIds,
	}

	resp, err := u.svc.Create(ctx, createReq)
	if err != nil {
		return nil, utils.ToGRPC(err)
	}

	createResp := &pb.CreateUserResponse{
		Id: resp.InsertedID,
	}
	return createResp, nil
}

func (u *userHandler) CreateMany(_ context.Context, req *pb.CreateManyUsersRequest) (*pb.CreateManyUsersResponse, error) {
	ctx, cancel := context.WithTimeout(u.ctx, u.cfg.Timeout.Duration)
	defer cancel()

	var createManyReq []models.CreateUserReq
	for _, user := range req.Users {
		createManyReq = append(createManyReq, models.CreateUserReq{
			Name:     user.Name,
			Surnames: user.Surnames,
			Email:    user.Email,
			Password: user.Password,
			ClaimIDs: user.ClaimIds,
		})
	}

	resp, err := u.svc.CreateMany(ctx, createManyReq)
	if err != nil {
		return nil, utils.ToGRPC(err)
	}

	createManyResp := &pb.CreateManyUsersResponse{
		Ids: resp.InsertedIDs,
	}
	return createManyResp, nil
}

func (u *userHandler) GetAll(_ context.Context, _ *emptypb.Empty) (*pb.GetAllUsersResponse, error) {
	ctx, cancel := context.WithTimeout(u.ctx, u.cfg.Timeout.Duration)
	defer cancel()

	resp, err := u.svc.GetAll(ctx)
	if err != nil {
		return nil, utils.ToGRPC(err)
	}

	var getAllRespList []*pb.GetUserResponse
	for _, user := range resp {
		getAllRespList = append(getAllRespList, &pb.GetUserResponse{
			Id:        user.ID,
			Name:      user.Name,
			Surnames:  user.Surnames,
			Email:     user.Email,
			ClaimIds:  user.ClaimIDs,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}

	getAllResp := &pb.GetAllUsersResponse{
		Users: getAllRespList,
	}
	return getAllResp, nil
}

func (u *userHandler) GetByEmail(_ context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserResponse, error) {
	ctx, cancel := context.WithTimeout(u.ctx, u.cfg.Timeout.Duration)
	defer cancel()

	resp, err := u.svc.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, utils.ToGRPC(err)
	}

	getByEmailResp := &pb.GetUserResponse{
		Id:        resp.ID,
		Name:      resp.Name,
		Surnames:  resp.Surnames,
		Email:     resp.Email,
		ClaimIds:  resp.ClaimIDs,
		CreatedAt: timestamppb.New(resp.CreatedAt),
		UpdatedAt: timestamppb.New(resp.CreatedAt),
	}
	return getByEmailResp, nil
}

func (u *userHandler) GetByID(_ context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserResponse, error) {
	ctx, cancel := context.WithTimeout(u.ctx, u.cfg.Timeout.Duration)
	defer cancel()

	resp, err := u.svc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, utils.ToGRPC(err)
	}

	getByIDResp := &pb.GetUserResponse{
		Id:        resp.ID,
		Name:      resp.Name,
		Surnames:  resp.Surnames,
		Email:     resp.Email,
		ClaimIds:  resp.ClaimIDs,
		CreatedAt: timestamppb.New(resp.CreatedAt),
		UpdatedAt: timestamppb.New(resp.CreatedAt),
	}
	return getByIDResp, nil
}

func (u *userHandler) Update(_ context.Context, req *pb.UpdateUserRequest) (*emptypb.Empty, error) {
	ctx, cancel := context.WithTimeout(u.ctx, u.cfg.Timeout.Duration)
	defer cancel()

	updateReq := models.UpdateUserReq{
		Name:        &req.Name,
		Surnames:    &req.Surnames,
		Email:       &req.Email,
		OldPassword: &req.OldPassword,
		NewPassword: &req.NewPassword,
		ClaimIDs:    &req.ClaimIds,
	}

	err := u.svc.Update(ctx, req.Id, updateReq)
	if err != nil {
		return nil, utils.ToGRPC(err)
	}
	return &emptypb.Empty{}, nil
}

func (u *userHandler) GetClaims(_ context.Context, _ *emptypb.Empty) (*pb.GetClaimsResponse, error) {
	ctx, cancel := context.WithTimeout(u.ctx, u.cfg.Timeout.Duration)
	defer cancel()

	resp := u.svc.GetUserClaims(ctx)

	var getClaimsList []*pb.Claim
	for id, value := range resp {
		getClaimsList = append(getClaimsList, &pb.Claim{
			Id:    int64(id),
			Value: value,
		})
	}
	getClaimsResp := &pb.GetClaimsResponse{
		Claims: getClaimsList,
	}
	return getClaimsResp, nil
}

func (u *userHandler) Delete(_ context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	ctx, cancel := context.WithTimeout(u.ctx, u.cfg.Timeout.Duration)
	defer cancel()

	err := u.svc.Delete(ctx, req.Id)
	if err != nil {
		return nil, utils.ToGRPC(err)
	}

	return &emptypb.Empty{}, nil
}

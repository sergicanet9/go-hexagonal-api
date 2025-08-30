package handlers

import (
	"context"

	"github.com/sergicanet9/go-hexagonal-api/config"
	pb "github.com/sergicanet9/go-hexagonal-api/proto/protogen"
)

type healthHandler struct {
	ctx context.Context
	cfg config.Config
	pb.UnimplementedHealthServiceServer
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(ctx context.Context, cfg config.Config) *healthHandler {
	return &healthHandler{
		ctx: ctx,
		cfg: cfg,
	}
}

// TODO:
func (h *healthHandler) HealthCheck(_ context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Version:     h.cfg.Version,
		Environment: h.cfg.Environment,
		Database:    h.cfg.Database,
		Port:        int32(h.cfg.Port),
	}, nil
}

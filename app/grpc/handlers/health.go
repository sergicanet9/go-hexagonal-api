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

// HealthCheck .
func (h *healthHandler) HealthCheck(_ context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	response := &pb.HealthCheckResponse{
		Version:     h.cfg.Version,
		Environment: h.cfg.Environment,
		Database:    h.cfg.Database,
		HttpPort:    int32(h.cfg.HTTPPort),
		GrpcPort:    int32(h.cfg.GRPCPort),
		Dsn:         "***FILTERED***",
	}
	if h.cfg.Environment == "local" {
		response.Dsn = h.cfg.DSN
	}

	return response, nil
}

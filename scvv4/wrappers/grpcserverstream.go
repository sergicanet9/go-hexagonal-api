package wrappers

import (
	"context"

	"google.golang.org/grpc"
)

type gRPCServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// NewGRPCServerStream creates a new wrapper for grpc.ServerStream
func NewGRPCServerStream(ctx context.Context) *gRPCServerStream {
	return &gRPCServerStream{ctx: ctx}
}

// Context returns the context associated to the stream
func (m *gRPCServerStream) Context() context.Context {
	return m.ctx
}

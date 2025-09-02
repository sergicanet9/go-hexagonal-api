package mocks

import (
	"context"

	"google.golang.org/grpc"
)

type MockGRPCServerStream struct {
	grpc.ServerStream
	Ctx context.Context
}

func (m *MockGRPCServerStream) Context() context.Context {
	return m.Ctx
}

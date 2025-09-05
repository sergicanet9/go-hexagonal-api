package interceptors

import (
	"context"
	"fmt"

	"github.com/sergicanet9/go-hexagonal-api/scvv4/utils"
	"google.golang.org/grpc"
)

// UnaryRecover is a gRPC unary interceptor that recovers from panics and returns a gRPC error for the incomming call.
func UnaryRecover() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoverFromPanic(info.FullMethod, r)
			}
		}()
		return handler(ctx, req)
	}
}

// StreamRecover is a gRPC stream interceptor that recovers from panics and returns a gRPC error for the incomming call.
func StreamRecover() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoverFromPanic(info.FullMethod, r)
			}
		}()
		return handler(srv, ss)
	}
}

func recoverFromPanic(methodName string, r interface{}) error {
	return utils.ToGRPC(fmt.Errorf("recovered from a panic during gRPC call for method: %s, Panic: %v", methodName, r))
}

package interceptors

import (
	"context"
	"time"

	"github.com/sergicanet9/go-hexagonal-api/scvv4/observability"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryLogger is a gRPC unary interceptor that logs details of the incomming call.
func UnaryLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		logResult("Unary", info.FullMethod, start, err, req, resp)

		return resp, err
	}
}

// StreamLogger is a gRPC stream interceptor that logs details of the incomming call.
func StreamLogger() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		err := handler(srv, ss)

		logResult("Stream", info.FullMethod, start, err, nil, nil)

		return err
	}
}

func logResult(callType, fullMethod string, start time.Time, err error, req interface{}, resp interface{}) {
	latency := time.Since(start)

	if err != nil {
		s, _ := status.FromError(err)
		observability.Logger().Printf("gRPC %s Call: %s - Request: %v - Status: %s - Latency: %s - Error: %v",
			callType, fullMethod, req, s.Code(), latency, err)
	} else {
		observability.Logger().Printf("gRPC %s Call: %s - Request: %v - Status: OK - Latency: %s - Response: %v",
			callType, fullMethod, req, latency, resp)
	}
}

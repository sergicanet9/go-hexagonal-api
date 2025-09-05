package interceptors

import (
	"context"
	"errors"
	"testing"

	"github.com/sergicanet9/go-hexagonal-api/scvv4/utils"
	"github.com/sergicanet9/go-hexagonal-api/scvv4/wrappers"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestUnaryLogger_HandlerOk checks that the unary interceptor preserves the response when the handler returns a successful response
func TestUnaryLogger_HandlerOk(t *testing.T) {
	method := "/TestService/TestMethod"

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	interceptor := UnaryLogger()
	resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: method,
	}, handler)

	assert.Nil(t, err)
	assert.Equal(t, "ok", resp)
}

// TestUnaryLogger_HandlerError checks that the unary interceptor preserves the status code and message when the handler returns a grpc error
func TestUnaryLogger_HandlerError(t *testing.T) {
	method := "/TestService/TestMethod"

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, utils.ToGRPC(errors.New("test error"))
	}

	interceptor := UnaryLogger()
	resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: method,
	}, handler)

	status, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, status.Code())
	assert.Equal(t, "test error", status.Message())
	assert.Nil(t, resp)
}

// TestStreamLogger_HandlerOk checks that the stream interceptor does not interfere when the handler does not return a grpc error
func TestStreamLogger_HandlerOk(t *testing.T) {
	method := "/TestService/TestStreamMethod"
	handler := func(srv interface{}, ss grpc.ServerStream) error {
		return nil
	}

	interceptor := StreamLogger()
	stream := wrappers.NewGRPCServerStream(context.Background())
	err := interceptor(nil, stream, &grpc.StreamServerInfo{FullMethod: method}, handler)

	assert.Nil(t, err)
}

// TestStreamLogger_HandlerError checks that the unary interceptor preserves the status code and message when the handler returns a grpc error
func TestStreamLogger_HandlerError(t *testing.T) {
	method := "/TestService/TestStreamMethod"
	handler := func(srv interface{}, ss grpc.ServerStream) error {
		return utils.ToGRPC(errors.New("test error"))
	}

	interceptor := StreamLogger()
	stream := wrappers.NewGRPCServerStream(context.Background())
	err := interceptor(nil, stream, &grpc.StreamServerInfo{FullMethod: method}, handler)

	status, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, status.Code())
	assert.Equal(t, "test error", status.Message())
}

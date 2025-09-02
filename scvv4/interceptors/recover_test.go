package interceptors

import (
	"context"
	"testing"

	"github.com/sergicanet9/go-hexagonal-api/scvv4/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestUnaryRecover_NoPanic checks that the inteceptor does not return an error when no panic happens in the handler
func TestUnaryRecover_NoPanic(t *testing.T) {
	// Arrange
	interceptor := UnaryRecover()
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok-response", nil
	}

	// Act
	resp, err := interceptor(context.Background(), "req", &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}, handler)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "ok-response", resp)
}

// TestUnaryRecover_Panic checks that the interceptor returns an error when there is a panic in the handler.
func TestUnaryRecover_Panic(t *testing.T) {
	// Arrange
	interceptor := UnaryRecover()
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("test panic")
	}

	// Act
	resp, err := interceptor(context.Background(), "req", &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}, handler)

	// Assert
	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "test panic")
}

// TestStreamRecover_NoPanic checks that the inteceptor does not return an error when no panic happens in the handler
func TestStreamRecover_NoPanic(t *testing.T) {
	// Arrange
	interceptor := StreamRecover()
	handler := func(srv interface{}, ss grpc.ServerStream) error {
		return nil
	}

	// Act
	err := interceptor("srv", &mocks.MockGRPCServerStream{}, &grpc.StreamServerInfo{FullMethod: "/test.Service/Stream"}, handler)

	// Assert
	assert.Nil(t, err)
}

// TestStreamRecover_Panic checks that the interceptor returns an error when there is a panic in the handler.
func TestStreamRecover_Panic(t *testing.T) {
	// Arrange
	interceptor := StreamRecover()
	handler := func(srv interface{}, ss grpc.ServerStream) error {
		panic("stream panic")
	}

	// Act
	err := interceptor("srv", &mocks.MockGRPCServerStream{}, &grpc.StreamServerInfo{FullMethod: "/test.Service/Stream"}, handler)

	// Assert
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "stream panic")
}

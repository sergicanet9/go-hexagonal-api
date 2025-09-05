package utils

import (
	"errors"
	"testing"

	wrappersv4 "github.com/sergicanet9/go-hexagonal-api/scvv4/wrappers"
	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestToGRPC verifies that different errors are mapped to the expected grpcStatus code
func TestToGRPC(t *testing.T) {
	tests := []struct {
		name     string
		inputErr error
		wantCode codes.Code
	}{
		{"nil error", nil, codes.OK},
		{"not found", wrappers.NonExistentErr, codes.NotFound},
		{"invalid input", wrappers.ValidationErr, codes.InvalidArgument},
		{"unauthorized", wrappers.UnauthorizedErr, codes.Unauthenticated},
		{"unauthenticated", wrappersv4.UnauthenticatedErr, codes.PermissionDenied},
		{"unknown error", errors.New("test-error"), codes.Internal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ToGRPC(tt.inputErr)

			if tt.inputErr == nil {
				if err != nil {
					t.Fatalf("expected nil, got %v", err)
				}
				return
			}

			st, ok := status.FromError(err)
			if !ok {
				t.Fatalf("expected status error, got %v", err)
			}
			if st.Code() != tt.wantCode {
				t.Errorf("expected code %v, got %v", tt.wantCode, st.Code())
			}
		})
	}
}

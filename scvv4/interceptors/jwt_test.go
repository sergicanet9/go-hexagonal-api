package interceptors

import (
	"context"
	"strings"
	"testing"

	"github.com/sergicanet9/go-hexagonal-api/scvv4/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// TestUnaryJWT checks that the unary JWT inteceptor correctly handles all expected scenarios
func TestUnaryJWT(t *testing.T) {
	secret := "test-secret"
	method := "/TestService/TestMethod"

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			md := metadata.New(nil)
			if tc.jwtToken != "" {
				md.Set("authorization", tc.jwtToken)
			}
			ctx := metadata.NewIncomingContext(context.Background(), md)

			interceptor := UnaryJWT(secret, []MethodPolicy{
				{MethodName: method, RequiredClaims: tc.requiredClaims},
			})

			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "", nil
			}

			resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: method}, handler)

			if tc.expectedCode == codes.OK {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedMsg, resp)
			} else {
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedCode, st.Code())
				assert.True(t, strings.Contains(st.Message(), tc.expectedMsg))
			}
		})
	}
}

// TestStreamJWT checks that the stream JWT inteceptor correctly handles all expected scenarios
func TestStreamJWT(t *testing.T) {
	secret := "test-secret"
	method := "/TestService/TestStreamMethod"

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			md := metadata.New(nil)
			if tc.jwtToken != "" {
				md.Set("authorization", tc.jwtToken)
			}
			ctx := metadata.NewIncomingContext(context.Background(), md)
			stream := &mocks.MockGRPCServerStream{Ctx: ctx}

			interceptor := StreamJWT(secret, []MethodPolicy{
				{MethodName: method, RequiredClaims: tc.requiredClaims},
			})

			handler := func(srv interface{}, ss grpc.ServerStream) error {
				return nil
			}

			err := interceptor(nil, stream, &grpc.StreamServerInfo{FullMethod: method}, handler)

			if tc.expectedCode == codes.OK {
				assert.Nil(t, err)
			} else {
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedCode, st.Code())
				assert.True(t, strings.Contains(st.Message(), tc.expectedMsg))
			}
		})
	}
}

var cases = []struct {
	name           string
	jwtToken       string
	requiredClaims []string
	expectedCode   codes.Code
	expectedMsg    string
}{
	{
		name:         "Valid token and claims",
		jwtToken:     "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek",
		expectedCode: codes.OK,
		expectedMsg:  "",
	},
	{
		name:         "Missing token",
		jwtToken:     "",
		expectedCode: codes.Unauthenticated,
		expectedMsg:  "authorization token is not provided",
	},
	{
		name:         "Malformed token",
		jwtToken:     "123",
		expectedCode: codes.Unauthenticated,
		expectedMsg:  "invalid token format",
	},
	{
		name:         "Invalid signin method",
		jwtToken:     "Bearer eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.e30.",
		expectedCode: codes.Unauthenticated,
		expectedMsg:  "signin method not valid",
	},
	{
		name:         "Invalid secret",
		jwtToken:     "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.df4nTfuWWdndfrlIxF0iWUrrcANrM4bzKdbYa9VeAj8",
		expectedCode: codes.Unauthenticated,
		expectedMsg:  "invalid token",
	},
	{
		name:           "Missing required claim",
		jwtToken:       "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.mpHl842O7xEZjgQ8CyX8xYLDoEORGVMnAxULkW-u8Ek",
		requiredClaims: []string{"role"},
		expectedCode:   codes.PermissionDenied,
		expectedMsg:    "required claim",
	},
}

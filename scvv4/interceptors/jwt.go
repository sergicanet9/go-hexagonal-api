package interceptors

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sergicanet9/go-hexagonal-api/scvv4/utils"
	wrappersv4 "github.com/sergicanet9/go-hexagonal-api/scvv4/wrappers"
	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type claimsCtxKey string

const ClaimsKey claimsCtxKey = "claims"

type MethodPolicy struct {
	MethodName     string
	RequiredClaims []string
}

// UnaryJWT is a configurable gRPC unary interceptor that validates the JWT tokens and its claims for the incomming call
func UnaryJWT(jwtSecret string, methods []MethodPolicy) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		policy, isProtected := findMethodPolicy(methods, info.FullMethod)
		if !isProtected {
			return handler(ctx, req)
		}

		newCtx, err := jwtValidator(ctx, jwtSecret, policy.RequiredClaims)
		if err != nil {
			return nil, err
		}

		return handler(newCtx, req)
	}
}

// StreamJWT is a configurable gRPC stream interceptor that validates the JWT tokens and its claims for the incomming call
func StreamJWT(jwtSecret string, methods []MethodPolicy) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		policy, isProtected := findMethodPolicy(methods, info.FullMethod)
		if !isProtected {
			return handler(srv, ss)
		}

		ctx := ss.Context()
		newCtx, err := jwtValidator(ctx, jwtSecret, policy.RequiredClaims)
		if err != nil {
			return err
		}
		wrappedStream := wrappersv4.NewGRPCServerStream(newCtx)
		wrappedStream.ServerStream = ss

		return handler(srv, wrappedStream)
	}
}

func findMethodPolicy(methods []MethodPolicy, fullMethod string) (MethodPolicy, bool) {
	for _, policy := range methods {
		if policy.MethodName == fullMethod {
			return policy, true
		}
	}
	return MethodPolicy{}, false
}

func jwtValidator(ctx context.Context, jwtSecret string, requiredClaims []string) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.ToGRPC(wrappers.NewUnauthorizedErr(errors.New("metadata is not provided")))
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return nil, utils.ToGRPC(wrappers.NewUnauthorizedErr(errors.New("authorization token is not provided")))
	}
	authorization := tokens[0]

	bearerToken := strings.Split(authorization, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		return nil, utils.ToGRPC(wrappers.NewUnauthorizedErr(errors.New("invalid token format, should be Bearer + {token}")))
	}
	tokenString := bearerToken[1]

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, utils.ToGRPC(wrappers.NewUnauthorizedErr(errors.New("signin method not valid")))
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, utils.ToGRPC(wrappers.NewUnauthorizedErr(fmt.Errorf("invalid token: %v", err)))
	}

	for _, requiredClaim := range requiredClaims {
		if _, ok := claims[requiredClaim]; !ok {
			return nil, utils.ToGRPC(wrappersv4.NewUnauthenticatedErr(fmt.Errorf("insufficient permissions: required claim '%s' not found", requiredClaim)))
		}
	}

	return context.WithValue(ctx, ClaimsKey, claims), nil
}

package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sergicanet9/go-hexagonal-api/scvv4/utils"
)

type claimsCtxKey string

const ClaimsKey claimsCtxKey = "claims"

// JWT is a configurable HTTP middleware that validates the JWT tokens and its claims for the incomming call
func JWT(jwtSecret string, requiredClaims ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				utils.ErrorResponse(w, http.StatusUnauthorized, "authorization token is not provided")
				return
			}

			bearerToken := strings.Split(authorization, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				utils.ErrorResponse(w, http.StatusUnauthorized, "invalid token format, should be Bearer + {token}")
				return
			}
			tokenString := bearerToken[1]

			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("signin method not valid")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				utils.ErrorResponse(w, http.StatusUnauthorized, fmt.Sprintf("invalid token. %v", err))
				return
			}

			for _, requiredClaim := range requiredClaims {
				if _, ok := claims[requiredClaim]; !ok {
					utils.ErrorResponse(w, http.StatusForbidden, fmt.Sprintf("insufficient permissions: required claim '%s' not found or incorrect", requiredClaim))
					return
				}
			}

			newCtx := context.WithValue(r.Context(), ClaimsKey, claims)
			r = r.WithContext(newCtx)
			next.ServeHTTP(w, r)
		})
	}
}

package http

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sinland/smtp-proxy-api/internal/domain"
	"net/http"
	"strings"
)

const (
	ctxKeyClaims = "authClaims"
)

func WithJWTClaims(ctx context.Context, claims *domain.JWTClaims) context.Context {
	return context.WithValue(ctx, ctxKeyClaims, claims)
}

func JWTClaimsFromContext(ctx context.Context) *domain.JWTClaims {
	val, ok := ctx.Value(ctxKeyClaims).(*domain.JWTClaims)
	if !ok {
		return nil
	}
	return val
}

// NewJWTMiddleware Middleware to authenticate by jwt token
func NewJWTMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing Authorization header", http.StatusUnauthorized)
				return
			}

			// Check if it's a Bearer token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				// No Bearer prefix found
				http.Error(w, "invalid token format", http.StatusUnauthorized)
				return
			}

			// Parse and validate the JWT
			claims := &domain.JWTClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// Attach user ID to context
			next.ServeHTTP(w, r.WithContext(WithJWTClaims(r.Context(), claims)))
		})
	}
}

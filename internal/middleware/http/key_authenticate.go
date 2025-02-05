package http

import (
	"net/http"
	"strings"
)

// NewApiKeyMiddleware Middleware to authenticate by api key
func NewApiKeyMiddleware(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing Authorization header", http.StatusUnauthorized)
				return
			}

			// Check if it's a Bearer token
			tokenString := strings.TrimPrefix(authHeader, "Token ")
			if tokenString == authHeader {
				// No Bearer prefix found
				http.Error(w, "invalid token format", http.StatusUnauthorized)
				return
			}

			if apiKey != tokenString {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// allow to process request
			next.ServeHTTP(w, r)
		})
	}
}

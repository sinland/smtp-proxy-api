package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// JWTClaims struct
type JWTClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

// GenerateToken Generate JWT Token
func GenerateToken(username string, jwtSecret string) (string, error) {
	cl := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
		Username: username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)

	return token.SignedString([]byte(jwtSecret))
}

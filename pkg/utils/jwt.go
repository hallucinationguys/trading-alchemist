package utils

import (
	"fmt"
	"time"
	"trading-alchemist/internal/domain/auth"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token
func GenerateJWT(user *auth.User, secret string, ttl time.Duration, issuer string) (string, error) {
	if user == nil {
		return "", fmt.Errorf("user cannot be nil")
	}

	claims := &Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT validates a JWT token and returns the user ID
func ValidateJWT(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
} 
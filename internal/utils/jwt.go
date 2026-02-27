package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// JWTClaims represents JWT claims
type JWTClaims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token
func GenerateToken(userID, email, secret string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken generates a refresh token (longer expiration)
func GenerateRefreshToken(userID, email, secret string) (string, error) {
	return GenerateToken(userID, email, secret, 7*24*time.Hour) // 7 days
}

// ValidateToken validates and parses a JWT token
func ValidateToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ParseTokenDuration parses duration string (e.g., "24h", "7d")
func ParseTokenDuration(durationStr string) (time.Duration, error) {
	// Handle common formats
	if durationStr == "" {
		return 24 * time.Hour, nil // default 24 hours
	}

	// Try parsing as Go duration
	duration, err := time.ParseDuration(durationStr)
	if err == nil {
		return duration, nil
	}

	// Try parsing custom formats like "7d"
	if len(durationStr) >= 2 {
		unit := durationStr[len(durationStr)-1:]
		valueStr := durationStr[:len(durationStr)-1]

		var value int64
		_, err := fmt.Sscanf(valueStr, "%d", &value)
		if err != nil {
			return 0, fmt.Errorf("invalid duration format: %s", durationStr)
		}

		switch unit {
		case "d":
			return time.Duration(value) * 24 * time.Hour, nil
		case "h":
			return time.Duration(value) * time.Hour, nil
		case "m":
			return time.Duration(value) * time.Minute, nil
		}
	}

	return 0, fmt.Errorf("invalid duration format: %s", durationStr)
}

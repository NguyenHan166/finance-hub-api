package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomToken generates a random token for verification/reset
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateVerificationToken generates a token for email verification
func GenerateVerificationToken() (string, error) {
	return GenerateRandomToken(32) // 64 character hex string
}

// GenerateResetToken generates a token for password reset
func GenerateResetToken() (string, error) {
	return GenerateRandomToken(32) // 64 character hex string
}

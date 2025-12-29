package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(userID uuid.UUID, orgID uuid.UUID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"org_id":  orgID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 3 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Get secret from environment variable
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "fallback_secret_for_dev_only" // Optional fallback
	}

	return token.SignedString([]byte(secret))
}

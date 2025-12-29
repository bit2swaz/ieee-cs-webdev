package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	return func(c fiber.Ctx) error {
		// 1. Get Token from Header (Authorization: Bearer <token>)
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized: No token provided"})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// 2. Parse Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Fetch secret with fallback to match utils/jwt.go
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				secret = "fallback_secret_for_dev_only"
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized: Invalid token"})
		}

		// 3. Extract Claims (User Data)
		claims := token.Claims.(jwt.MapClaims)

		// 4. Inject into Context (So handlers can access it)
		c.Locals("user_id", claims["user_id"])
		c.Locals("org_id", claims["org_id"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}
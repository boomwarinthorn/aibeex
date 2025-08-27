package middleware

import (
	"strings"

	"fiber-hello-world/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware validates JWT tokens in requests
func JWTMiddleware(jwtService *jwt.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Authorization header required",
			})
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Bearer token required",
			})
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return c.Status(401).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Token required",
			})
		}

		// Validate token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Invalid token",
			})
		}

		// Store user claims in context
		c.Locals("user", claims)
		return c.Next()
	}
}

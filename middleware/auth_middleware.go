package middleware

import (
	"os"
	"strings"
	"tugas_uas/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// [FR-002] RBAC Middleware
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.ErrorResponse{Success: false, Message: "Missing token"})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.ErrorResponse{Success: false, Message: "Invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.ErrorResponse{Success: false, Message: "Invalid claims"})
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}

// Permission Check
func PermissionCheck(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.ErrorResponse{Success: false, Message: "Role not found"})
		}

		roleStr := userRole.(string)
		if roleStr == "Admin" { return c.Next() } // Admin Bypass

		if roleStr != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(helper.ErrorResponse{Success: false, Message: "Forbidden access"})
		}
		return c.Next()
	}
}
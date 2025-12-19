package middleware

import (
	"os"
	"strings"
	"tugas_uas/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// [FR-002] RBAC Middleware: Validasi Token
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		
		// 1. Cek apakah header ada
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.ErrorResponse{Success: false, Message: "Missing Authorization Header"})
		}

		// 2. Bersihkan format "Bearer <token>"
		// Kode ini akan menghapus kata "Bearer " jika ada, tapi jika user cuma kirim token mentah, tetap diterima.
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.ErrorResponse{Success: false, Message: "Token string is empty"})
		}

		// 3. Parse & Validasi Token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Pastikan metode signing sesuai (HMAC)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.ErrorResponse{Success: false, Message: "Invalid or Expired Token"})
		}

		// 4. Ambil Claims (Data dalam token)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.ErrorResponse{Success: false, Message: "Invalid Token Claims"})
		}

		// 5. Simpan ke Context agar bisa dipakai di Controller
		// Pastikan casting tipe datanya aman
		if userID, ok := claims["user_id"].(string); ok {
			c.Locals("user_id", userID)
		}
		if role, ok := claims["role"].(string); ok {
			c.Locals("role", role)
		}

		return c.Next()
	}
}

// Permission Check: Validasi Role
func PermissionCheck(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok || userRole == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(helper.ErrorResponse{Success: false, Message: "Role not found in token"})
		}

		// Admin selalu boleh akses (Super User bypass)
		if userRole == "Admin" {
			return c.Next()
		}

		// Cek apakah role user sesuai
		if userRole != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(helper.ErrorResponse{Success: false, Message: "Forbidden access: You are not " + requiredRole})
		}
		
		return c.Next()
	}
}
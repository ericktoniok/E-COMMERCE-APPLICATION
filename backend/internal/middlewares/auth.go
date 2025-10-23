package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"mini-ecommerce/backend/internal/appctx"
)

const CtxUserID = "uid"
const CtxUserRole = "role"

func JWTProtect(ax *appctx.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := c.Get("Authorization")
		if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing bearer token"})
		}
		token := strings.TrimSpace(h[7:])
		claims, err := ax.JWT.Parse(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		c.Locals(CtxUserID, claims.UserID)
		c.Locals(CtxUserRole, claims.Role)
		return c.Next()
	}
}

func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if r, ok := c.Locals(CtxUserRole).(string); !ok || r != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		return c.Next()
	}
}

package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"mini-ecommerce/backend/internal/appctx"
	"mini-ecommerce/backend/internal/services"
)

type AuthController struct {
	Ax   *appctx.Context
	Auth *services.AuthService
}

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResp struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

func NewAuthController(ax *appctx.Context, auth *services.AuthService) *AuthController {
	return &AuthController{Ax: ax, Auth: auth}
}

func (h *AuthController) Register(c *fiber.Ctx) error {
	var req registerReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad request"})
	}
	u, err := h.Auth.Register(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	tok, err := h.Ax.JWT.Generate(u.ID, string(u.Role))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "token error"})
	}
	return c.JSON(authResp{Token: tok, Role: string(u.Role)})
}

func (h *AuthController) Login(c *fiber.Ctx) error {
	var req loginReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad request"})
	}
	u, err := h.Auth.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	tok, err := h.Ax.JWT.Generate(u.ID, string(u.Role))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "token error"})
	}
	return c.JSON(authResp{Token: tok, Role: string(u.Role)})
}

// JWT expiry is set in the manager at construction time
func defaultJWTExpiry() time.Duration { return 72 * time.Hour }

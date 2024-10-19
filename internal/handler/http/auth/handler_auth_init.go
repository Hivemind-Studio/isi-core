package user

import (
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService serviceAuthInterface
}

func NewAuthHandler(authService serviceAuthInterface) *AuthHandler {
	return &AuthHandler{
		authService,
	}
}

func (h *AuthHandler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1/auth")

	v1.Post("/login", h.Login)
}

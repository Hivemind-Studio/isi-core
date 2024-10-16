package user

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	authService serviceAuthInterface
}

func NewAuthHandler(authService serviceAuthInterface) *Handler {
	return &Handler{
		authService,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/auth")

	v1.Post("/login", h.Login)
}

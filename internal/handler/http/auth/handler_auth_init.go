package user

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	authService serviceAuthInterface
	userService serviceUserInterface
}

func NewAuthHandler(authService serviceAuthInterface, userService serviceUserInterface) *Handler {
	return &Handler{
		authService,
		userService,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1/auth")

	v1.Post("/login", h.Login)
	v1.Post("/register", h.Create)
	v1.Post("/verify-email", h.SendEmailVerification)
}

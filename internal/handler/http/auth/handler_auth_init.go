package user

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	authService  serviceAuthInterface
	userService  serviceUserInterface
	coachService serviceCoachInterface
}

func NewAuthHandler(authService serviceAuthInterface, userService serviceUserInterface, coachService serviceCoachInterface) *Handler {
	return &Handler{
		authService,
		userService,
		coachService,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1/auth")

	v1.Post("/login", h.Login)
	v1.Post("/register", h.Create)
	v1.Post("/verify-email", h.SendEmailVerification)
	v1.Patch("/coach/password", h.PatchPassword)
}

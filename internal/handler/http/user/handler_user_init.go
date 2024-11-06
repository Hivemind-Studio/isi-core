package user

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	userService serviceUserInterface
}

func NewUserHandler(userService serviceUserInterface) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	v1.Post("/users", h.Create)
	v1.Get("/users", h.GetUsers)
	v1.Get("/users/:id", h.GetUserByID)
}

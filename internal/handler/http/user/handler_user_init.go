package user

import (
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService serviceUserInterface
}

func NewUserHandler(userService serviceUserInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	v1.Post("/users", h.CreateUser)
	v1.Get("/users/:id", h.GetStorageUnitDetail)
}

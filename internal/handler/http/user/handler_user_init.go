package user

import (
	"github.com/Hivemind-Studio/isi-core/utils"
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
	v1 := app.Group("/api/v1", utils.RoleMiddleware("Admin"))

	v1.Post("/users", h.CreateUser)
	v1.Get("/users/:id", h.GetStorageUnitDetail)
}

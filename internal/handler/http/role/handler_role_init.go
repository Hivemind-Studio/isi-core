package role

import "github.com/gofiber/fiber/v2"

type Handler struct {
	createRoleUseCase CreateRoleUseCaseInterface
}

func NewRoleHandler(createRoleUseCase CreateRoleUseCaseInterface) *Handler {
	return &Handler{
		createRoleUseCase,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	v1.Post("/createrole", h.CreateRole)
}

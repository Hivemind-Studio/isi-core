package role

import "github.com/gofiber/fiber/v2"

type Handler struct {
	roleService serviceRoleInterface
}

func NewRoleHandler(roleService serviceRoleInterface) *Handler {
	return &Handler{
		roleService,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	v1.Post("/role", h.CreateRole)
}

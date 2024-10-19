package role

import "github.com/gofiber/fiber/v2"

type RoleHandler struct {
	roleService serviceRoleInterface
}

func NewRoleHandler(roleService serviceRoleInterface) *RoleHandler {
	return &RoleHandler{
		roleService,
	}
}

func (h *RoleHandler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	v1.Post("/role", h.CreateRole)
}

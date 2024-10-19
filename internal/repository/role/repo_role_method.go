package role

import (
	"github.com/gofiber/fiber/v2"
)

func (r *RoleRepository) CreateRole(ctx *fiber.Ctx, name string) (result string, err error) {
	return name, nil
}

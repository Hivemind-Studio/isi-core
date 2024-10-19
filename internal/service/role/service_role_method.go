package role

import (
	"github.com/gofiber/fiber/v2"
)

func (s *RoleService) CreateRole(ctx *fiber.Ctx, name string) (result string, err error) {
	return s.repoRole.CreateRole(ctx, name)
}

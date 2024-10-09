package role

import (
	"github.com/gofiber/fiber/v2"
)

type serviceRoleInterface interface {
	CreateRole(ctx *fiber.Ctx, name string) (res string, err error)
}

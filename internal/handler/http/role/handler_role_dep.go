package role

import (
	"github.com/gofiber/fiber/v2"
)

type CreateRoleUseCaseInterface interface {
	Execute(ctx *fiber.Ctx, name string) (res string, err error)
}

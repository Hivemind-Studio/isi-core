package role

import (
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/gofiber/fiber/v2"
)

type repoRoleInterface interface {
	dbtx.DBTXInterface
	CreateRole(ctx *fiber.Ctx, name string) (result string, err error)
}

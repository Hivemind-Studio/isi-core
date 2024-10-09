package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/gofiber/fiber/v2"
)

type repoUserInterface interface {
	dbtx.DBTXInterface
	GetTest(ctx *fiber.Ctx, id int) (result string, err error)
	Create(ctx *fiber.Ctx, body *user.RegisterDTO) (result *user.RegisterDTO, err error)
}

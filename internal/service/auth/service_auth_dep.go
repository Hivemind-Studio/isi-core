package auth

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/gofiber/fiber/v2"
)

type repoAuthInterface interface {
	dbtx.DBTXInterface
	Login(ctx *fiber.Ctx, body *user.LoginDTO) (result string, err error)
}

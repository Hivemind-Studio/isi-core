package auth

import (
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/gofiber/fiber/v2"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	FindByEmail(ctx *fiber.Ctx, email string) (user.User, error)
}

package auth

import (
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type repoAuthInterface interface {
	dbtx.DBTXInterface
	FindByEmail(ctx *fiber.Ctx, tx *sqlx.Tx, email string) (user.User, error)
}

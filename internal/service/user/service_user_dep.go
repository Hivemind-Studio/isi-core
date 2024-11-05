package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	Create(ctx *fiber.Ctx, tx *sqlx.Tx, body *user.RegistrationDTO, role int,
	) (result *user.RegisterResponse, err error)
}

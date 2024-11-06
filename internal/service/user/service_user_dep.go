package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"time"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	Create(ctx *fiber.Ctx, tx *sqlx.Tx, name string, email string, password string, roleId int64, phoneNumber string,
	) (err error)

	GetUsers(ctx *fiber.Ctx, name string, email string, startDate, endDate *time.Time,
		page int64, perPage int64,
	) ([]user.User, error)

	GetUserByID(ctx *fiber.Ctx, id int64) (user.User, error)
	//Update(ctx *fiber.Ctx, tx *sqlx.Tx, user user.Update, id string) (user.User, error)
}

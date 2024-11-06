package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/gofiber/fiber/v2"
	"time"
)

type serviceUserInterface interface {
	Create(ctx *fiber.Ctx, body *auth.RegistrationDTO) (result *auth.RegisterResponse, err error)
	GetUsers(ctx *fiber.Ctx, name string, email string, startDate, endDate *time.Time, page int64, perPage int64,
	) ([]user.UserDTO, error)
	GetUserByID(ctx *fiber.Ctx, id int64) (result *user.UserDTO, err error)
	//Update(ctx *fiber.Ctx)
}

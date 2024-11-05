package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/gofiber/fiber/v2"
	"time"
)

type serviceUserInterface interface {
	Create(ctx *fiber.Ctx, body *user.RegistrationDTO) (result *user.RegisterResponse, err error)
	GetUsers(ctx *fiber.Ctx, name string, email string, startDate, endDate *time.Time, page int64, perPage int64,
	) ([]user.UserDTO, error)
}

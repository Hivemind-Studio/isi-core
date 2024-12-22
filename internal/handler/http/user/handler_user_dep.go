package user

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"time"
)

type CreateUserUseCaseInterface interface {
	Execute(ctx context.Context, body *auth.RegistrationDTO) (result *auth.RegisterResponse, err error)
}

type GetUsersUseCaseInterface interface {
	Execute(ctx context.Context, name string, email string, startDate, endDate *time.Time, page int64, perPage int64) ([]user.UserDTO, error)
}

type GetUserByIDUseCaseInterface interface {
	Execute(ctx context.Context, id int64) (result *user.UserDTO, err error)
}

type UpdateUserStatusUseCaseInterface interface {
	Execute(ctx context.Context, ids []int64, status string) error
}

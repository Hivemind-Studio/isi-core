package profile

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
)

type GetProfileUserUseCaseInterface interface {
	Execute(ctx context.Context, id int64) (result *user.UserDTO, err error)
}

type UpdateProfileUserPasswordUseCaseInterface interface {
	Execute(ctx context.Context, id int64, currentPassword string, password string) (err error)
}

type UpdateProfileUserUseCaseInterface interface {
	Execute(ctx context.Context, id int64, name string, address string, gender string, phoneNumber string) (err error)
}

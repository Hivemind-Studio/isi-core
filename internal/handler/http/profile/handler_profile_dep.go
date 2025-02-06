package profile

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
)

type GetProfileUserUseCaseInterface interface {
	Execute(ctx context.Context, id int64, role string) (result *user.UserDTO, err error)
}

type UpdateProfileUserPasswordUseCaseInterface interface {
	Execute(ctx context.Context, id int64, currentPassword string, password string) (err error)
}

type UpdateProfileUserUseCaseInterface interface {
	Execute(ctx context.Context, id int64, role string, payload user.UpdateUserDTO) (result *user.UserDTO, err error)
}

type UpdatePhotoUseCaseInterface interface {
	Execute(ctx context.Context, id int64, photo string) error
}

type DeletePhotoUseCaseInterface interface {
	Execute(ctx context.Context, id int64) error
}

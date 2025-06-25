package profile

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
)

type GetProfileUserUseCaseInterface interface {
	GetProfileUser(ctx context.Context, id int64, role string) (result *user.UserDTO, err error)
}

type UpdateProfileUserPasswordUseCaseInterface interface {
	UpdateProfileUserPassword(ctx context.Context, id int64, currentPassword string, password string) (err error)
}

type UpdateProfileUserUseCaseInterface interface {
	UpdateProfileUser(ctx context.Context, id int64, role string, payload user.UpdateUserDTO) (result *user.UserDTO, err error)
}

type UpdatePhotoUseCaseInterface interface {
	updateCoachLevel(ctx context.Context, id int64, photo string) error
}

type DeletePhotoUseCaseInterface interface {
	DeletePhoto(ctx context.Context, id int64) error
}

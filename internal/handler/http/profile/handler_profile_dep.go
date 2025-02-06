package profile

import (
	"context"
	coachDto "github.com/Hivemind-Studio/isi-core/internal/dto/coach"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
)

type GetProfileUserUseCaseInterface interface {
	Execute(ctx context.Context, id int64) (result *user.UserDTO, err error)
}

type UpdateProfileUserPasswordUseCaseInterface interface {
	Execute(ctx context.Context, id int64, currentPassword string, password string) (err error)
}

type UpdateProfileUserUseCaseInterface interface {
	Execute(ctx context.Context, id int64, name string, address string, gender string,
		phoneNumber string, occupation string) (result *user.UserDTO, err error)
}

type UpdateProfileCoachUseCaseInterface interface {
	Execute(ctx context.Context, id int64, name string, address string, gender string,
		phoneNumber string, dateOfBirth string, title string, bio string, expertise string) (dto *coachDto.DTO, err error)
}

type UpdatePhotoUseCaseInterface interface {
	Execute(ctx context.Context, id int64, photo string) error
}

type DeletePhotoUseCaseInterface interface {
	Execute(ctx context.Context, id int64) error
}

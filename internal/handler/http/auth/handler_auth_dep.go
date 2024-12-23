package user

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
)

type LoginUseCaseInterface interface {
	Execute(ctx context.Context, body *auth.LoginDTO) (user dto.UserDTO, err error)
}

type SendVerificationUseCaseInterface interface {
	Execute(ctx context.Context, email string) error
}

type VerifyRegistrationTokenUseCaseInterface interface {
	Execute(ctx context.Context, registrationToken string, token string) (err error)
}

type CreateUserUseCaseInterface interface {
	Execute(ctx context.Context, body *auth.RegistrationDTO) (result *auth.RegisterResponse, err error)
}

type UpdateUserUseCaseInterface interface {
	Execute(ctx context.Context, password string, confirmPassword string, token string) (err error)
}

type UpdateCoachPasswordInterface interface {
	Execute(ctx context.Context, password string, confirmPassword string, token string) (err error)
}

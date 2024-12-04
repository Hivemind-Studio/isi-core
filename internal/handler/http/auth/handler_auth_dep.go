package user

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
)

type serviceAuthInterface interface {
	Login(ctx context.Context, body *auth.LoginDTO) (result dto.UserDTO, err error)
	SendEmailVerification(ctx context.Context, email string) error
	VerifyRegistrationToken(ctx context.Context, email string, token string) (err error)
}

type serviceUserInterface interface {
	CreateUser(ctx context.Context, body *auth.RegistrationDTO) (result *auth.RegisterResponse, err error)
}

type serviceCoachInterface interface {
	UpdateCoachPassword(ctx context.Context, password string, confirmPassword string, token string) (err error)
}

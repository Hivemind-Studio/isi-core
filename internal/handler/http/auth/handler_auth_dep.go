package user

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
)

type serviceAuthInterface interface {
	Login(ctx context.Context, body *auth.LoginDTO) (result dto.UserDTO, err error)
	SignUp(ctx context.Context, body *auth.SignUpDTO) (err error)
	SendEmailVerification(context context.Context, email string) (err error)
}

type serviceUserInterface interface {
	Create(ctx context.Context, body *auth.RegistrationDTO) (result *auth.RegisterResponse, err error)
}

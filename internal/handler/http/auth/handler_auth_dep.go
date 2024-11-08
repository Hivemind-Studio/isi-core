package user

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
)

type serviceAuthInterface interface {
	Login(context context.Context, body *auth.LoginDTO) (result dto.UserDTO, err error)
}

type serviceUserInterface interface {
	Create(context context.Context, body *auth.RegistrationDTO) (result *auth.RegisterResponse, err error)
}

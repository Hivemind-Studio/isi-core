package user

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/internal/dto/campaign"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/gofiber/fiber/v2"
)

type LoginUseCaseInterface interface {
	Login(ctx context.Context, body *auth.LoginDTO) (user dto.UserDTO, err error)
}

type SendVerificationUseCaseInterface interface {
	SendVerificationUseCase(ctx context.Context, email string) error
}

type VerifyRegistrationTokenUseCaseInterface interface {
	VerifyRegistrationToken(ctx context.Context, registrationToken string, token string) (err error)
}

type CreateUserUseCaseInterface interface {
	CreateUser(ctx context.Context, body *auth.RegistrationDTO) (result *auth.RegisterResponse, err error)
}

type UpdateCoachPasswordInterface interface {
	SendConfirmationChangeNewEmail(ctx context.Context, password string, confirmPassword string, token string) (err error)
}

type ForgotPasswordUseCaseInterface interface {
	SendVerificationUseCase(ctx context.Context, email string) (err error)
}

type GoogleLoginUseCaseInterface interface {
	GoogleLogin(*fiber.Ctx) string
}

type GoogleCallbackUseCaseInterface interface {
	GetUserDataFromGoogle(ctx context.Context, code string) (dto.UserDTO, error)
}

type CreateUserCampaign interface {
	CreateUser(ctx context.Context, payload campaign.UserCampaign) error
}

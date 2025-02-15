package user

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	loginUseCase                             LoginUseCaseInterface
	sendRegistrationEmailVerificationUseCase SendVerificationUseCaseInterface
	verifyRegistrationTokenUseCase           VerifyRegistrationTokenUseCaseInterface
	createUserUseCase                        CreateUserUseCaseInterface
	updateCoachPasswordUseCase               UpdateCoachPasswordInterface
	forgotPasswordUseCase                    ForgotPasswordUseCaseInterface
	googleLoginUseCase                       GoogleLoginUseCaseInterface
}

func NewAuthHandler(
	loginUseCase LoginUseCaseInterface,
	sendRegistrationEmailVerificationUseCase SendVerificationUseCaseInterface,
	verifyRegistrationTokenUseCase VerifyRegistrationTokenUseCaseInterface,
	createUserUseCase CreateUserUseCaseInterface,
	updateCoachPasswordUseCase UpdateCoachPasswordInterface,
	forgotPasswordUseCase ForgotPasswordUseCaseInterface,
	googleLoginUseCaseInterface GoogleLoginUseCaseInterface,
) *Handler {
	return &Handler{
		loginUseCase,
		sendRegistrationEmailVerificationUseCase,
		verifyRegistrationTokenUseCase,
		createUserUseCase,
		updateCoachPasswordUseCase,
		forgotPasswordUseCase,
		googleLoginUseCaseInterface,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1/auth")

	v1.Post("/login", h.Login)
	v1.Post("/register", h.Create)
	v1.Post("/email/verify", h.EmailVerification)
	v1.Patch("/password", h.PatchPassword)
	v1.Post("/password/recover", h.ForgotPassword)
	v1.Get("/google", h.GoogleLogin)
}

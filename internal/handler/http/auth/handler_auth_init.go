package user

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	loginUseCase                   LoginUseCaseInterface
	sendEmailVerificationUseCase   SendVerificationUseCaseInterface
	verifyRegistrationTokenUseCase VerifyRegistrationTokenUseCaseInterface
	createUserUseCase              CreateUserUseCaseInterface
	updateCoachPasswordUseCase     UpdateCoachPasswordInterface
	forgotPasswordUseCase          ForgotPasswordUseCaseInterface
}

func NewAuthHandler(
	loginUseCase LoginUseCaseInterface,
	sendVerificationUseCase SendVerificationUseCaseInterface,
	verifyRegistrationTokenUseCase VerifyRegistrationTokenUseCaseInterface,
	createUserUseCase CreateUserUseCaseInterface,
	updateCoachPasswordUseCase UpdateCoachPasswordInterface,
	forgotPasswordUseCase ForgotPasswordUseCaseInterface) *Handler {
	return &Handler{
		loginUseCase,
		sendVerificationUseCase,
		verifyRegistrationTokenUseCase,
		createUserUseCase,
		updateCoachPasswordUseCase,
		forgotPasswordUseCase,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1/auth")

	v1.Post("/login", h.Login)
	v1.Post("/register", h.Create)
	v1.Post("/email/verify", h.EmailVerification)
	v1.Patch("/password", h.UpdatePassword)
	v1.Post("/password/recover", h.ForgotPassword)
}

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
}

func NewAuthHandler(
	loginUseCase LoginUseCaseInterface,
	sendVerificationUseCase SendVerificationUseCaseInterface,
	verifyRegistrationTokenUseCase VerifyRegistrationTokenUseCaseInterface,
	createUserUseCase CreateUserUseCaseInterface,
	updateCoachPasswordUseCase UpdateCoachPasswordInterface) *Handler {
	return &Handler{
		loginUseCase,
		sendVerificationUseCase,
		verifyRegistrationTokenUseCase,
		createUserUseCase,
		updateCoachPasswordUseCase,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1/auth")

	v1.Post("/login", h.Login)
	v1.Post("/register", h.Create)
	v1.Post("/verify-email", h.SendEmailVerification)
	v1.Patch("/coach/password", h.PatchPassword)
}

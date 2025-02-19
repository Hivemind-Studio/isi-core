package user

import (
	"github.com/Hivemind-Studio/isi-core/pkg/session"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	sessionManager                           *session.SessionManager
	loginUseCase                             LoginUseCaseInterface
	sendRegistrationEmailVerificationUseCase SendVerificationUseCaseInterface
	verifyRegistrationTokenUseCase           VerifyRegistrationTokenUseCaseInterface
	createUserUseCase                        CreateUserUseCaseInterface
	updateCoachPasswordUseCase               UpdateCoachPasswordInterface
	forgotPasswordUseCase                    ForgotPasswordUseCaseInterface
	googleLoginUseCase                       GoogleLoginUseCaseInterface
	googleCallbackUseCase                    GoogleCallbackUseCaseInterface
}

func NewAuthHandler(
	sessionManager *session.SessionManager,
	loginUseCase LoginUseCaseInterface,
	sendRegistrationEmailVerificationUseCase SendVerificationUseCaseInterface,
	verifyRegistrationTokenUseCase VerifyRegistrationTokenUseCaseInterface,
	createUserUseCase CreateUserUseCaseInterface,
	updateCoachPasswordUseCase UpdateCoachPasswordInterface,
	forgotPasswordUseCase ForgotPasswordUseCaseInterface,
	googleLoginUseCase GoogleLoginUseCaseInterface,
	googleCallbackUseCase GoogleCallbackUseCaseInterface,
) *Handler {
	return &Handler{
		sessionManager:                           sessionManager,
		loginUseCase:                             loginUseCase,
		sendRegistrationEmailVerificationUseCase: sendRegistrationEmailVerificationUseCase,
		verifyRegistrationTokenUseCase:           verifyRegistrationTokenUseCase,
		createUserUseCase:                        createUserUseCase,
		updateCoachPasswordUseCase:               updateCoachPasswordUseCase,
		forgotPasswordUseCase:                    forgotPasswordUseCase,
		googleLoginUseCase:                       googleLoginUseCase,
		googleCallbackUseCase:                    googleCallbackUseCase,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App, sessionManager *session.SessionManager) {
	v1 := app.Group("/api/v1/auth")

	v1.Post("/login", h.Login)
	v1.Post("/register", h.Create)
	v1.Post("/email/verify", h.EmailVerification)
	v1.Patch("/password", h.PatchPassword)
	v1.Post("/password/recover", h.ForgotPassword)
	v1.Get("/google", h.GoogleLogin)
	v1.Get("/google/callback", h.GoogleCallback)
}

package user

import (
	"github.com/Hivemind-Studio/isi-core/pkg/session"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	createUserCampaign                       CreateUserCampaign
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
	createUserCampaign CreateUserCampaign,
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
		createUserCampaign:                       createUserCampaign,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App, sessionManager *session.SessionManager) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://dashboard.inspirasisatu.com", // ✅ Explicit frontend URL
		AllowCredentials: true,
		AllowMethods:     "GET,POST,OPTIONS",
		AllowHeaders:     "Content-Type, Authorization",
	}))

	v1 := app.Group("/api/v1/auth")

	v1.Post("/login", h.Login)
	v1.Post("/register", h.Create)
	v1.Post("/email/verify", h.EmailVerification)
	v1.Patch("/password", h.PatchPassword)
	v1.Post("/password/recover", h.ForgotPassword)
	v1.Get("/google", h.GoogleLogin)
	v1.Get("/google/callback", h.GoogleCallback)
}

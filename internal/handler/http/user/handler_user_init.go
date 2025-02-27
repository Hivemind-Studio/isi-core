package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/Hivemind-Studio/isi-core/pkg/session"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	createUserStaffUseCase                CreateUserStaffUseCaseInterface
	getUsersUseCase                       GetUsersUseCaseInterface
	getUserByIDUseCase                    GetUserByIDUseCaseInterface
	updateUserStatusUseCase               UpdateUserStatusUseCaseInterface
	updateUserRoleCase                    UpdateUserRoleUseCaseInterface
	updateUserEmailUseCase                UpdateUserEmailInterface
	sendChangeEmailVerificationUseCase    SendChangeEmailVerificationInterface
	sendConfirmationChangeNewEmailUseCase SendConfirmationChangeNewEmailInterface
}

func NewUserHandler(
	createUserStaffUseCase CreateUserStaffUseCaseInterface,
	getUsersUseCase GetUsersUseCaseInterface,
	getUserByIDUseCase GetUserByIDUseCaseInterface,
	updateUserStatusUseCase UpdateUserStatusUseCaseInterface,
	updateUserRoleCase UpdateUserRoleUseCaseInterface,
	updateUserEmailUseCase UpdateUserEmailInterface,
	sendChangeEmailVerificationUseCase SendChangeEmailVerificationInterface,
	sendConfirmationChangeNewEmailUseCase SendConfirmationChangeNewEmailInterface,
) *Handler {
	return &Handler{
		createUserStaffUseCase:                createUserStaffUseCase,
		getUsersUseCase:                       getUsersUseCase,
		getUserByIDUseCase:                    getUserByIDUseCase,
		updateUserStatusUseCase:               updateUserStatusUseCase,
		updateUserRoleCase:                    updateUserRoleCase,
		updateUserEmailUseCase:                updateUserEmailUseCase,
		sendChangeEmailVerificationUseCase:    sendChangeEmailVerificationUseCase,
		sendConfirmationChangeNewEmailUseCase: sendConfirmationChangeNewEmailUseCase,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App, sessionManager *session.SessionManager) {
	v1 := app.Group("/api/v1/users")

	accessControlRules := h.manageAccessControl()
	v1.Use(middleware.SessionAuthMiddleware(sessionManager, accessControlRules))

	v1.Post("/", h.Create)
	v1.Get("/", h.GetUsers)
	v1.Get("/:id", h.GetUserById)
	v1.Post("/email", h.SendChangeEmailVerification)
	v1.Post("/new-email", h.SendConfirmationNewEmail)
	v1.Patch("/email", h.ConfirmChangeNewEmail)
	v1.Patch("/status", h.UpdateStatusUser)
	v1.Patch("/:id/role", h.UpdateUserRole)
}

func (h *Handler) manageAccessControl() map[string]middleware.AccessControlRule {
	accessControlRules := map[string]middleware.AccessControlRule{
		"Admin": {
			Role: "Admin",
			AllowedMethod: map[string][]string{
				constant.V1 + "/users": {"GET", "POST", "DELETE", "PATCH"},
			},
		},
		"Staff": {
			Role: "Staff",
			AllowedMethod: map[string][]string{
				constant.V1 + "/users": {"GET"},
			},
		},
	}
	return accessControlRules
}

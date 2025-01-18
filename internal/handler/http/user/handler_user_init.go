package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	createUserStaffUseCase  CreateUserStaffUseCaseInterface
	getUsersUseCase         GetUsersUseCaseInterface
	getUserByIDUseCase      GetUserByIDUseCaseInterface
	updateUserStatusUseCase UpdateUserStatusUseCaseInterface
	updateUserRoleCase      UpdateUserRoleUseCaseInterface
}

func NewUserHandler(
	createUserStaffUseCase CreateUserStaffUseCaseInterface,
	getUsersUseCase GetUsersUseCaseInterface,
	getUserByIDUseCase GetUserByIDUseCaseInterface,
	updateUserStatusUseCase UpdateUserStatusUseCaseInterface,
	updateUserRoleCase UpdateUserRoleUseCaseInterface,
) *Handler {
	return &Handler{
		createUserStaffUseCase:  createUserStaffUseCase,
		getUsersUseCase:         getUsersUseCase,
		getUserByIDUseCase:      getUserByIDUseCase,
		updateUserStatusUseCase: updateUserStatusUseCase,
		updateUserRoleCase:      updateUserRoleCase,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	accessControlRules := h.manageAccessControl()
	v1.Use(middleware.JWTAuthMiddleware(accessControlRules))

	v1.Post("/users", h.Create)
	v1.Get("/users", h.GetUsers)
	v1.Get("/users/:id", h.GetUserById)
	v1.Patch("/users/:id/status", h.UpdateStatusUser)
	v1.Patch("/users/:id/role", h.UpdateUserRole)
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

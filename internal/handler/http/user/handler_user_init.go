package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	createUserUsecase       CreateUserUseCaseInterface
	getUsersUseCase         GetUsersUseCaseInterface
	getUserByIDUseCase      GetUserByIDUseCaseInterface
	updateUserStatusUseCase UpdateUserStatusUseCaseInterface
}

func NewUserHandler(
	createUserUsecase CreateUserUseCaseInterface,
	getUsersUseCase GetUsersUseCaseInterface,
	getUserByIDUseCase GetUserByIDUseCaseInterface,
	updateUserStatusUseCase UpdateUserStatusUseCaseInterface,
) *Handler {
	return &Handler{
		createUserUsecase:       createUserUsecase,
		getUsersUseCase:         getUsersUseCase,
		getUserByIDUseCase:      getUserByIDUseCase,
		updateUserStatusUseCase: updateUserStatusUseCase,
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
}

func (h *Handler) manageAccessControl() map[string]middleware.AccessControlRule {
	accessControlRules := map[string]middleware.AccessControlRule{
		"coachee": {
			Role: "Coachee",
			AllowedMethod: map[string][]string{
				constant.V1 + "/users": {"GET", "POST", "DELETE", "PATCH"},
			},
		},
	}
	return accessControlRules
}

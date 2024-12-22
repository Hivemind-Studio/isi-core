package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	getCoachUseCase    GetCoachesUseCaseInterface
	createCoachUseCase CreateCoachUseCaseInterface
}

func NewCoachHandler(
	getCoachUseCase GetCoachesUseCaseInterface,
	createCoachUseCase CreateCoachUseCaseInterface,
) *Handler {
	return &Handler{
		getCoachUseCase:    getCoachUseCase,
		createCoachUseCase: createCoachUseCase,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1/coach")

	accessControlRules := h.manageAccessControl()

	v1.Use(middleware.JWTAuthMiddleware(accessControlRules))

	v1.Get("/", h.GetCoaches)
	v1.Post("/", h.CreateCoach)
}

func (h *Handler) manageAccessControl() map[string]middleware.AccessControlRule {
	accessControlRules := map[string]middleware.AccessControlRule{
		"Admin": {
			Role: "Admin",
			AllowedMethod: map[string][]string{
				constant.V1 + "/coach": {"GET", "POST", "DELETE"},
			},
		},
	}
	return accessControlRules
}

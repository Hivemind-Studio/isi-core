package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	getCoacheesUseCase GetCoacheesUseCaseInterface
}

func NewCoacheeHandler(getCoacheesUseCase GetCoacheesUseCaseInterface) *Handler {
	return &Handler{
		getCoacheesUseCase: getCoacheesUseCase,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1/coachees")

	accessControlRules := h.manageAccessControl()
	v1.Use(middleware.JWTAuthMiddleware(accessControlRules))

	v1.Get("/", h.GetCoachees)
}

func (h *Handler) manageAccessControl() map[string]middleware.AccessControlRule {
	accessControlRules := map[string]middleware.AccessControlRule{
		"Staff": {
			Role: "Staff",
			AllowedMethod: map[string][]string{
				constant.V1 + "/users": {"GET", "POST", "DELETE", "PATCH"},
			},
		},
	}
	return accessControlRules
}

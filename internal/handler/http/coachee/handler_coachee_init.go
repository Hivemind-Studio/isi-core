package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	getCoacheesUseCase    GetCoacheesUseCaseInterface
	getCoacheeByIdUseCase GetCoacheeByIdUseCaseInterface
}

func NewCoacheeHandler(getCoacheesUseCase GetCoacheesUseCaseInterface,
	getCoacheeByIdUseCase GetCoacheeByIdUseCaseInterface) *Handler {
	return &Handler{
		getCoacheesUseCase,
		getCoacheeByIdUseCase,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1/coachees")

	//accessControlRules := h.manageAccessControl()
	//v1.Use(middleware.JWTAuthMiddleware(accessControlRules))

	v1.Get("/", h.GetCoachees)
	v1.Get("/:id", h.GetCoacheeById)
}

func (h *Handler) manageAccessControl() map[string]middleware.AccessControlRule {
	accessControlRules := map[string]middleware.AccessControlRule{
		"Admin": {
			Role: "Admin",
			AllowedMethod: map[string][]string{
				constant.V1 + "/coachees": {"GET", "POST", "DELETE", "PATCH"},
			},
		},
		"Staff": {
			Role: "Staff",
			AllowedMethod: map[string][]string{
				constant.V1 + "/coachees": {"GET", "POST", "DELETE", "PATCH"},
			},
		},
	}
	return accessControlRules
}

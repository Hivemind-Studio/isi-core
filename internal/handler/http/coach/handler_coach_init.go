package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	getCoachUseCase     GetCoachesUseCaseInterface
	createCoachUseCase  CreateCoachUseCaseInterface
	getCoachByIdUseCase GetCoachByIdUseCaseInterface
}

func NewCoachHandler(
	getCoachUseCase GetCoachesUseCaseInterface,
	createCoachUseCase CreateCoachUseCaseInterface,
	getCoachByIdUseCase GetCoachByIdUseCaseInterface,
) *Handler {
	return &Handler{
		getCoachUseCase:     getCoachUseCase,
		createCoachUseCase:  createCoachUseCase,
		getCoachByIdUseCase: getCoachByIdUseCase,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1/coaches")

	//accessControlRules := h.manageAccessControl()
	//v1.Use(middleware.JWTAuthMiddleware(accessControlRules))

	v1.Get("/", h.GetCoaches)
	v1.Get("/:id", h.GetCoachById)
	v1.Post("/", h.CreateCoach)
}

func (h *Handler) manageAccessControl() map[string]middleware.AccessControlRule {
	accessControlRules := map[string]middleware.AccessControlRule{
		"Admin": {
			Role: "Admin",
			AllowedMethod: map[string][]string{
				constant.V1 + "/coaches": {"GET", "POST", "DELETE", "PATCH"},
			},
		},
		"Staff": {
			Role: "Staff",
			AllowedMethod: map[string][]string{
				constant.V1 + "/coaches": {"GET", "POST", "DELETE", "PATCH"},
			},
		},
	}
	return accessControlRules
}

package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/Hivemind-Studio/isi-core/pkg/session"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	createCampaignUseCase CreateCampaignUseCaseInterface
}

func NewCampaignHandler(createCampaignUseCase CreateCampaignUseCaseInterface) *Handler {
	return &Handler{
		createCampaignUseCase,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App, sessionManager *session.SessionManager) {
	v1 := app.Group("/api/v1/campaign")

	accessControlRules := h.manageAccessControl()
	v1.Use(middleware.SessionAuthMiddleware(sessionManager, accessControlRules))

	v1.Get("/", h.Create)
}

func (h *Handler) manageAccessControl() map[string]middleware.AccessControlRule {
	accessControlRules := map[string]middleware.AccessControlRule{
		"Admin": {
			Role: "Admin",
			AllowedMethod: map[string][]string{
				constant.V1 + "/campaign": {"GET", "POST", "DELETE", "PATCH"},
			},
		},
		"Staff": {
			Role: "Staff",
			AllowedMethod: map[string][]string{
				constant.V1 + "/campaign": {"GET", "POST", "DELETE", "PATCH"},
			},
		},
	}
	return accessControlRules
}

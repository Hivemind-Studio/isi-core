package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/Hivemind-Studio/isi-core/pkg/session"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	createCampaignUseCase CreateCampaignUseCaseInterface
	getCampaignUseCase    GetCampaignUseCaseInterface
	updateStatusCampaign  UpdateStatusCampaignUseCaseInterface
	updateCampaign        UpdateCampaignUseCaseInterface
	getCampaignByID       GetCampaignByIDUseCaseInterface
}

func NewCampaignHandler(createCampaignUseCase CreateCampaignUseCaseInterface,
	getCampaignUseCase GetCampaignUseCaseInterface,
	updateStatusCampaign UpdateStatusCampaignUseCaseInterface,
	updateCampaign UpdateCampaignUseCaseInterface,
	getCampaignByID GetCampaignByIDUseCaseInterface) *Handler {
	return &Handler{
		createCampaignUseCase,
		getCampaignUseCase,
		updateStatusCampaign,
		updateCampaign,
		getCampaignByID,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App, sessionManager *session.SessionManager) {
	v1 := app.Group("/api/v1/campaign")
	accessControlRules := h.manageAccessControl()

	v1.Use(middleware.SessionAuthMiddleware(sessionManager, accessControlRules))

	v1.Get("/", h.Get)
	v1.Post("/", h.Create)
	v1.Get("/:id")
	v1.Put("/:id", h.Update)
	v1.Patch("/status", h.UpdateStatusCampaign)

}

func (h *Handler) manageAccessControl() map[string]middleware.AccessControlRule {
	accessControlRules := map[string]middleware.AccessControlRule{
		"Admin": {
			Role: "Admin",
			AllowedMethod: map[string][]string{
				constant.V1 + "/campaign": {"GET", "POST", "PUT", "DELETE", "PATCH"},
			},
		},
		"Marketing": {
			Role: "Marketing",
			AllowedMethod: map[string][]string{
				constant.V1 + "/campaign": {"GET", "POST", "PUT", "DELETE", "PATCH"},
			},
		},
	}
	return accessControlRules
}

package analytic

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/Hivemind-Studio/isi-core/pkg/session"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	getTotalRegistrant GetTotalRegistrantUseCase
}

func NewAnalyticHandler(getTotalRegistrant GetTotalRegistrantUseCase) *Handler {
	return &Handler{
		getTotalRegistrant: getTotalRegistrant,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App, sessionManager *session.SessionManager) {
	v1 := app.Group("/api/v1/analytics")
	accessControlRules := h.manageAccessControl()

	v1.Use(middleware.SessionAuthMiddleware(sessionManager, accessControlRules))

	v1.Get("/registrant", h.GetTotalRegistrant)

}

func (h *Handler) manageAccessControl() map[string]middleware.AccessControlRule {
	accessControlRules := map[string]middleware.AccessControlRule{
		"Admin": {
			Role: "Admin",
			AllowedMethod: map[string][]string{
				constant.V1 + "/analytics": {"GET", "POST", "PUT", "DELETE", "PATCH"},
			},
		},
		"Staff": {
			Role: "Staff",
			AllowedMethod: map[string][]string{
				constant.V1 + "/analytics": {"GET"},
			},
		},
		"Marketing": {
			Role: "Marketing",
			AllowedMethod: map[string][]string{
				constant.V1 + "/analytics": {"GET"},
			},
		},
	}
	return accessControlRules
}

package profile

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	getProfileUser GetProfileUserUseCaseInterface
}

func NewProfileHandler(getProfileUser GetProfileUserUseCaseInterface) *Handler {
	return &Handler{
		getProfileUser,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	accessControlRules := h.manageAccessControl()
	v1.Use(middleware.JWTAuthMiddleware(accessControlRules))

	v1.Get("/profile", h.GetProfileUser)
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
				constant.V1 + "/users": {"GET", "POST", "DELETE", "PATCH"},
			},
		},
	}
	return accessControlRules
}

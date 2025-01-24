package profile

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	getProfileUser        GetProfileUserUseCaseInterface
	updateProfilePassword UpdateProfileUserPasswordUseCaseInterface
}

func NewProfileHandler(getProfileUser GetProfileUserUseCaseInterface, updateProfilePassword UpdateProfileUserPasswordUseCaseInterface) *Handler {
	return &Handler{
		getProfileUser,
		updateProfilePassword,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	accessControlRules := h.manageAccessControl()
	v1.Use(middleware.JWTAuthMiddleware(accessControlRules))

	v1.Get("/profile", h.GetProfile)
	v1.Patch("/profile", h.UpdateProfilePassword)
}

func (h *Handler) manageAccessControl() map[string]middleware.AccessControlRule {
	accessControlRules := map[string]middleware.AccessControlRule{
		"Admin": {
			Role: "Admin",
			AllowedMethod: map[string][]string{
				constant.V1 + "/profile": {"GET", "PUT", "PATCH"},
			},
		},
		"Staff": {
			Role: "Staff",
			AllowedMethod: map[string][]string{
				constant.V1 + "/profile": {"GET", "PUT", "PATCH"},
			},
		},
		"Coach": {
			Role: "Coach",
			AllowedMethod: map[string][]string{
				constant.V1 + "/profile": {"GET", "PUT", "PATCH"},
			},
		},
		"Coachee": {
			Role: "Coachee",
			AllowedMethod: map[string][]string{
				constant.V1 + "/profile": {"GET", "PUT", "PATCH"},
			},
		},
	}
	return accessControlRules
}

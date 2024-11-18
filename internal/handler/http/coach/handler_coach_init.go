package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	coachService serviceCoachInterface
}

func NewCoachHandler(coachService serviceCoachInterface) *Handler {
	return &Handler{
		coachService: coachService,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	accessControlRules := map[string]middleware.AccessControlRule{
		"Admin": {
			Role: "Admin",
			AllowedMethod: map[string][]string{
				constant.V1 + "/coaches": {"GET", "POST", "DELETE"},
			},
		},
	}

	v1.Use(middleware.JWTAuthMiddleware(accessControlRules))

	v1.Post("/coaches", h.GetCoaches)
}

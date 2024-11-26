package user

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	userService serviceUserInterface
}

func NewUserHandler(userService serviceUserInterface) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	accessControlRules := map[string]middleware.AccessControlRule{
		"coachee": {
			Role: "Coachee",
			AllowedMethod: map[string][]string{
				constant.V1 + "/users": {"GET", "POST", "DELETE", "PATCH"},
			},
		},
	}

	v1.Use(middleware.JWTAuthMiddleware(accessControlRules))

	v1.Post("/users", h.Create)
	v1.Get("/users", h.GetUsers)
	v1.Get("/users/:id", h.GetUserById)
	v1.Patch("/users/status", h.UpdateStatusUser)
}

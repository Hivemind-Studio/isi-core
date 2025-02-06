package profile

import (
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	getProfileUser        GetProfileUserUseCaseInterface
	updateProfilePassword UpdateProfileUserPasswordUseCaseInterface
	updateProfile         UpdateProfileUserUseCaseInterface
	updatePhoto           UpdatePhotoUseCaseInterface
	deletePhoto           DeletePhotoUseCaseInterface
}

func NewProfileHandler(
	getProfileUser GetProfileUserUseCaseInterface,
	updateProfilePassword UpdateProfileUserPasswordUseCaseInterface,
	updateProfile UpdateProfileUserUseCaseInterface,
	updatePhoto UpdatePhotoUseCaseInterface,
	deletePhoto DeletePhotoUseCaseInterface) *Handler {
	return &Handler{
		getProfileUser,
		updateProfilePassword,
		updateProfile,
		updatePhoto,
		deletePhoto,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	accessControlRules := h.manageAccessControl()
	v1.Use(middleware.JWTAuthMiddleware(accessControlRules))

	v1.Get("/profile", h.GetProfile)
	v1.Put("/profile", h.UpdateProfile)
	v1.Patch("/profile", h.UpdateProfilePassword)
	v1.Patch("/profile/photo", h.UploadPhoto)
	v1.Delete("/profile/photo", h.DeletePhoto)
}

func (h *Handler) manageAccessControl() map[string]middleware.AccessControlRule {
	accessControlRules := map[string]middleware.AccessControlRule{
		"Admin": {
			Role: "Admin",
			AllowedMethod: map[string][]string{
				constant.V1 + "/profile": {"GET", "PUT", "PATCH", "DELETE"},
			},
		},
		"Staff": {
			Role: "Staff",
			AllowedMethod: map[string][]string{
				constant.V1 + "/profile": {"GET", "PUT", "PATCH", "DELETE"},
			},
		},
		"Coach": {
			Role: "Coach",
			AllowedMethod: map[string][]string{
				constant.V1 + "/profile": {"GET", "PUT", "PATCH", "DELETE"},
			},
		},
		"Coachee": {
			Role: "Coachee",
			AllowedMethod: map[string][]string{
				constant.V1 + "/profile": {"GET", "PUT", "PATCH", "DELETE"},
			},
		},
	}
	return accessControlRules
}

package userlogin

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/hash"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) Execute(ctx context.Context, body *auth.LoginDTO) (user dto.UserDTO, err error) {
	requestID := ctx.Value("request_id").(string)

	logger.Print("info", requestID, "user_login", "execute", "attempting login for email", body.Email)

	savedUser, err := uc.repoUser.FindByEmail(ctx, body.Email)
	if err != nil {
		logger.Print("error", requestID, "user_login", "execute", "failed to find user", err)
		return dto.UserDTO{}, err
	}

	if savedUser.Password == nil {
		logger.Print("error", requestID, "user_login", "execute", "user has not set password", body.Email)
		return dto.UserDTO{}, httperror.New(fiber.StatusBadRequest, "user has not set password")
	}

	isValidPassword, _ := hash.ComparePassword(body.Password, *savedUser.Password)
	if !isValidPassword {
		logger.Print("error", requestID, "user_login", "execute", "invalid password attempt", body.Email)
		return dto.UserDTO{}, httperror.New(fiber.StatusBadRequest, "invalid password")
	}

	logger.Print("info", requestID, "user_login", "execute", "login successful", body.Email)

	return dto.UserDTO{
		ID:    savedUser.ID,
		Name:  savedUser.Name,
		Role:  savedUser.RoleName,
		Email: savedUser.Email,
		Photo: savedUser.Photo,
	}, nil
}

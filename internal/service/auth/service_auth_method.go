package auth

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/hash"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

func (s *Service) Login(ctx *fiber.Ctx, body *auth.LoginDTO) (user dto.UserDTO, err error) {
	savedUser, err := s.repoAuth.FindByEmail(ctx, body.Email)
	if err != nil {
		return dto.UserDTO{}, err
	}

	isValidPassword, _ := hash.ComparePassword(body.Password, savedUser.Password)
	if !isValidPassword {
		return dto.UserDTO{}, httperror.New(fiber.StatusBadRequest, "invalid password")
	}
	if err != nil {
		return dto.UserDTO{}, err
	}

	return dto.UserDTO{
		Name:  savedUser.Name,
		Role:  savedUser.RoleName,
		Email: savedUser.Email,
		Photo: savedUser.Photo,
	}, nil
}

package auth

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/hash"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

func (s *Service) Login(ctx *fiber.Ctx, body *user.LoginDTO) (userId string, err error) {
	savedUser, err := s.repoAuth.FindByEmail(ctx, body.Email)
	if err != nil {
		return "", err
	}

	isValidPassword, _ := hash.ComparePassword(savedUser.Password, body.Password)
	if !isValidPassword {
		return "", httperror.New(fiber.StatusBadRequest, "invalid password")
	}
	if err != nil {
		return "", err
	}

	return savedUser.Email, nil
}

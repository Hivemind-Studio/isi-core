package auth

import (
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	userRepo "github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/hash"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

func (s *Service) Login(ctx *fiber.Ctx, body *user.LoginDTO) (user userRepo.Cookie, err error) {
	savedUser, err := s.repoAuth.FindByEmail(ctx, body.Email)
	if err != nil {
		return userRepo.Cookie{}, err
	}

	isValidPassword, _ := hash.ComparePassword(savedUser.Password, body.Password)
	if !isValidPassword {
		return userRepo.Cookie{}, httperror.New(fiber.StatusBadRequest, "invalid password")
	}
	if err != nil {
		return userRepo.Cookie{}, err
	}

	return savedUser, nil
}

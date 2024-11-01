package auth

import (
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
)

func (s *Service) Login(ctx *fiber.Ctx, body *user.LoginDTO) (userId string, err error) {
	err = s.repoAuth.StartTx()
	if err != nil {
		return userId, fmt.Errorf("failed to start transaction: %w", err)
	}

	tx, err := s.repoAuth.GetTx()
	defer dbtx.HandleRollback(tx)
	if err != nil {
		dbtx.HandleRollback(tx)
		return "", fmt.Errorf("failed to get transaction: %w", err)
	}

	savedUser, err := s.repoAuth.FindByEmail(ctx, tx, body)

	isValidPassword, _ := utils.ComparePassword(savedUser.Password, body.Password)
	if !isValidPassword {
		return "", fmt.Errorf("invalid password")
	}

	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return savedUser.Email, nil
}

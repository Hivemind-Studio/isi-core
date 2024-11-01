package user

import (
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/gofiber/fiber/v2"
)

func (s *Service) GetTest(ctx *fiber.Ctx, id int) (result string, err error) {
	return s.repoUser.GetTest(ctx, id)
}

func (s *Service) Create(ctx *fiber.Ctx, body *user.RegisterDTO) (result *user.RegisterResponse, err error) {
	err = s.repoUser.StartTx()
	if err != nil {
		return result, fmt.Errorf("failed to start transaction: %w", err)
	}

	tx, err := s.repoUser.GetTx()
	defer dbtx.HandleRollback(tx)
	if err != nil {
		dbtx.HandleRollback(tx)
		return result, fmt.Errorf("failed to get transaction: %w", err)
	}

	result, err = s.repoUser.Create(ctx, tx, body)
	if err != nil {
		dbtx.HandleRollback(tx)
		return result, fmt.Errorf("failed to get transaction: %w", err)
	}

	err = tx.Commit()

	if err != nil {
		return result, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, err
}

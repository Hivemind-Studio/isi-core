package user

import (
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/enum"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

func (s *Service) Create(ctx *fiber.Ctx, body *user.RegistrationDTO) (result *user.RegisterResponse, err error) {
	err = s.repoUser.StartTx()
	if err != nil {
		return result, httperror.New(fiber.StatusInternalServerError, "error when start transaction")
	}

	tx, err := s.repoUser.GetTx()
	defer dbtx.HandleRollback(tx)
	if err != nil {
		dbtx.HandleRollback(tx)
		return result, httperror.New(fiber.StatusInternalServerError, "error when get transaction")
	}

	result, err = s.repoUser.Create(ctx, tx, body, enum.CoacheeRoleId)
	if err != nil {
		dbtx.HandleRollback(tx)
		return result, err
	}

	err = tx.Commit()
	if err != nil {
		return result, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, err
}

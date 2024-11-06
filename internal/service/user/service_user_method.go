package user

import (
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/enum"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"time"
)

func (s *Service) Create(ctx *fiber.Ctx, body *user.RegistrationDTO) (result *user.RegisterResponse, err error) {
	tx, err := s.repoUser.StartTx()
	if err != nil {
		return nil, httperror.New(fiber.StatusInternalServerError, "error when start transaction")
	}
	defer dbtx.HandleRollback(tx)

	err = s.repoUser.Create(ctx, tx, body.Name, body.Email, body.Password, enum.CoacheeRoleId, body.PhoneNumber)
	if err != nil {
		dbtx.HandleRollback(tx)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &user.RegisterResponse{
		Name:  body.Name,
		Email: body.Email,
	}, err
}

func (s *Service) GetUsers(ctx *fiber.Ctx, name string, email string, startDate,
	endDate *time.Time, page int64, perPage int64,
) ([]user.UserDTO, error) {
	users, err := s.repoUser.GetUsers(ctx, name, email, startDate, endDate, page, perPage)
	if err != nil {
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve users")
	}
	userDTOs := user.ConvertUsersToDTOs(users)

	return userDTOs, nil
}

package user

import (
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/enum"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"time"
)

func (s *Service) Create(ctx *fiber.Ctx, body *auth.RegistrationDTO) (result *auth.RegisterResponse, err error) {
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

	return &auth.RegisterResponse{
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

func (s *Service) GetUserByID(ctx *fiber.Ctx, id int64) (*user.UserDTO, error) {
	res, err := s.repoUser.GetUserByID(ctx, id)

	if err != nil {
		return nil, httperror.New(fiber.StatusNotFound, "user not found")
	}

	userDto := user.ConvertUserToDTO(res)

	return &userDto, nil
}

func (s *Service) SuspendUsers(ctx *fiber.Ctx, ids []int64) error {
	tx, err := s.repoUser.StartTx()
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when start transaction")
	}
	defer dbtx.HandleRollback(tx)

	err = s.repoUser.SuspendUsers(ctx, tx, ids)
	if err != nil {
		dbtx.HandleRollback(tx)
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

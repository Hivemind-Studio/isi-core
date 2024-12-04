package user

import (
	"context"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/enum"
	user "github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"time"
)

func (s *Service) CreateUser(ctx context.Context, body *auth.RegistrationDTO) (result *auth.RegisterResponse, err error) {
	tx, err := s.repoUser.StartTx(ctx)
	requestId := ctx.Value("request_id").(string)
	logger.Print("info", requestId, "User service", "CreateUser", "function start", body)

	if err != nil {
		return nil, httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	_, err = s.repoUser.Create(ctx, tx, body.Name, body.Email, body.Password, enum.CoacheeRoleId, body.PhoneNumber, int(constant.ACTIVE))
	if err != nil {
		logger.Print("error", requestId, "User service", "CreateUser", err.Error(), body)
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

func (s *Service) GetUsers(ctx context.Context, name string, email string, startDate,
	endDate *time.Time, page int64, perPage int64,
) ([]dto.UserDTO, error) {
	params := dto.GetUsersDTO{Name: name,
		Email:     email,
		StartDate: startDate,
		EndDate:   endDate,
		Role:      nil,
	}
	users, err := s.repoUser.GetUsers(ctx, params, page, perPage)
	if err != nil {
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve users")
	}
	userDTOs := user.ConvertUsersToDTOs(users)

	return userDTOs, nil
}

func (s *Service) GetUserByID(ctx context.Context, id int64) (*dto.UserDTO, error) {
	res, err := s.repoUser.GetUserByID(ctx, id)
	if err != nil {
		return nil, httperror.New(fiber.StatusNotFound, "user not found")
	}

	userDto := user.ConvertUserToDTO(res)

	return &userDto, nil
}

func (s *Service) UpdateUserStatus(ctx context.Context, ids []int64, updatedStatus string) error {
	tx, err := s.repoUser.StartTx(ctx)
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	err = s.repoUser.UpdateUserStatus(ctx, tx, ids, updatedStatus)

	if err != nil {
		dbtx.HandleRollback(tx)
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

package createuser

import (
	"context"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/internal/enum"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (s *UseCase) Execute(ctx context.Context, body *auth.RegistrationDTO) (result *auth.RegisterResponse, err error) {
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

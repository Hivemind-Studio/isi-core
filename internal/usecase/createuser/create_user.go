package createuser

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
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

	u, err := s.repoUser.GetByVerificationToken(ctx, body.Token)
	if err != nil {
		return nil, httperror.New(fiber.StatusInternalServerError, "error when fetching email")
	}

	_, err = s.repoUser.Create(ctx, tx, user.CreateUserParams{
		Name:          body.Name,
		Email:         u.Email,
		Password:      &body.Password,
		RoleID:        constant.RoleIDCoachee,
		PhoneNumber:   &body.PhoneNumber,
		Gender:        body.Gender,
		Address:       body.Address,
		Status:        int(constant.ACTIVE),
		GoogleID:      nil,
		Photo:         nil,
		VerifiedEmail: false,
	})

	if err != nil {
		logger.Print("error", requestId, "User service", "CreateUser", err.Error(), body)
		dbtx.HandleRollback(tx)
		return nil, err
	}

	err = s.repoUser.DeleteEmailTokenVerificationByTokenAndType(ctx, tx, body.Token, constant.REGISTER)
	if err != nil {
		logger.Print("error", requestId, "User service", "CreateUser",
			"failed to delete verification token with error :"+err.Error(), body)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, httperror.New(fiber.StatusInternalServerError, "error when committing transaction")
	}

	return &auth.RegisterResponse{
		Name:  body.Name,
		Email: u.Email,
	}, err
}

package updateprofile

import (
	"context"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) Execute(ctx context.Context, id int64, name string, address string, gender string,
	phoneNumber string, occupation string) (*dto.UserDTO, error) {
	tx, err := uc.repoUser.StartTx(ctx)
	requestId := ctx.Value("request_id").(string)
	logger.Print("info", requestId, "Profile service", "UpdateProfile", "function start", name)

	if err != nil {
		return nil, httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	findUser, err := uc.repoUser.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res, err := uc.repoUser.UpdateUser(ctx, tx, id, name, address, gender, phoneNumber, occupation, findUser.Version)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, httperror.New(fiber.StatusInternalServerError, "Failed to update password")
	}

	result := user.ConvertUserToDTO(*res)

	return &result, nil

}

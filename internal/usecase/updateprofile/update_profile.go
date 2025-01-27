package updateprofile

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) Execute(ctx context.Context, id int64, name string, address string, gender string, phoneNumber string) (err error) {
	tx, err := uc.repoUser.StartTx(ctx)
	requestId := ctx.Value("request_id").(string)
	logger.Print("info", requestId, "Profile service", "UpdateProfile", "function start", name)

	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	user, err := uc.repoUser.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	err = uc.repoUser.UpdateUser(ctx, tx, id, name, address, gender, phoneNumber, user.Version)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "Failed to update password")
	}

	return nil

}

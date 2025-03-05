package updateprofilepassword

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/hash"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) Execute(ctx context.Context, id int64, currentPassword string, password string) (err error) {
	tx, err := uc.repoUser.StartTx(ctx)
	requestId := ctx.Value("request_id").(string)
	logger.Print("info", requestId, "Profile service", "UpdatePassword", "function start", id)

	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	userProfile, err := uc.repoUser.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	current, _ := hash.HashPassword(currentPassword)

	if current != *userProfile.Password {
		return httperror.New(fiber.StatusBadRequest, "password mismatch")
	}

	err = uc.repoUser.UpdatePassword(ctx, tx, password, userProfile.Email, userProfile.Version)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "Failed to update password")
	}

	return nil
}

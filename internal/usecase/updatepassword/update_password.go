package updatepassword

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) SendConfirmationChangeNewEmail(ctx context.Context, password string, confirmPassword string, token string) (err error) {
	tx, err := uc.repoUser.StartTx(ctx)
	requestId := ctx.Value("request_id").(string)
	logger.Print("info", requestId, "User service", "UpdatePassword", "function start", token)

	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	if password != confirmPassword {
		return httperror.New(fiber.StatusBadRequest, "password mismatch")
	}

	email, err := uc.repoUser.GetTokenEmailVerification(token)
	if err != nil {
		return err
	}

	user, err := uc.repoUser.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	err = uc.repoUser.UpdatePassword(ctx, tx, password, email, user.Version)

	if err != nil {
		return err
	}

	err = uc.repoUser.DeleteEmailTokenVerificationByToken(ctx, tx, token)
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "Failed to create user")
	}

	err = tx.Commit()
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "Failed to create user")
	}

	return nil

}

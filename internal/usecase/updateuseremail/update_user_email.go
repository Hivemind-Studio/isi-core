package updateuseremail

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) Execute(ctx context.Context, token string, newEmail string, oldEmail string) (err error) {
	tx, err := uc.repoUser.StartTx(ctx)
	defer dbtx.HandleRollback(tx)

	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}

	_, err = uc.repoUser.GetTokenEmailVerificationWithType(ctx, token, constant.CONFIRM_TO_CHANGED_EMAIL_UPDATE, newEmail)
	if err != nil {
		return httperror.New(fiber.StatusUnauthorized, "token is not valid")
	}

	err = uc.repoUser.UpdateUserEmail(ctx, tx, newEmail, oldEmail)
	if err != nil {
		return err
	}

	err = uc.repoUser.DeleteEmailTokenVerificationByToken(ctx, tx, token, constant.CONFIRM_TO_CHANGED_EMAIL_UPDATE)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "Failed to update password")
	}

	return nil
}

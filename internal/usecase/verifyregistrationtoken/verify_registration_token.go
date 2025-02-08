package verifyregistrationtoken

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"time"
)

func (uc *UseCase) Execute(ctx context.Context, email string, token string) (err error) {
	tx, err := uc.repoUser.StartTx(ctx)
	if err != nil {
		return err
	}
	defer dbtx.HandleRollback(tx)

	emailVerification, err := uc.repoUser.GetByVerificationToken(ctx, token)
	if err != nil {
		return err
	}
	if emailVerification == nil {
		return httperror.New(fiber.StatusBadRequest, "invalid token")
	}
	if time.Now().After(emailVerification.ExpiredAt) {
		return httperror.New(fiber.StatusBadRequest, "verification token has expired")
	}
	return nil
}

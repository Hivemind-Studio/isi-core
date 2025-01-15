package updatecoachpassword

import (
	"context"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) Execute(ctx context.Context, password string, confirmPassword string, token string) (err error) {
	fmt.Print("Masuk sini crot")
	tx, err := uc.repoCoach.StartTx(ctx)
	requestId := ctx.Value("request_id").(string)
	logger.Print("info", requestId, "Coach service", "UpdateCoachPassword", "function start", token)

	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	if password != confirmPassword {
		return httperror.New(fiber.StatusBadRequest, "password mismatch")
	}

	email, err := uc.repoCoach.GetTokenEmailVerification(token)
	if err != nil {
		return err
	}

	_, err = uc.repoUser.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	err = uc.repoCoach.UpdateCoachPassword(ctx, tx, password, email)

	if err != nil {
		return err
	}

	err = uc.repoUser.DeleteEmailTokenVerification(ctx, tx, email)
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "Failed to create coach")
	}

	err = tx.Commit()
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "Failed to create coach")
	}

	return nil

}

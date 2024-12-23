package sendverification

import (
	"context"
	"errors"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

func (uc *UseCase) Execute(ctx context.Context, email string) error {
	if valid := uc.validateEmail(ctx, email); !valid {
		return httperror.New(fiber.StatusBadRequest, "email already exists")
	}

	trial, err := uc.repoUser.GetEmailVerificationTrialRequestByDate(ctx, email, time.Now())
	if err != nil {
		return err
	}
	if *trial >= 2 {
		return httperror.New(fiber.StatusTooManyRequests, "email verification limit reached for today")
	}

	token, err := uc.handleTokenGeneration(ctx, email, *trial)
	if err != nil {
		return err
	}

	if err := uc.emailVerification(email, token, email); err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send email verification")
	}

	return nil
}

func (uc *UseCase) validateEmail(ctx context.Context, email string) bool {
	_, err := uc.repoUser.FindByEmail(ctx, email)
	if err != nil {
		var customErr *httperror.CustomError
		if !errors.As(err, &customErr) {
			return false
		}
		if customErr.Code == fiber.StatusNotFound {
			return true
		}
	}

	return false
}

func (uc *UseCase) handleTokenGeneration(ctx context.Context, email string, trial int8) (string, error) {
	tx, err := uc.repoUser.StartTx(ctx)
	if err != nil {
		return "", httperror.Wrap(fiber.StatusInternalServerError, err, "failed to start transaction")
	}
	defer dbtx.HandleRollback(tx)

	token, err := uc.generateAndSaveToken(ctx, tx, email, trial)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", httperror.Wrap(fiber.StatusInternalServerError, err, "failed to commit transaction")
	}

	return token, nil
}

func (uc *UseCase) generateAndSaveToken(ctx context.Context, tx *sqlx.Tx, email string, trial int8) (string, error) {
	token := utils.GenerateVerificationToken()
	expiredAt := time.Now().Add(1 * time.Hour)
	currentDate := time.Now().Format("2006-01-02")

	if trial == 0 {
		if err := uc.repoUser.InsertEmailVerificationTrial(ctx, tx, email, token, expiredAt); err != nil {
			return "", httperror.Wrap(fiber.StatusInternalServerError, err, "failed to insert verification record")
		}
	} else {
		if err := uc.repoUser.UpdateEmailVerificationTrial(ctx, tx, email, currentDate, token, expiredAt); err != nil {
			return "", httperror.Wrap(fiber.StatusInternalServerError, err, "failed to update verification record")
		}
	}

	return token, nil
}

func (uc *UseCase) emailVerification(name string, token string, email string) error {
	emailData := struct {
		Name            string
		VerificationURL string
		Year            int
	}{
		Name:            name,
		VerificationURL: fmt.Sprintf("%stoken=%s", os.Getenv("CALLBACK_VERIFICATION_URL"), token),
		Year:            time.Now().Year(),
	}

	err := uc.emailClient.SendMail(
		[]string{email},
		"Inspirasi Satu - Verify Your Email",
		"template/verification_email.html",
		emailData,
	)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send verification email")
	}

	return nil
}

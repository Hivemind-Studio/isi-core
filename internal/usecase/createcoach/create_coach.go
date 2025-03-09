package createcoach

import (
	"context"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/internal/dto/coach"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

func (uc *UseCase) Execute(ctx context.Context, payload coach.CreateCoachDTO) (err error) {
	tx, err := uc.repoCoach.StartTx(ctx)
	requestId := ctx.Value("request_id").(string)
	logger.Print("info", requestId, "Coach service", "CreateCoach", "function start", payload)

	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	generatePassword := time.Now().String()

	userId, err := uc.repoUser.Create(ctx, tx, payload.Name, payload.Email, &generatePassword, constant.RoleIDCoach, payload.PhoneNumber, payload.Gender, payload.Address, int(constant.PENDING), nil, nil, false)

	if err != nil {
		logger.Print("error", requestId, "Coach service", "CreateCoach", err.Error(), payload)
		dbtx.HandleRollback(tx)
		return err
	}

	err = uc.repoCoach.CreateCoach(ctx, tx, userId)
	if err != nil {
		logger.Print("error", requestId, "Coach service", "CreateCoach", err.Error(), payload)
		dbtx.HandleRollback(tx)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "Failed to create coach")
	}

	err = uc.sendEmailVerification(ctx, payload.Name, payload.Email)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send email verification")
	}

	return nil
}

func (uc *UseCase) sendEmailVerification(ctx context.Context, name string, email string) error {
	trial, err := uc.repoCoach.GetEmailVerificationTrialRequestByDate(ctx, email, time.Now())
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

	if err := uc.emailVerification(name, token, email); err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send email verification")
	}

	return nil
}

func (uc *UseCase) handleTokenGeneration(ctx context.Context, email string, trial int8) (string, error) {
	tx, err := uc.repoCoach.StartTx(ctx)
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
		if err := uc.repoCoach.InsertEmailVerificationTrial(ctx, tx, email, token, expiredAt); err != nil {
			return "", httperror.Wrap(fiber.StatusInternalServerError, err, "failed to insert verification record")
		}
	} else {
		if err := uc.repoCoach.UpdateEmailVerificationTrial(ctx, tx, email, currentDate, token, expiredAt); err != nil {
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
		VerificationURL: fmt.Sprintf("%sjoin?token=%s", os.Getenv("CALLBACK_VERIFICATION_URL"), token),
		Year:            time.Now().Year(),
	}

	err := uc.emailClient.SendMail(
		[]string{email},
		"Inspirasi Satu - Email Verification",
		"template/verification_email.html",
		emailData,
	)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send verification email")
	}

	return nil
}

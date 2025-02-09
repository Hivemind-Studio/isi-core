package useremail

import (
	"context"
	"errors"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"os"
	"strconv"
	"time"
)

func (s *Service) ValidateEmail(ctx context.Context, email string) bool {
	_, err := s.repoUser.FindByEmail(ctx, email)
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

func (s *Service) ValidateTrialByDate(ctx context.Context, email string, tokenType string) (*int8, error) {
	emailLimit := os.Getenv("EMAIL_LIMIT_NUM")
	limit, err := strconv.ParseInt(emailLimit, 10, 64)
	trial, err := s.repoUser.GetEmailVerificationTrialRequestByDate(ctx, email, time.Now(), tokenType)
	if err != nil {
		return trial, err
	}
	if *trial >= int8(limit) {
		return trial, httperror.New(fiber.StatusTooManyRequests, "user email verification limit reached for today")
	}

	return trial, nil
}

func (s *Service) HandleTokenGeneration(ctx context.Context, email string, trial int8, tokenType string) (string, error) {
	tx, err := s.repoUser.StartTx(ctx)
	if err != nil {
		return "", httperror.Wrap(fiber.StatusInternalServerError, err, "failed to start transaction")
	}
	defer dbtx.HandleRollback(tx)

	token, err := s.generateAndSaveToken(ctx, tx, email, trial, tokenType)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", httperror.Wrap(fiber.StatusInternalServerError, err, "failed to commit transaction")
	}

	return token, nil
}

func (s *Service) generateAndSaveToken(ctx context.Context, tx *sqlx.Tx, email string, trial int8, tokenType string) (string, error) {
	token := utils.GenerateVerificationToken()
	expiredAt := time.Now().Add(1 * time.Hour)
	currentDate := time.Now().Format("2006-01-02")
	if trial == 0 {
		if err := s.repoUser.InsertEmailVerificationTrial(ctx, tx, email, token, expiredAt, tokenType); err != nil {
			return "", httperror.Wrap(fiber.StatusInternalServerError, err, "failed to insert verification record")
		}
	} else {
		existingEmailVerification, err := s.repoUser.GetByEmail(ctx, email, tokenType)
		if err != nil {
			return "", err
		}
		if err := s.repoUser.UpdateEmailVerificationTrial(ctx, tx, email, currentDate, token,
			expiredAt, existingEmailVerification.Version, tokenType); err != nil {
			return "", httperror.Wrap(fiber.StatusInternalServerError, err, "failed to update verification record")
		}
	}

	return token, nil
}

func (s *Service) SendEmail(recipients []string, subject, templatePath string, emailData interface{}) error {
	err := s.emailClient.SendMail(recipients, subject, templatePath, emailData)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send email")
	}
	return nil
}

func (s *Service) GenerateTokenVerification(ctx context.Context, email string, tokenType string, validateExistingEmail bool) (string, error) {
	if validateExistingEmail {
		if valid := s.ValidateEmail(ctx, email); !valid {
			return "", httperror.New(fiber.StatusBadRequest, "user email already exists")
		}
	}

	trial, err := s.ValidateTrialByDate(ctx, email, tokenType)
	if err != nil {
		return "", err
	}

	token, err := s.HandleTokenGeneration(ctx, email, *trial, tokenType)
	if err != nil {
		return "", err
	}
	return token, nil
}

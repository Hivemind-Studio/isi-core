package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/hash"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

func (s *Service) Login(ctx context.Context, body *auth.LoginDTO) (user dto.UserDTO, err error) {
	savedUser, err := s.repoUser.FindByEmail(ctx, body.Email)
	if err != nil {
		return dto.UserDTO{}, err
	}

	isValidPassword, _ := hash.ComparePassword(body.Password, savedUser.Password)
	if !isValidPassword {
		return dto.UserDTO{}, httperror.New(fiber.StatusBadRequest, "invalid password")
	}
	if err != nil {
		return dto.UserDTO{}, err
	}

	return dto.UserDTO{
		Name:  savedUser.Name,
		Role:  savedUser.RoleName,
		Email: savedUser.Email,
		Photo: savedUser.Photo,
	}, nil
}

func (s *Service) SignUp(ctx context.Context, body *auth.SignUpDTO) (err error) {
	_, err = s.repoUser.FindByEmail(ctx, body.Email)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SendEmailVerification(ctx context.Context, email string) error {
	if err := s.validateEmail(ctx, email); err != nil {
		return err
	}

	trial, err := s.repoUser.GetEmailVerificationTrialRequestByDate(ctx, email, time.Now())
	if err != nil {
		return err
	}
	if *trial >= 2 {
		return httperror.New(fiber.StatusTooManyRequests, "email verification limit reached for today")
	}

	token, err := s.handleTokenGeneration(ctx, email, *trial)
	if err != nil {
		return err
	}

	if err := s.emailVerification(email, token, email); err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send email verification")
	}

	return nil
}

func (s *Service) validateEmail(ctx context.Context, email string) error {
	_, err := s.repoUser.FindByEmail(ctx, email)
	if err != nil {
		var customErr *httperror.CustomError
		if !errors.As(err, &customErr) {
			return err
		}
		if customErr.Code != fiber.StatusNotFound {
			return err
		}
	}
	return nil
}

func (s *Service) handleTokenGeneration(ctx context.Context, email string, trial int8) (string, error) {
	tx, err := s.repoUser.StartTx(ctx)
	if err != nil {
		return "", httperror.Wrap(fiber.StatusInternalServerError, err, "failed to start transaction")
	}
	defer dbtx.HandleRollback(tx)

	token, err := s.generateAndSaveToken(ctx, tx, email, trial)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", httperror.Wrap(fiber.StatusInternalServerError, err, "failed to commit transaction")
	}

	return token, nil
}

func (s *Service) generateAndSaveToken(ctx context.Context, tx *sqlx.Tx, email string, trial int8) (string, error) {
	token := utils.GenerateVerificationToken()
	expiredAt := time.Now().Add(1 * time.Hour)
	currentDate := time.Now().Format("2006-01-02")

	if trial == 0 {
		if err := s.repoUser.InsertEmailVerificationTrial(ctx, tx, email, token, expiredAt); err != nil {
			return "", httperror.Wrap(fiber.StatusInternalServerError, err, "failed to insert verification record")
		}
	} else {
		if err := s.repoUser.UpdateEmailVerificationTrial(ctx, tx, email, currentDate, token, expiredAt); err != nil {
			return "", httperror.Wrap(fiber.StatusInternalServerError, err, "failed to update verification record")
		}
	}

	return token, nil
}

func (s *Service) emailVerification(name string, token string, email string) error {
	emailData := struct {
		Name            string
		VerificationURL string
		Year            int
	}{
		Name:            name,
		VerificationURL: fmt.Sprintf("%stoken=%s", os.Getenv("CALLBACK_VERIFICATION_URL"), token),
		Year:            time.Now().Year(),
	}

	err := s.emailClient.SendMail(
		[]string{email},
		"Verify Your Email",
		"template/verification_email.html",
		emailData,
	)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send verification email")
	}

	return nil
}

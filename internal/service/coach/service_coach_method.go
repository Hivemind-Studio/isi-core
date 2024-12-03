package coach

import (
	"context"
	"errors"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/internal/dto/coach"
	userdto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

func (s *Service) GetCoaches(ctx context.Context, name string, email string, startDate,
	endDate *time.Time, page int64, perPage int64,
) ([]userdto.UserDTO, error) {
	coachRoleId := constant.RoleIDCoach
	params := userdto.GetUsersDTO{Name: name,
		Email:     email,
		StartDate: startDate,
		EndDate:   endDate,
		Role:      &coachRoleId,
	}
	users, err := s.repoCoach.GetUsers(ctx, params, page, perPage)
	if err != nil {
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve users")
	}
	userDTOs := user.ConvertUsersToDTOs(users)

	return userDTOs, nil
}

func (s *Service) CreateCoach(ctx context.Context, payload coach.CreateCoachDTO) (err error) {
	tx, err := s.repoCoach.StartTx(ctx)
	requestId := ctx.Value("request_id").(string)
	logger.Print("info", requestId, "Coach service", "CreateCoach", "function start", payload)

	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	err = s.repoCoach.CreateCoach(ctx, tx, payload.Name, payload.Email, payload.PhoneNumber, payload.Gender, payload.Address)
	if err != nil {
		logger.Print("error", requestId, "Coach service", "CreateCoach", err.Error(), payload)
		dbtx.HandleRollback(tx)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "Failed to create coach")
	}

	err = s.SendEmailVerification(ctx, payload.Email)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send email verification")
	}

	return nil
}

func (s *Service) SendEmailVerification(ctx context.Context, email string) error {
	//if valid := s.validateEmail(ctx, email); !valid {
	//	return httperror.New(fiber.StatusBadRequest, "email already exists")
	//}

	trial, err := s.repoCoach.GetEmailVerificationTrialRequestByDate(ctx, email, time.Now())
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

func (s *Service) validateEmail(ctx context.Context, email string) bool {
	_, err := s.repoCoach.FindByEmail(ctx, email)
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

func (s *Service) handleTokenGeneration(ctx context.Context, email string, trial int8) (string, error) {
	tx, err := s.repoCoach.StartTx(ctx)
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
		if err := s.repoCoach.InsertEmailVerificationTrial(ctx, tx, email, token, expiredAt); err != nil {
			return "", httperror.Wrap(fiber.StatusInternalServerError, err, "failed to insert verification record")
		}
	} else {
		if err := s.repoCoach.UpdateEmailVerificationTrial(ctx, tx, email, currentDate, token, expiredAt); err != nil {
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
		VerificationURL: fmt.Sprintf("%scoach/token=%s", os.Getenv("CALLBACK_VERIFICATION_URL"), token),
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

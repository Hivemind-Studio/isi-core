package createstaff

import (
	"context"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/dto/auth"
	user2 "github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"os"
	"time"
)

func (uc *UseCase) Execute(ctx context.Context, body auth.RegistrationStaffDTO) (err error) {
	tx, err := uc.repoUser.StartTx(ctx)
	requestId := ctx.Value("request_id").(string)
	logger.Print("info", requestId, "User service", "CreateUser", "function start", body)

	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	generatePassword := time.Now().String()
	user := user2.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: generatePassword,
		Gender:   body.Gender,
		Address:  body.Address,
	}
	_, err = uc.repoUser.CreateStaff(ctx, tx, user)

	if err != nil {
		logger.Print("error", requestId, "User service", "CreateUser", err.Error(), body)
		dbtx.HandleRollback(tx)
		return nil
	}

	err = tx.Commit()
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "Transaction commit failed")
	}

	// sendEmail
	err = uc.sendEmailVerification(ctx, body.Name, body.Email)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send useremail verification")
	}

	return nil
}

func (uc *UseCase) sendEmailVerification(ctx context.Context, name string, email string) error {
	trial, err := uc.userEmailService.GetEmailVerificationTrialRequestByDate(ctx, email, time.Now())
	if err != nil {
		return err
	}
	if *trial >= 2 {
		return httperror.New(fiber.StatusTooManyRequests, "useremail verification limit reached for today")
	}

	token, err := uc.userEmailService.HandleTokenGeneration(ctx, email, *trial)
	if err != nil {
		return err
	}

	if err := uc.emailVerification(name, token, email); err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send useremail verification")
	}

	return nil
}

func (uc *UseCase) emailVerification(name string, token string, email string) error {
	emailData := struct {
		Name            string
		VerificationURL string
		Year            int
	}{
		Name:            name,
		VerificationURL: fmt.Sprintf("%s/register/password/token=%s", os.Getenv("CALLBACK_VERIFICATION_URL"), token),
		Year:            time.Now().Year(),
	}

	err := uc.emailClient.SendMail(
		[]string{email},
		"Inspirasi Satu - Verify Your Email",
		"template/verification_email.html",
		emailData,
	)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send verification user email")
	}

	return nil
}

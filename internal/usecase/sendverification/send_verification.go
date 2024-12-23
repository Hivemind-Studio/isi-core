package sendverification

import (
	"context"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"os"
	"time"
)

func (uc *UseCase) Execute(ctx context.Context, email string) error {
	if valid := uc.userEmailService.ValidateEmail(ctx, email); !valid {
		return httperror.New(fiber.StatusBadRequest, "useremail already exists")
	}

	trial, err := uc.repoUser.GetEmailVerificationTrialRequestByDate(ctx, email, time.Now())
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

	if err := uc.emailVerification(email, token, email); err != nil {
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
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send verification user email")
	}

	return nil
}

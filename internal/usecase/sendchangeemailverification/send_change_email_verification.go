package sendchangeemailverification

import (
	"context"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"os"
	"time"
)

func (uc *UseCase) Execute(ctx context.Context, email string) error {
	//if valid := uc.userEmailService.ValidateEmail(ctx, email); !valid {
	//	return httperror.New(fiber.StatusBadRequest, "user email already exists")
	//}

	trial, err := uc.userEmailService.ValidateTrialByDate(ctx, email)
	if err != nil {
		return err
	}

	token, err := uc.userEmailService.HandleTokenGeneration(ctx, email, *trial, constant.EMAIL_UPDATE)
	if err != nil {
		return err
	}

	if err := uc.emailVerification(email, token, email); err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send user email verification")
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

	err := uc.userEmailService.SendEmail([]string{email},
		"Inspirasi Satu - Change Email Verification",
		"template/change_email.html",
		emailData,
	)

	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send verification user email")
	}

	return nil
}

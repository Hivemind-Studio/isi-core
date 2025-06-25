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

func (uc *UseCase) SendVerificationUseCase(ctx context.Context, email string) error {
	token, err := uc.userEmailService.GenerateTokenVerification(ctx, email, constant.EMAIL_UPDATE, false)
	if err != nil {
		return err
	}

	if err := uc.sendEmailVerification(email, token, email); err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send user email verification")
	}

	return nil
}

func (uc *UseCase) sendEmailVerification(name string, token string, email string) error {
	emailData := struct {
		Name            string
		VerificationURL string
		Year            int
	}{
		Name: name,
		VerificationURL: fmt.Sprintf("%s%stoken=%s", os.Getenv("CALLBACK_VERIFICATION_URL"),
			"change-email?", token),
		Year: time.Now().Year(),
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

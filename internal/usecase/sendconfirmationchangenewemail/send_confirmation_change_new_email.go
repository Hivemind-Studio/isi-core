package sendconfirmationchangenewemail

import (
	"context"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"os"
	"time"
)

func (uc *UseCase) Execute(ctx context.Context, token string, newEmail string, oldEmail string) (err error) {
	tx, err := uc.repoUser.StartTx(ctx)
	defer dbtx.HandleRollback(tx)

	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}

	_, err = uc.repoUser.GetTokenEmailVerificationWithType(ctx, token, constant.EMAIL_UPDATE, oldEmail)
	if err != nil {
		return err
	}

	err = uc.repoUser.DeleteEmailTokenVerificationByTokenAndType(ctx, tx, token, constant.EMAIL_UPDATE)
	if err != nil {
		return err
	}

	newToken, err := uc.emailService.GenerateTokenVerification(ctx, newEmail, constant.CONFIRM_TO_CHANGED_EMAIL_UPDATE, true)
	if err != nil {
		return err
	}

	err = uc.sendEmailVerification(newEmail, newToken, newEmail)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "Failed to update password")
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
			"new-email-confirm?", token),
		Year: time.Now().Year(),
	}

	err := uc.emailService.SendEmail([]string{email},
		"Inspirasi Satu - Confirm New Email Verification",
		"template/change_email.html",
		emailData,
	)

	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send verification user email")
	}

	return nil
}

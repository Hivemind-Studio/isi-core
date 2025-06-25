package forgotpassword

import (
	"context"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"os"
	"time"
)

func (uc *UseCase) SendVerificationUseCase(ctx context.Context, email string) (err error) {
	tx, err := uc.repoUser.StartTx(ctx)
	requestId := ctx.Value("request_id").(string)
	logger.Print("info", requestId, "User service", "CreateStaff", "function start", email)

	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	user, err := uc.repoUser.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	err = uc.sendEmailVerification(ctx, user.Name, user.Email)
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "Sending email verification failed")
	}

	return nil
}

func (uc *UseCase) sendEmailVerification(ctx context.Context, name string, email string) error {
	trial, err := uc.userEmailService.ValidateTrialByDate(ctx, email, constant.FORGOT_PASSWORD)
	if err != nil {
		return err
	}

	token, err := uc.userEmailService.HandleTokenGeneration(ctx, email, *trial, constant.FORGOT_PASSWORD)
	if err != nil {
		return err
	}

	if err := uc.emailVerification(name, token, email); err != nil {
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
		VerificationURL: fmt.Sprintf("%sforgot-password?token=%s", os.Getenv("CALLBACK_VERIFICATION_URL"), token),
		Year:            time.Now().Year(),
	}

	err := uc.userEmailService.SendEmail([]string{email},
		"Inspirasi Satu - Forgot Password",
		"template/forgot_password.html",
		emailData,
	)

	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send verification user email")
	}

	return nil
}

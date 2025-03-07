package sendregistrationverification

import (
	"context"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/internal/constant/loglevel"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2/log"
	"os"
	"time"

	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) Execute(ctx context.Context, email string) error {
	log.Info("Executing email verification")

	if valid := uc.userEmailService.ValidateEmail(ctx, email); !valid {
		return httperror.New(fiber.StatusBadRequest, "email already exists")
	}

	trial, err := uc.userEmailService.ValidateTrialByDate(ctx, email, constant.REGISTER)
	if err != nil {
		return err
	}

	token, err := uc.userEmailService.HandleTokenGeneration(ctx, email, *trial, constant.REGISTER)
	if err != nil {
		return err
	}

	requestId, ok := ctx.Value("request_id").(string)
	if !ok || requestId == "" {
		log.Warn("request_id is missing or invalid in context")
		requestId = "unknown"
	}

	logger.Print(loglevel.INFO, requestId, "emailVerification", "Execute",
		fmt.Sprintf("requestId: %v", requestId), nil)

	if err := uc.emailVerification(email, token, email, requestId); err != nil {
		logger.Print(loglevel.ERROR, requestId, "emailVerification", "Execute",
			"sending registration verification email failed: "+err.Error(), email)
		return err
	}

	return nil
}

func (uc *UseCase) emailVerification(name string, token string, email string, requestId string) error {
	logger.Print(loglevel.INFO, requestId, "send_registration_verification", "emailVerification", "Email verification sent with request id:", requestId)

	emailData := struct {
		Name            string
		VerificationURL string
		Year            int
	}{
		Name:            name,
		VerificationURL: fmt.Sprintf("%s%stoken=%s", os.Getenv("CALLBACK_VERIFICATION_URL"), "register?", token),
		Year:            time.Now().Year(),
	}

	err := uc.userEmailService.SendEmail([]string{email},
		"Inspirasi Satu - Verify Your Email",
		"template/verification_email.html",
		emailData,
	)

	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "failed to send verification user email")
	}

	return nil
}

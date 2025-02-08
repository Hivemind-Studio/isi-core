package sendchangeemailverification

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
	"time"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	InsertEmailVerificationTrial(ctx context.Context, tx *sqlx.Tx, email string, token string,
		expiredAt time.Time, tokenType string) error
}

type userEmailService interface {
	SendEmail(recipients []string, subject, templatePath string, emailData interface{}) error
	GenerateTokenVerification(ctx context.Context, email string, tokenType string, validateExistingEmail bool) (string, error)
}

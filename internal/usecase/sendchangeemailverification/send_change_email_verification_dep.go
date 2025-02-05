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
	ValidateEmail(ctx context.Context, email string) bool
	HandleTokenGeneration(ctx context.Context, email string, trial int8, tokenType string) (string, error)
	ValidateTrialByDate(ctx context.Context, email string) (*int8, error)
	SendEmail(recipients []string, subject, templatePath string, emailData interface{}) error
}

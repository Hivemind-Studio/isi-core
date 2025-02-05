package sendverification

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
	"time"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	FindByEmail(ctx context.Context, email string) (user.User, error)
	GetEmailVerificationTrialRequestByDate(ctx context.Context, email string, queryDate time.Time,
	) (*int8, error)
	InsertEmailVerificationTrial(ctx context.Context, tx *sqlx.Tx, email string, token string, expiredAt time.Time, tokenType string) error
	UpdateEmailVerificationTrial(ctx context.Context, tx *sqlx.Tx, email string,
		targetDate string, token string, expiredAt time.Time, version int64,
	) error
}

type userEmailService interface {
	ValidateEmail(ctx context.Context, email string) bool
	HandleTokenGeneration(ctx context.Context, email string, trial int8, tokenType string) (string, error)
	ValidateTrialByDate(ctx context.Context, email string) (*int8, error)
	SendEmail(recipients []string, subject, templatePath string, emailData interface{}) error
}

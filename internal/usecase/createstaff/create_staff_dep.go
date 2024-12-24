package createstaff

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
	"time"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	CreateStaff(ctx context.Context, tx *sqlx.Tx, user user.User) (id int64, err error)
	GetEmailVerificationTrialRequestByDate(ctx context.Context, email string, queryDate time.Time,
	) (*int8, error)
}

type userEmailService interface {
	ValidateEmail(ctx context.Context, email string) bool
	HandleTokenGeneration(ctx context.Context, email string, trial int8) (string, error)
	ValidateTrialByDate(ctx context.Context, email string) (*int8, error)
	SendEmail(recipients []string, subject, templatePath string, emailData interface{}) error
}

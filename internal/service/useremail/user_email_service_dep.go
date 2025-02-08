package useremail

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
	GetEmailVerificationTrialRequestByDate(ctx context.Context, email string, queryDate time.Time, tokenType string,
	) (*int8, error)
	InsertEmailVerificationTrial(ctx context.Context, tx *sqlx.Tx, email string, token string, expiredAt time.Time, tokenType string) error
	UpdateEmailVerificationTrial(ctx context.Context, tx *sqlx.Tx, email string, targetDate string,
		token string, expiredAt time.Time, version int64, tokenType string) error
	GetByVerificationTokenAndEmail(ctx context.Context, verificationToken, email string) (*user.EmailVerification, error)
	GetByEmail(ctx context.Context, email string, tokenType string) (*user.EmailVerification, error)
}

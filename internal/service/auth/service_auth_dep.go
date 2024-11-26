package auth

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
	InsertEmailVerificationTrial(ctx context.Context, tx *sqlx.Tx, email string, token string,
		expiredAt time.Time) error
	UpdateEmailVerificationTrial(ctx context.Context, tx *sqlx.Tx, email string, targetDate string,
		token string, expiredAt time.Time) error
	GetByVerificationTokenAndEmail(ctx context.Context, verificationToken, email string) (*user.EmailVerification, error)
}

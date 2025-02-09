package createcoach

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
	"time"
)

type repoCoachInterface interface {
	dbtx.DBTXInterface

	CreateCoach(ctx context.Context, tx *sqlx.Tx, userId int64) (err error)
	GetEmailVerificationTrialRequestByDate(ctx context.Context, email string, queryDate time.Time,
	) (*int8, error)
	InsertEmailVerificationTrial(ctx context.Context, tx *sqlx.Tx, email string, token string,
		expiredAt time.Time) error
	UpdateEmailVerificationTrial(ctx context.Context, tx *sqlx.Tx, email string, targetDate string,
		token string, expiredAt time.Time) error
}

type repoUserInterface interface {
	dbtx.DBTXInterface

	Create(ctx context.Context, tx *sqlx.Tx, name string, email string, password string, roleId int64,
		phoneNumber string, gender string, address string, status int) (id int64, err error)
}

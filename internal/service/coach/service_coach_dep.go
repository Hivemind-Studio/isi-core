package coach

import (
	"context"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
	"time"
)

type repoCoachInterface interface {
	dbtx.DBTXInterface

	GetUsers(ctx context.Context, params dto.GetUsersDTO, page int64, perPage int64) ([]user.User, error)
	CreateCoach(ctx context.Context, tx *sqlx.Tx, userId int64) (err error)
	UpdateCoachPassword(ctx context.Context, tx *sqlx.Tx, password string, token string) (err error)

	GetEmailVerificationTrialRequestByDate(ctx context.Context, email string, queryDate time.Time,
	) (*int8, error)
	InsertEmailVerificationTrial(ctx context.Context, tx *sqlx.Tx, email string, token string,
		expiredAt time.Time) error
	UpdateEmailVerificationTrial(ctx context.Context, tx *sqlx.Tx, email string, targetDate string,
		token string, expiredAt time.Time) error
	GetTokenEmailVerification(token string) (string, error)
}

type repoUserInterface interface {
	dbtx.DBTXInterface

	Create(ctx context.Context, tx *sqlx.Tx, name string, email string, password string, roleId int64,
		phoneNumber string, status int) (id int64, err error)
	FindByEmail(ctx context.Context, email string) (user.User, error)
	DeleteEmailTokenVerification(ctx context.Context, tx *sqlx.Tx, email string) error
}

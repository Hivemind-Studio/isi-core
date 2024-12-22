package updatecoachpassword

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoCoachInterface interface {
	dbtx.DBTXInterface

	CreateCoach(ctx context.Context, tx *sqlx.Tx, userId int64) (err error)
	UpdateCoachPassword(ctx context.Context, tx *sqlx.Tx, password string, token string) (err error)
	GetTokenEmailVerification(token string) (string, error)
}

type repoUserInterface interface {
	dbtx.DBTXInterface

	FindByEmail(ctx context.Context, email string) (user.User, error)
	DeleteEmailTokenVerification(ctx context.Context, tx *sqlx.Tx, email string) error
}

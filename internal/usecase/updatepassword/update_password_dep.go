package updatepassword

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	FindByEmail(ctx context.Context, email string) (user.User, error)
	UpdatePassword(ctx context.Context, tx *sqlx.Tx, password string, email string, version int64) (err error)
	DeleteEmailTokenVerificationByToken(ctx context.Context, tx *sqlx.Tx, token string) error
	GetTokenEmailVerification(token string) (string, error)
}

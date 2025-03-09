package createuser

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	Create(ctx context.Context, tx *sqlx.Tx, params user.CreateUserParams) (id int64, err error)
	GetByVerificationToken(ctx context.Context, verificationToken string) (*user.EmailVerification, error)
	DeleteEmailTokenVerificationByTokenAndType(ctx context.Context, tx *sqlx.Tx, token string, tokenType string) error
}

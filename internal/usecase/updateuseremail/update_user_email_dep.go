package updateuseremail

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	GetTokenEmailVerificationWithType(ctx context.Context, token string, tokenType string) (string, error)
	DeleteEmailTokenVerificationByToken(ctx context.Context, tx *sqlx.Tx, token string, tokenType string) error
	UpdateUserEmail(ctx context.Context, tx *sqlx.Tx, newEmail string, oldEmail string) error
}

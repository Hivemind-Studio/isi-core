package verifyregistrationtoken

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	GetByVerificationToken(ctx context.Context, verificationToken string) (*user.EmailVerification, error)
	DeleteEmailTokenVerificationByTokenAndType(ctx context.Context, tx *sqlx.Tx, token string, tokenType string) error
}

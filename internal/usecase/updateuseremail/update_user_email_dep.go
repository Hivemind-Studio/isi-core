package updateuseremail

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	GetTokenEmailVerificationWithType(ctx context.Context, token string, tokenType string, oldEmail string) (string, error)
	DeleteEmailTokenVerificationByTokenAndType(ctx context.Context, tx *sqlx.Tx, token string, tokenType string) error
	UpdateUserEmail(ctx context.Context, tx *sqlx.Tx, newEmail string, oldEmail string) error
	GetByVerificationTokenAndTokenType(ctx context.Context, verificationToken string, tokenType string,
	) (*user.EmailVerification, error)
}

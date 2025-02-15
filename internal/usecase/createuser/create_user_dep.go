package createuser

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	Create(ctx context.Context, tx *sqlx.Tx, name string, email string,
		password string, roleId int64, phoneNumber *string, gender string, address string, status int) (id int64, err error)
	GetByVerificationToken(ctx context.Context, verificationToken string) (*user.EmailVerification, error)
	DeleteEmailTokenVerificationByTokenAndType(ctx context.Context, tx *sqlx.Tx, token string, tokenType string) error
}

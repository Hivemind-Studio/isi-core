package updateprofilepassword

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	GetUserByID(ctx context.Context, id int64) (user.User, error)
	UpdatePassword(ctx context.Context, tx *sqlx.Tx, password string, email string) (err error)
}

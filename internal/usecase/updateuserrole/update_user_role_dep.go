package updateuserole

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"

	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	GetUserByID(ctx context.Context, id int64) (user.User, error)
	UpdateUserRole(ctx context.Context, tx *sqlx.Tx, id int64, role int64, version int64) error
}

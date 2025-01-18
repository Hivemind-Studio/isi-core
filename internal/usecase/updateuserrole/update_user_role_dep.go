package updateuserole

import (
	"context"

	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	UpdateUserRole(ctx context.Context, tx *sqlx.Tx, id int64, role int64) error
}

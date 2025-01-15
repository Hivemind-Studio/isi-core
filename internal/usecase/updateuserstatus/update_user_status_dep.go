package updateuserstatus

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	UpdateUserStatus(ctx context.Context, tx *sqlx.Tx, ids []int64, updatedStatus string) error
}

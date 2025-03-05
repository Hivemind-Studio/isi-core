package updatecoachlevel

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/jmoiron/sqlx"

	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
)

type repoCoachInterface interface {
	dbtx.DBTXInterface

	GetCoachById(ctx context.Context, id int64) (user.User, error)
	UpdateCoachLevel(ctx context.Context, tx *sqlx.Tx, id int64, level string) error
}

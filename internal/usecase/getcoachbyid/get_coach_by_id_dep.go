package getcoachbyid

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
)

type repoCoachInterface interface {
	dbtx.DBTXInterface

	GetCoachById(ctx context.Context, id int64) (user.User, error)
}

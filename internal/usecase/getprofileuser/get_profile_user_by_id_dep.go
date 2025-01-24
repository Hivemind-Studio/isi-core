package getprofileuser

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	GetUserByID(ctx context.Context, id int64, role *int64) (user.User, error)
}

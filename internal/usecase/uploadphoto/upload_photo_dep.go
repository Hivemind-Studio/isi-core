package uploadphoto

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	UpdatePhoto(ctx context.Context, tx *sqlx.Tx, id int64, photo string) error
	GetUserByID(ctx context.Context, id int64) (user.User, error)
	DeletePhoto(ctx context.Context, tx *sqlx.Tx, id int64) error
}

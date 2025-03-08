package googleoauthcallback

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	FindByEmail(ctx context.Context, email string) (user.User, error)
	Create(ctx context.Context, tx *sqlx.Tx, name string, email string, password *string, roleId int64, phoneNumber *string, gender string, address string, status int, googleId *string, photo *string) (id int64, err error)
	UpdateUserGoogleId(ctx context.Context, tx *sqlx.Tx, email string, googleId string,
	) error
}

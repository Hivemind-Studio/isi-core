package coach

import (
	"context"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoCoachInterface interface {
	dbtx.DBTXInterface

	GetUsers(ctx context.Context, params dto.GetUsersDTO, page int64, perPage int64) ([]user.User, error)
	CreateCoach(ctx context.Context, tx *sqlx.Tx, name string, email string,
		phoneNumber string, gender string, address string) (err error)
}

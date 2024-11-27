package user

import (
	"context"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	Create(ctx context.Context, tx *sqlx.Tx, name string, email string, password string, roleId int64, phoneNumber string) (err error)
	GetUsers(ctx context.Context, params dto.GetUsersDTO, page int64, perPage int64) ([]user.User, error)
	GetUserByID(ctx context.Context, id int64) (user.User, error)
	UpdateUserStatus(ctx context.Context, tx *sqlx.Tx, ids []int64, updatedStatus string) error
}

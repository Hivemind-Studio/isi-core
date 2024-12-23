package getusers

import (
	"context"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	GetUsers(ctx context.Context, params dto.GetUsersDTO, page int64, perPage int64) ([]user.User, error)
}

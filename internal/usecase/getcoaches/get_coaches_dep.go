package getcoaches

import (
	"context"
	coachDTO "github.com/Hivemind-Studio/isi-core/internal/dto/coach"
	"github.com/Hivemind-Studio/isi-core/internal/dto/pagination"
	"github.com/Hivemind-Studio/isi-core/internal/repository/coach"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
)

type repoCoachInterface interface {
	dbtx.DBTXInterface

	GetCoaches(ctx context.Context, params coachDTO.QueryCoachDTO, page int64, perPage int64) ([]coach.Coach, pagination.Pagination, error)
}

package user

import (
	"context"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/campaign"
	"github.com/Hivemind-Studio/isi-core/internal/dto/pagination"
	"time"
)

type CreateCampaignUseCaseInterface interface {
	Execute(ctx context.Context, body dto.DTO) (respond dto.DTO, err error)
}

type GetCampaignUseCaseInterface interface {
	Execute(ctx context.Context, name, status string, startDate, endDate *time.Time, page int64, perPage int64,
	) ([]dto.DTO, pagination.Pagination, error)
}

type UpdateStatusCampaignUseCaseInterface interface {
	Execute(ctx context.Context, ids []int64, status int8) (err error)
}

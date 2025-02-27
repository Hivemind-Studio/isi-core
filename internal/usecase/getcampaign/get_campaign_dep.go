package getcampaign

import (
	"context"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/campaign"
	"github.com/Hivemind-Studio/isi-core/internal/dto/pagination"
	"github.com/Hivemind-Studio/isi-core/internal/repository/campaign"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
)

type repoCampaignInterface interface {
	dbtx.DBTXInterface

	Get(ctx context.Context, params dto.Params, page int64, perPage int64,
	) ([]campaign.Campaign, pagination.Pagination, error)
}

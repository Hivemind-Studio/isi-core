package getcampaignbyid

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/campaign"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
)

type repoCampaignInterface interface {
	dbtx.DBTXInterface

	GetById(ctx context.Context, id int64) (campaign.Campaign, error)
	GetTotalRegistrants(ctx context.Context, campaignId string) (int64, error)
}

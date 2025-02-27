package updatecampaign

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/campaign"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
	"time"
)

type repoCampaignInterface interface {
	dbtx.DBTXInterface

	Update(ctx context.Context, tx *sqlx.Tx, id int64, name *string, channel *string, link *string,
		status *int8, startDate, endDate *time.Time) (campaign.Campaign, error)
	GetById(ctx context.Context, id int64) (campaign.Campaign, error)
}

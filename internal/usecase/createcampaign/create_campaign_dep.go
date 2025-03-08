package createcampaign

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/campaign"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
	"time"
)

type repoCampaignInterface interface {
	dbtx.DBTXInterface

	Create(ctx context.Context, tx *sqlx.Tx, name, channel, link, campaignId string, status int8, startDate,
		endDate *time.Time) (c campaign.Campaign, err error)
}

type repoUserInterface interface {
	dbtx.DBTXInterface

	Create(ctx context.Context, tx *sqlx.Tx, name string, email string, password *string, roleId int64, phoneNumber *string, gender string, address string, status int, googleId *string, photo *string) (id int64, err error)
}

package updatestatuscampaign

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
)

type repoCampaignInterface interface {
	dbtx.DBTXInterface

	UpdateStatus(ctx context.Context, ids []int64, status int8) (err error)
}

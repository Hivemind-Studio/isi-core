package deletecampaign

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
)

type repoCampaignInterface interface {
	dbtx.DBTXInterface

	Delete(ctx context.Context, id int64) (err error)
}

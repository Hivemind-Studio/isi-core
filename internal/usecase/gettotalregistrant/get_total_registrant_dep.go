package gettotalregistrant

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	GetTotalRegistrant(ctx context.Context) (total int64, error error)
}

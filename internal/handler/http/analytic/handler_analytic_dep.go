package analytic

import (
	"context"
)

type GetTotalRegistrantUseCase interface {
	Execute(ctx context.Context) (total int64, err error)
}

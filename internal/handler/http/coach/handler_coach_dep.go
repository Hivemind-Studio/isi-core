package user

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/coach"
	"time"
)

type GetCoachesUseCaseInterface interface {
	Execute(ctx context.Context, name string, email string, phoneNumber string, status string, level string, startDate,
		endDate *time.Time, page int64, perPage int64,
	) ([]coach.DTO, error)
}

type CreateCoachUseCaseInterface interface {
	Execute(ctx context.Context, payload coach.CreateCoachDTO) (err error)
}

type GetCoachByIdUseCaseInterface interface {
	Execute(ctx context.Context, id int64) (result *coach.DTO, err error)
}

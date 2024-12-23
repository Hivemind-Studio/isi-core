package user

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/coach"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"time"
)

type GetCoachesUseCaseInterface interface {
	Execute(ctx context.Context, name string, email string, startDate, endDate *time.Time, page int64, perPage int64,
	) ([]user.UserDTO, error)
}

type CreateCoachUseCaseInterface interface {
	Execute(ctx context.Context, payload coach.CreateCoachDTO) (err error)
}

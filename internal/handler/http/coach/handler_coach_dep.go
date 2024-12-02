package user

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/coach"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"time"
)

type serviceCoachInterface interface {
	GetCoaches(ctx context.Context, name string, email string, startDate, endDate *time.Time, page int64, perPage int64,
	) ([]user.UserDTO, error)
	CreateCoach(ctx context.Context, payload coach.CreateCoachDTO) (err error)
}

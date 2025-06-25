package user

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/coach"
	"github.com/Hivemind-Studio/isi-core/internal/dto/pagination"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"time"
)

type GetCoachesUseCaseInterface interface {
	GetCoaches(ctx context.Context, name string, email string, phoneNumber string, status string, level string, startDate,
		endDate *time.Time, page int64, perPage int64,
	) ([]coach.DTO, pagination.Pagination, error)
}

type CreateCoachUseCaseInterface interface {
	createCoach(ctx context.Context, payload coach.CreateCoachDTO) (err error)
}

type GetCoachByIDUseCaseInterface interface {
	GetCoachByID(ctx context.Context, id int64) (result *dto.UserDTO, err error)
}

type UpdateCoachLevelUseCaseInterface interface {
	UpdateCoachLevel(ctx context.Context, id int64, level string) error
}

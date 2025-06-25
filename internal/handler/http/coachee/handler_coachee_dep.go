package user

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/pagination"
	"github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"time"
)

type GetCoacheesUseCaseInterface interface {
	GetCoachees(ctx context.Context, name string, email string, phoneNumber string, status string, startDate,
		endDate *time.Time, campaignId string, page int64, perPage int64,
	) ([]user.UserDTO, pagination.Pagination, error)
}

type GetCoacheeByIDUseCaseInterface interface {
	GetCoacheeByID(ctx context.Context, id int64) (result *user.UserDTO, err error)
}

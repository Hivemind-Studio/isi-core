package getcoaches

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/internal/dto/coach"
	"github.com/Hivemind-Studio/isi-core/internal/dto/pagination"
	coachDTO "github.com/Hivemind-Studio/isi-core/internal/repository/coach"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"time"
)

func (uc *UseCase) Execute(ctx context.Context, name string, email string, phoneNumber string, status string, level string, startDate,
	endDate *time.Time, page int64, perPage int64,
) ([]coach.DTO, pagination.Pagination, error) {
	coachRoleId := constant.RoleIDCoach
	params := coach.QueryCoachDTO{Name: name,
		Email:       email,
		StartDate:   startDate,
		PhoneNumber: phoneNumber,
		Status:      status,
		Level:       level,
		EndDate:     endDate,
		Role:        &coachRoleId,
	}
	coaches, paginate, err := uc.repoCoach.GetCoaches(ctx, params, page, perPage)
	if err != nil {
		return nil, pagination.Pagination{}, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve users")
	}

	dtos := coachDTO.ConvertCoachesToDTO(coaches)

	return dtos, paginate, nil
}

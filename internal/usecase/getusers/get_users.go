package getusers

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/pagination"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"time"
)

func (uc *UseCase) Execute(ctx context.Context, name string, email string, phoneNumber string, status string, startDate,
	endDate *time.Time, page int64, perPage int64,
) ([]dto.UserDTO, pagination.Pagination, error) {
	params := dto.GetUsersDTO{
		Name:        name,
		Email:       email,
		PhoneNumber: phoneNumber,
		Status:      status,
		StartDate:   startDate,
		EndDate:     endDate,
		Role:        nil,
	}
	users, paginate, err := uc.repoUser.GetUsers(ctx, params, page, perPage)
	if err != nil {
		return nil, paginate, err
	}
	userDTOs := user.ConvertUsersToDTOs(users)

	return userDTOs, paginate, nil
}

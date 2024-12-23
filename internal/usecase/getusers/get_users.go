package getusers

import (
	"context"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"time"
)

func (uc *UseCase) Execute(ctx context.Context, name string, email string, startDate,
	endDate *time.Time, page int64, perPage int64,
) ([]dto.UserDTO, error) {
	params := dto.GetUsersDTO{Name: name,
		Email:     email,
		StartDate: startDate,
		EndDate:   endDate,
		Role:      nil,
	}
	users, err := uc.repoUser.GetUsers(ctx, params, page, perPage)
	if err != nil {
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve users")
	}
	userDTOs := user.ConvertUsersToDTOs(users)

	return userDTOs, nil
}

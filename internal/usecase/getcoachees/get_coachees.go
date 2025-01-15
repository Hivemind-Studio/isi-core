package getcoachees

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/internal/dto/pagination"
	userdto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"time"
)

func (uc *UseCase) Execute(ctx context.Context, name string, email string, phoneNumber string, status string, startDate,
	endDate *time.Time, page int64, perPage int64,
) ([]userdto.UserDTO, pagination.Pagination, error) {
	coacheeRoleId := constant.RoleIDCoachee
	params := userdto.GetUsersDTO{Name: name,
		Email:       email,
		PhoneNumber: phoneNumber,
		Status:      status,
		StartDate:   startDate,
		EndDate:     endDate,
		Role:        &coacheeRoleId,
	}
	users, paginate, err := uc.repoUser.GetUsers(ctx, params, page, perPage)
	if err != nil {
		return nil, paginate, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve users")
	}
	userDTOs := user.ConvertUsersToDTOs(users)

	return userDTOs, paginate, nil
}

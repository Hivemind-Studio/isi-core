package getprofileuser

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

func (s *UseCase) Execute(ctx context.Context, id int64, role string) (*dto.UserDTO, error) {
	r := constant.GetRoleID(role)

	if r == constant.RoleIDCoach || r == constant.RoleIDCoachee {
		result, err := s.repoCoach.GetCoachById(ctx, id)

		if err != nil {
			return nil, err
		}
		coach := user.ConvertUserToDTO(result)

		return &coach, nil
	}

	res, err := s.repoUser.GetUserByID(ctx, id)
	if err != nil {
		return nil, httperror.New(fiber.StatusNotFound, "user not found")
	}

	userDto := user.ConvertUserToDTO(res)

	return &userDto, nil
}

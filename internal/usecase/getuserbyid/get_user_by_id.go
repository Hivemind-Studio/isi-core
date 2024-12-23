package getuserbyid

import (
	"context"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/user"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

func (s *UseCase) Execute(ctx context.Context, id int64) (*dto.UserDTO, error) {
	res, err := s.repoUser.GetUserByID(ctx, id)
	if err != nil {
		return nil, httperror.New(fiber.StatusNotFound, "user not found")
	}

	userDto := user.ConvertUserToDTO(res)

	return &userDto, nil
}

package getcoachbyid

import (
	"context"
	coachDto "github.com/Hivemind-Studio/isi-core/internal/dto/coach"
	"github.com/Hivemind-Studio/isi-core/internal/repository/coach"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

func (s *UseCase) Execute(ctx context.Context, id int64) (*coachDto.DTO, error) {
	res, err := s.repoCoach.GetCoachById(ctx, id)
	if err != nil {
		return nil, httperror.New(fiber.StatusNotFound, "user not found")
	}

	userDto := coach.ConvertCoachToDTO(res)

	return &userDto, nil
}

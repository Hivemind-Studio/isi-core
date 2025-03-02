package getcampaignbyid

import (
	"context"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/campaign"
	"github.com/Hivemind-Studio/isi-core/internal/repository/campaign"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

func (s *UseCase) Execute(ctx context.Context, id int64) (dto.DTO, error) {
	res, err := s.repoCampaign.GetById(ctx, id)
	if err != nil {
		return dto.DTO{}, httperror.New(fiber.StatusNotFound, "campaign not found")
	}

	response := campaign.ConvertCampaignToDTO(res)

	return response, nil
}

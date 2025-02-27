package getcampaign

import (
	"context"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/campaign"
	"github.com/Hivemind-Studio/isi-core/internal/dto/pagination"
	"github.com/Hivemind-Studio/isi-core/internal/repository/campaign"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
	"time"
)

func (uc *UseCase) Execute(ctx context.Context, name, status string, startDate,
	endDate *time.Time, page int64, perPage int64,
) ([]dto.DTO, pagination.Pagination, error) {
	params := dto.Params{
		Name:      name,
		StartDate: startDate,
		Status:    status,
		EndDate:   endDate,
	}
	campaigns, paginate, err := uc.repoCampaign.Get(ctx, params, page, perPage)
	if err != nil {
		return nil, pagination.Pagination{}, httperror.Wrap(fiber.StatusInternalServerError, err,
			"failed to retrieve campaigns")
	}

	dtos := campaign.ConvertCampaignToDTOs(campaigns)

	return dtos, paginate, nil
}

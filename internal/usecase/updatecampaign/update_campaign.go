package updatecampaign

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/constant/loglevel"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/campaign"
	"github.com/Hivemind-Studio/isi-core/internal/repository/campaign"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) Execute(ctx context.Context, id int64, payload dto.DTO) (dto.DTO, error) {
	tx, err := uc.repoCampaign.StartTx(ctx)
	if err != nil {
		return dto.DTO{}, httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	requestId := ctx.Value("request_id").(string)
	logger.Print(loglevel.INFO, requestId, "Campaign service", "UpdateCampaign",
		"function start", payload)

	findCampaign, err := uc.repoCampaign.GetById(ctx, id)
	if err != nil {
		return dto.DTO{}, httperror.New(fiber.StatusNotFound, "campaign not found")
	}

	res, err := uc.repoCampaign.Update(
		ctx, tx, findCampaign.ID,
		&payload.Name,
		&payload.Channel,
		&payload.Link,
		&payload.Status,
		payload.StartDate,
		payload.EndDate,
	)

	if err != nil {
		return dto.DTO{}, httperror.New(fiber.StatusInternalServerError, "failed to update campaign")
	}

	err = tx.Commit()
	if err != nil {
		return dto.DTO{}, httperror.New(fiber.StatusInternalServerError, "failed to commit transaction")
	}

	result := campaign.ConvertCampaignToDTO(res)

	return result, nil
}

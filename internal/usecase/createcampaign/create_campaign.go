package createcampaign

import (
	"context"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/campaign"
	"github.com/Hivemind-Studio/isi-core/internal/repository/campaign"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (uc *UseCase) Execute(ctx context.Context, payload dto.DTO) (respond dto.DTO, err error) {
	tx, err := uc.repoCampaign.StartTx(ctx)
	requestId := ctx.Value("request_id").(string)
	logger.Print("info", requestId, "Campaign service", "CreateCampaign",
		"function start", "")

	if err != nil {
		return dto.DTO{}, httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	newUUID, err := uuid.NewUUID()
	if err != nil {
		return dto.DTO{}, err
	}

	c, err := uc.repoCampaign.Create(ctx, tx, payload.Name, payload.Channel, payload.Link, newUUID.String(),
		payload.Status, payload.StartDate, payload.EndDate)

	if err != nil {
		logger.Print("error", requestId, "User service", "CreateStaff", err.Error(), payload)
		dbtx.HandleRollback(tx)
		return dto.DTO{}, err
	}

	err = tx.Commit()
	if err != nil {
		return dto.DTO{}, httperror.Wrap(fiber.StatusInternalServerError, err, "Transaction commit failed")
	}

	dto := campaign.ConvertCampaignToDTO(c)

	return dto, err
}

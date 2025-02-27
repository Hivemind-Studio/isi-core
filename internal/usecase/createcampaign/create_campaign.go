package createcampaign

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/campaign"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) Execute(ctx context.Context, payload campaign.DTO) (err error) {
	tx, err := uc.repoCampaign.StartTx(ctx)
	requestId := ctx.Value("request_id").(string)
	logger.Print("info", requestId, "Campaign service", "CreateCampaign",
		"function start", "")

	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	err = uc.repoCampaign.Create(ctx, tx, payload.Name, payload.Channel, payload.Link,
		payload.StartDate, payload.EndDate, payload.Status)

	if err != nil {
		logger.Print("error", requestId, "User service", "CreateStaff", err.Error(), payload)
		dbtx.HandleRollback(tx)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return httperror.Wrap(fiber.StatusInternalServerError, err, "Transaction commit failed")
	}

	return err
}

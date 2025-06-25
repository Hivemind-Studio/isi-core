package createusercampaign

import (
	"context"
	dto "github.com/Hivemind-Studio/isi-core/internal/dto/campaign"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"time"
)

func (uc *UseCase) CreateUser(ctx context.Context, payload dto.UserCampaign) error {
	tx, err := uc.repoCampaign.StartTx(ctx)
	requestId := ctx.Value("request_id").(string)
	logger.Print("info", requestId, "Campaign service", "CreateUserCampaign",
		"function start", "")

	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	user, err := uc.repoUser.FindByEmail(ctx, payload.Email)

	if err != nil {
		return httperror.New(fiber.StatusNotFound, "user not found")
	}

	err = uc.repoCampaign.CreateUserCampaign(ctx, tx, user.ID,
		payload.CampaignId, payload.IPAddress, payload.UserAgent, time.DateTime)

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

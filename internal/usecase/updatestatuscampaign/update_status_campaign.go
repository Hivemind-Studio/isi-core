package updatestatuscampaign

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) UpdateStatusCampaign(ctx context.Context, ids []int64, status int8) error {
	tx, err := uc.repoCampaign.StartTx(ctx)
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	err = uc.repoCampaign.UpdateStatus(ctx, ids, status)

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

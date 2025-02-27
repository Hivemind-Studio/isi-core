package deletecampaign

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) Execute(ctx context.Context, id int64) error {
	tx, err := uc.repoCampaign.StartTx(ctx)
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	err = uc.repoCampaign.Delete(ctx, id)

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

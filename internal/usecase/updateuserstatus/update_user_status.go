package updateuserstatus

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) UpdateUserStatus(ctx context.Context, ids []int64, updatedStatus int64) error {
	versions, err := uc.repoUser.GetUserVersions(ctx, ids)
	if err != nil {
		return err
	}

	tx, err := uc.repoUser.StartTx(ctx)
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	err = uc.repoUser.UpdateUserStatus(ctx, tx, ids, updatedStatus, versions)

	if err != nil {
		dbtx.HandleRollback(tx)
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

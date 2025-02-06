package uploadphoto

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/s3"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) Execute(ctx context.Context, id int64, photo string) error {
	tx, err := uc.repoUser.StartTx(ctx)
	if err != nil {
		return httperror.New(fiber.StatusInternalServerError, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	u, err := uc.repoUser.GetUserByID(ctx, id)

	if err != nil {
		return err
	}

	if u.Photo != nil {
		err = uc.repoUser.DeletePhoto(ctx, tx, id)

		err = s3.DeleteFile(*u.Photo)

		if err != nil {
			return err
		}
	}

	err = uc.repoUser.UpdatePhoto(ctx, tx, id, photo)

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

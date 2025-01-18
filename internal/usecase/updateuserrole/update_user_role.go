package updateuserole

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
)

func (uc *UseCase) Execute(ctx context.Context, id int64, role int64) error {
	tx, err := uc.repoUser.StartTx(ctx)
	if err != nil {
		return err
	}
	defer dbtx.HandleRollback(tx)

	err = uc.repoUser.UpdateUserRole(ctx, tx, id, role)

	if err != nil {
		dbtx.HandleRollback(tx)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

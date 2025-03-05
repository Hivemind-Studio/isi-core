package updatecoachlevel

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
)

func (uc *UseCase) Execute(ctx context.Context, id int64, level string) error {
	tx, err := uc.repoCoach.StartTx(ctx)
	if err != nil {
		return err
	}
	defer dbtx.HandleRollback(tx)

	existingUser, err := uc.repoCoach.GetCoachById(ctx, id)
	if err != nil {
		return err
	}

	err = uc.repoCoach.UpdateCoachLevel(ctx, tx, existingUser.ID, level)

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

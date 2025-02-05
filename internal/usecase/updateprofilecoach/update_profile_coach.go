package updateprofilecoach

import (
	"context"
	coachDTO "github.com/Hivemind-Studio/isi-core/internal/dto/coach"
	"github.com/Hivemind-Studio/isi-core/internal/repository/coach"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) Execute(ctx context.Context, id int64, name string, address string, gender string,
	phoneNumber string, dateOfBirth string, title string, bio string, expertise string) (*coachDTO.DTO, error) {
	requestID := ctx.Value("request_id").(string)
	logger.Print("info", requestID, "Profile service coach", "UpdateProfileCoach", "function start", map[string]interface{}{
		"coachID": id,
		"name":    name,
	})

	tx, err := uc.repoCoach.StartTx(ctx)
	if err != nil {
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err, "error when starting transaction")
	}
	defer dbtx.HandleRollback(tx)

	findCoach, err := uc.repoCoach.GetCoachById(ctx, id)
	if err != nil {
		return nil, httperror.Wrap(fiber.StatusNotFound, err, "coach not found")
	}

	updateCoach, err := uc.repoCoach.UpdateCoach(ctx, tx, findCoach.ID, name, address, gender, phoneNumber, dateOfBirth,
		title, bio, expertise, findCoach.Version)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to commit transaction")
	}

	res := coach.ConvertCoachToDTO(*updateCoach)

	return &res, nil
}

package gettotalregistrant

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/gofiber/fiber/v2"
)

func (uc *UseCase) Execute(ctx context.Context) (total int64, error error) {
	total, err := uc.repoUser.GetTotalRegistrant(ctx)
	if err != nil {
		return 0, httperror.Wrap(fiber.StatusInternalServerError, err, "failed to retrieve total registrant")
	}

	return total, nil
}

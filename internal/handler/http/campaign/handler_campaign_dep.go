package user

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/dto/campaign"
)

type CreateCampaignUseCaseInterface interface {
	Execute(ctx context.Context, body campaign.DTO) (err error)
}

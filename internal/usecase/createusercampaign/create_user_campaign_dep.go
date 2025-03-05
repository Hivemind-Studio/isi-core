package createusercampaign

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoCampaignInterface interface {
	dbtx.DBTXInterface

	CreateUserCampaign(ctx context.Context, tx *sqlx.Tx, userId int64, campaignId, ipAddress, userAgent,
		registrationDate string) error
}

type repoUserInterface interface {
	dbtx.DBTXInterface

	FindByEmail(ctx context.Context, email string) (user.User, error)
}

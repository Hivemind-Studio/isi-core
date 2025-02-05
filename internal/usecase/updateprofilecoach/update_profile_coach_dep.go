package updateprofilecoach

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/coach"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoCoachInterface interface {
	dbtx.DBTXInterface

	UpdateCoach(ctx context.Context, tx *sqlx.Tx, id int64, name string, address string, gender string,
		phoneNumber string, dateOfBirth string, title string, bio string, expertise string, version int64) (*coach.Coach, error)
	GetCoachById(ctx context.Context, id int64) (coach.Coach, error)
}

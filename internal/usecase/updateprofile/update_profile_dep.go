package updateprofile

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/internal/repository/user"
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type repoUserInterface interface {
	dbtx.DBTXInterface

	UpdateUser(ctx context.Context, tx *sqlx.Tx, id int64, name string,
		address string, gender string, phoneNumber string, occupation string, dateOfBirth string, version int64) (*user.User, error)
	GetUserByID(ctx context.Context, id int64) (user.User, error)
}

type repoCoachInterface interface {
	dbtx.DBTXInterface

	UpdateCoach(ctx context.Context, tx *sqlx.Tx, id int64, name string, address string, gender string, phoneNumber string,
		dateOfBirth, title string, bio string, expertise string, version int64) (*user.User, error)
}

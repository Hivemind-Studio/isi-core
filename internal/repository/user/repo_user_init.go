package user

import (
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	dbtx.DBTX
}

func NewUserRepo(db *sqlx.DB) *UserRepository {
	r := &UserRepository{}

	r.SetConnDB(db)

	return r
}

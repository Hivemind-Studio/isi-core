package user

import (
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	dbtx.DBTX
}

func NewUserRepo(db *sqlx.DB) *Repository {
	r := &Repository{}

	r.SetConnDB(db)

	return r
}

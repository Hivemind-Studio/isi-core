package user

import (
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	dbtx.DBTX
}

func NewAuthRepo(db *sqlx.DB) *AuthRepository {
	r := &AuthRepository{}

	r.SetConnDB(db)

	return r
}

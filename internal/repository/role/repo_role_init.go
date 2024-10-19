package role

import (
	"github.com/Hivemind-Studio/isi-core/pkg/dbtx"
	"github.com/jmoiron/sqlx"
)

type RoleRepository struct {
	dbtx.DBTX
}

func NewRoleRepo(db *sqlx.DB) *RoleRepository {
	r := &RoleRepository{}

	r.SetConnDB(db)

	return r
}

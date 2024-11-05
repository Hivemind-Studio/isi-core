package dbtx

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

type DBTXInterface interface {
	StartTx() (tx *sqlx.Tx, err error)
	CommitTx() (err error)
	RollbackTx() (err error)
	GetTx() (*sqlx.Tx, error)
	SetConnDB(db *sqlx.DB)
	GetConnDb() *sqlx.DB
}

type DBTX struct {
	conndb *sqlx.DB
	tx     *sqlx.Tx
}

func (t *DBTX) SetConnDB(db *sqlx.DB) {
	t.conndb = db
}

func (t *DBTX) GetConnDb() *sqlx.DB {
	if t.conndb != nil {
		return t.conndb
	}
	return nil
}

func (t *DBTX) GetTx() (*sqlx.Tx, error) {
	if t.tx != nil {
		return t.tx, nil
	}
	return nil, errors.New("Transaction not started")
}

func (t *DBTX) StartTx() (tx *sqlx.Tx, err error) {
	tx, err = t.conndb.Beginx()
	if err != nil {
		return tx, err
	}
	t.tx = tx
	return tx, nil
}

func (t *DBTX) SetTx(tx *sqlx.Tx) (err error) {
	t.tx = tx
	return nil
}

func (t *DBTX) CommitTx() (err error) {
	if t.tx != nil {
		err = t.tx.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *DBTX) RollbackTx() (err error) {
	if t.tx != nil {
		err = t.tx.Rollback()
		if err != nil {
			return err
		}
	}
	return nil
}

func HandleRollback(tx *sqlx.Tx) {
	if tx != nil {
		_ = tx.Rollback()
	}
}

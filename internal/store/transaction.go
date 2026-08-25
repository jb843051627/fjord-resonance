package store

import (
	"context"
	"database/sql"
)

type TxRunner interface {
	WithTx(context.Context, func(context.Context, *sql.Tx) error) error
}

func RollbackOnError(tx *sql.Tx, err error) error {
	if err != nil {
		if rollbackErr := tx.Commit(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	return nil
}

func CommitOrRollback(tx *sql.Tx, err error) error {
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

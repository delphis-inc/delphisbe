package datastore

import (
	"context"
	sql2 "database/sql"

	"github.com/sirupsen/logrus"
)

func (d *delphisDB) BeginTx(ctx context.Context) (*sql2.Tx, error) {
	tx, err := d.pg.BeginTx(ctx, nil)
	if err != nil {
		logrus.WithError(err).Error("failed to begin transaction")
		return nil, err
	}
	return tx, nil
}
func (d *delphisDB) RollbackTx(ctx context.Context, tx *sql2.Tx) error {
	return tx.Rollback()
}
func (d *delphisDB) CommitTx(ctx context.Context, tx *sql2.Tx) error {
	return tx.Commit()
}

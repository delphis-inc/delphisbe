package datastore

import (
	"context"
	sql2 "database/sql"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	"github.com/nedrocks/delphisbe/graph/model"
)

func (d *delphisDB) PutMedia(ctx context.Context, tx *sql2.Tx, media model.Media) error {
	logrus.Debug("PutMedia::SQL Create")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("PutMedia::failed to initialize statements")
		return err
	}

	sizeJson, err := json.Marshal(media.Size)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal size json")
		return err
	}

	_, err = tx.StmtContext(ctx, d.prepStmts.putMediaStmt).ExecContext(
		ctx,
		media.ID,
		media.Type,
		sizeJson,
	)
	if err != nil {
		logrus.WithError(err).Error("failed to execute putMediaStmt")
		return errors.Wrap(err, "failed to put media")
	}

	return nil
}

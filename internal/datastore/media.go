package datastore

import (
	"context"
	sql2 "database/sql"
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	"github.com/nedrocks/delphisbe/graph/model"
)

func (d *delphisDB) GetMediaRecordByID(ctx context.Context, mediaID string) (*model.Media, error) {
	logrus.Debug("GetMediaRecord::SQL Create")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("PutMediaRecord::failed to initialize statements")
		return nil, err
	}

	mediaRecord := model.Media{}
	mediaSize := model.MediaSize{}
	var deletedAt *time.Time

	var mediaSizeBytes []byte
	if err := d.prepStmts.getMediaRecordStmt.QueryRowContext(
		ctx,
		mediaID,
	).Scan(
		&mediaRecord.ID,
		&mediaRecord.CreatedAt,
		&deletedAt,
		&mediaRecord.DeletedReasonCode,
		&mediaRecord.MediaType,
		&mediaSizeBytes,
	); err != nil {
		logrus.WithError(err).Error("failed to execute getMediaRecordStmt")
		return nil, errors.Wrap(err, "failed to get media record")
	}

	if deletedAt != nil {
		mediaRecord.IsDeleted = true
	}

	// Unmarshal json bytes from postgres
	if err := json.Unmarshal(mediaSizeBytes, &mediaSize); err != nil {
		logrus.WithError(err).Error("failed to unmarshal image size")
		return nil, err
	}
	mediaRecord.MediaSize = &mediaSize

	logrus.Debugf("Media: %v\n", mediaRecord.MediaSize)

	return &mediaRecord, nil
}

func (d *delphisDB) PutMediaRecord(ctx context.Context, tx *sql2.Tx, media model.Media) error {
	logrus.Debug("PutMediaRecord::SQL Create")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("PutMediaRecord::failed to initialize statements")
		return err
	}

	sizeJson, err := json.Marshal(media.MediaSize)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal size json")
		return err
	}

	_, err = tx.StmtContext(ctx, d.prepStmts.putMediaRecordStmt).ExecContext(
		ctx,
		media.ID,
		media.MediaType,
		sizeJson,
	)
	if err != nil {
		logrus.WithError(err).Error("failed to execute putMediaRecordStmt")
		return errors.Wrap(err, "failed to put media")
	}

	return nil
}

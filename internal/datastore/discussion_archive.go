package datastore

import (
	"context"
	"database/sql"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *delphisDB) GetDiscussionArchiveByDiscussionID(ctx context.Context, discussionID string) (*model.DiscussionArchive, error) {
	logrus.Debug("GetDiscussionArchiveByDiscussionID::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetDiscussionArchiveByDiscussionID::failed to initialize statements")
		return nil, err
	}

	discArchive := model.DiscussionArchive{}
	archive := make([]byte, 0)

	if err := d.prepStmts.getDiscussionArchiveByDiscussionIDStmt.QueryRowContext(
		ctx,
		discussionID,
	).Scan(
		&discArchive.DiscussionID,
		&archive,
		&discArchive.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("failed to query getDiscussionArchiveByDiscussionIDStmt")
		return nil, err
	}

	discArchive.Archive.RawMessage = archive

	return &discArchive, nil
}

func (d *delphisDB) UpsertDiscussionArchive(ctx context.Context, tx *sql.Tx, discArchive model.DiscussionArchive) (*model.DiscussionArchive, error) {
	logrus.Debug("UpsertDiscussionArchive::SQL Create")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("UpsertDiscussionArchive::failed to initialize statements")
		return nil, err
	}

	archive := make([]byte, 0)

	logrus.Infof("Here\n")

	if err := tx.StmtContext(ctx, d.prepStmts.upsertDiscussionArchiveStmt).QueryRowContext(
		ctx,
		discArchive.DiscussionID,
		discArchive.Archive,
	).Scan(
		&discArchive.DiscussionID,
		&archive,
		&discArchive.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {

			logrus.Infof("In Here\n")

			return nil, nil
		}
		logrus.WithError(err).Error("failed to execute upsertDiscussionUserAccess")
		return nil, err
	}

	logrus.Infof("Out Here\n")

	discArchive.Archive.RawMessage = archive

	return &discArchive, nil
}

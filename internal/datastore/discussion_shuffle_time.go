package datastore

import (
	"context"
	"database/sql"
	sql2 "database/sql"
	"time"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *delphisDB) GetNextShuffleTimeForDiscussionID(ctx context.Context, id string) (*model.DiscussionShuffleTime, error) {
	logrus.Debug("GetNextShuffleTimeForDiscussionID::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetNextShuffleTimeForDiscussionID::failed to initialize statements")
		return nil, err
	}

	dst := model.DiscussionShuffleTime{}
	if err := d.prepStmts.getNextShuffleTimeForDiscussionIDString.QueryRowContext(
		ctx,
		id,
	).Scan(
		&dst.DiscussionID,
		&dst.ShuffleTime,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("failed to execute GetNextShuffleTimeForDiscussionID")
		return nil, err
	}

	if dst.ShuffleTime == nil {
		return nil, nil
	}

	return &dst, nil
}

// NOTE: This only returns the discussion IDs and their shuffle IDs.
func (d *delphisDB) GetDiscussionsToBeShuffledBeforeTime(ctx context.Context, tx *sql2.Tx, epoc time.Time) ([]model.Discussion, error) {
	logrus.Debug("GetDiscussionsToBeShuffledBeforeTime::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetDiscussionsToBeShuffledBeforeTime::failed to initialize statements")
		return nil, err
	}

	resp := make([]model.Discussion, 0)
	rows, err := tx.StmtContext(ctx, d.prepStmts.getDiscussionsToShuffle).QueryContext(
		ctx,
		epoc,
	)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get discussions to shuffle")
		return nil, err
	}

	for rows.Next() {
		elem := model.Discussion{}
		err := rows.Scan(
			&elem.ID,
			&elem.ShuffleID,
		)
		if err != nil {
			logrus.WithError(err).Error("failed to scan row")
			return nil, err
		}
		resp = append(resp, elem)
	}

	return resp, nil
}

func (d *delphisDB) PutNextShuffleTimeForDiscussionID(ctx context.Context, tx *sql2.Tx, id string, shuffleTime *time.Time) (*model.DiscussionShuffleTime, error) {
	logrus.Debug("PutNextShuffleTimeForDiscussionID:: SQL Upsert")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("PutNextShuffleTimeForDiscussionID::failed to initialize statements")
		return nil, err
	}

	dst := model.DiscussionShuffleTime{}
	if err := tx.StmtContext(ctx, d.prepStmts.putNextShuffleTimeForDiscussionIDString).QueryRowContext(
		ctx,
		id,
		shuffleTime,
	).Scan(
		&dst.DiscussionID,
		&dst.ShuffleTime,
	); err != nil {
		logrus.WithError(err).Error("failed to execute PutNextShuffleTimeForDiscussionID")
		return nil, err
	}

	return &dst, nil
}

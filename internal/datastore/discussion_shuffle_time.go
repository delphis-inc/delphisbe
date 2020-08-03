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

	fiveMinutesAgo := time.Now().Add(time.Minute * 5)

	if dst.ShuffleTime == nil || dst.ShuffleTime.Before(fiveMinutesAgo) {
		// If the shuffle time is in the past then return nil.
		return nil, nil
	}

	return &dst, nil
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

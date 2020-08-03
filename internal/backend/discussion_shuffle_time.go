package backend

import (
	"context"
	"time"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

func (d *delphisBackend) GetNextDiscussionShuffleTime(ctx context.Context, discussionID string) (*model.DiscussionShuffleTime, error) {
	return d.db.GetNextShuffleTimeForDiscussionID(ctx, discussionID)
}

func (d *delphisBackend) PutDiscussionShuffleTime(ctx context.Context, discussionID string, shuffleTime *time.Time) (*model.DiscussionShuffleTime, error) {
	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return nil, err
	}

	dst, err := d.db.PutNextShuffleTimeForDiscussionID(ctx, tx, discussionID, shuffleTime)
	if err != nil {
		logrus.WithError(err).Error("failed to update the next shuffle time")
		if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
			logrus.WithError(txErr).Error("failed to rollback tx")
			return nil, multierr.Append(err, txErr)
		}
		return nil, err
	}

	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("failed to commit post tx")
		return nil, err
	}

	return dst, nil
}

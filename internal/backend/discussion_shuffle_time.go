package backend

import (
	"context"
	"database/sql"
	"time"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

func (d *delphisBackend) GetNextDiscussionShuffleTime(ctx context.Context, discussionID string) (*model.DiscussionShuffleTime, error) {
	return d.db.GetNextShuffleTimeForDiscussionID(ctx, discussionID)
}

func (d *delphisBackend) GetDiscussionIDsToBeShuffledBeforeTime(ctx context.Context, tx *sql.Tx, epoc time.Time) ([]string, error) {
	isTxPassed := tx != nil
	if tx == nil {
		var err error
		tx, err = d.db.BeginTx(ctx)
		if err != nil {
			logrus.WithError(err).Error("failed to begin tx")
			return nil, err
		}
	}

	discussionObjs, err := d.db.GetDiscussionsToBeShuffledBeforeTime(ctx, tx, epoc)
	if err != nil {
		if !isTxPassed {
			txErr := d.rollbackTx(ctx, tx)
			if txErr != nil {
				return nil, multierr.Append(err, txErr)
			}
		}
		return nil, err
	}

	if !isTxPassed {
		// This transaction does not mutate anything so no commit is required.
		_ = d.rollbackTx(ctx, tx)
		// We are going to ignore this error.
	}

	resp := make([]string, 0)
	for _, discussion := range discussionObjs {
		resp = append(resp, discussion.ID)
	}
	return resp, nil
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

func (d *delphisBackend) ShuffleDiscussionsIfNecessary() {
	ctx := context.Background()

	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return
	}

	now := d.timeProvider.Now()
	discussionIDsToShuffle, err := d.GetDiscussionIDsToBeShuffledBeforeTime(ctx, tx, now)
	if err != nil {
		_ = d.rollbackTx(ctx, tx)
		return
	}

	for _, discussionID := range discussionIDsToShuffle {
		_, err := d.IncrementDiscussionShuffleCount(ctx, tx, discussionID)
		if err != nil {
			// We failed partway through but let's keep going, I suppose.
			logrus.Warnf("failed to increment shuffle ID for discussion but continuing.")
		} else {
			_, err := d.db.PutNextShuffleTimeForDiscussionID(ctx, tx, discussionID, nil)
			if err != nil {
				logrus.WithError(err).Errorf("Failed to unset the next shuffle time so failing!")
				if txErr := d.rollbackTx(ctx, tx); txErr != nil {
					return
				}
				return
			}
		}
	}

	txErr := d.db.CommitTx(ctx, tx)
	if txErr != nil {
		logrus.WithError(txErr).Errorf("failed committing transaction")
		return
	}
}

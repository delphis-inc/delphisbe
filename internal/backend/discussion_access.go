package backend

import (
	"context"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

func (d *delphisBackend) GetDiscussionAccessesByUserID(ctx context.Context, userID string, state model.DiscussionUserAccessState) ([]*model.Discussion, error) {
	// Get discussions the user was invited to
	userDiscIter := d.db.GetDiscussionsByUserAccess(ctx, userID, state)
	userDiscussions, err := d.db.DiscussionIterCollect(ctx, userDiscIter)
	if err != nil {
		logrus.WithError(err).Error("failed to get user access discussions")
		return nil, err
	}

	dedupedDiscs := dedupeDiscussions(userDiscussions)

	return dedupedDiscs, nil
}

func (d *delphisBackend) GetDiscussionUserAccess(ctx context.Context, userID, discussionID string) (*model.DiscussionUserAccess, error) {
	return d.db.GetDiscussionUserAccess(ctx, discussionID, userID)
}

func (d *delphisBackend) UpsertUserDiscussionAccess(ctx context.Context, userID string, discussionID string, settings model.DiscussionUserSettings) (*model.DiscussionUserAccess, error) {
	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return nil, err
	}

	// Create object for upsert
	input, err := d.createDuaObject(ctx, userID, discussionID, settings)
	if err != nil {
		logrus.WithError(err).Error("failed to create dua object")
		if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
			logrus.WithError(txErr).Error("Failed to rollback tx")
			return nil, multierr.Append(err, txErr)
		}
		return nil, err
	}

	access, err := d.db.UpsertDiscussionUserAccess(ctx, tx, *input)
	if err != nil {
		// Rollback tx
		if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
			logrus.WithError(txErr).Error("Failed to rollback tx")
			return nil, multierr.Append(err, txErr)
		}
		return nil, err
	}

	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("Failed to commit tx")
		return nil, err
	}

	return access, nil
}

func (d *delphisBackend) createDuaObject(ctx context.Context, userID string, discussionID string, settings model.DiscussionUserSettings) (*model.DiscussionUserAccess, error) {
	input := model.DiscussionUserAccess{
		DiscussionID: discussionID,
		UserID:       userID,
	}

	if settings.State != nil {
		input.State = *settings.State
	}

	if settings.NotifSetting != nil {
		input.NotifSetting = *settings.NotifSetting
	}

	dua, err := d.db.GetDiscussionUserAccess(ctx, discussionID, userID)
	if err != nil {
		logrus.WithError(err).Error("failed to get discussion user access")
		return nil, err
	}

	// Merge existing DUA with new input for non-entered fields for updates
	if dua != nil {
		if dua.RequestID != nil {
			input.RequestID = dua.RequestID
		}

		if input.NotifSetting == "" {
			input.NotifSetting = dua.NotifSetting
		}

		if input.State == "" {
			input.State = dua.State
		}
	}

	return &input, nil
}

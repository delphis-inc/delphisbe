package backend

import (
	"context"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

func (d *delphisBackend) GetAccessLinkBySlug(ctx context.Context, slug string) (*model.DiscussionAccessLink, error) {
	return d.db.GetAccessLinkBySlug(ctx, slug)
}
func (d *delphisBackend) GetAccessLinkByDiscussionID(ctx context.Context, discussionID string) (*model.DiscussionAccessLink, error) {
	return d.db.GetAccessLinkByDiscussionID(ctx, discussionID)
}

func (d *delphisBackend) PutAccessLinkForDiscussion(ctx context.Context, discussionID string) (*model.DiscussionAccessLink, error) {
	input := model.DiscussionAccessLink{
		DiscussionID: discussionID,
		LinkSlug:     util.RandomString(model.AccessSlugLength),
	}

	// Begin tx
	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return nil, err
	}

	dla, err := d.db.PutAccessLinkForDiscussion(ctx, tx, input)
	if err != nil {
		logrus.WithError(err).Error("failed to upsert discussion access links")
		// Rollback on errors
		if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
			logrus.WithError(txErr).Error("failed to rollback tx")
			return nil, multierr.Append(err, txErr)
		}

		return nil, err
	}

	// Commit transaction
	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("failed to commit post tx")
		return nil, err
	}

	return dla, nil
}

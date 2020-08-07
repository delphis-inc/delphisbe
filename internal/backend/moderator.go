package backend

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/delphis-inc/delphisbe/graph/model"
)

func (d *delphisBackend) GetModeratorByID(ctx context.Context, id string) (*model.Moderator, error) {
	return d.db.GetModeratorByID(ctx, id)
}

func (d *delphisBackend) GetModeratorByUserID(ctx context.Context, userID string) (*model.Moderator, error) {
	return d.db.GetModeratorByUserID(ctx, userID)
}

func (d *delphisBackend) GetModeratorByUserIDAndDiscussionID(ctx context.Context, userID, discussionID string) (*model.Moderator, error) {
	return d.db.GetModeratorByUserIDAndDiscussionID(ctx, userID, discussionID)
}

func (d *delphisBackend) GetModeratedDiscussionsByUserID(ctx context.Context, userID string) ([]*model.Discussion, error) {
	iter := d.db.GetModeratedDiscussionsByUserID(ctx, userID)
	return d.db.DiscussionIterCollect(ctx, iter)
}

func (d *delphisBackend) CheckIfModerator(ctx context.Context, userID string) (bool, error) {
	mod, err := d.GetModeratorByUserID(ctx, userID)
	if err != nil {
		logrus.WithError(err).Error("failed to get moderator by userID")
		return false, err
	}

	return mod != nil, nil
}

func (d *delphisBackend) CheckIfModeratorForDiscussion(ctx context.Context, userID string, discussionID string) (bool, error) {
	mod, err := d.GetModeratorByUserIDAndDiscussionID(ctx, userID, discussionID)
	if err != nil {
		logrus.WithError(err).Error("failed to get moderator by userID")
		return false, err
	}

	return mod != nil, nil
}

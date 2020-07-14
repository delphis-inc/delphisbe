package backend

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/delphis-inc/delphisbe/graph/model"
)

func (b *delphisBackend) GetModeratorByID(ctx context.Context, id string) (*model.Moderator, error) {
	return b.db.GetModeratorByID(ctx, id)
}

func (b *delphisBackend) GetModeratorByUserID(ctx context.Context, userID string) (*model.Moderator, error) {
	return b.db.GetModeratorByUserID(ctx, userID)
}

func (b *delphisBackend) GetModeratorByUserIDAndDiscussionID(ctx context.Context, userID, discussionID string) (*model.Moderator, error) {
	return b.db.GetModeratorByUserIDAndDiscussionID(ctx, userID, discussionID)
}

func (b *delphisBackend) CheckIfModerator(ctx context.Context, userID string) (bool, error) {
	mod, err := b.GetModeratorByUserID(ctx, userID)
	if err != nil {
		logrus.WithError(err).Error("failed to get moderator by userID")
		return false, err
	}

	return mod != nil, nil
}

func (b *delphisBackend) CheckIfModeratorForDiscussion(ctx context.Context, userID string, discussionID string) (bool, error) {
	mod, err := b.GetModeratorByUserIDAndDiscussionID(ctx, userID, discussionID)
	if err != nil {
		logrus.WithError(err).Error("failed to get moderator by userID")
		return false, err
	}

	return mod != nil, nil
}

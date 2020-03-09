package backend

import (
	"context"

	"github.com/nedrocks/delphisbe/graph/model"
)

func (d *delphisBackend) GetUserProfileByID(ctx context.Context, id string) (*model.UserProfile, error) {
	return d.db.GetUserProfileByID(ctx, id)
}

func (d *delphisBackend) AddModeratedDiscussionToUserProfile(ctx context.Context, userProfileID string, discussionID string) (*model.UserProfile, error) {
	return d.db.AddModeratedDiscussionToUserProfile(ctx, userProfileID, discussionID)
}

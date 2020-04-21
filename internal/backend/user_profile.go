package backend

import (
	"context"

	"github.com/nedrocks/delphisbe/graph/model"
)

func (d *delphisBackend) GetUserProfileByUserID(ctx context.Context, userID string) (*model.UserProfile, error) {
	return d.db.GetUserProfileByUserID(ctx, userID)
}

func (d *delphisBackend) GetUserProfileByID(ctx context.Context, id string) (*model.UserProfile, error) {
	return d.db.GetUserProfileByID(ctx, id)
}

func (d *delphisBackend) GetSocialInfosByUserProfileID(ctx context.Context, userProfileID string) ([]model.SocialInfo, error) {
	return d.db.GetSocialInfosByUserProfileID(ctx, userProfileID)
}

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

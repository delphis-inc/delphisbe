package backend

import (
	"context"

	"github.com/nedrocks/delphisbe/graph/model"
)

func (d *daoManager) GetUserProfileByID(ctx context.Context, id string) (*model.UserProfile, error) {
	return d.db.GetUserProfileByID(ctx, id)
}

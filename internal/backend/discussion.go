package backend

import (
	"context"

	"github.com/nedrocks/delphisbe/graph/model"
)

func (d *daoManager) GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error) {
	return d.db.GetDiscussionByID(ctx, id)
}

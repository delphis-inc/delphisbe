package backend

import (
	"context"

	"github.com/nedrocks/delphisbe/graph/model"
)

func (b *delphisBackend) GetModeratorByID(ctx context.Context, id string) (*model.Moderator, error) {
	return b.db.GetModeratorByID(ctx, id)
}

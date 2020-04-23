package backend

import (
	"context"

	"github.com/nedrocks/delphisbe/graph/model"
)

func (b *delphisBackend) GetPostContentByID(ctx context.Context, id string) (*model.PostContent, error) {
	return b.db.GetPostContentByID(ctx, id)
}

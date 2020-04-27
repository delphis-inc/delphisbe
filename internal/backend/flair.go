package backend

import (
	"context"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (d *delphisBackend) GetFlairByID(ctx context.Context, id string) (*model.Flair, error) {
	return d.db.GetFlairByID(ctx, id)
}

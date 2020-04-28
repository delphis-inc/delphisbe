package backend

import (
	"context"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (d *delphisBackend) GetFlairByID(ctx context.Context, id string) (*model.Flair, error) {
	return d.db.GetFlairByID(ctx, id)
}

func (d *delphisBackend) GetFlairByUserIDFlairID(ctx context.Context, userID string, flairID string) (*model.Flair, error) {
	return d.db.GetFlairByUserIDFlairID(ctx, userID, flairID)
}

func (d *delphisBackend) GetFlairsByUserID(ctx context.Context, userID string) ([]model.Flair, error) {
	return d.db.GetFlairsByUserID(ctx, userID)
}

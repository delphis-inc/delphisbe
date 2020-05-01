package backend

import (
	"context"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/util"
)

func (d *delphisBackend) CreateFlair(ctx context.Context, userID string, templateID string) (*model.Flair, error) {
	flair := model.Flair{
		ID:         util.UUIDv4(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		UserID:     userID,
		TemplateID: templateID,
	}

	_, err := d.db.UpsertFlair(ctx, flair)
	if err != nil {
		return nil, err
	}

	return &flair, nil
}

func (d *delphisBackend) GetFlairByID(ctx context.Context, id string) (*model.Flair, error) {
	return d.db.GetFlairByID(ctx, id)
}

func (d *delphisBackend) GetFlairsByUserID(ctx context.Context, userID string) ([]*model.Flair, error) {
	return d.db.GetFlairsByUserID(ctx, userID)
}

func (d *delphisBackend) RemoveFlair(ctx context.Context, flair model.Flair) (*model.Flair, error) {
	return d.db.RemoveFlair(ctx, flair)
}

package backend

import (
	"context"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/util"
)

func (d *delphisBackend) CreateNewFlair(ctx context.Context, displayName *string, imageURL *string, source string) (*model.Flair, error) {
	flairID := util.UUIDv4()
	flairObj := model.Flair{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ID:          flairID,
		DisplayName: displayName,
		ImageURL:    imageURL,
		Source:      source,
	}

	_, err := d.db.UpsertFlair(ctx, flairObj)

	if err != nil {
		return nil, err
	}

	return &flairObj, nil
}


func (d *delphisBackend) GetFlairByID(ctx context.Context, id string) (*model.Flair, error) {
	return d.db.GetFlairByID(ctx, id)
}

func (d *delphisBackend) GetFlairByUserIDFlairID(ctx context.Context, userID string, flairID string) (*model.Flair, error) {
	return d.db.GetFlairByUserIDFlairID(ctx, userID, flairID)
}

func (d *delphisBackend) GetFlairsByUserID(ctx context.Context, userID string) ([]model.Flair, error) {
	return d.db.GetFlairsByUserID(ctx, userID)
}

package backend

import (
	"context"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/util"
)

func (d *delphisBackend) CreateFlairTemplate(ctx context.Context, displayName *string, imageURL *string, source string) (*model.FlairTemplate, error) {
	flairTemplateObj := model.FlairTemplate{
		ID:          util.UUIDv4(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DisplayName: displayName,
		ImageURL:    imageURL,
		Source:      source,
	}

	_, err := d.db.UpsertFlairTemplate(ctx, flairTemplateObj)
	if err != nil {
		return nil, err
	}

	return &flairTemplateObj, nil
}

func (d *delphisBackend) GetFlairTemplateByID(ctx context.Context, id string) (*model.FlairTemplate, error) {
	return d.db.GetFlairTemplateByID(ctx, id)
}

func (d *delphisBackend) RemoveFlairTemplate(ctx context.Context, flairTemplate model.FlairTemplate) (*model.FlairTemplate, error) {
	return d.db.RemoveFlairTemplate(ctx, flairTemplate)
}

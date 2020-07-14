package backend

import (
	"context"
	"time"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/util"
)

func (d *delphisBackend) CreateViewerForDiscussion(ctx context.Context, discussionID string, userID string) (*model.Viewer, error) {
	viewerObj := model.Viewer{
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ID:           util.UUIDv4(),
		DiscussionID: &discussionID,
		UserID:       &userID,
	}

	_, err := d.db.UpsertViewer(ctx, viewerObj)

	if err != nil {
		return nil, err
	}

	return &viewerObj, err
}

func (d *delphisBackend) GetViewersByIDs(ctx context.Context, viewerIDs []string) (map[string]*model.Viewer, error) {
	return d.db.GetViewersByIDs(ctx, viewerIDs)
}

func (d *delphisBackend) GetViewerByID(ctx context.Context, viewerID string) (*model.Viewer, error) {
	viewers, err := d.GetViewersByIDs(ctx, []string{viewerID})

	if err != nil {
		return nil, err
	}

	return viewers[viewerID], nil
}

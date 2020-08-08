package backend

import (
	"context"

	"github.com/delphis-inc/delphisbe/graph/model"
)

func (d *delphisBackend) CreateViewerForDiscussion(ctx context.Context, discussionID string, userID string) (*model.Viewer, error) {
	return d.GetViewerForDiscussion(ctx, discussionID, userID, true)
}

func (d *delphisBackend) SetViewerLastPostViewed(ctx context.Context, viewerID, postID string) (*model.Viewer, error) {
	now := d.timeProvider.Now()

	return d.db.SetViewerLastPostViewed(ctx, viewerID, postID, now)
}

func (d *delphisBackend) GetViewerForDiscussion(ctx context.Context, discussionID, userID string, createIfNotFound bool) (*model.Viewer, error) {
	viewer, err := d.db.GetViewerForDiscussion(ctx, discussionID, userID)
	if err != nil {
		return nil, err
	}

	if viewer == nil && createIfNotFound {
		return d.CreateViewerForDiscussion(ctx, discussionID, userID)
	}

	return viewer, nil
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

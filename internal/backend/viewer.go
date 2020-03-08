package backend

import (
	"context"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/util"
)

func (d *daoManager) CreateViewerForDiscussion(ctx context.Context, discussionID string, userID string) (*model.Viewer, error) {
	viewerObj := model.Viewer{
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ID:           util.UUIDv4(),
		DiscussionID: discussionID,
		UserID:       userID,
	}

	_, err := d.db.PutViewer(ctx, viewerObj)

	if err != nil {
		return nil, err
	}

	return &viewerObj, err
}

func (d *daoManager) GetViewersByIDs(ctx context.Context, discussionViewerKeys []model.DiscussionViewerKey) (map[model.DiscussionViewerKey]*model.Viewer, error) {
	return d.db.GetViewersByIDs(ctx, discussionViewerKeys)
}

func (d *daoManager) GetViewerByID(ctx context.Context, discussionViewerKey model.DiscussionViewerKey) (*model.Viewer, error) {
	viewers, err := d.GetViewersByIDs(ctx, []model.DiscussionViewerKey{discussionViewerKey})

	if err != nil {
		return nil, err
	}

	return viewers[discussionViewerKey], nil
}

package backend

import (
	"context"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/util"
)

func (d *delphisBackend) CreateNewDiscussion(ctx context.Context, creatingUser *model.User, anonymityType model.AnonymityType) (*model.Discussion, error) {
	discussionID := util.UUIDv4()
	discussionObj := model.Discussion{
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ID:            discussionID,
		AnonymityType: anonymityType,
		Moderator: model.Moderator{
			ID:            discussionID,
			UserProfileID: creatingUser.UserProfileID,
			DiscussionID:  discussionID,
		},
	}

	_, err := d.db.PutDiscussion(ctx, discussionObj)

	if err != nil {
		return nil, err
	}

	return &discussionObj, nil
}

func (d *delphisBackend) GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error) {
	return d.db.GetDiscussionByID(ctx, id)
}

func (d *delphisBackend) GetDiscussionsByIDs(ctx context.Context, ids []string) (map[string]*model.Discussion, error) {
	return d.db.GetDiscussionsByIDs(ctx, ids)
}

func (d *delphisBackend) ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error) {
	return d.db.ListDiscussions(ctx)
}

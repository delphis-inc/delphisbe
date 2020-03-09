package backend

import (
	"context"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/util"
)

func (d *delphisBackend) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	return d.db.GetUserByID(ctx, userID)
}

// AddParticipantToUser updates a user object to refer to the participant
// and viewer.
func (d *delphisBackend) AddParticipantAndViewerToUser(ctx context.Context, userID string, participantID int, discussionID string, viewerID string) (*model.User, error) {
	_, err := d.db.AddParticipantToUser(ctx, userID, model.DiscussionParticipantKey{
		DiscussionID:  discussionID,
		ParticipantID: participantID,
	})

	if err != nil {
		return nil, err
	}

	user, err := d.AddViewerToUser(ctx, userID, discussionID, viewerID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *delphisBackend) AddViewerToUser(ctx context.Context, userID, discussionID, viewerID string) (*model.User, error) {
	return d.db.AddViewerToUser(ctx, userID, model.DiscussionViewerKey{
		DiscussionID: discussionID,
		ViewerID:     viewerID,
	})
}

func (d *delphisBackend) CreateUser(ctx context.Context) (*model.User, error) {
	userObj := model.User{
		ID:        util.UUIDv4(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	resp, err := d.db.PutUser(ctx, userObj)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

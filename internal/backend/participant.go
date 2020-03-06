package backend

import (
	"context"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
)

func (d *daoManager) CreateParticipantForDiscussion(ctx context.Context, discussionID string, userID string) (*model.Participant, error) {
	allParticipants, err := d.GetParticipantsByDiscussionID(ctx, discussionID)

	if err != nil {
		return nil, err
	}

	viewerObj, err := d.CreateViewerForDiscussion(ctx, discussionID, userID)

	if err != nil {
		return nil, err
	}

	participantObj := model.Participant{
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ParticipantID: len(allParticipants),
		DiscussionID:  discussionID,
		// TODO: Create the viewer and add it here.
		ViewerID: viewerObj.ID,
		UserID:   userID,
	}

	return d.db.PutParticipant(ctx, participantObj)
}

func (d *daoManager) GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error) {
	return d.db.GetParticipantsByDiscussionID(ctx, id)
}

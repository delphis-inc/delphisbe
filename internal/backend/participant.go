package backend

import (
	"context"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
)

func (d *delphisBackend) CreateParticipantForDiscussion(ctx context.Context, discussionID string, userID string) (*model.Participant, error) {
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

	res, err := d.db.PutParticipant(ctx, participantObj)

	if err != nil {
		return nil, err
	}

	_, err = d.AddParticipantAndViewerToUser(ctx, userID, res.ParticipantID, discussionID, viewerObj.ID)

	if err != nil {
		return nil, err
	}

	return &participantObj, nil
}

func (d *delphisBackend) GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error) {
	return d.db.GetParticipantsByDiscussionID(ctx, id)
}

func (d *delphisBackend) GetParticipantByID(ctx context.Context, discussionParticipantKey model.DiscussionParticipantKey) (*model.Participant, error) {
	respMap, err := d.GetParticipantsByIDs(ctx, []model.DiscussionParticipantKey{discussionParticipantKey})

	if err != nil {
		return nil, err
	}

	return respMap[discussionParticipantKey], nil
}

func (d *delphisBackend) GetParticipantsByIDs(ctx context.Context, discussionParticipantKeys []model.DiscussionParticipantKey) (map[model.DiscussionParticipantKey]*model.Participant, error) {
	return d.db.GetParticipantsByIDs(ctx, discussionParticipantKeys)
}

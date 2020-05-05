package backend

import (
	"context"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/util"
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
		ID:            util.UUIDv4(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ParticipantID: len(allParticipants),
		DiscussionID:  &discussionID,
		ViewerID:      &viewerObj.ID,
		UserID:        &userID,
	}

	_, err = d.db.PutParticipant(ctx, participantObj)

	if err != nil {
		return nil, err
	}

	return &participantObj, nil
}

func (d *delphisBackend) GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error) {
	return d.db.GetParticipantsByDiscussionID(ctx, id)
}

func (d *delphisBackend) GetParticipantByDiscussionIDUserID(ctx context.Context, discussionID string, userID string) (*model.Participant, error) {
	return d.db.GetParticipantByDiscussionIDUserID(ctx, discussionID, userID)
}

func (d *delphisBackend) GetParticipantByID(ctx context.Context, id string) (*model.Participant, error) {
	participant, err := d.db.GetParticipantByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return participant, nil
}

func (d *delphisBackend) AssignFlair(ctx context.Context, participant model.Participant, flairID string) (*model.Participant, error) {
	return d.db.AssignFlair(ctx, participant, &flairID)
}

func (d *delphisBackend) UnassignFlair(ctx context.Context, participant model.Participant) (*model.Participant, error) {
	return d.db.AssignFlair(ctx, participant, nil)
}

func (d *delphisBackend) GetTotalParticipantCountByDiscussionID(ctx context.Context, discussionID string) int {
	return d.db.GetTotalParticipantCountByDiscussionID(ctx, discussionID)
}

func (d *delphisBackend) CopyAndUpdateParticipant(ctx context.Context, orig model.Participant, input model.UpdateParticipantInput) (*model.Participant, error) {
	participantCount := d.GetTotalParticipantCountByDiscussionID(ctx, *orig.DiscussionID)
	now := time.Now()
	copiedObj := orig
	copiedObj.ID = util.UUIDv4()
	copiedObj.CreatedAt = now
	copiedObj.UpdatedAt = now
	if input.GradientColor != nil || (input.IsUnsetGradient != nil && *input.IsUnsetGradient) {
		copiedObj.GradientColor = input.GradientColor
	}
	if input.FlairID != nil || (input.IsUnsetFlairID != nil && *input.IsUnsetFlairID) {
		copiedObj.FlairID = input.FlairID
	}
	if input.IsAnonymous != nil {
		copiedObj.IsAnonymous = *input.IsAnonymous
	}
	copiedObj.ParticipantID = participantCount

	return d.db.PutParticipant(ctx, copiedObj)
}

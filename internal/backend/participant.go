package backend

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/util"
)

func (d *delphisBackend) CreateParticipantForDiscussion(ctx context.Context, discussionID string, userID string, discussionParticipantInput model.AddDiscussionParticipantInput) (*model.Participant, error) {
	userObj, err := d.GetUserByID(ctx, userID)
	if err != nil || userObj == nil {
		if userObj == nil {
			err = fmt.Errorf("Could not find User with ID %s so failing creation of Participant", userID)
		}
		return nil, err
	}

	allParticipantCount := d.GetTotalParticipantCountByDiscussionID(ctx, discussionID)

	participantObj := model.Participant{
		ID:            util.UUIDv4(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ParticipantID: allParticipantCount,
		DiscussionID:  &discussionID,
		UserID:        &userID,
	}

	if discussionParticipantInput.GradientColor != nil {
		participantObj.GradientColor = discussionParticipantInput.GradientColor
	} else {
		gradientColor := model.GradientColorUnknown
		for gradientColor == model.GradientColorUnknown {
			gradientColor = model.AllGradientColor[rand.Intn(len(model.AllGradientColor))]
		}
		// TODO: We need to create a unique gradient color / name pairing once we have names.
		participantObj.GradientColor = &gradientColor
	}

	if discussionParticipantInput.FlairID != nil {
		if userObj.Flairs == nil {
			userObj.Flairs, err = d.GetFlairsByUserID(ctx, userID)
			if err == nil {
				return nil, err
			}
		}
		if len(userObj.Flairs) > 0 {
			for _, elem := range userObj.Flairs {
				if elem != nil && elem.ID == *discussionParticipantInput.FlairID {
					participantObj.FlairID = discussionParticipantInput.FlairID
				}
			}
		}
	}

	participantObj.HasJoined = discussionParticipantInput.HasJoined != nil && *discussionParticipantInput.HasJoined
	participantObj.IsAnonymous = discussionParticipantInput.IsAnonymous

	viewerObj, err := d.CreateViewerForDiscussion(ctx, discussionID, userID)

	if err != nil {
		return nil, err
	}

	participantObj.ViewerID = &viewerObj.ID

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
	if input.HasJoined != nil {
		copiedObj.HasJoined = *input.HasJoined
	}
	copiedObj.ParticipantID = participantCount

	return d.db.PutParticipant(ctx, copiedObj)
}

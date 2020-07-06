package backend

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/util"
)

type UserDiscussionParticipants struct {
	Anon    *model.Participant
	NonAnon *model.Participant
}

// TODO: Should this be gated to only allow two participants from the same UserID (one anonymous and the other not)?
// Where are we gating whether a user can add another participant?
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
			if err != nil {
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

	_, err = d.db.UpsertParticipant(ctx, participantObj)
	if err != nil {
		return nil, err
	}

	if _, err := d.CreateAlertPost(ctx, discussionID, userObj, discussionParticipantInput.IsAnonymous); err != nil {
		// Don't return err. Just log
		logrus.WithError(err).Error("failed to create alert post")
	}

	return &participantObj, nil
}

func (d *delphisBackend) BanParticipant(ctx context.Context, discussionID string, participantID string, requestingUserID string) (*model.Participant, error) {
	discussionObj, err := d.GetDiscussionByID(ctx, discussionID)
	if err != nil || discussionObj == nil {
		return nil, fmt.Errorf("Failed to retrieve discussion")
	}

	if discussionObj.ModeratorID == nil || requestingUserID != *discussionObj.ModeratorID {
		return nil, fmt.Errorf("Only the moderator may ban users")
	}

	participantObj, err := d.GetParticipantByID(ctx, participantID)
	if err != nil || participantObj == nil {
		return nil, fmt.Errorf("Failed to retreive participant")
	}

	if participantObj.DiscussionID == nil || *participantObj.DiscussionID != discussionID {
		return nil, fmt.Errorf("Participant is not part of this discussion")
	}

	if participantObj.IsBanned {
		return participantObj, nil
	}

	participantObj.IsBanned = true
	updatedParticipant, err := d.db.UpsertParticipant(ctx, *participantObj)

	if err != nil {
		logrus.WithError(err).Error("Failed to update participant")
		return nil, err
	}

	_, err = d.db.DeleteAllParticipantPosts(ctx, discussionID, participantID, model.PostDeletedReasonParticipantRemoved)

	if err != nil {
		logrus.WithError(err).Error("Failed to delete participant posts.")
		// Do not return error here.
	}

	return updatedParticipant, nil
}

func (d *delphisBackend) GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error) {
	return d.db.GetParticipantsByDiscussionID(ctx, id)
}

func (d *delphisBackend) GetParticipantsByDiscussionIDUserID(ctx context.Context, discussionID string, userID string) (*UserDiscussionParticipants, error) {
	participants, err := d.db.GetParticipantsByDiscussionIDUserID(ctx, discussionID, userID)
	if err != nil {
		return nil, err
	}

	participantResponse := &UserDiscussionParticipants{}

	for i, participant := range participants {

		if participant.IsAnonymous && participantResponse.Anon == nil {
			participantResponse.Anon = &participants[i]
		}
		if !participant.IsAnonymous && participantResponse.NonAnon == nil {
			participantResponse.NonAnon = &participants[i]
		}
	}
	return participantResponse, nil
}

func (d *delphisBackend) GetParticipantByID(ctx context.Context, id string) (*model.Participant, error) {
	participant, err := d.db.GetParticipantByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return participant, nil
}

func (d *delphisBackend) GetParticipantsByIDs(ctx context.Context, ids []string) (map[string]*model.Participant, error) {
	return d.db.GetParticipantsByIDs(ctx, ids)
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

func (d *delphisBackend) UpdateParticipant(ctx context.Context, participants UserDiscussionParticipants, currentParticipantID string, input model.UpdateParticipantInput) (*model.Participant, error) {
	var currentParticipantObj *model.Participant
	var otherParticipantObj *model.Participant
	if participants.Anon != nil && participants.Anon.ID == currentParticipantID {
		currentParticipantObj = participants.Anon
		otherParticipantObj = participants.NonAnon
	} else if participants.NonAnon != nil && participants.NonAnon.ID == currentParticipantID {
		currentParticipantObj = participants.NonAnon
		otherParticipantObj = participants.Anon
	}

	if currentParticipantObj == nil {
		return nil, fmt.Errorf("Failed to find participant with ID %s", currentParticipantID)
	}

	if input.IsAnonymous != nil && *input.IsAnonymous != currentParticipantObj.IsAnonymous {
		// We are changing the participant here. Potentially creating a new one.
		if otherParticipantObj == nil {
			// We have to create a new one.
			participantCount := d.GetTotalParticipantCountByDiscussionID(ctx, *currentParticipantObj.DiscussionID)
			now := time.Now()
			copiedObj := *currentParticipantObj
			copiedObj.ParticipantID = participantCount
			copiedObj.ID = util.UUIDv4()
			copiedObj.CreatedAt = now
			copiedObj.UpdatedAt = now
			copiedObj.IsAnonymous = *input.IsAnonymous
			currentParticipantObj = &copiedObj
		} else {
			// In this case we can use the other participant object.
			currentParticipantObj = otherParticipantObj
		}
	}
	if input.GradientColor != nil || (input.IsUnsetGradient != nil && *input.IsUnsetGradient) {
		currentParticipantObj.GradientColor = input.GradientColor
	}
	if input.FlairID != nil || (input.IsUnsetFlairID != nil && *input.IsUnsetFlairID) {
		currentParticipantObj.FlairID = input.FlairID
	}
	if input.HasJoined != nil {
		// Cannot unjoin a conversation.
		if !currentParticipantObj.HasJoined {
			currentParticipantObj.HasJoined = *input.HasJoined
		}
	}

	return d.db.UpsertParticipant(ctx, *currentParticipantObj)
}

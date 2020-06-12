package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/auth"
	"github.com/sirupsen/logrus"
)

func (r *mutationResolver) AddDiscussionParticipant(ctx context.Context, discussionID string, userID string, discussionParticipantInput model.AddDiscussionParticipantInput) (*model.Participant, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	if authedUser.UserID != userID {
		// Check if moderator is trying to add discussion participant
		modCheck, err := r.DAOManager.CheckIfModeratorForDiscussion(ctx, authedUser.UserID, discussionID)
		if err != nil || !modCheck {
			return nil, fmt.Errorf("unauthorized")
		}
	}

	participantObj, err := r.DAOManager.CreateParticipantForDiscussion(ctx, discussionID, authedUser.UserID, discussionParticipantInput)
	if err != nil {
		return nil, err
	}

	return participantObj, nil
}

func (r *mutationResolver) AddPost(ctx context.Context, discussionID string, participantID string, postContent model.PostContentInput) (*model.Post, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	// Note: This is here mainly to ensure the discussion is not (soft) deleted
	discussion, err := r.DAOManager.GetDiscussionByID(ctx, discussionID)
	if discussion == nil || err != nil {
		return nil, fmt.Errorf("Discussion with ID %s not found", discussionID)
	}

	//participant, err := r.DAOManager.GetParticipantByDiscussionIDUserID(ctx, discussionID, authedUser.UserID)
	participant, err := r.DAOManager.GetParticipantByID(ctx, participantID)
	if err != nil {
		return nil, err
	} else if participant == nil {
		return nil, fmt.Errorf("Could not find Participant with ID %s", participantID)
	}

	// Verify that the posting participant belongs to the logged-in user
	if *participant.UserID != authedUser.UserID {
		return nil, fmt.Errorf("Unauthorized")
	}

	createdPost, err := r.DAOManager.CreatePost(ctx, discussionID, participant.ID, postContent)
	if err != nil {
		return nil, fmt.Errorf("Failed to create post")
	}

	err = r.DAOManager.NotifySubscribersOfCreatedPost(ctx, createdPost, discussionID)
	if err != nil {
		// Silently ignore this
		logrus.Warnf("Failed to notify subscribers of created post")
	}

	return createdPost, nil
}

func (r *mutationResolver) PostImportedContent(ctx context.Context, discussionID string, participantID string, contentID string) (*model.Post, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	// Only allow the mod to manual post imported content from the drip
	modCheck, err := r.DAOManager.CheckIfModeratorForDiscussion(ctx, authedUser.UserID, discussionID)
	if err != nil || !modCheck {
		return nil, fmt.Errorf("unauthorized")
	}

	participant, err := r.DAOManager.GetParticipantByID(ctx, participantID)
	if err != nil {
		return nil, err
	} else if participant == nil {
		return nil, fmt.Errorf("Could not find Participant with ID %s", participantID)
	}

	// Verify that the posting participant belongs to the logged-in user
	if *participant.UserID != authedUser.UserID {
		return nil, fmt.Errorf("Unauthorized")
	}

	now := time.Now()
	return r.DAOManager.PostImportedContent(ctx, participantID, discussionID, contentID, &now, nil, model.ManualDrip)
}

func (r *mutationResolver) ScheduleImportedContent(ctx context.Context, discussionID string, contentID string) (*model.ContentQueueRecord, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	// Only allow the mod to manual post imported content from the drip
	modCheck, err := r.DAOManager.CheckIfModeratorForDiscussion(ctx, authedUser.UserID, discussionID)
	if err != nil || !modCheck {
		return nil, fmt.Errorf("unauthorized")
	}

	return r.DAOManager.PutImportedContentQueue(ctx, discussionID, contentID, nil, nil, "")
}

func (r *mutationResolver) CreateDiscussion(ctx context.Context, anonymityType model.AnonymityType, title string) (*model.Discussion, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}
	if authedUser.User == nil {
		var err error
		authedUser.User, err = r.DAOManager.GetUserByID(ctx, authedUser.UserID)
		if err != nil {
			return nil, fmt.Errorf("Error fetching user with ID (%s)", authedUser.UserID)
		}
	}

	discussionObj, err := r.DAOManager.CreateNewDiscussion(ctx, authedUser.User, anonymityType, title)

	if err != nil {
		return nil, err
	}

	trueObj := true
	_, err = r.DAOManager.CreateParticipantForDiscussion(ctx, discussionObj.ID, authedUser.UserID, model.AddDiscussionParticipantInput{HasJoined: &trueObj})

	if err != nil {
		return nil, err
	}

	return discussionObj, nil
}

func (r *mutationResolver) CreateFlair(ctx context.Context, userID string, templateID string) (*model.Flair, error) {
	if userID == "" || templateID == "" {
		return nil, fmt.Errorf("parameters cannot be empty")
	}

	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	// This API takes in userID so that an admin or moderator
	// could potentially create flair on behalf of another user.
	if authedUser.UserID != userID {
		// Check if a moderator is trying to addFlair
		modCheck, err := r.DAOManager.CheckIfModerator(ctx, authedUser.UserID)
		if err != nil || !modCheck {
			return nil, fmt.Errorf("unauthorized")
		}
	}

	// TODO: Add validation and verification mechanisms.
	//
	// flairTemplate, err := r.DAOManager.GetFlairTemplateByID(ctx, templateID)
	// if err != nil {
	// 	return nil, err
	// } else if flairTemplate == nil {
	// 	return nil, fmt.Errorf("Error fetching flair template with ID (%s)", templateID)
	// }
	//
	// Validation should ensure that this flair template makes sense to add to
	// this user. i.e. if the  template is only for companies or the user
	// currently has a contradictory flair, then this would not be valid.
	// flairTemplate.validate(userID)
	//
	// Verification should use a third-party (someone/something other than the
	// user receiving the flair) to authenticate the claim. This could be a
	// synchronous API request to some verification service, or an asynchronous
	// request that fills in the `verified_at` field on this flair hours later.
	// flairTemplate.verify(userID)

	return r.DAOManager.CreateFlair(ctx, userID, templateID)
}

func (r *mutationResolver) RemoveFlair(ctx context.Context, id string) (*model.Flair, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	flair, err := r.DAOManager.GetFlairByID(ctx, id)
	if err != nil {
		return nil, err
	} else if flair == nil {
		return nil, fmt.Errorf("Already removed")
	}

	// An admin or moderator could potentially remove flair on
	// behalf of another user.
	if authedUser.UserID != flair.UserID {
		// Check if a moderator is trying to removeFlair
		modCheck, err := r.DAOManager.CheckIfModerator(ctx, authedUser.UserID)
		if err != nil || !modCheck {
			return nil, fmt.Errorf("unauthorized")
		}
	}

	return r.DAOManager.RemoveFlair(ctx, *flair)
}

func (r *mutationResolver) AssignFlair(ctx context.Context, participantID string, flairID string) (*model.Participant, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	flair, err := r.DAOManager.GetFlairByID(ctx, flairID)
	if err != nil {
		return nil, err
	} else if flair == nil {
		return nil, fmt.Errorf("Error fetching flair with ID (%s)", flairID)
	}

	participant, err := r.DAOManager.GetParticipantByID(ctx, participantID)
	if err != nil {
		return nil, err
	} else if participant == nil {
		return nil, fmt.Errorf("Error fetching participant with ID (%s)", participantID)
	}

	if *participant.UserID != flair.UserID {
		// Flair does not belong to this participant, bad request.
		// Note: This error message is intentionally ambiguous as to not reveal
		// the lack of flair for a participant.
		return nil, fmt.Errorf("Unathorized")
	}

	// This API takes in userID so that an admin or moderator
	// could potentially assign flair on behalf of another participant.
	if authedUser.UserID != *participant.UserID {
		// Check if a moderator is trying to assignFlair
		modCheck, err := r.DAOManager.CheckIfModerator(ctx, authedUser.UserID)
		if err != nil || !modCheck {
			return nil, fmt.Errorf("unauthorized")
		}
	}

	return r.DAOManager.AssignFlair(ctx, *participant, flairID)
}

func (r *mutationResolver) UnassignFlair(ctx context.Context, participantID string) (*model.Participant, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	participant, err := r.DAOManager.GetParticipantByID(ctx, participantID)
	if err != nil {
		return nil, err
	} else if participant == nil {
		return nil, fmt.Errorf("Error fetching participant with ID (%s)", participantID)
	}

	// This API takes in userID so that an admin or
	// moderator could potentially unassign flair on behalf of another
	// participant.
	if authedUser.UserID != *participant.UserID {
		// Check if a moderator is trying to assignFlair
		modCheck, err := r.DAOManager.CheckIfModerator(ctx, authedUser.UserID)
		if err != nil || !modCheck {
			return nil, fmt.Errorf("unauthorized")
		}
	}

	return r.DAOManager.UnassignFlair(ctx, *participant)
}

func (r *mutationResolver) CreateFlairTemplate(ctx context.Context, displayName *string, imageURL *string, source string) (*model.FlairTemplate, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}
	// TODO: Only admins should be able to create flair templates

	return r.DAOManager.CreateFlairTemplate(ctx, displayName, imageURL, source)
}

func (r *mutationResolver) RemoveFlairTemplate(ctx context.Context, id string) (*model.FlairTemplate, error) {
	currentUser := auth.GetAuthedUser(ctx)
	if currentUser == nil {
		return nil, fmt.Errorf("Need auth")
	}
	// TODO: Only admins should be able to remove flair templates

	flairTemplate, err := r.DAOManager.GetFlairTemplateByID(ctx, id)
	if err != nil {
		return nil, err
	} else if flairTemplate == nil {
		return nil, fmt.Errorf("Already removed")
	}

	return r.DAOManager.RemoveFlairTemplate(ctx, *flairTemplate)
}

func (r *mutationResolver) UpdateParticipant(ctx context.Context, discussionID string, participantID string, updateInput model.UpdateParticipantInput) (*model.Participant, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	participantResponse, err := r.DAOManager.GetParticipantsByDiscussionIDUserID(ctx, discussionID, authedUser.UserID)
	if err != nil {
		return nil, err
	}
	if participantResponse == nil {
		return nil, fmt.Errorf("Failed to find participant with ID %s", participantID)
	}

	// Verify that the updating participant belongs to the logged-in user
	var nonAnonUserID, anonUserID string
	if participantResponse.NonAnon != nil && participantResponse.NonAnon.ID == participantID {
		nonAnonUserID = *participantResponse.NonAnon.UserID
	}
	if participantResponse.Anon != nil && participantResponse.Anon.ID == participantID {
		anonUserID = *participantResponse.Anon.UserID
	}
	if authedUser.UserID != nonAnonUserID && authedUser.UserID != anonUserID {
		return nil, fmt.Errorf("Unauthorized")
	}

	res, err := r.DAOManager.UpdateParticipant(ctx, *participantResponse, participantID, updateInput)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *mutationResolver) UpsertUserDevice(ctx context.Context, userID *string, platform model.Platform, deviceID string, token *string) (*model.UserDevice, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	if authedUser.UserID != *userID {
		// Check if the mod is trying to upsert the device
		modCheck, err := r.DAOManager.CheckIfModerator(ctx, authedUser.UserID)
		if err != nil || !modCheck {
			return nil, fmt.Errorf("unauthorized")
		}
	}

	resp, err := r.DAOManager.UpsertUserDevice(ctx, deviceID, &authedUser.UserID, platform.String(), token)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *mutationResolver) UpdateDiscussion(ctx context.Context, discussionID string, input model.DiscussionInput) (*model.Discussion, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	// Only allow the mod to update discussion
	modCheck, err := r.DAOManager.CheckIfModeratorForDiscussion(ctx, authedUser.UserID, discussionID)
	if err != nil || !modCheck {
		return nil, fmt.Errorf("unauthorized")
	}

	return r.DAOManager.UpdateDiscussion(ctx, discussionID, input)
}

func (r *mutationResolver) AddDiscussionTags(ctx context.Context, discussionID string, tags []string) ([]*model.Tag, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	// Only allow the mod to add discussion tags
	modCheck, err := r.DAOManager.CheckIfModeratorForDiscussion(ctx, authedUser.UserID, discussionID)
	if err != nil || !modCheck {
		return nil, fmt.Errorf("unauthorized")
	}

	return r.DAOManager.PutDiscussionTags(ctx, discussionID, tags)
}

func (r *mutationResolver) DeleteDiscussionTags(ctx context.Context, discussionID string, tags []string) ([]*model.Tag, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	// Only allow the mod to delete discussion tags
	modCheck, err := r.DAOManager.CheckIfModeratorForDiscussion(ctx, authedUser.UserID, discussionID)
	if err != nil || !modCheck {
		return nil, fmt.Errorf("unauthorized")
	}

	return r.DAOManager.DeleteDiscussionTags(ctx, discussionID, tags)
}

func (r *queryResolver) Discussion(ctx context.Context, id string) (*model.Discussion, error) {
	return r.resolveDiscussionByID(ctx, id)
}

func (r *queryResolver) ListDiscussions(ctx context.Context) ([]*model.Discussion, error) {
	connection, err := r.DAOManager.ListDiscussions(ctx)
	if err != nil {
		return nil, err
	}

	discussions := make([]*model.Discussion, 0)
	for i, edge := range connection.Edges {
		if edge != nil {
			discussions = append(discussions, connection.Edges[i].Node)
		}
	}
	return discussions, nil
}

func (r *queryResolver) FlairTemplates(ctx context.Context, query *string) ([]*model.FlairTemplate, error) {
	return r.DAOManager.ListFlairTemplates(ctx, query)
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	userObj, err := r.DAOManager.GetUserByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return userObj, err
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}
	if authedUser.User == nil {
		var err error
		authedUser.User, err = r.DAOManager.GetUserByID(ctx, authedUser.UserID)
		if err != nil {
			return nil, fmt.Errorf("Error fetching user with ID (%s)", authedUser.UserID)
		}
	}

	return authedUser.User, nil
}

func (r *subscriptionResolver) PostAdded(ctx context.Context, discussionID string) (<-chan *model.Post, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}
	events := make(chan *model.Post, 1)

	go func() {
		<-ctx.Done()
		err := r.DAOManager.UnSubscribeFromDiscussion(ctx, authedUser.UserID, discussionID)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to unsubscribe from discussion")
		}
		close(events)
	}()

	err := r.DAOManager.SubscribeToDiscussion(ctx, authedUser.UserID, events, discussionID)
	if err != nil {
		close(events)
		return nil, err
	}

	return events, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

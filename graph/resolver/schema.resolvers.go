// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"
	"fmt"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/auth"
	"github.com/sirupsen/logrus"
)

func (r *mutationResolver) AddDiscussionParticipant(ctx context.Context, discussionID string, userID string) (*model.Participant, error) {
	participantObj, err := r.DAOManager.CreateParticipantForDiscussion(ctx, discussionID, userID)

	// TODO: Only the current user can join the conversation
	if err != nil {
		return nil, err
	}

	return participantObj, nil
}

func (r *mutationResolver) AddPost(ctx context.Context, discussionID string, postContent string) (*model.Post, error) {
	creatingUser := auth.GetAuthedUser(ctx)
	if creatingUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	if creatingUser.User == nil {
		user, err := r.DAOManager.GetUserByID(ctx, creatingUser.UserID)
		if err != nil {
			return nil, fmt.Errorf("Error finding user with ID: %s", creatingUser.UserID)
		}
		creatingUser.User = user
	}

	// TODO: Don't need to look up discussion here since it's already handled in the participant lookup below
	discussion, err := r.DAOManager.GetDiscussionByID(ctx, discussionID)
	if discussion == nil || err != nil {
		return nil, fmt.Errorf("Discussion with ID %s not found", discussionID)
	}

	participant, err := r.DAOManager.GetParticipantByDiscussionIDUserID(ctx, discussionID, creatingUser.UserID)
	if err != nil {
		return nil, err
	} else if participant == nil {
		return nil, fmt.Errorf("Current user not a participant in this discussion")
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

func (r *mutationResolver) CreateDiscussion(ctx context.Context, anonymityType model.AnonymityType, title string) (*model.Discussion, error) {
	creatingUser := auth.GetAuthedUser(ctx)
	if creatingUser == nil {
		return nil, fmt.Errorf("Need auth")
	}
	if creatingUser.User == nil {
		var err error
		creatingUser.User, err = r.DAOManager.GetUserByID(ctx, creatingUser.UserID)
		if err != nil {
			return nil, fmt.Errorf("Error fetching user with ID (%s)", creatingUser.UserID)
		}
	}

	discussionObj, err := r.DAOManager.CreateNewDiscussion(ctx, creatingUser.User, anonymityType, title)

	if err != nil {
		return nil, err
	}

	_, err = r.DAOManager.CreateParticipantForDiscussion(ctx, discussionObj.ID, creatingUser.UserID)

	if err != nil {
		return nil, err
	}

	return discussionObj, nil
}

func (r *mutationResolver) CreateFlair(ctx context.Context, userID string, templateID string) (*model.Flair, error) {
	currentUser := auth.GetAuthedUser(ctx)
	if currentUser == nil {
		return nil, fmt.Errorf("Need auth")
	}
	// TODO: Currently this only allows the requesting user to create flair for
	// themselves. But this API takes in userID so that an admin or moderator
	// could potentially create flair on behalf of another user. To add this
	// functionality remove the following line and add authorization/permission
	// logic.
	if currentUser.UserID != userID {
		return nil, fmt.Errorf("Unauthorized")
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
	currentUser := auth.GetAuthedUser(ctx)
	if currentUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	flair, err := r.DAOManager.GetFlairByID(ctx, id)
	if err != nil {
		return nil, err
	} else if flair == nil {
		return nil, fmt.Errorf("Already removed")
	}

	// TODO: Currently this only allows the requesting user to remove flair for
	// themselves. But an admin or moderator could potentially remove flair on
	// behalf of another user. To add this functionality remove the following
	// line and add authorization/permission logic.
	if currentUser.UserID != flair.UserID {
		return nil, fmt.Errorf("Unauthorized")
	}

	return r.DAOManager.RemoveFlair(ctx, *flair)
}

func (r *mutationResolver) AssignFlair(ctx context.Context, participantID string, flairID string) (*model.Participant, error) {
	currentUser := auth.GetAuthedUser(ctx)
	if currentUser == nil {
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

	// TODO: Currently this only allows the requesting user to assign flair to
	// their own participant. But this API takes in userID so that an admin or
	// moderator could potentially assign flair on behalf of another
	// participant. To add this functionality remove the following line and
	// add authorization/permission logic.
	if currentUser.UserID != *participant.UserID {
		return nil, fmt.Errorf("Unauthorized")
	}

	return r.DAOManager.AssignFlair(ctx, *participant, flairID)
}

func (r *mutationResolver) UnassignFlair(ctx context.Context, participantID string) (*model.Participant, error) {
	currentUser := auth.GetAuthedUser(ctx)
	if currentUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	participant, err := r.DAOManager.GetParticipantByID(ctx, participantID)
	if err != nil {
		return nil, err
	} else if participant == nil {
		return nil, fmt.Errorf("Error fetching participant with ID (%s)", participantID)
	}

	// TODO: Currently this only allows the requesting user to unassign flair
	// from their own participant. But this API takes in userID so that an admin
	// or moderator could potentially unassign flair on behalf of another
	// participant. To add this functionality remove the following line and
	// add authorization/permission logic.
	if currentUser.UserID != *participant.UserID {
		return nil, fmt.Errorf("Unauthorized")
	}

	return r.DAOManager.UnassignFlair(ctx, *participant)
}

func (r *mutationResolver) CreateFlairTemplate(ctx context.Context, displayName *string, imageURL *string, source string) (*model.FlairTemplate, error) {
	currentUser := auth.GetAuthedUser(ctx)
	if currentUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	return r.DAOManager.CreateFlairTemplate(ctx, displayName, imageURL, source)
}

func (r *mutationResolver) RemoveFlairTemplate(ctx context.Context, id string) (*model.FlairTemplate, error) {
	currentUser := auth.GetAuthedUser(ctx)
	if currentUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	flairTemplate, err := r.DAOManager.GetFlairTemplateByID(ctx, id)
	if err != nil {
		return nil, err
	} else if flairTemplate == nil {
		return nil, fmt.Errorf("Already removed")
	}

	return r.DAOManager.RemoveFlairTemplate(ctx, *flairTemplate)
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
	for _, edge := range connection.Edges {
		if edge != nil {
			discussions = append(discussions, edge.Node)
		}
	}
	return discussions, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	userObj, err := r.DAOManager.GetUserByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return userObj, err
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	creatingUser := auth.GetAuthedUser(ctx)
	if creatingUser == nil {
		return nil, fmt.Errorf("Need auth")
	}
	if creatingUser.User == nil {
		var err error
		creatingUser.User, err = r.DAOManager.GetUserByID(ctx, creatingUser.UserID)
		if err != nil {
			return nil, fmt.Errorf("Error fetching user with ID (%s)", creatingUser.UserID)
		}
	}

	return creatingUser.User, nil
}

func (r *subscriptionResolver) PostAdded(ctx context.Context, discussionID string) (<-chan *model.Post, error) {
	currentUser := auth.GetAuthedUser(ctx)
	if currentUser == nil {
		return nil, fmt.Errorf("Need auth")
	}
	events := make(chan *model.Post, 1)

	go func() {
		<-ctx.Done()
		err := r.DAOManager.UnSubscribeFromDiscussion(ctx, currentUser.UserID, discussionID)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to unsubscribe from discussion")
		}
		close(events)
	}()

	err := r.DAOManager.SubscribeToDiscussion(ctx, currentUser.UserID, events, discussionID)
	if err != nil {
		close(events)
		return nil, err
	}

	return events, nil
}

func (r *Resolver) Mutation() generated.MutationResolver         { return &mutationResolver{r} }
func (r *Resolver) Query() generated.QueryResolver               { return &queryResolver{r} }
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

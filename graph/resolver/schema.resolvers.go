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

func (r *mutationResolver) AddFlair(ctx context.Context, discussionID string, flairID string) (*model.Participant, error) {
	creatingUser := auth.GetAuthedUser(ctx)
	if creatingUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	participant, err := r.DAOManager.GetParticipantByDiscussionIDUserID(ctx, discussionID, creatingUser.UserID)
	if err != nil {
		return nil, err
	} else if participant == nil {
		return nil, fmt.Errorf("Current user not a participant in this discussion")
	}

	// Ensure this flair exists and belongs to the creating user
	flair, err := r.DAOManager.GetFlairByUserIDFlairID(ctx, creatingUser.UserID, flairID)
	if err != nil || flair == nil {
		return nil, fmt.Errorf("Flair with ID %s not found", flairID)
	}

	// Add the flair to the participant
	return r.DAOManager.AddFlair(ctx, participant, flairID)
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
		// Need to add auth logic here
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

	if err != nil {
		return nil, err
	}

	_, err = r.DAOManager.CreateParticipantForDiscussion(ctx, discussionObj.ID, creatingUser.UserID)

	if err != nil {
		return nil, err
	}

	return discussionObj, nil
}

func (r *mutationResolver) RemoveFlair(ctx context.Context, discussionID string) (*model.Participant, error) {
	creatingUser := auth.GetAuthedUser(ctx)
	if creatingUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	participant, err := r.DAOManager.GetParticipantByDiscussionIDUserID(ctx, discussionID, creatingUser.UserID)
	if err != nil {
		return nil, err
	} else if participant == nil {
		return nil, fmt.Errorf("Current user not a participant in this discussion")
	}

	// Remove any flair by assigning nil
	return r.DAOManager.RemoveFlair(ctx, participant)
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
		// Need to add auth logic here
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

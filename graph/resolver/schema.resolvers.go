// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"
	"fmt"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/auth"
)

func (r *mutationResolver) CreateDiscussion(ctx context.Context, anonymityType model.AnonymityType) (*model.Discussion, error) {
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
	discussionObj, err := r.DAOManager.CreateNewDiscussion(ctx, creatingUser.User, anonymityType)

	if err != nil {
		return nil, err
	}

	_, err = r.DAOManager.AddModeratedDiscussionToUserProfile(ctx, creatingUser.User.UserProfileID, discussionObj.ID)

	if err != nil {
		return nil, err
	}

	return discussionObj, nil
}

func (r *mutationResolver) AddDiscussionParticipant(ctx context.Context, discussionID string, userID string) (*model.Participant, error) {
	participantObj, err := r.DAOManager.CreateParticipantForDiscussion(ctx, discussionID, userID)

	// TODO: Only the current user can join the conversation
	if err != nil {
		return nil, err
	}

	return participantObj, nil
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

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() generated.QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

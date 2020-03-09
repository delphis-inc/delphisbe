// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"
	"fmt"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *discussionResolver) Posts(ctx context.Context, obj *model.Discussion) ([]*model.Post, error) {
	posts, err := r.DAOManager.GetPostsByDiscussionID(ctx, obj.ID)

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *discussionResolver) Participants(ctx context.Context, obj *model.Discussion) ([]*model.Participant, error) {
	if obj.Participants == nil {
		participants, err := r.DAOManager.GetParticipantsByDiscussionID(ctx, obj.ID)

		if err != nil {
			return nil, err
		}

		particPointers := make([]*model.Participant, 0)
		for _, elem := range participants {
			elem.Discussion = obj
			particPointers = append(particPointers, &elem)
		}

		obj.Participants = particPointers
	}
	return obj.Participants, nil
}

func (r *Resolver) Discussion() generated.DiscussionResolver { return &discussionResolver{r} }

type discussionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *discussionResolver) Moderator(ctx context.Context, obj *model.Discussion) (*model.Moderator, error) {
	panic(fmt.Errorf("not implemented"))
}

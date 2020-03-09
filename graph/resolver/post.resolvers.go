// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *postResolver) IsDeleted(ctx context.Context, obj *model.Post) (bool, error) {
	return obj.DeletedAt != nil, nil
}

func (r *postResolver) Content(ctx context.Context, obj *model.Post) (string, error) {
	return obj.PostContent.Content, nil
}

func (r *Resolver) Post() generated.PostResolver { return &postResolver{r} }

type postResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *postResolver) Discussion(ctx context.Context, obj *model.Post) (*model.Discussion, error) {
	if obj.Discussion == nil {
		res, err := r.DAOManager.GetDiscussionByID(ctx, obj.DiscussionID)

		if err != nil {
			return nil, err
		}
		obj.Discussion = res
	}
	return obj.Discussion, nil
}
func (r *postResolver) Participant(ctx context.Context, obj *model.Post) (*model.Participant, error) {
	if obj.Participant == nil {
		participant, err := r.DAOManager.GetParticipantByID(ctx, model.DiscussionParticipantKey{
			DiscussionID:  obj.DiscussionID,
			ParticipantID: obj.ParticipantID,
		})

		if err != nil {
			return nil, err
		}

		obj.Participant = participant
	}
	return obj.Participant, nil
}

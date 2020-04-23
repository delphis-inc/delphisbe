// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"
	"time"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *discussionResolver) Moderator(ctx context.Context, obj *model.Discussion) (*model.Moderator, error) {
	if obj.Moderator == nil && obj.ModeratorID != nil {
		moderator, err := r.DAOManager.GetModeratorByID(ctx, *obj.ModeratorID)
		if err != nil {
			return nil, err
		}
		obj.Moderator = moderator
	}
	return obj.Moderator, nil
}

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
		for i, elem := range participants {
			elem.Discussion = obj
			particPointers = append(particPointers, &participants[i])
		}

		obj.Participants = particPointers
	}
	return obj.Participants, nil
}

func (r *discussionResolver) CreatedAt(ctx context.Context, obj *model.Discussion) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

func (r *discussionResolver) UpdatedAt(ctx context.Context, obj *model.Discussion) (string, error) {
	return obj.UpdatedAt.Format(time.RFC3339), nil
}

func (r *Resolver) Discussion() generated.DiscussionResolver { return &discussionResolver{r} }

type discussionResolver struct{ *Resolver }

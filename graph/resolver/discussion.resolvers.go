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
	panic(fmt.Errorf("not implemented"))
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

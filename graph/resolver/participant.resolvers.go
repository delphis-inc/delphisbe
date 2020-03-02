// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"
	"fmt"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *participantResolver) Discussion(ctx context.Context, obj *model.Participant) (*model.Discussion, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *participantResolver) Viewer(ctx context.Context, obj *model.Participant) (*model.Viewer, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *participantResolver) DiscussionNotificationPreferences(ctx context.Context, obj *model.Participant) (model.DiscussionNotificationPreferences, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *participantResolver) Posts(ctx context.Context, obj *model.Participant, first *int, after *string) (*model.PostsConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Participant() generated.ParticipantResolver { return &participantResolver{r} }

type participantResolver struct{ *Resolver }

package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/delphis-inc/delphisbe/graph/generated"
	"github.com/delphis-inc/delphisbe/graph/model"
)

func (r *participantsConnectionResolver) Edges(ctx context.Context, obj *model.ParticipantsConnection) ([]*model.ParticipantsEdge, error) {
	panic(fmt.Errorf("not implemented"))
}

// ParticipantsConnection returns generated.ParticipantsConnectionResolver implementation.
func (r *Resolver) ParticipantsConnection() generated.ParticipantsConnectionResolver {
	return &participantsConnectionResolver{r}
}

type participantsConnectionResolver struct{ *Resolver }

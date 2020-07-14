package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/delphis-inc/delphisbe/graph/generated"
	"github.com/delphis-inc/delphisbe/graph/model"
)

func (r *viewersConnectionResolver) Edges(ctx context.Context, obj *model.ViewersConnection) ([]*model.ViewersEdge, error) {
	panic(fmt.Errorf("not implemented"))
}

// ViewersConnection returns generated.ViewersConnectionResolver implementation.
func (r *Resolver) ViewersConnection() generated.ViewersConnectionResolver {
	return &viewersConnectionResolver{r}
}

type viewersConnectionResolver struct{ *Resolver }

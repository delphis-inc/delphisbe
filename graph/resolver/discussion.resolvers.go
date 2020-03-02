// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"
	"fmt"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *discussionResolver) Posts(ctx context.Context, obj *model.Discussion, first *int, after *string) (*model.PostsConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *discussionsConnectionResolver) Edges(ctx context.Context, obj *model.DiscussionsConnection) ([]*model.DiscussionsEdge, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Discussion() generated.DiscussionResolver { return &discussionResolver{r} }
func (r *Resolver) DiscussionsConnection() generated.DiscussionsConnectionResolver {
	return &discussionsConnectionResolver{r}
}

type discussionResolver struct {
	*Resolver
}
type discussionsConnectionResolver struct{ *Resolver }

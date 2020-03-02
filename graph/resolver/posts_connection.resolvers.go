// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"context"
	"fmt"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *postsConnectionResolver) Edges(ctx context.Context, obj *model.PostsConnection) ([]*model.PostsEdge, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PostsConnection() generated.PostsConnectionResolver {
	return &postsConnectionResolver{r}
}

type postsConnectionResolver struct{ *Resolver }

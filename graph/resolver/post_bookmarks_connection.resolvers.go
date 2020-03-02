// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"
	"fmt"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *postBookmarksConnectionResolver) Edges(ctx context.Context, obj *model.PostBookmarksConnection) ([]*model.PostBookmarksEdge, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) PostBookmarksConnection() generated.PostBookmarksConnectionResolver {
	return &postBookmarksConnectionResolver{r}
}

type postBookmarksConnectionResolver struct{ *Resolver }

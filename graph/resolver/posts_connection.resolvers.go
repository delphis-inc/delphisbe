package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/delphis-inc/delphisbe/graph/model"
)

func (r *postsConnectionResolver) Edges(ctx context.Context, obj *model.PostsConnection) ([]*model.PostsEdge, error) {
	return obj.Edges, nil
}

type postsConnectionResolver struct{ *Resolver }

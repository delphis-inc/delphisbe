package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *postBookmarkResolver) Discussion(ctx context.Context, obj *model.PostBookmark) (*model.Discussion, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *postBookmarkResolver) Post(ctx context.Context, obj *model.PostBookmark) (*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

// PostBookmark returns generated.PostBookmarkResolver implementation.
func (r *Resolver) PostBookmark() generated.PostBookmarkResolver { return &postBookmarkResolver{r} }

type postBookmarkResolver struct{ *Resolver }

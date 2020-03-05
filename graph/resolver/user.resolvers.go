// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"
	"fmt"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *userResolver) Participants(ctx context.Context, obj *model.User) ([]*model.Participant, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Viewers(ctx context.Context, obj *model.User) ([]*model.Viewer, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Bookmarks(ctx context.Context, obj *model.User) ([]*model.PostBookmark, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"
	"fmt"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *flairResolver) ImageURL(ctx context.Context, obj *model.Flair) (*model.URL, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) Flair() generated.FlairResolver { return &flairResolver{r} }

type flairResolver struct{ *Resolver }

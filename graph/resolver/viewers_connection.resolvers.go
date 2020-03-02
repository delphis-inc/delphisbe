// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"
	"fmt"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *viewersConnectionResolver) Edges(ctx context.Context, obj *model.ViewersConnection) ([]*model.ViewersEdge, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *Resolver) ViewersConnection() generated.ViewersConnectionResolver {
	return &viewersConnectionResolver{r}
}

type viewersConnectionResolver struct{ *Resolver }

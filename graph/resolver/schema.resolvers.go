// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *queryResolver) Discussion(ctx context.Context, id string) (*model.Discussion, error) {
	discussionObj, err := r.DAOManager.GetDiscussionByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return discussionObj, nil
}

func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
package resolver

import (
	"context"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/backend"
)

type Resolver struct {
	DAOManager backend.DAOManager
}

func (r *Resolver) resolveDiscussionByID(ctx context.Context, id string) (*model.Discussion, error) {
	discussionObj, err := r.DAOManager.GetDiscussionByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return discussionObj, nil
}

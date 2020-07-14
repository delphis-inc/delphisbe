// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
package resolver

import (
	"context"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/backend"
)

type Resolver struct {
	DAOManager backend.DelphisBackend
}

func (r *Resolver) resolveDiscussionByID(ctx context.Context, id string) (*model.Discussion, error) {
	discussionObj, err := r.DAOManager.GetDiscussionByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return discussionObj, nil
}

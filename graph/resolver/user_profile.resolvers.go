// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"
	"sort"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *userProfileResolver) ModeratedDiscussions(ctx context.Context, obj *model.UserProfile) ([]*model.Discussion, error) {
	moderatedDiscussionMap, err := r.DAOManager.GetDiscussionsByIDs(ctx, obj.ModeratedDiscussionsIDs)

	if err != nil {
		return nil, err
	}

	moderatedDiscussions := make([]*model.Discussion, 0)
	for _, v := range moderatedDiscussionMap {
		if v != nil {
			moderatedDiscussions = append(moderatedDiscussions, v)
		}
	}

	sort.Slice(moderatedDiscussions, func(i, j int) bool {
		return moderatedDiscussions[i].CreatedAt.Before(moderatedDiscussions[j].CreatedAt)
	})

	return moderatedDiscussions, nil
}

func (r *Resolver) UserProfile() generated.UserProfileResolver { return &userProfileResolver{r} }

type userProfileResolver struct{ *Resolver }

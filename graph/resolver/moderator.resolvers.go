// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *moderatorResolver) Discussion(ctx context.Context, obj *model.Moderator) (*model.Discussion, error) {
	if obj.Discussion == nil {
		discussion, err := r.DAOManager.GetDiscussionByModeratorID(ctx, obj.ID)

		if err != nil {
			return nil, err
		}

		obj.Discussion = discussion
	}
	return obj.Discussion, nil
}

func (r *moderatorResolver) UserProfile(ctx context.Context, obj *model.Moderator) (*model.UserProfile, error) {
	if obj.UserProfile == nil && obj.UserProfileID != nil {
		userProfile, err := r.DAOManager.GetUserProfileByID(ctx, *obj.UserProfileID)

		if err != nil {
			return nil, err
		}

		obj.UserProfile = userProfile
	}
	return obj.UserProfile, nil
}

func (r *Resolver) Moderator() generated.ModeratorResolver { return &moderatorResolver{r} }

type moderatorResolver struct{ *Resolver }

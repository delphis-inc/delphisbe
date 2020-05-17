package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
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
	logrus.Debugf("Have moderator: %+v", obj)
	if obj.UserProfile == nil && obj.UserProfileID != nil {
		userProfile, err := r.DAOManager.GetUserProfileByID(ctx, *obj.UserProfileID)
		logrus.Debugf("User Profile: %+v; err: %+v", userProfile, err)

		if err != nil {
			return nil, err
		}

		obj.UserProfile = userProfile
	}
	return obj.UserProfile, nil
}

// Moderator returns generated.ModeratorResolver implementation.
func (r *Resolver) Moderator() generated.ModeratorResolver { return &moderatorResolver{r} }

type moderatorResolver struct{ *Resolver }

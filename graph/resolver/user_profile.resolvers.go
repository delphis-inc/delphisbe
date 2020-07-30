package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/delphis-inc/delphisbe/graph/generated"
	"github.com/delphis-inc/delphisbe/graph/model"
)

func (r *userProfileResolver) ProfileImageURL(ctx context.Context, obj *model.UserProfile) (string, error) {
	if len(obj.SocialInfos) == 0 {
		si, err := r.DAOManager.GetSocialInfosByUserProfileID(ctx, obj.ID)
		if err != nil {
			return "", err
		}
		obj.SocialInfos = si
	}
	return obj.SocialInfos[0].ProfileImageURL, nil
}

func (r *userProfileResolver) AuthenticatedWithTwitter(ctx context.Context, obj *model.UserProfile) (bool, error) {
	if len(obj.SocialInfos) == 0 {
		si, err := r.DAOManager.GetSocialInfosByUserProfileID(ctx, obj.ID)
		if err != nil {
			return false, err
		}
		for _, s := range si {
			if s.Network == "twitter" {
				return true, nil
			}
		}
	}
	return false, nil
}

// UserProfile returns generated.UserProfileResolver implementation.
func (r *Resolver) UserProfile() generated.UserProfileResolver { return &userProfileResolver{r} }

type userProfileResolver struct{ *Resolver }

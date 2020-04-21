// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
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

func (r *Resolver) UserProfile() generated.UserProfileResolver { return &userProfileResolver{r} }

type userProfileResolver struct{ *Resolver }

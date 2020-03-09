package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/util"
)

type LoginWithTwitterInput struct {
	User              *twitter.User
	AccessToken       string
	AccessTokenSecret string
}

func (t LoginWithTwitterInput) ID() string {
	return fmt.Sprintf("%s:%s", "twitter", t.User.IDStr)
}

func (b *delphisBackend) GetOrCreateUser(ctx context.Context, input LoginWithTwitterInput) (*model.User, error) {
	userProfileObj := &model.UserProfile{
		ID:            input.ID(),
		DisplayName:   input.User.Name,
		TwitterHandle: input.User.ScreenName,
		TwitterInfo: model.SocialInfo{
			AccessToken:       input.AccessToken,
			AccessTokenSecret: input.AccessTokenSecret,
			UserID:            input.User.IDStr,
			ProfileImageURL:   input.User.ProfileImageURLHttps,
			ScreenName:        input.User.ScreenName,
			IsVerified:        input.User.Verified,
		},
	}

	userProfileObj, isCreated, err := b.db.CreateOrUpdateUserProfile(ctx, *userProfileObj)
	if err != nil {
		return nil, err
	}

	var userObj *model.User
	if isCreated || userProfileObj.UserID == "" {
		userObj = &model.User{
			ID:            util.UUIDv4(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			UserProfileID: userProfileObj.ID,
		}
		userObj, err = b.db.PutUser(ctx, *userObj)

		if err != nil {
			logrus.WithError(err).Errorf("Failed putting user object: %+v", userObj)
			return nil, err
		}

		_, err = b.db.UpdateUserProfileUserID(ctx, userProfileObj.ID, userObj.ID)
		if err != nil {
			return nil, err
		}
	} else {
		userObj, err = b.GetUserByID(ctx, userProfileObj.UserID)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to find user")
			return nil, err
		}
	}

	return userObj, nil
}

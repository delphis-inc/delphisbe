package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/dghubble/go-twitter/twitter"
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
	}

	userProfileObj, isCreated, err := b.db.CreateOrUpdateUserProfile(ctx, *userProfileObj)
	if err != nil {
		return nil, err
	}

	socialInfoObj := &model.SocialInfo{
		Network:           "twitter",
		AccessToken:       input.AccessToken,
		AccessTokenSecret: input.AccessTokenSecret,
		UserID:            input.User.IDStr,
		ProfileImageURL:   input.User.ProfileImageURLHttps,
		ScreenName:        input.User.ScreenName,
		IsVerified:        input.User.Verified,
		UserProfileID:     userProfileObj.ID,
	}
	_, err = b.db.UpsertSocialInfo(ctx, *socialInfoObj)
	if err != nil {
		return nil, err
	}

	var userObj *model.User
	logrus.Debugf("isCreated? %t, userProfileObj: %+v", isCreated, userProfileObj)
	if isCreated || userProfileObj.UserID == nil {
		userObj = &model.User{
			ID:        util.UUIDv4(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		userObj, err = b.db.UpsertUser(ctx, *userObj)

		if err != nil {
			logrus.WithError(err).Errorf("Failed putting user object: %+v", userObj)
			return nil, err
		}

		userProfileObj.UserID = &userObj.ID
		_, _, err = b.db.CreateOrUpdateUserProfile(ctx, *userProfileObj)
		if err != nil {
			return nil, err
		}
	} else {
		userObj, err = b.GetUserByID(ctx, *userProfileObj.UserID)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to find user")
			return nil, err
		}
	}

	return userObj, nil
}

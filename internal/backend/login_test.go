package backend

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/dghubble/go-twitter/twitter"

	"github.com/delphis-inc/delphisbe/internal/backend/test_utils"

	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/cache"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/delphis-inc/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestDelphisBackend_GetOrCreateUser(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	profileObj := test_utils.TestUserProfile()
	socialObj := test_utils.TestSocialInfo()
	userObj := test_utils.TestUser()
	twitterInput := LoginWithTwitterInput{
		User: &twitter.User{
			Email: "test@email.com",
		},
		AccessToken:       test_utils.Token,
		AccessTokenSecret: test_utils.TokenSecret,
	}

	Convey("RespondToRequestAccess", t, func() {
		cacheObj := cache.NewInMemoryCache()
		authObj := auth.NewDelphisAuth(nil)
		mockDB := &mocks.Datastore{}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            authObj,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when CreateOrUpdateUserProfile errors out ", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("CreateOrUpdateUserProfile", ctx, mock.Anything).Return(nil, false, expectedError)

			resp, err := backendObj.GetOrCreateUser(ctx, twitterInput)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when UpsertSocialInfo errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("CreateOrUpdateUserProfile", ctx, mock.Anything).Return(&profileObj, false, nil)
			mockDB.On("UpsertSocialInfo", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetOrCreateUser(ctx, twitterInput)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when user needs to be created", func() {
			Convey("when UpsertUser errors out", func() {
				expectedError := fmt.Errorf("Some Error")
				mockDB.On("CreateOrUpdateUserProfile", ctx, mock.Anything).Return(&profileObj, true, nil)
				mockDB.On("UpsertSocialInfo", ctx, mock.Anything).Return(&socialObj, nil)
				mockDB.On("UpsertUser", ctx, mock.Anything).Return(nil, expectedError)

				resp, err := backendObj.GetOrCreateUser(ctx, twitterInput)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})

			Convey("when call is successful", func() {
				mockDB.On("CreateOrUpdateUserProfile", ctx, mock.Anything).Return(&profileObj, true, nil)
				mockDB.On("UpsertSocialInfo", ctx, mock.Anything).Return(&socialObj, nil)
				mockDB.On("UpsertUser", ctx, mock.Anything).Return(&userObj, nil)
				mockDB.On("CreateOrUpdateUserProfile", ctx, mock.Anything).Return(&profileObj, true, nil)

				resp, err := backendObj.GetOrCreateUser(ctx, twitterInput)

				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
			})
		})

		Convey("when user does not need to be created", func() {
			Convey("when GetUserByID errors out", func() {
				expectedError := fmt.Errorf("Some Error")
				mockDB.On("CreateOrUpdateUserProfile", ctx, mock.Anything).Return(&profileObj, false, nil)
				mockDB.On("UpsertSocialInfo", ctx, mock.Anything).Return(&socialObj, nil)
				mockDB.On("GetUserByID", ctx, mock.Anything).Return(nil, expectedError)

				resp, err := backendObj.GetOrCreateUser(ctx, twitterInput)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})

			Convey("when call is successful", func() {
				mockDB.On("CreateOrUpdateUserProfile", ctx, mock.Anything).Return(&profileObj, false, nil)
				mockDB.On("UpsertSocialInfo", ctx, mock.Anything).Return(&socialObj, nil)
				mockDB.On("GetUserByID", ctx, mock.Anything).Return(&userObj, nil)

				resp, err := backendObj.GetOrCreateUser(ctx, twitterInput)

				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
			})
		})
	})
}

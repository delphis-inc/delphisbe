package backend

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"

	"github.com/nedrocks/delphisbe/internal/auth"
	"github.com/nedrocks/delphisbe/internal/backend/test_utils"
	"github.com/nedrocks/delphisbe/internal/cache"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/nedrocks/delphisbe/internal/util"
	"github.com/nedrocks/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDelphisBackend_GetUserProfileByUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := test_utils.UserID

	profileObj := test_utils.TestUserProfile()

	Convey("GetUserProfileByUserID", t, func() {
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

		Convey("when the query errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetUserProfileByUserID", ctx, userID).Return(nil, expectedError)

			resp, err := backendObj.GetUserProfileByUserID(ctx, userID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetUserProfileByUserID", ctx, userID).Return(&profileObj, nil)

			resp, err := backendObj.GetUserProfileByUserID(ctx, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &profileObj)
		})
	})
}

func TestDelphisBackend_GetUserProfileByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	profileID := test_utils.ProfileID

	profileObj := test_utils.TestUserProfile()

	Convey("GetUserProfileByID", t, func() {
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

		Convey("when the query errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetUserProfileByID", ctx, profileID).Return(nil, expectedError)

			resp, err := backendObj.GetUserProfileByID(ctx, profileID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetUserProfileByID", ctx, profileID).Return(&profileObj, nil)

			resp, err := backendObj.GetUserProfileByID(ctx, profileID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &profileObj)
		})
	})
}

func TestDelphisBackend_GetSocialInfosByUserProfileID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	profileID := test_utils.ProfileID

	socialObj := test_utils.TestSocialInfo()

	Convey("GetUserProfileByID", t, func() {
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

		Convey("when the query errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetSocialInfosByUserProfileID", ctx, profileID).Return(nil, expectedError)

			resp, err := backendObj.GetSocialInfosByUserProfileID(ctx, profileID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetSocialInfosByUserProfileID", ctx, profileID).Return([]model.SocialInfo{socialObj}, nil)

			resp, err := backendObj.GetSocialInfosByUserProfileID(ctx, profileID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []model.SocialInfo{socialObj})
		})
	})
}

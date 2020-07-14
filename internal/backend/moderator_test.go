package backend

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/backend/test_utils"
	"github.com/delphis-inc/delphisbe/internal/cache"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/delphis-inc/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDelphisBackend_GetModeratorByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	modID := test_utils.ModeratorID

	modObj := test_utils.TestModerator()

	Convey("GetModeratorByID", t, func() {
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

		Convey("when GetModeratorByID errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetModeratorByID", ctx, modID).Return(nil, expectedError)

			resp, err := backendObj.GetModeratorByID(ctx, modID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when query succeeds", func() {
			mockDB.On("GetModeratorByID", ctx, modID).Return(&modObj, nil)

			resp, err := backendObj.GetModeratorByID(ctx, modID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &modObj)
		})
	})
}

func TestDelphisBackend_GetModeratorByUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := test_utils.UserID

	modObj := test_utils.TestModerator()

	Convey("GetModeratorByUserID", t, func() {
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

		Convey("when GetModeratorByUserID errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetModeratorByUserID", ctx, userID).Return(nil, expectedError)

			resp, err := backendObj.GetModeratorByUserID(ctx, userID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when query succeeds", func() {
			mockDB.On("GetModeratorByUserID", ctx, userID).Return(&modObj, nil)

			resp, err := backendObj.GetModeratorByUserID(ctx, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &modObj)
		})
	})
}

func TestDelphisBackend_GetModeratorByUserIDAndDiscussionID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := test_utils.UserID
	discussionID := test_utils.DiscussionID

	modObj := test_utils.TestModerator()

	Convey("GetModeratorByUserID", t, func() {
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

		Convey("when GetModeratorByUserIDAndDiscussionID errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetModeratorByUserIDAndDiscussionID", ctx, userID, discussionID).Return(nil, expectedError)

			resp, err := backendObj.GetModeratorByUserIDAndDiscussionID(ctx, userID, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when query succeeds", func() {
			mockDB.On("GetModeratorByUserIDAndDiscussionID", ctx, userID, discussionID).Return(&modObj, nil)

			resp, err := backendObj.GetModeratorByUserIDAndDiscussionID(ctx, userID, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &modObj)
		})
	})
}

func TestDelphisBackend_CheckIfModerator(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := test_utils.UserID

	modObj := test_utils.TestModerator()

	Convey("CheckIfModerator", t, func() {
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

		Convey("when GetModeratorByUserID errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetModeratorByUserID", ctx, userID).Return(nil, expectedError)

			resp, err := backendObj.CheckIfModerator(ctx, userID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldEqual, false)
		})

		Convey("when GetModeratorByUserID does not return a row; the user is not a mod", func() {
			mockDB.On("GetModeratorByUserID", ctx, userID).Return(nil, nil)

			resp, err := backendObj.CheckIfModerator(ctx, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldEqual, false)
		})

		Convey("when query succeeds", func() {
			mockDB.On("GetModeratorByUserID", ctx, userID).Return(&modObj, nil)

			resp, err := backendObj.CheckIfModerator(ctx, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, true)
		})
	})
}

func TestDelphisBackend_CheckIfModeratorForDiscussion(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := test_utils.UserID
	discussionID := test_utils.DiscussionID

	modObj := test_utils.TestModerator()

	Convey("CheckIfModerator", t, func() {
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

		Convey("when GetModeratorByUserIDAndDiscussionID errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetModeratorByUserIDAndDiscussionID", ctx, userID, discussionID).Return(nil, expectedError)

			resp, err := backendObj.CheckIfModeratorForDiscussion(ctx, userID, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldEqual, false)
		})

		Convey("when GetModeratorByUserIDAndDiscussionID does not return a row; the user is not a mod", func() {
			mockDB.On("GetModeratorByUserIDAndDiscussionID", ctx, userID, discussionID).Return(nil, nil)

			resp, err := backendObj.CheckIfModeratorForDiscussion(ctx, userID, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldEqual, false)
		})

		Convey("when query succeeds", func() {
			mockDB.On("GetModeratorByUserIDAndDiscussionID", ctx, userID, discussionID).Return(&modObj, nil)

			resp, err := backendObj.CheckIfModeratorForDiscussion(ctx, userID, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, true)
		})
	})
}

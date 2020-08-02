package backend

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/backend/test_utils"
	"github.com/delphis-inc/delphisbe/internal/cache"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/delphis-inc/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDelphisBackend_GetAccessLinkBySlug(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	slug := test_utils.LinkSlug

	dalObj := test_utils.TestDiscussionAccessLink()

	Convey("GetAccessLinkBySlug", t, func() {
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
			mockDB.On("GetAccessLinkBySlug", ctx, slug).Return(nil, expectedError)

			resp, err := backendObj.GetAccessLinkBySlug(ctx, slug)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetAccessLinkBySlug", ctx, slug).Return(&dalObj, nil)

			resp, err := backendObj.GetAccessLinkBySlug(ctx, slug)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &dalObj)
		})
	})
}

func TestDelphisBackend_GetAccessLinkByDiscussionID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	discussionID := test_utils.DiscussionID

	dalObj := test_utils.TestDiscussionAccessLink()

	Convey("GetAccessLinkByDiscussionID", t, func() {
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
			mockDB.On("GetAccessLinkByDiscussionID", ctx, discussionID).Return(nil, expectedError)

			resp, err := backendObj.GetAccessLinkByDiscussionID(ctx, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetAccessLinkByDiscussionID", ctx, discussionID).Return(&dalObj, nil)

			resp, err := backendObj.GetAccessLinkByDiscussionID(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &dalObj)
		})
	})
}

func TestDelphisBackend_PutAccessLinkForDiscussion(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	discussionID := test_utils.DiscussionID

	dalObj := test_utils.TestDiscussionAccessLink()

	tx := sql.Tx{}

	Convey("PutAccessLinkForDiscussion", t, func() {
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

		Convey("when BeginTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(nil, expectedError)

			resp, err := backendObj.PutAccessLinkForDiscussion(ctx, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when PutAccessLinkForDiscussion errors out and Rollback fails", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutAccessLinkForDiscussion", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.PutAccessLinkForDiscussion(ctx, discussionID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when PutAccessLinkForDiscussion errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutAccessLinkForDiscussion", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.PutAccessLinkForDiscussion(ctx, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when CommitTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutAccessLinkForDiscussion", ctx, mock.Anything, mock.Anything).Return(&dalObj, nil)

			mockDB.On("CommitTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.PutAccessLinkForDiscussion(ctx, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when response succeeds", func() {
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutAccessLinkForDiscussion", ctx, mock.Anything, mock.Anything).Return(&dalObj, nil)

			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.PutAccessLinkForDiscussion(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

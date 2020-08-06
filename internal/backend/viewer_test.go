package backend

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/backend/test_utils"
	"github.com/delphis-inc/delphisbe/internal/cache"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/delphis-inc/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_GetViewersByIDs(t *testing.T) {
	ctx := context.Background()

	viewerIDs := []string{"view1", "view2"}

	Convey("GetViewersByIDs", t, func() {
		now := time.Now()
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

		Convey("when db returns an error it is returned unchanged", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetViewersByIDs", ctx, viewerIDs).Return(nil, expectedError)

			resp, err := backendObj.GetViewersByIDs(ctx, viewerIDs)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when db returns response it is returned", func() {
			dbResponse := map[string]*model.Viewer{
				"view1": &model.Viewer{ID: "view1"},
				"view2": &model.Viewer{ID: "view2"},
			}

			mockDB.On("GetViewersByIDs", ctx, viewerIDs).Return(dbResponse, nil)

			resp, err := backendObj.GetViewersByIDs(ctx, viewerIDs)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, dbResponse)
		})
	})
}

func TestDelphisBackend_SetViewerLastPostViewed(t *testing.T) {
	ctx := context.Background()
	viewerObj := test_utils.TestViewer()
	postID := test_utils.PostID

	Convey("SetViewerLastPostViewed", t, func() {
		now := time.Now()
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
		Convey("when db returns an error", func() {
			mockDB.On("SetViewerLastPostViewed", ctx, viewerObj.ID, postID, now).Return(nil, fmt.Errorf("sth"))

			resp, err := backendObj.SetViewerLastPostViewed(ctx, viewerObj.ID, postID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when db returns nil", func() {
			mockDB.On("SetViewerLastPostViewed", ctx, viewerObj.ID, postID, now).Return(nil, nil)

			resp, err := backendObj.SetViewerLastPostViewed(ctx, viewerObj.ID, postID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when db returns viewer object", func() {
			mockDB.On("SetViewerLastPostViewed", ctx, viewerObj.ID, postID, now).Return(&viewerObj, nil)

			resp, err := backendObj.SetViewerLastPostViewed(ctx, viewerObj.ID, postID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &viewerObj)
		})
	})
}

func TestDelphisBackend_GetViewerForDiscussion(t *testing.T) {
	ctx := context.Background()
	discussionID := test_utils.DiscussionID
	viewerObj := test_utils.TestViewer()
	userID := test_utils.UserID

	Convey("GetViewerForDiscussion", t, func() {
		now := time.Now()
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

		Convey("when db returns error", func() {
			mockDB.On("GetViewerForDiscussion", ctx, discussionID, userID).Return(nil, fmt.Errorf("sth"))

			resp, err := backendObj.GetViewerForDiscussion(ctx, discussionID, userID, false)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when db returns nil and create is false", func() {
			mockDB.On("GetViewerForDiscussion", ctx, discussionID, userID).Return(nil, nil)

			resp, err := backendObj.GetViewerForDiscussion(ctx, discussionID, userID, false)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when db returns nil and create is true", func() {
			mockDB.On("GetViewerForDiscussion", ctx, discussionID, userID).Return(nil, nil)
			mockDB.On("UpsertViewer", ctx, mock.Anything).Return(&viewerObj, nil)

			resp, err := backendObj.GetViewerForDiscussion(ctx, discussionID, userID, true)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &viewerObj)
		})

		Convey("when db returns a viewer object", func() {
			mockDB.On("GetViewerForDiscussion", ctx, discussionID, userID).Return(&viewerObj, nil)

			resp, err := backendObj.GetViewerForDiscussion(ctx, discussionID, userID, true)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &viewerObj)
		})
	})
}

func Test_GetViewerByID(t *testing.T) {
	ctx := context.Background()

	viewerID := "view1"

	Convey("GetViewerByID", t, func() {
		now := time.Now()
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

		Convey("when db returns an error it is returned unchanged", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetViewersByIDs", ctx, []string{viewerID}).Return(nil, expectedError)

			resp, err := backendObj.GetViewerByID(ctx, viewerID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when db returns response it is returned", func() {
			dbResponse := map[string]*model.Viewer{
				"view1": &model.Viewer{ID: "view1"},
			}

			mockDB.On("GetViewersByIDs", ctx, []string{viewerID}).Return(dbResponse, nil)

			resp, err := backendObj.GetViewerByID(ctx, viewerID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, dbResponse[viewerID])
		})
	})
}

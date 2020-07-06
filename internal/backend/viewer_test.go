package backend

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/auth"
	"github.com/nedrocks/delphisbe/internal/cache"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/nedrocks/delphisbe/internal/util"
	"github.com/nedrocks/delphisbe/mocks"
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

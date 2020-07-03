package backend

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nedrocks/delphisbe/internal/backend/test_utils"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/auth"
	"github.com/nedrocks/delphisbe/internal/cache"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/nedrocks/delphisbe/internal/util"
	"github.com/nedrocks/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestDelphisBackend_CreateFlair(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := test_utils.UserID
	templateID := test_utils.FlairTemplateID

	flairObj := test_utils.TestFlair()

	Convey("CreateFlair", t, func() {
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

		Convey("when UpsertFlair errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("UpsertFlair", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.CreateFlair(ctx, userID, templateID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when call succeeds", func() {
			mockDB.On("UpsertFlair", ctx, mock.Anything).Return(&flairObj, nil)

			resp, err := backendObj.CreateFlair(ctx, userID, templateID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_GetFlairByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	flairID := test_utils.FlairID

	flairObj := test_utils.TestFlair()

	Convey("GetFlairByID", t, func() {
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

		Convey("when the query errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetFlairByID", ctx, flairID).Return(nil, expectedError)

			resp, err := backendObj.GetFlairByID(ctx, flairID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetFlairByID", ctx, flairID).Return(&flairObj, nil)

			resp, err := backendObj.GetFlairByID(ctx, flairID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &flairObj)
		})
	})
}

func TestDelphisBackend_GetFlairsByUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := test_utils.UserID

	flairObj := test_utils.TestFlair()

	Convey("GetFlairsByUserID", t, func() {
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

		Convey("when the query errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetFlairsByUserID", ctx, userID).Return(nil, expectedError)

			resp, err := backendObj.GetFlairsByUserID(ctx, userID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetFlairsByUserID", ctx, userID).Return([]*model.Flair{&flairObj}, nil)

			resp, err := backendObj.GetFlairsByUserID(ctx, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.Flair{&flairObj})
		})
	})
}

func TestDelphisBackend_RemoveFlair(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	flairObj := test_utils.TestFlair()

	Convey("RemoveFlair", t, func() {
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

		Convey("when the query errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("RemoveFlair", ctx, flairObj).Return(nil, expectedError)

			resp, err := backendObj.RemoveFlair(ctx, flairObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("RemoveFlair", ctx, flairObj).Return(&flairObj, nil)

			resp, err := backendObj.RemoveFlair(ctx, flairObj)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &flairObj)
		})
	})
}

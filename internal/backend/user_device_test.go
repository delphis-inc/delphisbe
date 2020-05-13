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

func Test_UpsertUserDevice(t *testing.T) {
	ctx := context.Background()
	Convey("UpsertUserDevice", t, func() {
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

		Convey("calls db.UpsertUserDevice and returns response", func() {
			var expectedErr error = nil
			expectedResponse := &model.UserDevice{ID: util.UUIDv4()}
			deviceID := "12345"
			userID := "54321"
			platform := "ios"
			token := "11aa1"
			mockDB.On("UpsertUserDevice", ctx, model.UserDevice{
				ID:       deviceID,
				Platform: platform,
				LastSeen: now,
				Token:    &token,
				UserID:   &userID,
			}).Return(expectedResponse, expectedErr)

			resp, err := backendObj.UpsertUserDevice(ctx, deviceID, &userID, platform, &token)

			So(err, ShouldEqual, expectedErr)
			So(resp, ShouldEqual, expectedResponse)
			So(len(mockDB.Calls), ShouldEqual, 1)
		})

		Convey("when an error is returned it passes back the exact error", func() {
			var expectedErr error = fmt.Errorf("some weird error")
			deviceID := "12345"
			userID := "54321"
			platform := "ios"
			token := "11aa1"
			mockDB.On("UpsertUserDevice", ctx, model.UserDevice{
				ID:       deviceID,
				Platform: platform,
				LastSeen: now,
				Token:    &token,
				UserID:   &userID,
			}).Return(nil, expectedErr)

			resp, err := backendObj.UpsertUserDevice(ctx, deviceID, &userID, platform, &token)

			So(err, ShouldEqual, expectedErr)
			So(resp, ShouldBeNil)
			So(len(mockDB.Calls), ShouldEqual, 1)
		})
	})
}

func Test_GetUserDeviceByUserIDPlatform(t *testing.T) {
	ctx := context.Background()
	Convey("GetUserDeviceByUserIDPlatform", t, func() {
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
			userID := "54321"
			platform := "ios"
			mockDB.On("GetUserDevicesByUserID", ctx, userID).Return(nil, expectedError)

			resp, err := backendObj.GetUserDeviceByUserIDPlatform(ctx, userID, platform)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when db has success but returns 0 items it returns nil", func() {
			dbResponse := []model.UserDevice{}
			userID := "54321"
			platform := "ios"
			mockDB.On("GetUserDevicesByUserID", ctx, userID).Return(dbResponse, nil)

			resp, err := backendObj.GetUserDeviceByUserIDPlatform(ctx, userID, platform)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when db is successful", func() {
			userID := "54321"
			dbResponse := []model.UserDevice{
				{
					ID:       "11111",
					Platform: "ios",
				},
				{
					ID:       "2222",
					Platform: "android",
				},
			}
			mockDB.On("GetUserDevicesByUserID", ctx, userID).Return(dbResponse, nil)
			Convey("when platform is found in response it returns the object", func() {
				resp, err := backendObj.GetUserDeviceByUserIDPlatform(ctx, userID, "ios")

				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(resp.Platform, ShouldEqual, "ios")
			})

			Convey("when platform is not found in response it returns nil", func() {
				resp, err := backendObj.GetUserDeviceByUserIDPlatform(ctx, userID, "web")

				So(err, ShouldBeNil)
				So(resp, ShouldBeNil)
			})
		})
	})
}

func Test_GetUserDevicesByUserID(t *testing.T) {
	ctx := context.Background()
	Convey("GetUserDevicesByUserID", t, func() {
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
			userID := "54321"
			mockDB.On("GetUserDevicesByUserID", ctx, userID).Return(nil, expectedError)

			resp, err := backendObj.GetUserDevicesByUserID(ctx, userID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when db returns response it is returned", func() {
			dbResponse := []model.UserDevice{
				{
					ID: "1111",
				},
				{
					ID: "2222",
				},
			}
			userID := "54321"
			mockDB.On("GetUserDevicesByUserID", ctx, userID).Return(dbResponse, nil)

			resp, err := backendObj.GetUserDevicesByUserID(ctx, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, dbResponse)
		})
	})
}

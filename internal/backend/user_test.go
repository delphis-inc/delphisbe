package backend

import (
	"context"
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

func TestDelphisBackend_GetUserByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := test_utils.UserID

	userObj := test_utils.TestUser()

	Convey("GetUserByID", t, func() {
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
			mockDB.On("GetUserByID", ctx, userID).Return(nil, expectedError)

			resp, err := backendObj.GetUserByID(ctx, userID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetUserByID", ctx, userID).Return(&userObj, nil)

			resp, err := backendObj.GetUserByID(ctx, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &userObj)
		})
	})
}

func TestDelphisBackend_CreateUser(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userObj := test_utils.TestUser()

	Convey("CreateUser", t, func() {
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
			mockDB.On("UpsertUser", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.CreateUser(ctx)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("UpsertUser", ctx, mock.Anything).Return(&userObj, nil)

			resp, err := backendObj.CreateUser(ctx)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &userObj)
		})
	})
}

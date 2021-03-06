package backend

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/delphis-inc/delphisbe/internal/backend/test_utils"

	"github.com/delphis-inc/delphisbe/internal/auth"

	"github.com/delphis-inc/delphisbe/internal/cache"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/delphis-inc/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewAccessToken(t *testing.T) {
	ctx := context.Background()

	userID := test_utils.UserID

	Convey("NewAccessToken", t, func() {
		now := time.Now()
		cacheObj := cache.NewInMemoryCache()
		mockAuth := &mocks.DelphisAuth{}
		mockDB := &mocks.Datastore{}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            mockAuth,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when auth returns an error", func() {
			expectedError := fmt.Errorf("Some Error")
			mockAuth.On("NewAccessToken", userID).Return(nil, expectedError)

			resp, err := backendObj.NewAccessToken(ctx, userID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when db returns response it is returned", func() {
			authResponse := &auth.DelphisAccessToken{
				Claims:      nil,
				TokenString: "token",
			}

			mockAuth.On("NewAccessToken", userID).Return(authResponse, nil)

			resp, err := backendObj.NewAccessToken(ctx, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, authResponse)
		})
	})
}

func Test_ValidateAccessToken(t *testing.T) {
	ctx := context.Background()

	token := test_utils.Token

	Convey("ValidateAccessToken", t, func() {
		now := time.Now()
		cacheObj := cache.NewInMemoryCache()
		mockAuth := &mocks.DelphisAuth{}
		mockDB := &mocks.Datastore{}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            mockAuth,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when auth returns an error", func() {
			expectedError := fmt.Errorf("Some Error")
			mockAuth.On("ValidateAccessToken", ctx, token).Return(nil, expectedError)

			resp, err := backendObj.ValidateAccessToken(ctx, token)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when db returns response it is returned", func() {
			authResponse := &auth.DelphisAuthedUser{
				UserID: "userID",
				User:   nil,
			}

			mockAuth.On("ValidateAccessToken", ctx, token).Return(authResponse, nil)

			resp, err := backendObj.ValidateAccessToken(ctx, token)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, authResponse)
		})
	})
}

func Test_ValidateRefreshToken(t *testing.T) {
	ctx := context.Background()

	token := test_utils.Token

	Convey("ValidateRefreshToken", t, func() {
		now := time.Now()
		cacheObj := cache.NewInMemoryCache()
		mockAuth := &mocks.DelphisAuth{}
		mockDB := &mocks.Datastore{}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            mockAuth,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when auth returns an error", func() {
			expectedError := fmt.Errorf("Some Error")
			mockAuth.On("ValidateRefreshToken", ctx, token).Return(nil, expectedError)

			resp, err := backendObj.ValidateRefreshToken(ctx, token)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when db returns response it is returned", func() {
			authResponse := &auth.DelphisRefreshTokenUser{
				UserID: "userID",
				User:   nil,
			}

			mockAuth.On("ValidateRefreshToken", ctx, token).Return(authResponse, nil)

			resp, err := backendObj.ValidateRefreshToken(ctx, token)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, authResponse)
		})
	})
}

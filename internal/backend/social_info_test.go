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

func TestDelphisBackend_UpsertSocialInfo(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	socialObj := test_utils.TestSocialInfo()

	Convey("UpsertSocialInfo", t, func() {
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
			mockDB.On("UpsertSocialInfo", ctx, socialObj).Return(nil, expectedError)

			resp, err := backendObj.UpsertSocialInfo(ctx, socialObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("UpsertSocialInfo", ctx, socialObj).Return(&socialObj, nil)

			resp, err := backendObj.UpsertSocialInfo(ctx, socialObj)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &socialObj)
		})
	})
}

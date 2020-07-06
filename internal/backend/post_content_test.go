package backend

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nedrocks/delphisbe/internal/auth"
	"github.com/nedrocks/delphisbe/internal/backend/test_utils"
	"github.com/nedrocks/delphisbe/internal/cache"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/nedrocks/delphisbe/internal/util"
	"github.com/nedrocks/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDelphisBackend_GetPostContentByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	postContentID := test_utils.PostContentID

	postContentObj := test_utils.TestPostContent()

	Convey("GetPostContentByID", t, func() {
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
			mockDB.On("GetPostContentByID", ctx, postContentID).Return(nil, expectedError)

			resp, err := backendObj.GetPostContentByID(ctx, postContentID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetPostContentByID", ctx, postContentID).Return(&postContentObj, nil)

			resp, err := backendObj.GetPostContentByID(ctx, postContentID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &postContentObj)
		})
	})
}

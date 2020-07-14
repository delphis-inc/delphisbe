package backend

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/delphis-inc/delphisbe/internal/backend/test_utils"

	"github.com/stretchr/testify/mock"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/cache"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/delphis-inc/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDelphisBackend_ListFlairTemplates(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	ftObj := test_utils.TestFlairTemplate()
	query := "query"

	Convey("ListFlairTemplates", t, func() {
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
			mockDB.On("ListFlairTemplates", ctx, &query).Return(nil, expectedError)

			resp, err := backendObj.ListFlairTemplates(ctx, &query)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("ListFlairTemplates", ctx, &query).Return([]*model.FlairTemplate{&ftObj}, nil)

			resp, err := backendObj.ListFlairTemplates(ctx, &query)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.FlairTemplate{&ftObj})
		})
	})
}

func TestDelphisBackend_CreateFlairTemplate(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	source := test_utils.Source
	displayName := test_utils.DisplayName
	imageURL := test_utils.ImageURL

	ftObj := test_utils.TestFlairTemplate()

	Convey("CreateFlairTemplate", t, func() {
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

		Convey("when both display name and imageURL are nil", func() {

			resp, err := backendObj.CreateFlairTemplate(ctx, nil, nil, source)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when the UpsertFlairTemplate errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("UpsertFlairTemplate", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.CreateFlairTemplate(ctx, &displayName, &imageURL, source)

			So(err, ShouldResemble, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("UpsertFlairTemplate", ctx, mock.Anything).Return(&ftObj, nil)

			resp, err := backendObj.CreateFlairTemplate(ctx, &displayName, &imageURL, source)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_GetFlairTemplateByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	templateID := test_utils.FlairTemplateID

	ftObj := test_utils.TestFlairTemplate()

	Convey("GetFlairTemplateByID", t, func() {
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
			mockDB.On("GetFlairTemplateByID", ctx, templateID).Return(nil, expectedError)

			resp, err := backendObj.GetFlairTemplateByID(ctx, templateID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetFlairTemplateByID", ctx, templateID).Return(&ftObj, nil)

			resp, err := backendObj.GetFlairTemplateByID(ctx, templateID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &ftObj)
		})
	})
}

func TestDelphisBackend_RemoveFlairTemplate(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	ftObj := test_utils.TestFlairTemplate()

	Convey("RemoveFlairTemplate", t, func() {
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
			mockDB.On("RemoveFlairTemplate", ctx, ftObj).Return(nil, expectedError)

			resp, err := backendObj.RemoveFlairTemplate(ctx, ftObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("RemoveFlairTemplate", ctx, ftObj).Return(&ftObj, nil)

			resp, err := backendObj.RemoveFlairTemplate(ctx, ftObj)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &ftObj)
		})
	})
}

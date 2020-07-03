package backend

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nedrocks/delphisbe/internal/backend/test_utils"

	"github.com/stretchr/testify/mock"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/auth"
	"github.com/nedrocks/delphisbe/internal/cache"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/nedrocks/delphisbe/internal/util"
	"github.com/nedrocks/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

type mockImportedContentIter struct{}

func (m *mockImportedContentIter) Next(content *model.ImportedContent) bool { return true }
func (m *mockImportedContentIter) Close() error                             { return fmt.Errorf("error") }

func TestDelphisBackend_GetUpcomingImportedContentByDiscussionID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	discussionID := test_utils.DiscussionID
	limit := test_utils.Limit

	icObj := test_utils.TestImportedContent()

	Convey("GetUpcomingImportedContentByDiscussionID", t, func() {
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

		Convey("when GetScheduledImportedContentByDiscussionID errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetScheduledImportedContentByDiscussionID", ctx, discussionID).Return(&mockImportedContentIter{})
			mockDB.On("ContentIterCollect", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetUpcomingImportedContentByDiscussionID(ctx, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when query succeeds", func() {
			mockDB.On("GetScheduledImportedContentByDiscussionID", ctx, discussionID).Return(&mockImportedContentIter{})
			mockDB.On("ContentIterCollect", ctx, mock.Anything).Return([]*model.ImportedContent{&icObj}, nil)
			mockDB.On("GetImportedContentByDiscussionID", ctx, discussionID, limit).Return(&mockImportedContentIter{})
			mockDB.On("ContentIterCollect", ctx, mock.Anything).Return([]*model.ImportedContent{&icObj}, nil)

			resp, err := backendObj.GetUpcomingImportedContentByDiscussionID(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.ImportedContent{&icObj, &icObj})
		})
	})
}

func TestDelphisBackend_GetImportedContentByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	icID := test_utils.ContentID

	icObj := test_utils.TestImportedContent()

	Convey("GetImportedContentByID", t, func() {
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

		Convey("when query errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetImportedContentByID", ctx, icID).Return(nil, expectedError)

			resp, err := backendObj.GetImportedContentByID(ctx, icID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when query succeeds", func() {
			mockDB.On("GetImportedContentByID", ctx, icID).Return(&icObj, nil)

			resp, err := backendObj.GetImportedContentByID(ctx, icID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &icObj)
		})
	})
}

func TestDelphisBackend_GetMatchingsTags(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	icID := test_utils.ContentID
	discussionID := test_utils.DiscussionID
	tags := []string{"tag1", "tag2"}

	Convey("GetMatchingTags", t, func() {
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

		Convey("when query errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetMatchingTags", ctx, discussionID, icID).Return(nil, expectedError)

			resp, err := backendObj.GetMatchingTags(ctx, discussionID, icID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when query succeeds", func() {
			mockDB.On("GetMatchingTags", ctx, discussionID, icID).Return(tags, nil)

			resp, err := backendObj.GetMatchingTags(ctx, discussionID, icID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, tags)
		})
	})
}

func TestDelphisBackend_PutImportedContentAndTags(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tags := []string{"tag1"}

	inputObj := test_utils.TestImportedContentInput()
	icObj := test_utils.TestImportedContent()
	tagObj := test_utils.TestContentTag()

	inputObj.Tags = tags[0]
	icObj.Tags = tags

	tx := sql.Tx{}

	Convey("PutImportedContentAndTags", t, func() {
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

		Convey("when BeginTx errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(nil, expectedError)

			resp, err := backendObj.PutImportedContentAndTags(ctx, inputObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when PutImportedContent errors outs and RollbackTx fails", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutImportedContent", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.PutImportedContentAndTags(ctx, inputObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when PutImportedContent errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutImportedContent", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.PutImportedContentAndTags(ctx, inputObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when PutImportedContentTags errors out and RollbackTx fails", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutImportedContent", ctx, mock.Anything, mock.Anything).Return(&icObj, nil)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)
			mockDB.On("PutImportedContentTags", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.PutImportedContentAndTags(ctx, inputObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when PutImportedContentTags errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutImportedContent", ctx, mock.Anything, mock.Anything).Return(&icObj, nil)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)
			mockDB.On("PutImportedContentTags", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.PutImportedContentAndTags(ctx, inputObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when CommitTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutImportedContent", ctx, mock.Anything, mock.Anything).Return(&icObj, nil)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)
			mockDB.On("PutImportedContentTags", ctx, mock.Anything, mock.Anything).Return(&tagObj, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.PutImportedContentAndTags(ctx, inputObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when call succeeds", func() {
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutImportedContent", ctx, mock.Anything, mock.Anything).Return(&icObj, nil)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)
			mockDB.On("PutImportedContentTags", ctx, mock.Anything, mock.Anything).Return(&tagObj, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.PutImportedContentAndTags(ctx, inputObj)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

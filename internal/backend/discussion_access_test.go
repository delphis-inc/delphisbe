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

type mockDiscussionIter struct{}

func (m *mockDiscussionIter) Next(discussion *model.Discussion) bool { return true }
func (m *mockDiscussionIter) Close() error                           { return fmt.Errorf("error") }

type mockDFAIter struct{}

func (m *mockDFAIter) Next(dfa *model.DiscussionFlairTemplateAccess) bool { return true }
func (m *mockDFAIter) Close() error                                       { return fmt.Errorf("error") }

func TestDelphisBackend_GetDiscussionAccessByUserID(t *testing.T) {
	ctx := context.Background()
	userID := test_utils.UserID

	discObj := test_utils.TestDiscussion()

	Convey("GetDiscussionAccessByUserID", t, func() {
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

		Convey("when GetPublicDiscussions errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetPublicDiscussions", ctx).Return(&mockDiscussionIter{})
			mockDB.On("DiscussionIterCollect", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetDiscussionAccessByUserID(ctx, userID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when GetDiscussionsForFlairTemplateByUserID errors outs", func() {
			publicIter := mockDiscussionIter{}

			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetPublicDiscussions", ctx).Return(&publicIter)
			mockDB.On("DiscussionIterCollect", ctx, publicIter).Return([]*model.Discussion{&discObj}, nil)
			mockDB.On("GetDiscussionsForFlairTemplateByUserID", ctx, userID).Return(&mockDiscussionIter{})
			mockDB.On("DiscussionIterCollect", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetDiscussionAccessByUserID(ctx, userID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when GetDiscussionsForUserAccessByUserID errors outs", func() {
			publicIter := mockDiscussionIter{}
			flairIter := mockDiscussionIter{}

			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetPublicDiscussions", ctx).Return(&publicIter)
			mockDB.On("DiscussionIterCollect", ctx, publicIter).Return([]*model.Discussion{&discObj}, nil)
			mockDB.On("GetDiscussionsForFlairTemplateByUserID", ctx, userID).Return(&flairIter)
			mockDB.On("DiscussionIterCollect", ctx, flairIter).Return([]*model.Discussion{&discObj}, nil)
			mockDB.On("GetDiscussionsForUserAccessByUserID", ctx, userID).Return(&mockDiscussionIter{})
			mockDB.On("DiscussionIterCollect", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetDiscussionAccessByUserID(ctx, userID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetPublicDiscussions", ctx).Return(&mockDiscussionIter{})
			mockDB.On("DiscussionIterCollect", ctx, mock.Anything).Return([]*model.Discussion{&discObj}, nil)
			mockDB.On("GetDiscussionsForFlairTemplateByUserID", ctx, userID).Return(&mockDiscussionIter{})
			mockDB.On("DiscussionIterCollect", ctx, mock.Anything).Return([]*model.Discussion{&discObj}, nil)
			mockDB.On("GetDiscussionsForUserAccessByUserID", ctx, userID).Return(&mockDiscussionIter{})
			mockDB.On("DiscussionIterCollect", ctx, mock.Anything).Return(nil, nil)

			resp, err := backendObj.GetDiscussionAccessByUserID(ctx, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.Discussion{&discObj})
		})
	})
}

func TestDelphisBackend_GetDiscussionFlairTemplateAccessByDiscussionID(t *testing.T) {
	ctx := context.Background()

	discussionID := test_utils.DiscussionID

	ftObj := test_utils.TestFlairTemplate()

	Convey("GetDiscussionFlairTemplateAccessByDiscussionID", t, func() {
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

		Convey("when the query errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetDiscussionFlairTemplatesAccessByDiscussionID", ctx, discussionID).Return(&mockDFAIter{})
			mockDB.On("FlairTemplatesIterCollect", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetDiscussionFlairTemplateAccessByDiscussionID(ctx, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetDiscussionFlairTemplatesAccessByDiscussionID", ctx, discussionID).Return(&mockDFAIter{})
			mockDB.On("FlairTemplatesIterCollect", ctx, mock.Anything).Return([]*model.FlairTemplate{&ftObj}, nil)

			resp, err := backendObj.GetDiscussionFlairTemplateAccessByDiscussionID(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.FlairTemplate{&ftObj})
		})
	})
}

func TestDelphisBackend_PutDiscussionFlairTemplatesAccess(t *testing.T) {
	ctx := context.Background()

	discussionID := test_utils.DiscussionID
	userID := test_utils.UserID

	ftObj := test_utils.TestFlairTemplate()
	flairObj := test_utils.TestFlair()
	dfaObj := test_utils.TestDiscussionFlairTemplateAccess()

	flairTemplateIDs := []string{ftObj.ID}

	tx := sql.Tx{}

	Convey("PutDiscussionFlairTemplatesAccess", t, func() {
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

		Convey("when no template IDs are passed in", func() {
			resp, err := backendObj.PutDiscussionFlairTemplatesAccess(ctx, userID, discussionID, nil)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when validateFlairTemplatesToAdd errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.PutDiscussionFlairTemplatesAccess(ctx, userID, discussionID, flairTemplateIDs)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when BeginTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return([]*model.Flair{&flairObj}, nil)
			mockDB.On("BeginTx", ctx).Return(nil, expectedError)

			resp, err := backendObj.PutDiscussionFlairTemplatesAccess(ctx, userID, discussionID, flairTemplateIDs)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when UpsertDiscussionFlairTemplatesAccess and RollbackTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return([]*model.Flair{&flairObj}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpsertDiscussionFlairTemplatesAccess", ctx, mock.Anything, discussionID, ftObj.ID).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.PutDiscussionFlairTemplatesAccess(ctx, userID, discussionID, flairTemplateIDs)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when UpsertDiscussionFlairTemplatesAccess errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return([]*model.Flair{&flairObj}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpsertDiscussionFlairTemplatesAccess", ctx, mock.Anything, discussionID, ftObj.ID).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.PutDiscussionFlairTemplatesAccess(ctx, userID, discussionID, flairTemplateIDs)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns and CommitTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return([]*model.Flair{&flairObj}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpsertDiscussionFlairTemplatesAccess", ctx, mock.Anything, discussionID, ftObj.ID).Return(&dfaObj, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.PutDiscussionFlairTemplatesAccess(ctx, userID, discussionID, flairTemplateIDs)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return([]*model.Flair{&flairObj}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpsertDiscussionFlairTemplatesAccess", ctx, mock.Anything, discussionID, ftObj.ID).Return(&dfaObj, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.PutDiscussionFlairTemplatesAccess(ctx, userID, discussionID, flairTemplateIDs)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.DiscussionFlairTemplateAccess{&dfaObj})
		})
	})
}

func TestDelphisBackend_DeleteDiscussionFlairTemplatesAccess(t *testing.T) {
	ctx := context.Background()

	discussionID := test_utils.DiscussionID

	ftObj := test_utils.TestFlairTemplate()
	flairObj := test_utils.TestFlair()
	dfaObj := test_utils.TestDiscussionFlairTemplateAccess()

	flairTemplateIDs := []string{ftObj.ID}

	tx := sql.Tx{}

	Convey("DeleteDiscussionFlairTemplatesAccess", t, func() {
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

		Convey("when no template IDs are passed in", func() {
			resp, err := backendObj.DeleteDiscussionFlairTemplatesAccess(ctx, discussionID, nil)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when BeginTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(nil, expectedError)

			resp, err := backendObj.DeleteDiscussionFlairTemplatesAccess(ctx, discussionID, flairTemplateIDs)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when UpsertDiscussionFlairTemplatesAccess and RollbackTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return([]*model.Flair{&flairObj}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("DeleteDiscussionFlairTemplatesAccess", ctx, mock.Anything, discussionID, ftObj.ID).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.DeleteDiscussionFlairTemplatesAccess(ctx, discussionID, flairTemplateIDs)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when UpsertDiscussionFlairTemplatesAccess errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return([]*model.Flair{&flairObj}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("DeleteDiscussionFlairTemplatesAccess", ctx, mock.Anything, discussionID, ftObj.ID).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.DeleteDiscussionFlairTemplatesAccess(ctx, discussionID, flairTemplateIDs)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns and CommitTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return([]*model.Flair{&flairObj}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("DeleteDiscussionFlairTemplatesAccess", ctx, mock.Anything, discussionID, ftObj.ID).Return(&dfaObj, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.DeleteDiscussionFlairTemplatesAccess(ctx, discussionID, flairTemplateIDs)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return([]*model.Flair{&flairObj}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("DeleteDiscussionFlairTemplatesAccess", ctx, mock.Anything, discussionID, ftObj.ID).Return(&dfaObj, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.DeleteDiscussionFlairTemplatesAccess(ctx, discussionID, flairTemplateIDs)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.DiscussionFlairTemplateAccess{&dfaObj})
		})
	})
}

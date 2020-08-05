package backend

import (
	"context"
	"database/sql"
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

type mockDiscussionIter struct{}

func (m *mockDiscussionIter) Next(discussion *model.Discussion) bool { return true }
func (m *mockDiscussionIter) Close() error                           { return fmt.Errorf("error") }

func TestDelphisBackend_GetDiscussionAccessByUserID(t *testing.T) {
	ctx := context.Background()
	userID := test_utils.UserID
	state := model.DiscussionUserAccessStateActive

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

		Convey("when GetDiscussionsByUserAccess errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetDiscussionsByUserAccess", ctx, userID, state).Return(&mockDiscussionIter{})
			mockDB.On("DiscussionIterCollect", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetDiscussionAccessByUserID(ctx, userID, state)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetDiscussionsByUserAccess", ctx, userID, state).Return(&mockDiscussionIter{})
			mockDB.On("DiscussionIterCollect", ctx, mock.Anything).Return([]*model.Discussion{&discObj}, nil)

			resp, err := backendObj.GetDiscussionAccessByUserID(ctx, userID, state)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.Discussion{&discObj})
		})
	})
}

func TestDelphisBackend_GrantUserDiscussionAccess(t *testing.T) {
	ctx := context.Background()

	discussionID := test_utils.DiscussionID
	userID := test_utils.UserID
	state := model.DiscussionUserAccessStateActive

	discussionUserAccess := test_utils.TestDiscussionUserAccess()

	tx := sql.Tx{}

	Convey("UpsertUserDiscussionAccess", t, func() {
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

		Convey("when creating a transaction fails", func() {
			mockDB.On("BeginTx", ctx).Return(nil, fmt.Errorf("sth"))

			resp, err := backendObj.UpsertUserDiscussionAccess(ctx, userID, discussionID, state)

			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		mockDB.On("BeginTx", ctx).Return(&tx, nil)

		Convey("when getting discussion user access errors out", func() {
			mockDB.On("GetDiscussionUserAccess", ctx, discussionID, userID).Return(nil, fmt.Errorf("error"))

			resp, err := backendObj.UpsertUserDiscussionAccess(ctx, userID, discussionID, state)

			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		Convey("when upserting user access fails", func() {
			mockDB.On("GetDiscussionUserAccess", ctx, discussionID, userID).Return(&discussionUserAccess, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, &tx, discussionUserAccess).Return(nil, fmt.Errorf("sth"))

			Convey("when rolling back transaction fails", func() {
				mockDB.On("RollbackTx", ctx, &tx).Return(fmt.Errorf("sth"))

				resp, err := backendObj.UpsertUserDiscussionAccess(ctx, userID, discussionID, state)

				So(resp, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})

			Convey("when rolling back transaction succeeds", func() {
				mockDB.On("RollbackTx", ctx, &tx).Return(nil)

				resp, err := backendObj.UpsertUserDiscussionAccess(ctx, userID, discussionID, state)

				So(resp, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})

		mockDB.On("GetDiscussionUserAccess", ctx, discussionID, userID).Return(&discussionUserAccess, nil)
		mockDB.On("UpsertDiscussionUserAccess", ctx, &tx, discussionUserAccess).Return(&discussionUserAccess, nil)

		Convey("when committing transaction fails", func() {
			mockDB.On("CommitTx", ctx, &tx).Return(fmt.Errorf("sth"))

			resp, err := backendObj.UpsertUserDiscussionAccess(ctx, userID, discussionID, state)

			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		mockDB.On("CommitTx", ctx, &tx).Return(nil)

		Convey("when successful", func() {
			resp, err := backendObj.UpsertUserDiscussionAccess(ctx, userID, discussionID, state)

			So(resp, ShouldResemble, &discussionUserAccess)
			So(err, ShouldBeNil)
		})
	})
}

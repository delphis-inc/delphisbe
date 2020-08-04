package backend

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/backend/test_utils"
	"github.com/delphis-inc/delphisbe/internal/cache"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/delphis-inc/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDelphisBackend_GetDiscussionIDsToBeShuffledBeforeTime(t *testing.T) {
	ctx := context.Background()
	discussionObj := test_utils.TestDiscussion()

	Convey("GetDiscussionIDsToBeShuffledBeforeTime", t, func() {
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

		Convey("when no tx is passed", func() {
			tx := sql.Tx{}
			Convey("when creating tx fails", func() {
				mockDB.On("BeginTx", ctx).Return(nil, fmt.Errorf("sth"))

				resp, err := backendObj.GetDiscussionIDsToBeShuffledBeforeTime(ctx, nil, now)

				So(resp, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})

			Convey("when creating tx succeeds", func() {
				mockDB.On("BeginTx", ctx).Return(&tx, nil)
				Convey("when getting discussions errors", func() {
					Convey("when rolling back fails", func() {
						mockDB.On("GetDiscussionsToBeShuffledBeforeTime", ctx, &tx, now).Return(nil, fmt.Errorf("sth"))
						mockDB.On("RollbackTx", ctx, &tx).Return(fmt.Errorf("sth"))

						resp, err := backendObj.GetDiscussionIDsToBeShuffledBeforeTime(ctx, nil, now)

						So(resp, ShouldBeNil)
						So(err, ShouldNotBeNil)
					})

					Convey("when rolling back succeeds", func() {
						mockDB.On("GetDiscussionsToBeShuffledBeforeTime", ctx, &tx, now).Return(nil, fmt.Errorf("sth"))
						mockDB.On("RollbackTx", ctx, &tx).Return(nil)

						resp, err := backendObj.GetDiscussionIDsToBeShuffledBeforeTime(ctx, nil, now)

						So(resp, ShouldBeNil)
						So(err, ShouldNotBeNil)
					})
				})

				Convey("when getting discussions succeeds", func() {
					Convey("when rolling back fails", func() {
						mockDB.On("GetDiscussionsToBeShuffledBeforeTime", ctx, &tx, now).Return([]model.Discussion{discussionObj}, nil)
						mockDB.On("RollbackTx", ctx, &tx).Return(fmt.Errorf("sth"))

						resp, err := backendObj.GetDiscussionIDsToBeShuffledBeforeTime(ctx, nil, now)

						So(err, ShouldBeNil)
						So(resp, ShouldResemble, []string{discussionObj.ID})
					})
				})
			})
		})
	})
}

func TestDelphisBackend_PutDiscussionShuffleTime(t *testing.T) {
	ctx := context.Background()
	discussionObj := test_utils.TestDiscussion()
	discussionShuffleTime := test_utils.TestDiscussionShuffleTime()

	Convey("PutDiscussionShuffleTime", t, func() {
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

		Convey("when begin transaction fails", func() {
			mockDB.On("BeginTx", ctx).Return(nil, fmt.Errorf("sth"))

			resp, err := backendObj.PutDiscussionShuffleTime(ctx, discussionObj.ID, &now)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when begin transaction succeeds", func() {
			tx := sql.Tx{}
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			Convey("when db call fails", func() {
				Convey("when rollback fails", func() {
					mockDB.On("PutNextShuffleTimeForDiscussionID", ctx, &tx, discussionObj.ID, &now).Return(nil, fmt.Errorf("sth"))
					mockDB.On("RollbackTx", ctx, &tx).Return(fmt.Errorf("sth"))

					resp, err := backendObj.PutDiscussionShuffleTime(ctx, discussionObj.ID, &now)

					So(err, ShouldNotBeNil)
					So(resp, ShouldBeNil)
				})
				Convey("when rollback succeeds", func() {
					mockDB.On("PutNextShuffleTimeForDiscussionID", ctx, &tx, discussionObj.ID, &now).Return(nil, fmt.Errorf("sth"))
					mockDB.On("RollbackTx", ctx, &tx).Return(nil)

					resp, err := backendObj.PutDiscussionShuffleTime(ctx, discussionObj.ID, &now)

					So(err, ShouldNotBeNil)
					So(resp, ShouldBeNil)
				})
			})
			Convey("when db call succeeds", func() {
				Convey("and commit fails", func() {
					mockDB.On("PutNextShuffleTimeForDiscussionID", ctx, &tx, discussionObj.ID, &now).Return(&discussionShuffleTime, nil)
					mockDB.On("CommitTx", ctx, &tx).Return(fmt.Errorf("sth"))

					resp, err := backendObj.PutDiscussionShuffleTime(ctx, discussionObj.ID, &now)

					So(err, ShouldNotBeNil)
					So(resp, ShouldBeNil)
				})
				Convey("and commit succeeds", func() {
					mockDB.On("PutNextShuffleTimeForDiscussionID", ctx, &tx, discussionObj.ID, &now).Return(&discussionShuffleTime, nil)
					mockDB.On("CommitTx", ctx, &tx).Return(nil)

					resp, err := backendObj.PutDiscussionShuffleTime(ctx, discussionObj.ID, &now)

					So(err, ShouldBeNil)
					So(resp, ShouldResemble, &discussionShuffleTime)
				})
			})
		})
	})
}

// NOTE: This test is odd because nothing is returned so we're really
// just testing the mocks.
func TestDelphisBackend_ShuffleDiscussionsIfNecessary(t *testing.T) {
	ctx := context.Background()
	discussionObj := test_utils.TestDiscussion()
	// discussionShuffleTime := test_utils.TestDiscussionShuffleTime()

	Convey("ShuffleDiscussionsIfNecessary", t, func() {
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

		Convey("when begin transaction fails", func() {
			mockDB.On("BeginTx", ctx).Return(nil, fmt.Errorf("sth"))

			backendObj.ShuffleDiscussionsIfNecessary()
		})

		Convey("when begin transaction succeeds", func() {
			tx := sql.Tx{}
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			Convey("when getting discussions errors", func() {
				Convey("and rollback errors", func() {
					mockDB.On("GetDiscussionsToBeShuffledBeforeTime", ctx, &tx, now).Return(nil, fmt.Errorf("sth"))
					mockDB.On("RollbackTx", ctx, &tx).Return(fmt.Errorf("sth"))

					backendObj.ShuffleDiscussionsIfNecessary()
				})
				Convey("and rollback succeeds", func() {
					mockDB.On("GetDiscussionsToBeShuffledBeforeTime", ctx, &tx, now).Return(nil, fmt.Errorf("sth"))
					mockDB.On("RollbackTx", ctx, &tx).Return(nil)

					backendObj.ShuffleDiscussionsIfNecessary()
				})
			})
			Convey("when getting discussions returns an empty array", func() {
				mockDB.On("GetDiscussionsToBeShuffledBeforeTime", ctx, &tx, now).Return([]model.Discussion{}, nil)

				Convey("and committing fails", func() {
					mockDB.On("CommitTx", ctx, &tx).Return(fmt.Errorf("sth"))

					backendObj.ShuffleDiscussionsIfNecessary()
				})
				Convey("and commit succeeds", func() {
					mockDB.On("CommitTx", ctx, &tx).Return(nil)

					backendObj.ShuffleDiscussionsIfNecessary()
				})
			})
			Convey("when getting discussions returns a non-empty array", func() {
				var nilTime *time.Time
				Convey("and all things work", func() {
					mockDB.On("GetDiscussionsToBeShuffledBeforeTime", ctx, &tx, now).Return([]model.Discussion{discussionObj}, nil)
					mockDB.On("IncrementDiscussionShuffleCount", ctx, &tx, discussionObj.ID).Return(nil, nil)
					mockDB.On("PutNextShuffleTimeForDiscussionID", ctx, &tx, discussionObj.ID, nilTime).Return(nil, nil)
					mockDB.On("CommitTx", ctx, &tx).Return(nil)

					backendObj.ShuffleDiscussionsIfNecessary()
				})
			})
		})
	})
}

func TestDelphisBackend_GetNextDiscussionShuffleTime(t *testing.T) {
	ctx := context.Background()
	discussionID := "1"
	Convey("GetNextDiscussionShuffleTime", t, func() {
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

		Convey("when underlying call returns an error", func() {
			mockDB.On("GetNextShuffleTimeForDiscussionID", ctx, discussionID).Return(nil, fmt.Errorf("sth"))

			resp, err := backendObj.GetNextDiscussionShuffleTime(ctx, discussionID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when underlying call does not error", func() {
			shuffleTime := model.DiscussionShuffleTime{}
			mockDB.On("GetNextShuffleTimeForDiscussionID", ctx, discussionID).Return(&shuffleTime, nil)

			resp, err := backendObj.GetNextDiscussionShuffleTime(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldEqual, &shuffleTime)
		})
	})
}

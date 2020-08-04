package datastore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/delphis-inc/delphisbe/graph/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestDelphisDB_GetNextShuffleTimeForDiscussionID(t *testing.T) {
	ctx := context.Background()
	discussionID := "discussion1"
	now := time.Now()

	Convey("GetNextShuffleTimeForDiscussionID", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		assert.Nil(t, err, "Failed setting up sqlmock db")

		gormDB, _ := gorm.Open("postgres", db)
		mockDatastore := &delphisDB{
			dbConfig:  config.TablesConfig{},
			sql:       gormDB,
			pg:        db,
			prepStmts: &dbPrepStmts{},
			dynamo:    nil,
			encoder:   nil,
		}
		defer db.Close()

		Convey("when preparing statements returns an error", func() {
			mockPreparedStatementsWithError(mock)

			resp, err := mockDatastore.GetNextShuffleTimeForDiscussionID(ctx, discussionID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getNextShuffleTimeForDiscussionIDString).WithArgs(discussionID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetNextShuffleTimeForDiscussionID(ctx, discussionID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns data", func() {
			expected := model.DiscussionShuffleTime{
				DiscussionID: discussionID,
				ShuffleTime:  &now,
			}
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"discussion_id", "shuffle_time"}).
				AddRow(discussionID, now)

			mock.ExpectQuery(getNextShuffleTimeForDiscussionIDString).WithArgs(discussionID).WillReturnRows(rs)

			resp, err := mockDatastore.GetNextShuffleTimeForDiscussionID(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &expected)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and nil shuffle time", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"discussion_id", "shuffle_time"}).
				AddRow(discussionID, nil)

			mock.ExpectQuery(getNextShuffleTimeForDiscussionIDString).WithArgs(discussionID).WillReturnRows(rs)

			resp, err := mockDatastore.GetNextShuffleTimeForDiscussionID(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetDiscussionsToBeShuffledBeforeTime(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	Convey("GetDiscussionsToBeShuffledBeforeTime", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		assert.Nil(t, err, "Failed setting up sqlmock db")

		gormDB, _ := gorm.Open("postgres", db)
		mockDatastore := &delphisDB{
			dbConfig:  config.TablesConfig{},
			sql:       gormDB,
			pg:        db,
			prepStmts: &dbPrepStmts{},
			dynamo:    nil,
			encoder:   nil,
		}
		defer db.Close()

		Convey("when preparing statements returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatementsWithError(mock)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.GetDiscussionsToBeShuffledBeforeTime(ctx, tx, now)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(getDiscussionsToShuffle)
			mock.ExpectQuery(getDiscussionsToShuffle).WithArgs(now).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.GetDiscussionsToBeShuffledBeforeTime(ctx, tx, now)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query succeeds and returns 2 objects", func() {
			expected := []model.Discussion{
				{
					ID:           "discussion1",
					ShuffleCount: 0,
				},
				{
					ID:           "discussion2",
					ShuffleCount: 1,
				},
			}
			rs := sqlmock.NewRows([]string{"discussion_id", "shuffle_count"}).
				AddRow("discussion1", 0).
				AddRow("discussion2", 1)

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(getDiscussionsToShuffle)
			mock.ExpectQuery(getDiscussionsToShuffle).WithArgs(now).WillReturnRows(rs)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.GetDiscussionsToBeShuffledBeforeTime(ctx, tx, now)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, expected)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_PutNextShuffleTimeForDiscussionID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"

	Convey("PutNextShuffleTimeForDiscussionID", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		assert.Nil(t, err, "Failed setting up sqlmock db")

		gormDB, _ := gorm.Open("postgres", db)
		mockDatastore := &delphisDB{
			dbConfig:  config.TablesConfig{},
			sql:       gormDB,
			pg:        db,
			prepStmts: &dbPrepStmts{},
			dynamo:    nil,
			encoder:   nil,
		}
		defer db.Close()

		Convey("when preparing statements returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatementsWithError(mock)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.PutNextShuffleTimeForDiscussionID(ctx, tx, discussionID, &now)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putNextShuffleTimeForDiscussionIDString)
			mock.ExpectQuery(putNextShuffleTimeForDiscussionIDString).WithArgs(discussionID, &now).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.PutNextShuffleTimeForDiscussionID(ctx, tx, discussionID, &now)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query succeeds and returns 2 objects", func() {
			expected := model.DiscussionShuffleTime{
				DiscussionID: discussionID,
				ShuffleTime:  &now,
			}
			rs := sqlmock.NewRows([]string{"discussion_id", "shuffle_count"}).
				AddRow("discussion1", &now)

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putNextShuffleTimeForDiscussionIDString)
			mock.ExpectQuery(putNextShuffleTimeForDiscussionIDString).WithArgs(discussionID, &now).WillReturnRows(rs)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.PutNextShuffleTimeForDiscussionID(ctx, tx, discussionID, &now)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &expected)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

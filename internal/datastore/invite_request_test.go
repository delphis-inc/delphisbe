package datastore

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"

	"time"
)

func TestDelphisDB_GetDiscussionRequestAccessByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	requestID := "request1"
	requestObj := model.DiscussionAccessRequest{
		ID:           requestID,
		UserID:       "user1",
		DiscussionID: "discussion1",
		CreatedAt:    now.Format(time.RFC3339),
		UpdatedAt:    now.Format(time.RFC3339),
		IsDeleted:    false,
		Status:       model.InviteRequestStatusPending,
	}

	Convey("GetDiscussionRequestAccessByID", t, func() {
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

			resp, err := mockDatastore.GetDiscussionRequestAccessByID(ctx, requestID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getDiscussionRequestAccessByIDString).WithArgs(requestID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetDiscussionRequestAccessByID(ctx, requestID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns imported content", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"id", "user_id", "discussion_id", "created_at",
				"updated_at", "status"}).
				AddRow(requestObj.ID, requestObj.UserID, requestObj.DiscussionID, requestObj.CreatedAt,
					requestObj.UpdatedAt, requestObj.Status)

			mock.ExpectQuery(getDiscussionRequestAccessByIDString).WithArgs(requestID).WillReturnRows(rs)

			resp, err := mockDatastore.GetDiscussionRequestAccessByID(ctx, requestID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &requestObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetDiscussionAccessRequestByDiscussionIDUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	userID := "user1"
	requestObj := model.DiscussionAccessRequest{
		ID:           "request1",
		UserID:       "user1",
		DiscussionID: discussionID,
		CreatedAt:    now.Format(time.RFC3339),
		UpdatedAt:    now.Format(time.RFC3339),
		IsDeleted:    false,
		Status:       model.InviteRequestStatusPending,
	}

	Convey("GetDiscussionAccessRequestByDiscussionIDUserID", t, func() {
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

			resp, err := mockDatastore.GetDiscussionAccessRequestByDiscussionIDUserID(ctx, discussionID, userID)

			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getDiscussionAccessRequestByUserIDString).WithArgs(discussionID, userID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetDiscussionAccessRequestByDiscussionIDUserID(ctx, discussionID, userID)

			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns no rows", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"id", "user_id", "discussion_id", "created_at",
				"updated_at", "status"})
			mock.ExpectQuery(getDiscussionAccessRequestByUserIDString).WithArgs(discussionID, userID).WillReturnRows(rs)

			resp, err := mockDatastore.GetDiscussionAccessRequestByDiscussionIDUserID(ctx, discussionID, userID)

			So(resp, ShouldBeNil)
			So(err, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns a row", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"id", "user_id", "discussion_id", "created_at",
				"updated_at", "status"}).AddRow(
				requestObj.ID, requestObj.UserID, requestObj.DiscussionID, requestObj.CreatedAt,
				requestObj.UpdatedAt, requestObj.Status,
			)
			mock.ExpectQuery(getDiscussionAccessRequestByUserIDString).WithArgs(discussionID, userID).WillReturnRows(rs)

			resp, err := mockDatastore.GetDiscussionAccessRequestByDiscussionIDUserID(ctx, discussionID, userID)

			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &requestObj)
			So(err, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetDiscussionAccessRequestsByDiscussionID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	requestObj := model.DiscussionAccessRequest{
		ID:           "request1",
		UserID:       "user1",
		DiscussionID: discussionID,
		CreatedAt:    now.Format(time.RFC3339),
		UpdatedAt:    now.Format(time.RFC3339),
		IsDeleted:    false,
		Status:       model.InviteRequestStatusPending,
	}

	emptyRequest := model.DiscussionAccessRequest{}

	Convey("GetDiscussionAccessRequestsByDiscussionID", t, func() {
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

			iter := mockDatastore.GetDiscussionAccessRequestsByDiscussionID(ctx, discussionID)

			So(iter.Next(&emptyRequest), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getDiscussionAccessRequestsString).WithArgs(discussionID).WillReturnError(fmt.Errorf("error"))

			iter := mockDatastore.GetDiscussionAccessRequestsByDiscussionID(ctx, discussionID)

			So(iter.Next(&emptyRequest), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns posts", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"id", "user_id", "discussion_id", "created_at",
				"updated_at", "status"}).
				AddRow(requestObj.ID, requestObj.UserID, requestObj.DiscussionID, requestObj.CreatedAt,
					requestObj.UpdatedAt, requestObj.Status).
				AddRow(requestObj.ID, requestObj.UserID, requestObj.DiscussionID, requestObj.CreatedAt,
					requestObj.UpdatedAt, requestObj.Status)

			mock.ExpectQuery(getDiscussionAccessRequestsString).WithArgs(discussionID).WillReturnRows(rs)

			iter := mockDatastore.GetDiscussionAccessRequestsByDiscussionID(ctx, discussionID)

			So(iter.Next(&emptyRequest), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetSentDiscussionAccessRequestsByUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	userID := "user1"
	requestObj := model.DiscussionAccessRequest{
		ID:           "request1",
		UserID:       userID,
		DiscussionID: "discussion1",
		CreatedAt:    now.Format(time.RFC3339),
		UpdatedAt:    now.Format(time.RFC3339),
		IsDeleted:    false,
		Status:       model.InviteRequestStatusPending,
	}

	emptyRequest := model.DiscussionAccessRequest{}

	Convey("GetSentDiscussionAccessRequestsByUserID", t, func() {
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

			iter := mockDatastore.GetSentDiscussionAccessRequestsByUserID(ctx, userID)

			So(iter.Next(&emptyRequest), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getSentDiscussionAccessRequestsForUserString).WithArgs(userID).WillReturnError(fmt.Errorf("error"))

			iter := mockDatastore.GetSentDiscussionAccessRequestsByUserID(ctx, userID)

			So(iter.Next(&emptyRequest), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns posts", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"id", "user_id", "discussion_id", "created_at",
				"updated_at", "status"}).
				AddRow(requestObj.ID, requestObj.UserID, requestObj.DiscussionID, requestObj.CreatedAt,
					requestObj.UpdatedAt, requestObj.Status).
				AddRow(requestObj.ID, requestObj.UserID, requestObj.DiscussionID, requestObj.CreatedAt,
					requestObj.UpdatedAt, requestObj.Status)

			mock.ExpectQuery(getSentDiscussionAccessRequestsForUserString).WithArgs(userID).WillReturnRows(rs)

			iter := mockDatastore.GetSentDiscussionAccessRequestsByUserID(ctx, userID)

			So(iter.Next(&emptyRequest), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_PutDiscussionAccessRequestRecord(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	requestID := "request1"
	requestObj := model.DiscussionAccessRequest{
		ID:           requestID,
		UserID:       "user1",
		DiscussionID: "discussion1",
		CreatedAt:    now.Format(time.RFC3339),
		UpdatedAt:    now.Format(time.RFC3339),
		IsDeleted:    false,
		Status:       model.InviteRequestStatusPending,
	}

	Convey("PutDiscussionAccessRequestRecord", t, func() {
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
			resp, err := mockDatastore.PutDiscussionAccessRequestRecord(ctx, tx, requestObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putDiscussionAccessRequestString)
			mock.ExpectQuery(putDiscussionAccessRequestString).WithArgs(requestObj.ID, requestObj.UserID, requestObj.DiscussionID,
				requestObj.Status).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.PutDiscussionAccessRequestRecord(ctx, tx, requestObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when put post succeeds and returns an object", func() {
			rs := sqlmock.NewRows([]string{"id", "user_id", "discussion_id", "created_at",
				"updated_at", "status"}).
				AddRow(requestObj.ID, requestObj.UserID, requestObj.DiscussionID, requestObj.CreatedAt,
					requestObj.UpdatedAt, requestObj.Status)

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putDiscussionAccessRequestString)
			mock.ExpectQuery(putDiscussionAccessRequestString).WithArgs(requestObj.ID, requestObj.UserID, requestObj.DiscussionID,
				requestObj.Status).WillReturnRows(rs)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.PutDiscussionAccessRequestRecord(ctx, tx, requestObj)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &requestObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_UpdateDiscussionAccessRequestRecord(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	requestID := "request1"
	requestObj := model.DiscussionAccessRequest{
		ID:           requestID,
		UserID:       "user1",
		DiscussionID: "discussion1",
		CreatedAt:    now.Format(time.RFC3339),
		UpdatedAt:    now.Format(time.RFC3339),
		IsDeleted:    false,
		Status:       model.InviteRequestStatusPending,
	}

	Convey("UpdateDiscussionAccessRequestRecord", t, func() {
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
			resp, err := mockDatastore.UpdateDiscussionAccessRequestRecord(ctx, tx, requestObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(updateDiscussionAccessRequestString)
			mock.ExpectQuery(updateDiscussionAccessRequestString).WithArgs(requestObj.ID,
				requestObj.Status).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.UpdateDiscussionAccessRequestRecord(ctx, tx, requestObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when put post succeeds and returns an object", func() {
			rs := sqlmock.NewRows([]string{"id", "user_id", "discussion_id", "created_at",
				"updated_at", "status"}).
				AddRow(requestObj.ID, requestObj.UserID, requestObj.DiscussionID, requestObj.CreatedAt,
					requestObj.UpdatedAt, requestObj.Status)

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(updateDiscussionAccessRequestString)
			mock.ExpectQuery(updateDiscussionAccessRequestString).WithArgs(requestObj.ID,
				requestObj.Status).WillReturnRows(rs)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.UpdateDiscussionAccessRequestRecord(ctx, tx, requestObj)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &requestObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDiscussionAccessRequestIter_Next(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	requestID := "request1"
	requestObj := model.DiscussionAccessRequest{
		ID:           requestID,
		UserID:       "user1",
		DiscussionID: "discussion1",
		CreatedAt:    now.Format(time.RFC3339),
		UpdatedAt:    now.Format(time.RFC3339),
		IsDeleted:    false,
		Status:       model.InviteRequestStatusPending,
	}

	emptyRequest := model.DiscussionAccessRequest{}

	Convey("DiscussionAccessRequestIter_Next", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := discussionAccessRequestIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Next(&emptyRequest), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has a context error passed in", func() {
			ctx1, cancelFunc := context.WithCancel(ctx)
			cancelFunc()
			iter := discussionAccessRequestIter{
				ctx: ctx1,
			}

			So(iter.Next(&emptyRequest), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has no more rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"id", "user_id", "discussion_id", "created_at",
				"updated_at", "status"})

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := discussionAccessRequestIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyRequest), ShouldBeFalse)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on scan", func() {
			rs := sqlmock.NewRows([]string{"id", "user_id", "discussion_id", "created_at",
				"updated_at"}).
				AddRow(requestObj.ID, requestObj.UserID, requestObj.DiscussionID, requestObj.CreatedAt,
					requestObj.UpdatedAt)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := discussionAccessRequestIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyRequest), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"id", "user_id", "discussion_id", "created_at",
				"updated_at", "status"}).
				AddRow(requestObj.ID, requestObj.UserID, requestObj.DiscussionID, requestObj.CreatedAt,
					requestObj.UpdatedAt, requestObj.Status)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := discussionAccessRequestIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyRequest), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDiscussionAccessRequestIter_Close(t *testing.T) {
	ctx := context.Background()

	Convey("DiscussionAccessRequestIter_Close", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := discussionAccessRequestIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on rows.Close", func() {
			rs := sqlmock.NewRows([]string{"id", "user_id", "discussion_id", "created_at",
				"updated_at", "status"}).CloseError(fmt.Errorf("error"))

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := discussionAccessRequestIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_AccessRequestIterCollect(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	requestID := "request1"
	requestObj := model.DiscussionAccessRequest{
		ID:           requestID,
		UserID:       "user1",
		DiscussionID: "discussion1",
		CreatedAt:    now.Format(time.RFC3339),
		UpdatedAt:    now.Format(time.RFC3339),
		IsDeleted:    false,
		Status:       model.InviteRequestStatusPending,
	}

	Convey("AccessRequestIterCollect", t, func() {
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

		Convey("when the iterator fails to close", func() {
			iter := &discussionAccessRequestIter{
				err: fmt.Errorf("error"),
			}

			resp, err := mockDatastore.AccessRequestIterCollect(ctx, iter)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has results and returns slice of DiscussionAccessRequests", func() {
			rs := sqlmock.NewRows([]string{"id", "user_id", "discussion_id", "created_at",
				"updated_at", "status"}).
				AddRow(requestObj.ID, requestObj.UserID, requestObj.DiscussionID, requestObj.CreatedAt,
					requestObj.UpdatedAt, requestObj.Status).
				AddRow(requestObj.ID, requestObj.UserID, requestObj.DiscussionID, requestObj.CreatedAt,
					requestObj.UpdatedAt, requestObj.Status)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := &discussionAccessRequestIter{
				ctx:  ctx,
				rows: rs1,
			}

			resp, err := mockDatastore.AccessRequestIterCollect(ctx, iter)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []*model.DiscussionAccessRequest{&requestObj, &requestObj})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

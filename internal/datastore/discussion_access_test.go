package datastore

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/delphis-inc/delphisbe/internal/backend/test_utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

func TestDelphisDB_GetDiscussionsForUserAccess(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	userID := "userID"
	state := model.DiscussionUserAccessStateActive
	emptyString := ""
	discObj := model.Discussion{
		ID:            "discussion1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
		Title:         "test",
		AnonymityType: "",
		ModeratorID:   &modID,
		IconURL:       &emptyString,
	}

	emptyDisc := model.Discussion{}

	Convey("GetDiscussionsByUserAccess", t, func() {
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

			iter := mockDatastore.GetDiscussionsByUserAccess(ctx, userID, state)

			So(iter.Next(&emptyDisc), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getDiscussionsByUserAccessString).WithArgs(userID, state).WillReturnError(fmt.Errorf("error"))

			iter := mockDatastore.GetDiscussionsByUserAccess(ctx, userID, state)

			So(iter.Next(&emptyDisc), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns discussions", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title",
				"anonymity_type", "moderator_id", "icon_url", "description", "title_history",
				"description_history", "discussion_joinability", "last_post_id", "last_post_created_at",
				"shuffle_count", "lock_status"}).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt,
					discObj.Title, discObj.AnonymityType, discObj.ModeratorID, discObj.IconURL,
					discObj.Description, discObj.TitleHistory, discObj.DescriptionHistory,
					discObj.DiscussionJoinability, discObj.LastPostID, discObj.LastPostCreatedAt,
					discObj.ShuffleCount, discObj.LockStatus)

			mock.ExpectQuery(getDiscussionsByUserAccessString).WithArgs(userID, state).WillReturnRows(rs)

			iter := mockDatastore.GetDiscussionsByUserAccess(ctx, userID, state)

			So(iter.Next(&emptyDisc), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetDiscussionUserAccess(t *testing.T) {
	ctx := context.Background()
	userID := "userID"
	discussionID := "discussionID"
	requestID := "requestID"
	duaObj := test_utils.TestDiscussionUserAccess()

	duaObj.RequestID = &requestID

	Convey("GetDiscussionUserAccess", t, func() {
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

			resp, err := mockDatastore.GetDiscussionUserAccess(ctx, discussionID, userID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getDiscussionUserAccessString).WithArgs(discussionID, userID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetDiscussionUserAccess(ctx, discussionID, userID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when there is no data returned", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getDiscussionUserAccessString).WithArgs(discussionID, userID).WillReturnError(sql.ErrNoRows)

			resp, err := mockDatastore.GetDiscussionUserAccess(ctx, discussionID, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "user_id", "state", "request_id",
				"notif_setting", "created_at", "updated_at", "deleted_at"}).
				AddRow(duaObj.DiscussionID, duaObj.UserID, duaObj.State, duaObj.RequestID,
					duaObj.NotifSetting, duaObj.CreatedAt, duaObj.UpdatedAt, duaObj.DeletedAt)

			mockPreparedStatements(mock)
			mock.ExpectQuery(getDiscussionUserAccessString).WithArgs(discussionID, userID).WillReturnRows(rs)

			resp, err := mockDatastore.GetDiscussionUserAccess(ctx, discussionID, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &duaObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetDUAForEverythingNotifications(t *testing.T) {
	ctx := context.Background()
	userID := test_utils.UserID
	discussionID := test_utils.DiscussionID
	duaObj := test_utils.TestDiscussionUserAccess()

	emptyDuaObj := model.DiscussionUserAccess{}

	Convey("GetDUAForEverythingNotifications", t, func() {
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

			iter := mockDatastore.GetDUAForEverythingNotifications(ctx, discussionID, userID)

			So(iter.Next(&emptyDuaObj), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getDUAForEverythingNotificationsString).WithArgs(discussionID, userID).WillReturnError(fmt.Errorf("error"))

			iter := mockDatastore.GetDUAForEverythingNotifications(ctx, discussionID, userID)

			So(iter.Next(&emptyDuaObj), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns discussions", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"discussion_id", "user_id", "state", "request_id",
				"notif_setting", "created_at", "updated_at", "deleted_at"}).
				AddRow(duaObj.DiscussionID, duaObj.UserID, duaObj.State, duaObj.RequestID,
					duaObj.NotifSetting, duaObj.CreatedAt, duaObj.UpdatedAt, duaObj.DeletedAt)

			mock.ExpectQuery(getDUAForEverythingNotificationsString).WithArgs(discussionID, userID).WillReturnRows(rs)

			iter := mockDatastore.GetDUAForEverythingNotifications(ctx, discussionID, userID)

			So(iter.Next(&emptyDuaObj), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetDUAForMentionNotifications(t *testing.T) {
	ctx := context.Background()
	discussionID := test_utils.DiscussionID
	duaObj := test_utils.TestDiscussionUserAccess()

	postingUser := "postingUser"
	mentionedUsers := []string{test_utils.UserID}

	emptyDuaObj := model.DiscussionUserAccess{}

	Convey("GetDUAForMentionNotifications", t, func() {
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

			iter := mockDatastore.GetDUAForMentionNotifications(ctx, discussionID, postingUser, mentionedUsers)

			So(iter.Next(&emptyDuaObj), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getDUAForMentionNotificationsString).WithArgs(discussionID, postingUser, pq.Array(mentionedUsers)).WillReturnError(fmt.Errorf("error"))

			iter := mockDatastore.GetDUAForMentionNotifications(ctx, discussionID, postingUser, mentionedUsers)

			So(iter.Next(&emptyDuaObj), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns discussions", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"discussion_id", "user_id", "state", "request_id",
				"notif_setting", "created_at", "updated_at", "deleted_at"}).
				AddRow(duaObj.DiscussionID, duaObj.UserID, duaObj.State, duaObj.RequestID,
					duaObj.NotifSetting, duaObj.CreatedAt, duaObj.UpdatedAt, duaObj.DeletedAt)

			mock.ExpectQuery(getDUAForMentionNotificationsString).WithArgs(discussionID, postingUser, pq.Array(mentionedUsers)).WillReturnRows(rs)

			iter := mockDatastore.GetDUAForMentionNotifications(ctx, discussionID, postingUser, mentionedUsers)

			So(iter.Next(&emptyDuaObj), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_UpsertDiscussionUserAccess(t *testing.T) {
	ctx := context.Background()
	userID := "userID"
	discussionID := "discussionID"
	state := model.DiscussionUserAccessStateActive
	notifSetting := model.DiscussionUserNotificationSettingEverything
	requestID := "requestID"
	duaObj := test_utils.TestDiscussionUserAccess()

	duaObj.RequestID = &requestID

	Convey("UpsertDiscussionUserAccess", t, func() {
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

			tx, _ := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.UpsertDiscussionUserAccess(ctx, tx, duaObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(upsertDiscussionUserAccessString)
			mock.ExpectQuery(upsertDiscussionUserAccessString).WithArgs(discussionID, userID, state, &requestID, notifSetting).WillReturnError(fmt.Errorf("error"))

			tx, _ := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.UpsertDiscussionUserAccess(ctx, tx, duaObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when there is no new data to upsert", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(upsertDiscussionUserAccessString)
			mock.ExpectQuery(upsertDiscussionUserAccessString).WithArgs(discussionID, userID, state, &requestID, notifSetting).WillReturnError(sql.ErrNoRows)

			tx, _ := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.UpsertDiscussionUserAccess(ctx, tx, duaObj)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &model.DiscussionUserAccess{})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns discussions", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "user_id", "state", "request_id",
				"notif_setting", "created_at", "updated_at", "deleted_at"}).
				AddRow(duaObj.DiscussionID, duaObj.UserID, duaObj.State, duaObj.RequestID,
					duaObj.NotifSetting, duaObj.CreatedAt, duaObj.UpdatedAt, duaObj.DeletedAt)

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(upsertDiscussionUserAccessString)
			mock.ExpectQuery(upsertDiscussionUserAccessString).WithArgs(discussionID, userID, state, &requestID, notifSetting).WillReturnRows(rs)

			tx, _ := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.UpsertDiscussionUserAccess(ctx, tx, duaObj)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &duaObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_DeleteDiscussionUserAccess(t *testing.T) {
	ctx := context.Background()
	userID := "userID"
	discussionID := "discussionID"
	duaObj := model.DiscussionUserAccess{
		DiscussionID: discussionID,
		UserID:       userID,
	}

	Convey("DeleteDiscussionUserAccess", t, func() {
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
			resp, err := mockDatastore.DeleteDiscussionUserAccess(ctx, tx, discussionID, userID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(deleteDiscussionUserAccessString)
			mock.ExpectQuery(deleteDiscussionUserAccessString).WithArgs(discussionID, userID).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.DeleteDiscussionUserAccess(ctx, tx, discussionID, userID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns discussions", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "user_id", "created_at", "updated_at", "deleted_at"}).
				AddRow(duaObj.DiscussionID, duaObj.UserID, duaObj.CreatedAt, duaObj.UpdatedAt, duaObj.DeletedAt)

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(deleteDiscussionUserAccessString)
			mock.ExpectQuery(deleteDiscussionUserAccessString).WithArgs(discussionID, userID).WillReturnRows(rs)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.DeleteDiscussionUserAccess(ctx, tx, discussionID, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &duaObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_DuaIterCollect(t *testing.T) {
	ctx := context.Background()

	duaObj := test_utils.TestDiscussionUserAccess()

	Convey("DuaIterCollect", t, func() {
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
			iter := &duaIter{
				err: fmt.Errorf("error"),
			}

			resp, err := mockDatastore.DuaIterCollect(ctx, iter)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has results and returns slice of Discussions", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "user_id", "state", "request_id",
				"notif_setting", "created_at", "updated_at", "deleted_at"}).
				AddRow(duaObj.DiscussionID, duaObj.UserID, duaObj.State, duaObj.RequestID,
					duaObj.NotifSetting, duaObj.CreatedAt, duaObj.UpdatedAt, duaObj.DeletedAt).
				AddRow(duaObj.DiscussionID, duaObj.UserID, duaObj.State, duaObj.RequestID,
					duaObj.NotifSetting, duaObj.CreatedAt, duaObj.UpdatedAt, duaObj.DeletedAt)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := &duaIter{
				ctx:  ctx,
				rows: rs1,
			}

			resp, err := mockDatastore.DuaIterCollect(ctx, iter)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []*model.DiscussionUserAccess{&duaObj, &duaObj})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDiscussionUserAccessIter_Next(t *testing.T) {
	ctx := context.Background()
	duaObj := test_utils.TestDiscussionUserAccess()

	emptyDuaObj := model.DiscussionUserAccess{}

	Convey("DiscussionIter_Next", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := duaIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Next(&emptyDuaObj), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has a context error passed in", func() {
			ctx1, cancelFunc := context.WithCancel(ctx)
			cancelFunc()
			iter := duaIter{
				ctx: ctx1,
			}

			So(iter.Next(&emptyDuaObj), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has no more rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "user_id", "state", "request_id",
				"notif_setting", "created_at", "updated_at", "deleted_at"})

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := duaIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyDuaObj), ShouldBeFalse)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on scan", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "user_id", "state", "request_id",
				"notif_setting", "created_at", "updated_at"}).
				AddRow(duaObj.DiscussionID, duaObj.UserID, duaObj.State, duaObj.RequestID,
					duaObj.NotifSetting, duaObj.CreatedAt, duaObj.UpdatedAt)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := duaIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyDuaObj), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "user_id", "state", "request_id",
				"notif_setting", "created_at", "updated_at", "deleted_at"}).
				AddRow(duaObj.DiscussionID, duaObj.UserID, duaObj.State, duaObj.RequestID,
					duaObj.NotifSetting, duaObj.CreatedAt, duaObj.UpdatedAt, duaObj.DeletedAt).
				AddRow(duaObj.DiscussionID, duaObj.UserID, duaObj.State, duaObj.RequestID,
					duaObj.NotifSetting, duaObj.CreatedAt, duaObj.UpdatedAt, duaObj.DeletedAt)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := duaIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyDuaObj), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDiscussionUserAccessIter_Close(t *testing.T) {
	ctx := context.Background()

	Convey("DiscussionIter_Close", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := duaIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on rows.Close", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "user_id", "state", "request_id",
				"notif_setting", "created_at", "updated_at", "deleted_at"}).CloseError(fmt.Errorf("error"))

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := duaIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

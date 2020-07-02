package datastore

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/config"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

func TestDelphisDB_GetPublicDiscussions(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	discObj := model.Discussion{
		ID:            "discussion1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
		Title:         "test",
		AnonymityType: "",
		ModeratorID:   &modID,
		AutoPost:      false,
		IdleMinutes:   120,
		PublicAccess:  false,
		IconURL:       "",
	}

	emptyDisc := model.Discussion{}

	Convey("GetPublicDiscussions", t, func() {
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

			iter := mockDatastore.GetPublicDiscussions(ctx)

			So(iter.Next(&emptyDisc), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getPublicDiscussionsString).WillReturnError(fmt.Errorf("error"))

			iter := mockDatastore.GetPublicDiscussions(ctx)

			So(iter.Next(&emptyDisc), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns discussions", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"id", "created_at", "title", "anonymity_type",
				"moderator_id", "auto_post", "icon_url", "idle_minutes", "public_access"}).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.Title, discObj.AnonymityType, discObj.ModeratorID,
					discObj.AutoPost, discObj.IconURL, discObj.IdleMinutes, discObj.PublicAccess)

			mock.ExpectQuery(getPublicDiscussionsString).WillReturnRows(rs)

			iter := mockDatastore.GetPublicDiscussions(ctx)

			So(iter.Next(&emptyDisc), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetDiscussionsForFlairTemplateByUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	userID := "userID"
	discObj := model.Discussion{
		ID:            "discussion1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
		Title:         "test",
		AnonymityType: "",
		ModeratorID:   &modID,
		AutoPost:      false,
		IdleMinutes:   120,
		PublicAccess:  false,
		IconURL:       "",
	}

	emptyDisc := model.Discussion{}

	Convey("GetDiscussionsForFlairTemplateByUserID", t, func() {
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

			iter := mockDatastore.GetDiscussionsForFlairTemplateByUserID(ctx, userID)

			So(iter.Next(&emptyDisc), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getDiscussionsByFlairTemplateForUserString).WithArgs(userID).WillReturnError(fmt.Errorf("error"))

			iter := mockDatastore.GetDiscussionsForFlairTemplateByUserID(ctx, userID)

			So(iter.Next(&emptyDisc), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns discussions", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"id", "created_at", "title", "anonymity_type",
				"moderator_id", "auto_post", "icon_url", "idle_minutes", "public_access"}).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.Title, discObj.AnonymityType, discObj.ModeratorID,
					discObj.AutoPost, discObj.IconURL, discObj.IdleMinutes, discObj.PublicAccess)

			mock.ExpectQuery(getDiscussionsByFlairTemplateForUserString).WithArgs(userID).WillReturnRows(rs)

			iter := mockDatastore.GetDiscussionsForFlairTemplateByUserID(ctx, userID)

			So(iter.Next(&emptyDisc), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetDiscussionsForUserAccessByUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	userID := "userID"
	discObj := model.Discussion{
		ID:            "discussion1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
		Title:         "test",
		AnonymityType: "",
		ModeratorID:   &modID,
		AutoPost:      false,
		IdleMinutes:   120,
		PublicAccess:  false,
		IconURL:       "",
	}

	emptyDisc := model.Discussion{}

	Convey("GetDiscussionsForUserAccessByUserID", t, func() {
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

			iter := mockDatastore.GetDiscussionsForUserAccessByUserID(ctx, userID)

			So(iter.Next(&emptyDisc), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getDiscussionsByUserAccessForUserString).WithArgs(userID).WillReturnError(fmt.Errorf("error"))

			iter := mockDatastore.GetDiscussionsForUserAccessByUserID(ctx, userID)

			So(iter.Next(&emptyDisc), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns discussions", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"id", "created_at", "title", "anonymity_type",
				"moderator_id", "auto_post", "icon_url", "idle_minutes", "public_access"}).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.Title, discObj.AnonymityType, discObj.ModeratorID,
					discObj.AutoPost, discObj.IconURL, discObj.IdleMinutes, discObj.PublicAccess)

			mock.ExpectQuery(getDiscussionsByUserAccessForUserString).WithArgs(userID).WillReturnRows(rs)

			iter := mockDatastore.GetDiscussionsForUserAccessByUserID(ctx, userID)

			So(iter.Next(&emptyDisc), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetDiscussionFlairTemplatesAccessByDiscussionID(t *testing.T) {
	ctx := context.Background()
	flairTemplateID := "flairID"
	discussionID := "discussionID"
	dfaObj := model.DiscussionFlairTemplateAccess{
		DiscussionID:    discussionID,
		FlairTemplateID: flairTemplateID,
	}

	emptyDFA := model.DiscussionFlairTemplateAccess{}

	Convey("GetDiscussionFlairTemplatesAccessByDiscussionID", t, func() {
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

			iter := mockDatastore.GetDiscussionFlairTemplatesAccessByDiscussionID(ctx, discussionID)

			So(iter.Next(&emptyDFA), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getDiscussionFlairAccessString).WithArgs(discussionID).WillReturnError(fmt.Errorf("error"))

			iter := mockDatastore.GetDiscussionFlairTemplatesAccessByDiscussionID(ctx, discussionID)

			So(iter.Next(&emptyDFA), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns discussions", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"discussion_id", "flair_template_id", "created_at", "updated_at"}).
				AddRow(dfaObj.DiscussionID, dfaObj.FlairTemplateID, dfaObj.CreatedAt, dfaObj.UpdatedAt)

			mock.ExpectQuery(getDiscussionFlairAccessString).WithArgs(discussionID).WillReturnRows(rs)

			iter := mockDatastore.GetDiscussionFlairTemplatesAccessByDiscussionID(ctx, discussionID)

			So(iter.Next(&emptyDFA), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_UpsertDiscussionFlairTemplatesAccess(t *testing.T) {
	ctx := context.Background()
	flairTemplateID := "flairID"
	discussionID := "discussionID"
	dfaObj := model.DiscussionFlairTemplateAccess{
		DiscussionID:    discussionID,
		FlairTemplateID: flairTemplateID,
	}

	Convey("UpsertDiscussionFlairTemplatesAccess", t, func() {
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
			resp, err := mockDatastore.UpsertDiscussionFlairTemplatesAccess(ctx, tx, discussionID, flairTemplateID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(upsertDiscussionFlairAccessString)
			mock.ExpectQuery(upsertDiscussionFlairAccessString).WithArgs(discussionID, flairTemplateID).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.UpsertDiscussionFlairTemplatesAccess(ctx, tx, discussionID, flairTemplateID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when there is no new data to upsert", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(upsertDiscussionFlairAccessString)
			mock.ExpectQuery(upsertDiscussionFlairAccessString).WithArgs(discussionID, flairTemplateID).WillReturnError(sql.ErrNoRows)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.UpsertDiscussionFlairTemplatesAccess(ctx, tx, discussionID, flairTemplateID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &model.DiscussionFlairTemplateAccess{})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns discussions", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "flair_template_id", "created_at", "updated_at"}).
				AddRow(dfaObj.DiscussionID, dfaObj.FlairTemplateID, dfaObj.CreatedAt, dfaObj.UpdatedAt)

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(upsertDiscussionFlairAccessString)
			mock.ExpectQuery(upsertDiscussionFlairAccessString).WithArgs(discussionID, flairTemplateID).WillReturnRows(rs)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.UpsertDiscussionFlairTemplatesAccess(ctx, tx, discussionID, flairTemplateID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &dfaObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_UpsertDiscussionUserAccess(t *testing.T) {
	ctx := context.Background()
	userID := "userID"
	discussionID := "discussionID"
	duaObj := model.DiscussionUserAccess{
		DiscussionID: discussionID,
		UserID:       userID,
	}

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

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.UpsertDiscussionUserAccess(ctx, tx, discussionID, userID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(upsertDiscussionUserAccessString)
			mock.ExpectQuery(upsertDiscussionUserAccessString).WithArgs(discussionID, userID).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.UpsertDiscussionUserAccess(ctx, tx, discussionID, userID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when there is no new data to upsert", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(upsertDiscussionUserAccessString)
			mock.ExpectQuery(upsertDiscussionUserAccessString).WithArgs(discussionID, userID).WillReturnError(sql.ErrNoRows)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.UpsertDiscussionUserAccess(ctx, tx, discussionID, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &model.DiscussionUserAccess{})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns discussions", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "user_id", "created_at", "updated_at"}).
				AddRow(duaObj.DiscussionID, duaObj.UserID, duaObj.CreatedAt, duaObj.UpdatedAt)

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(upsertDiscussionUserAccessString)
			mock.ExpectQuery(upsertDiscussionUserAccessString).WithArgs(discussionID, userID).WillReturnRows(rs)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.UpsertDiscussionUserAccess(ctx, tx, discussionID, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &duaObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_DeleteDiscussionFlairTemplatesAccess(t *testing.T) {
	ctx := context.Background()
	flairTemplateID := "flairID"
	discussionID := "discussionID"
	dfaObj := model.DiscussionFlairTemplateAccess{
		DiscussionID:    discussionID,
		FlairTemplateID: flairTemplateID,
	}

	Convey("DeleteDiscussionFlairTemplatesAccess", t, func() {
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
			resp, err := mockDatastore.DeleteDiscussionFlairTemplatesAccess(ctx, tx, discussionID, flairTemplateID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(deleteDiscussionFlairAccessString)
			mock.ExpectQuery(deleteDiscussionFlairAccessString).WithArgs(discussionID, flairTemplateID).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.DeleteDiscussionFlairTemplatesAccess(ctx, tx, discussionID, flairTemplateID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns discussions", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "flair_template_id", "created_at", "updated_at", "deleted_at"}).
				AddRow(dfaObj.DiscussionID, dfaObj.FlairTemplateID, dfaObj.CreatedAt, dfaObj.UpdatedAt, dfaObj.DeletedAt)

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(deleteDiscussionFlairAccessString)
			mock.ExpectQuery(deleteDiscussionFlairAccessString).WithArgs(discussionID, flairTemplateID).WillReturnRows(rs)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.DeleteDiscussionFlairTemplatesAccess(ctx, tx, discussionID, flairTemplateID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &dfaObj)
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

func TestDiscussionIter_Next(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	discObj := model.Discussion{
		ID:            "discussion1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
		Title:         "test",
		AnonymityType: "",
		ModeratorID:   &modID,
		AutoPost:      false,
		IdleMinutes:   120,
		PublicAccess:  false,
		IconURL:       "",
	}
	emptyDisc := model.Discussion{}

	Convey("DiscussionIter_Next", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := discussionIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Next(&emptyDisc), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has a context error passed in", func() {
			ctx1, cancelFunc := context.WithCancel(ctx)
			cancelFunc()
			iter := discussionIter{
				ctx: ctx1,
			}

			So(iter.Next(&emptyDisc), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has no more rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "title", "anonymity_type",
				"moderator_id", "auto_post", "icon_url", "idle_minutes", "public_access"})

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := discussionIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyDisc), ShouldBeFalse)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on scan", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "title", "anonymity_type",
				"moderator_id", "auto_post", "icon_url", "idle_minutes"}).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.Title, discObj.AnonymityType, discObj.ModeratorID,
					discObj.AutoPost, discObj.IconURL, discObj.IdleMinutes)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := discussionIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyDisc), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "title", "anonymity_type",
				"moderator_id", "auto_post", "icon_url", "idle_minutes", "public_access"}).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.Title, discObj.AnonymityType, discObj.ModeratorID,
					discObj.AutoPost, discObj.IconURL, discObj.IdleMinutes, discObj.PublicAccess).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.Title, discObj.AnonymityType, discObj.ModeratorID,
					discObj.AutoPost, discObj.IconURL, discObj.IdleMinutes, discObj.PublicAccess)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := discussionIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyDisc), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDiscussionIter_Close(t *testing.T) {
	ctx := context.Background()

	Convey("DiscussionIter_Close", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := discussionIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on rows.Close", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "title", "anonymity_type",
				"moderator_id", "auto_post", "icon_url", "idle_minutes", "public_access"}).CloseError(fmt.Errorf("error"))

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := discussionIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDfaIter_Next(t *testing.T) {
	ctx := context.Background()
	flairTemplateID := "flairID"
	discussionID := "discussionID"
	dfaObj := model.DiscussionFlairTemplateAccess{
		DiscussionID:    discussionID,
		FlairTemplateID: flairTemplateID,
	}

	emptyDFA := model.DiscussionFlairTemplateAccess{}

	Convey("DfaIter_Next", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := dfaIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Next(&emptyDFA), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has a context error passed in", func() {
			ctx1, cancelFunc := context.WithCancel(ctx)
			cancelFunc()
			iter := dfaIter{
				ctx: ctx1,
			}

			So(iter.Next(&emptyDFA), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has no more rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "flair_template_id", "created_at", "updated_at"})

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := dfaIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyDFA), ShouldBeFalse)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on scan", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "flair_template_id", "created_at"}).
				AddRow(dfaObj.DiscussionID, dfaObj.FlairTemplateID, dfaObj.CreatedAt)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := dfaIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyDFA), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "flair_template_id", "created_at", "updated_at"}).
				AddRow(dfaObj.DiscussionID, dfaObj.FlairTemplateID, dfaObj.CreatedAt, dfaObj.UpdatedAt).
				AddRow(dfaObj.DiscussionID, dfaObj.FlairTemplateID, dfaObj.CreatedAt, dfaObj.UpdatedAt)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := dfaIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyDFA), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDfaIter_Close(t *testing.T) {
	ctx := context.Background()

	Convey("DfaIter_Close", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := dfaIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on rows.Close", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "flair_template_id", "created_at", "updated_at"}).CloseError(fmt.Errorf("error"))

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := dfaIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

package datastore

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/delphis-inc/delphisbe/internal/datastore/tests"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestDelphisDB_GetAccessLinkBySlug(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	slug := "slug"
	discussionID := tests.Discussion1ID
	linkObject := model.DiscussionAccessLink{
		DiscussionID: discussionID,
		LinkSlug:     "slug",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	Convey("GetAccessLinkBySlug", t, func() {
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

			resp, err := mockDatastore.GetAccessLinkByDiscussionID(ctx, slug)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getAccessLinkBySlugString).WithArgs(slug).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetAccessLinkBySlug(ctx, slug)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution does not find a record", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getAccessLinkBySlugString).WithArgs(slug).WillReturnError(sql.ErrNoRows)

			resp, err := mockDatastore.GetAccessLinkBySlug(ctx, slug)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &model.DiscussionAccessLink{})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns access link", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"discussion_id", "link_slug", "created_at", "updated_at", "deleted_at"}).
				AddRow(linkObject.DiscussionID, linkObject.LinkSlug, linkObject.CreatedAt,
					linkObject.UpdatedAt, linkObject.DeletedAt)

			mock.ExpectQuery(getAccessLinkBySlugString).WithArgs(slug).WillReturnRows(rs)

			resp, err := mockDatastore.GetAccessLinkBySlug(ctx, slug)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &linkObject)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetAccessLinkByDiscussionID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := tests.Discussion1ID
	linkObject := model.DiscussionAccessLink{
		DiscussionID: discussionID,
		LinkSlug:     "slug",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	Convey("GetAccessLinkByDiscussionID", t, func() {
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

			resp, err := mockDatastore.GetAccessLinkByDiscussionID(ctx, discussionID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getAccessLinkByDiscussionIDString).WithArgs(discussionID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetAccessLinkByDiscussionID(ctx, discussionID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution does not find a record", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getAccessLinkByDiscussionIDString).WithArgs(discussionID).WillReturnError(sql.ErrNoRows)

			resp, err := mockDatastore.GetAccessLinkByDiscussionID(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &model.DiscussionAccessLink{})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns access link", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"discussion_id", "link_slug", "created_at", "updated_at", "deleted_at"}).
				AddRow(linkObject.DiscussionID, linkObject.LinkSlug, linkObject.CreatedAt,
					linkObject.UpdatedAt, linkObject.DeletedAt)

			mock.ExpectQuery(getAccessLinkByDiscussionIDString).WithArgs(discussionID).WillReturnRows(rs)

			resp, err := mockDatastore.GetAccessLinkByDiscussionID(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &linkObject)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_PutAccessLinkForDiscussion(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	slug := "slug"
	linkObject := model.DiscussionAccessLink{
		DiscussionID: discussionID,
		LinkSlug:     slug,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	Convey("PutAccessLinkForDiscussion", t, func() {
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
			resp, err := mockDatastore.PutAccessLinkForDiscussion(ctx, tx, linkObject)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putAccessLinkForDiscussionString)
			mock.ExpectQuery(putAccessLinkForDiscussionString).WithArgs(
				linkObject.DiscussionID, linkObject.LinkSlug,
			).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.PutAccessLinkForDiscussion(ctx, tx, linkObject)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns imported content", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "link_slug", "created_at", "updated_at", "deleted_at"}).
				AddRow(linkObject.DiscussionID, linkObject.LinkSlug, linkObject.CreatedAt,
					linkObject.UpdatedAt, linkObject.DeletedAt)

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putAccessLinkForDiscussionString)
			mock.ExpectQuery(putAccessLinkForDiscussionString).WithArgs(
				linkObject.DiscussionID, linkObject.LinkSlug,
			).WillReturnRows(rs)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.PutAccessLinkForDiscussion(ctx, tx, linkObject)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &linkObject)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

package datastore

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/lib/pq"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/config"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestDelphisDB_GetImportedContentByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	icID := "ic1"
	icObject := model.ImportedContent{
		ID:          icID,
		CreatedAt:   now,
		ContentName: "name",
		ContentType: "type",
		Link:        "http://",
		Overview:    "overview",
		Source:      "my source",
		Tags:        nil,
	}

	Convey("GetImportedContentByID", t, func() {
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

			resp, err := mockDatastore.GetImportedContentByID(ctx, icID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getImportedContentByIDString).WithArgs(icID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetImportedContentByID(ctx, icID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns imported content", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"id", "created_at", "content_name", "content_type", "link", "overview", "source"}).
				AddRow(icObject.ID, icObject.CreatedAt, icObject.ContentName, icObject.ContentType, icObject.Link,
					icObject.Overview, icObject.Source)

			mock.ExpectQuery(getImportedContentByIDString).WithArgs(icID).WillReturnRows(rs)

			resp, err := mockDatastore.GetImportedContentByID(ctx, icID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &icObject)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetImportedContentTags(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	icID := "ic1"
	tagObject := model.Tag{
		ID:        icID,
		Tag:       "testTag",
		CreatedAt: now,
		DeletedAt: nil,
	}

	emptyTag := model.Tag{}

	Convey("GetImportedContentTags", t, func() {
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

			iter := mockDatastore.GetImportedContentTags(ctx, icID)

			So(iter.Next(&emptyTag), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getImportedContentTagsString).WithArgs(icID).WillReturnError(fmt.Errorf("error"))

			iter := mockDatastore.GetImportedContentTags(ctx, icID)

			So(iter.Next(&emptyTag), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns posts", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"imported_content_id", "tag", "created_at"}).
				AddRow(tagObject.ID, tagObject.Tag, tagObject.CreatedAt).
				AddRow(tagObject.ID, tagObject.Tag, tagObject.CreatedAt)

			mock.ExpectQuery(getImportedContentTagsString).WithArgs(icID).WillReturnRows(rs)

			iter := mockDatastore.GetImportedContentTags(ctx, icID)

			So(iter.Next(&emptyTag), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetImportedContentByDiscussionID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	icID := "ic1"
	limit := 10
	icObject := model.ImportedContent{
		ID:          icID,
		CreatedAt:   now,
		ContentName: "name",
		ContentType: "type",
		Link:        "http://",
		Overview:    "overview",
		Source:      "my source",
		Tags:        nil,
	}

	emptyIC := model.ImportedContent{}

	Convey("GetImportedContentByDiscussionID", t, func() {
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

			iter := mockDatastore.GetImportedContentByDiscussionID(ctx, icID, limit)

			So(iter.Next(&emptyIC), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getImportedContentForDiscussionString).WithArgs(icID, limit).WillReturnError(fmt.Errorf("error"))

			iter := mockDatastore.GetImportedContentByDiscussionID(ctx, icID, limit)

			So(iter.Next(&emptyIC), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns posts", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"i.id", "i.created_at", "i.content_name", "i.content_type", "i.link", "i.overview",
				"i.source", "d.matching_tags"}).
				AddRow(icObject.ID, icObject.CreatedAt, icObject.ContentName, icObject.ContentType, icObject.Link,
					icObject.Overview, icObject.Source, pq.Array(icObject.Tags)).
				AddRow(icObject.ID, icObject.CreatedAt, icObject.ContentName, icObject.ContentType, icObject.Link,
					icObject.Overview, icObject.Source, pq.Array(icObject.Tags))

			mock.ExpectQuery(getImportedContentForDiscussionString).WithArgs(icID, limit).WillReturnRows(rs)

			iter := mockDatastore.GetImportedContentByDiscussionID(ctx, icID, limit)

			So(iter.Next(&emptyIC), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetScheduledImportedContentByDiscussionID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	icID := "ic1"
	icObject := model.ImportedContent{
		ID:          icID,
		CreatedAt:   now,
		ContentName: "name",
		ContentType: "type",
		Link:        "http://",
		Overview:    "overview",
		Source:      "my source",
		Tags:        nil,
	}

	emptyIC := model.ImportedContent{}

	Convey("GetScheduledImportedContentByDiscussionID", t, func() {
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

			iter := mockDatastore.GetScheduledImportedContentByDiscussionID(ctx, discussionID)

			So(iter.Next(&emptyIC), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getScheduledImportedContentByDiscussionIDString).WithArgs(discussionID).WillReturnError(fmt.Errorf("error"))

			iter := mockDatastore.GetScheduledImportedContentByDiscussionID(ctx, discussionID)

			So(iter.Next(&emptyIC), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns posts", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"i.id", "i.created_at", "i.content_name", "i.content_type", "i.link", "i.overview",
				"i.source", "q.matching_tags"}).
				AddRow(icObject.ID, icObject.CreatedAt, icObject.ContentName, icObject.ContentType, icObject.Link,
					icObject.Overview, icObject.Source, pq.Array(icObject.Tags)).
				AddRow(icObject.ID, icObject.CreatedAt, icObject.ContentName, icObject.ContentType, icObject.Link,
					icObject.Overview, icObject.Source, pq.Array(icObject.Tags))

			mock.ExpectQuery(getScheduledImportedContentByDiscussionIDString).WithArgs(discussionID).WillReturnRows(rs)

			iter := mockDatastore.GetScheduledImportedContentByDiscussionID(ctx, discussionID)

			So(iter.Next(&emptyIC), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_PutImportedContent(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	icID := "ic1"
	icObject := model.ImportedContent{
		ID:          icID,
		CreatedAt:   now,
		ContentName: "name",
		ContentType: "type",
		Link:        "http://",
		Overview:    "overview",
		Source:      "my source",
		Tags:        nil,
	}

	Convey("PutImportedContent", t, func() {
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
			resp, err := mockDatastore.PutImportedContent(ctx, tx, icObject)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putImportedContentString)
			mock.ExpectQuery(putImportedContentString).WithArgs(icObject.ID, icObject.ContentName, icObject.ContentType, icObject.Link,
				icObject.Overview, icObject.Source).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.PutImportedContent(ctx, tx, icObject)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when put post succeeds and returns an object", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "content_name", "content_type", "link", "overview", "source"}).
				AddRow(icObject.ID, icObject.CreatedAt, icObject.ContentName, icObject.ContentType, icObject.Link, icObject.Overview, icObject.Source)

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putImportedContentString)
			mock.ExpectQuery(putImportedContentString).WithArgs(icObject.ID, icObject.ContentName, icObject.ContentType, icObject.Link,
				icObject.Overview, icObject.Source).WillReturnRows(rs)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.PutImportedContent(ctx, tx, icObject)

			logrus.Infof("Resp: %+v\n", resp)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &icObject)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_PutImportedContentTags(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	icID := "ic1"
	tagObject := model.Tag{
		ID:        icID,
		Tag:       "testTag",
		CreatedAt: now,
		DeletedAt: nil,
	}

	Convey("PutImportedContentTags", t, func() {
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
			resp, err := mockDatastore.PutImportedContentTags(ctx, tx, tagObject)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putImportedContentTagsString)
			mock.ExpectQuery(putImportedContentTagsString).WithArgs(tagObject.ID, tagObject.Tag).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.PutImportedContentTags(ctx, tx, tagObject)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when put post succeeds and returns an object", func() {
			rs := sqlmock.NewRows([]string{"id", "tag", "created_at"}).
				AddRow(tagObject.ID, tagObject.Tag, tagObject.CreatedAt)

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putImportedContentTagsString)
			mock.ExpectQuery(putImportedContentTagsString).WithArgs(tagObject.ID, tagObject.Tag).WillReturnRows(rs)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.PutImportedContentTags(ctx, tx, tagObject)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &tagObject)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetMatchingTags(t *testing.T) {
	ctx := context.Background()
	discussionID := "discussion1"
	icID := "ic1"
	tags := []string{"test", "test1"}

	Convey("GetMatchingTags", t, func() {
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

			resp, err := mockDatastore.GetMatchingTags(ctx, discussionID, icID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getMatchingTagsString).WithArgs(discussionID, icID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetMatchingTags(ctx, discussionID, icID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution doesn't find any matching tags", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getMatchingTagsString).WithArgs(discussionID, icID).WillReturnError(sql.ErrNoRows)

			resp, err := mockDatastore.GetMatchingTags(ctx, discussionID, icID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns posts", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"matching_tags"}).
				AddRow(pq.Array(tags))

			mock.ExpectQuery(getMatchingTagsString).WithArgs(discussionID, icID).WillReturnRows(rs)

			resp, err := mockDatastore.GetMatchingTags(ctx, discussionID, icID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, tags)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_PutImportedContentDiscussionQueue(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	icID := "ic1"
	tags := []string{"test", "test1"}
	contentQueueObj := model.ContentQueueRecord{
		DiscussionID:      discussionID,
		ImportedContentID: icID,
		PostedAt:          &now,
		MatchingTags:      tags,
	}

	Convey("PutImportedContentDiscussionQueue", t, func() {
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

			resp, err := mockDatastore.PutImportedContentDiscussionQueue(ctx, discussionID, icID, &now, tags)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(putImportedContentDiscussionQueueString).WithArgs(contentQueueObj.DiscussionID, contentQueueObj.ImportedContentID,
				contentQueueObj.PostedAt, pq.Array(contentQueueObj.MatchingTags)).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.PutImportedContentDiscussionQueue(ctx, discussionID, icID, &now, tags)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when put post succeeds and returns an object", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "imported_content_id", "created_at", "updated_at", "deleted_at", "posted_at", "matching_tags"}).
				AddRow(contentQueueObj.DiscussionID, contentQueueObj.ImportedContentID, contentQueueObj.CreatedAt, contentQueueObj.UpdatedAt,
					contentQueueObj.DeletedAt, contentQueueObj.PostedAt, pq.Array(contentQueueObj.MatchingTags))

			mockPreparedStatements(mock)
			mock.ExpectQuery(putImportedContentDiscussionQueueString).WithArgs(contentQueueObj.DiscussionID, contentQueueObj.ImportedContentID,
				contentQueueObj.PostedAt, pq.Array(contentQueueObj.MatchingTags)).WillReturnRows(rs)

			resp, err := mockDatastore.PutImportedContentDiscussionQueue(ctx, discussionID, icID, &now, tags)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &contentQueueObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_UpdateImportedContentDiscussionQueue(t *testing.T) {
	ctx := context.Background()
	postedAt := time.Now()
	discussionID := "discussion1"
	icID := "ic1"
	contentQueueObj := model.ContentQueueRecord{
		DiscussionID:      discussionID,
		ImportedContentID: icID,
		PostedAt:          &postedAt,
	}

	Convey("UpdateImportedContentDiscussionQueue", t, func() {
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

			resp, err := mockDatastore.UpdateImportedContentDiscussionQueue(ctx, discussionID, icID, &postedAt)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(updateImportedContentDiscussionQueueString).WithArgs(contentQueueObj.DiscussionID, contentQueueObj.ImportedContentID,
				&contentQueueObj.PostedAt).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.UpdateImportedContentDiscussionQueue(ctx, discussionID, icID, &postedAt)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when put post succeeds and returns an object", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "imported_content_id", "created_at", "updated_at", "deleted_at", "posted_at", "matching_tags"}).
				AddRow(contentQueueObj.DiscussionID, contentQueueObj.ImportedContentID, contentQueueObj.CreatedAt, contentQueueObj.UpdatedAt,
					contentQueueObj.DeletedAt, contentQueueObj.PostedAt, pq.Array(contentQueueObj.MatchingTags))

			mockPreparedStatements(mock)
			mock.ExpectQuery(updateImportedContentDiscussionQueueString).WithArgs(contentQueueObj.DiscussionID, contentQueueObj.ImportedContentID,
				&contentQueueObj.PostedAt).WillReturnRows(rs)

			resp, err := mockDatastore.UpdateImportedContentDiscussionQueue(ctx, discussionID, icID, &postedAt)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &contentQueueObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestTagIter_Next(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	icID := "ic1"
	tagObject := model.Tag{
		ID:        icID,
		Tag:       "testTag",
		CreatedAt: now,
		DeletedAt: nil,
	}

	emptyTag := model.Tag{}

	Convey("TagIter Next", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := tagIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Next(&emptyTag), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has a context error passed in", func() {
			ctx1, cancelFunc := context.WithCancel(ctx)
			cancelFunc()
			iter := tagIter{
				ctx: ctx1,
			}

			So(iter.Next(&emptyTag), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has no more rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"imported_content_id", "tag", "created_at"})

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := tagIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyTag), ShouldBeFalse)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on scan", func() {
			rs := sqlmock.NewRows([]string{"imported_content_id", "tag"}).
				AddRow(tagObject.ID, tagObject.Tag)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := tagIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyTag), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"imported_content_id", "tag", "created_at"}).
				AddRow(tagObject.ID, tagObject.Tag, tagObject.CreatedAt).
				AddRow(tagObject.ID, tagObject.Tag, tagObject.CreatedAt)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := tagIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyTag), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestTagIter_Close(t *testing.T) {
	ctx := context.Background()

	Convey("TagIter Close", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := tagIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on rows.Close", func() {
			rs := sqlmock.NewRows([]string{"imported_content_id", "tag", "created_at"}).CloseError(fmt.Errorf("error"))

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := tagIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestContentIter_Next(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	icID := "ic1"
	icObject := model.ImportedContent{
		ID:          icID,
		CreatedAt:   now,
		ContentName: "name",
		ContentType: "type",
		Link:        "http://",
		Overview:    "overview",
		Source:      "my source",
		Tags:        nil,
	}

	emptyIC := model.ImportedContent{}

	Convey("ContentIter Next", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := contentIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Next(&emptyIC), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has a context error passed in", func() {
			ctx1, cancelFunc := context.WithCancel(ctx)
			cancelFunc()
			iter := contentIter{
				ctx: ctx1,
			}

			So(iter.Next(&emptyIC), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has no more rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"i.id", "i.created_at", "i.content_name", "i.content_type", "i.link", "i.overview",
				"i.source", "d.matching_tags"})

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := contentIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyIC), ShouldBeFalse)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on scan", func() {
			rs := sqlmock.NewRows([]string{"i.id", "i.created_at", "i.content_name", "i.content_type", "i.link", "i.overview",
				"i.source"}).
				AddRow(icObject.ID, icObject.CreatedAt, icObject.ContentName, icObject.ContentType, icObject.Link,
					icObject.Overview, icObject.Source)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := contentIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyIC), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"i.id", "i.created_at", "i.content_name", "i.content_type", "i.link", "i.overview",
				"i.source", "d.matching_tags"}).
				AddRow(icObject.ID, icObject.CreatedAt, icObject.ContentName, icObject.ContentType, icObject.Link,
					icObject.Overview, icObject.Source, pq.Array(icObject.Tags)).
				AddRow(icObject.ID, icObject.CreatedAt, icObject.ContentName, icObject.ContentType, icObject.Link,
					icObject.Overview, icObject.Source, pq.Array(icObject.Tags))

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := contentIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyIC), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestContentIter_Close(t *testing.T) {
	ctx := context.Background()

	Convey("ContentIter Close", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := contentIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on rows.Close", func() {
			rs := sqlmock.NewRows([]string{"i.id", "i.created_at", "i.content_name", "i.content_type", "i.link", "i.overview",
				"i.source", "d.matching_tags"}).CloseError(fmt.Errorf("error"))

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := contentIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_ContentIterCollect(t *testing.T) {
	ctx := context.Background()
	discussionID := "discussion1"

	icObj := model.ImportedContent{
		ID:          discussionID,
		ContentName: "name",
		ContentType: "type",
	}

	Convey("ContentIterCollect", t, func() {
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
			iter := &contentIter{
				err: fmt.Errorf("error"),
			}

			resp, err := mockDatastore.ContentIterCollect(ctx, iter)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has results and returns slice of ImportedContent", func() {
			rs := sqlmock.NewRows([]string{"i.id", "i.created_at", "i.content_name", "i.content_type", "i.link", "i.overview",
				"i.source", "q.matching_tags"}).
				AddRow(icObj.ID, icObj.CreatedAt, icObj.ContentName, icObj.ContentType, icObj.Link,
					icObj.Overview, icObj.Source, pq.Array(icObj.Tags)).
				AddRow(icObj.ID, icObj.CreatedAt, icObj.ContentName, icObj.ContentType, icObj.Link,
					icObj.Overview, icObj.Source, pq.Array(icObj.Tags))

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := &contentIter{
				ctx:  ctx,
				rows: rs1,
			}

			resp, err := mockDatastore.ContentIterCollect(ctx, iter)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []*model.ImportedContent{&icObj, &icObj})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

	})
}

func TestDelphisDB_TagIterCollect(t *testing.T) {
	ctx := context.Background()
	discussionID := "discussion1"

	tagObj := model.Tag{
		ID:  discussionID,
		Tag: "tag",
	}

	Convey("TagIterCollect", t, func() {
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
			iter := &tagIter{
				err: fmt.Errorf("error"),
			}

			resp, err := mockDatastore.TagIterCollect(ctx, iter)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has results and returns slice of Tags", func() {
			rs := sqlmock.NewRows([]string{"discussion_id", "tag", "created_at"}).
				AddRow(tagObj.ID, tagObj.Tag, tagObj.CreatedAt).
				AddRow(tagObj.ID, tagObj.Tag, tagObj.CreatedAt)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := &tagIter{
				ctx:  ctx,
				rows: rs1,
			}

			resp, err := mockDatastore.TagIterCollect(ctx, iter)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []*model.Tag{&tagObj, &tagObj})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

	})
}

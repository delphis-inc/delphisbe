package datastore

import (
	"context"
	"fmt"
	"testing"

	"github.com/lib/pq"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/config"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestDelphisDB_PutPostContent(t *testing.T) {
	ctx := context.Background()
	postContentObject := model.PostContent{
		ID:                "content1",
		Content:           "hello world",
		MentionedEntities: nil,
	}

	Convey("PutPostContent", t, func() {
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
			err = mockDatastore.PutPostContent(ctx, tx, postContentObject)

			So(err, ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putPostContentsString)
			mock.ExpectExec(putPostContentsString).WithArgs(postContentObject.ID, postContentObject.Content, pq.Array(postContentObject.MentionedEntities)).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			err = mockDatastore.PutPostContent(ctx, tx, postContentObject)

			So(err, ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when put post succeeds and returns an object", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putPostContentsString)
			mock.ExpectExec(putPostContentsString).WithArgs(postContentObject.ID, postContentObject.Content, pq.Array(postContentObject.MentionedEntities)).WillReturnResult(sqlmock.NewResult(0, 0))

			tx, err := mockDatastore.BeginTx(ctx)
			err = mockDatastore.PutPostContent(ctx, tx, postContentObject)

			So(err, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetPostContentByID(t *testing.T) {
	ctx := context.Background()
	postContentID := "content1"
	postContentObject := model.PostContent{
		ID:                postContentID,
		Content:           "hello world",
		MentionedEntities: nil,
	}

	Convey("GetPostByID", t, func() {
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

		expectedQueryString := `SELECT * FROM "post_contents"  WHERE ("post_contents"."id" = $1) ORDER BY "post_contents"."id" ASC LIMIT 1`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(postContentID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetPostContentByID(ctx, postContentID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a postContents", func() {
			rs := sqlmock.NewRows([]string{"id", "content", "created_at", "updated_at", "mentioned_entities"}).
				AddRow(postContentObject.ID, postContentObject.Content, postContentObject.CreatedAt, postContentObject.UpdatedAt, pq.Array(postContentObject.MentionedEntities))

			mock.ExpectQuery(expectedQueryString).WithArgs(postContentID).WillReturnRows(rs)

			resp, err := mockDatastore.GetPostContentByID(ctx, postContentID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &postContentObject)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

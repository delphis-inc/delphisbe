package datastore

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/config"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestDelphisDB_PutActivity(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	participantID := "participant1"
	postID := "post1"
	postObject := model.Post{
		ID:            "post1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		QuotedPostID:  &postID,
		PostType:      model.PostTypeStandard,
	}

	postContentObj := model.PostContent{
		ID:                postID,
		Content:           "<0> <1>",
		MentionedEntities: []string{"discussion:1234", "participant:1234"},
	}

	entity0 := strings.Split(postContentObj.MentionedEntities[0], ":")
	entity1 := strings.Split(postContentObj.MentionedEntities[1], ":")

	postObject.PostContent = &postContentObj

	Convey("PutActivity", t, func() {
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
			err = mockDatastore.PutActivity(ctx, tx, &postObject)

			So(err, ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putActivityString)
			mock.ExpectExec(putActivityString).WithArgs(postObject.ParticipantID, postObject.PostContent.ID,
				entity0[1], entity0[0]).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			err = mockDatastore.PutActivity(ctx, tx, &postObject)

			So(err, ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when put post succeeds and returns an object", func() {

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putActivityString)
			mock.ExpectExec(putActivityString).WithArgs(postObject.ParticipantID, postObject.PostContent.ID,
				entity0[1], entity0[0]).WillReturnResult(sqlmock.NewResult(0, 0))
			//mock.ExpectPrepare(putActivityString)
			mock.ExpectExec(putActivityString).WithArgs(postObject.ParticipantID, postObject.PostContent.ID,
				entity1[1], entity1[0]).WillReturnResult(sqlmock.NewResult(0, 0))

			tx, err := mockDatastore.BeginTx(ctx)
			err = mockDatastore.PutActivity(ctx, tx, &postObject)

			So(err, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

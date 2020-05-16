package datastore

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/config"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestGetPostsByDiscussionID(t *testing.T) {
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
		PostContentID: &postID,
	}

	Convey("GetPostsByDiscussionID", t, func() {
		db, mock, err := sqlmock.New()

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

		Convey("when the objects are found", func() {
			if err := mockDatastore.initializeStatements(ctx); err != nil {
				logrus.WithError(err).Error(":faield")
			}

			mock.ExpectPrepare(regexp.QuoteMeta(putPostString))

			mock.ExpectExec("INSERT into posts").WithArgs(2, 3, 4, 5, 6, 7, 8).WillReturnResult(sqlmock.NewResult(1, 1))
			//mock.ExpectPrepare(string(mockDatastore.putPostStmt))
			//mock.ExpectExec(regexp.QuoteMeta(mockDatastore.putPostStmt)).WithArgs(discussionID).WillReturnRows(
			//	sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "deleted_reason_code", "discussion_id", "participant_id", "post_content_id"}).
			//		AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DeletedAt, postObject.DeletedReasonCode, postObject.DiscussionID, postObject.ParticipantID, postObject.PostContentID))

			resp, err := mockDatastore.PutPost(ctx, postObject)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []model.Post{postObject})
			//So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

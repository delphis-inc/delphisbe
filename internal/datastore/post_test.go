package datastore

import (
	"context"
	"regexp"
	"testing"
	"time"

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
		QuotedPostID:  &postID,
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

		// Prepare statements
		mock.ExpectPrepare(regexp.QuoteMeta(putPostString))
		mock.ExpectPrepare(regexp.QuoteMeta(getPostsByDiscussionIDString))
		mock.ExpectPrepare(regexp.QuoteMeta(putPostContentsString))

		Convey("when the works", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "discussion_id", "participant_id", "post_content_id", "quoted_post_id"}).
				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DiscussionID, postObject.ParticipantID, postObject.PostContentID, postObject.QuotedPostID)

			mock.ExpectQuery(regexp.QuoteMeta(putPostString)).WithArgs(postObject.ID, postObject.DiscussionID, postObject.ParticipantID, postObject.PostContentID, postObject.QuotedPostID).WillReturnRows(rs)

			resp, err := mockDatastore.PutPost(ctx, postObject)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &postObject)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

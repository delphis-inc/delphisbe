package datastore

import (
	"testing"
	//. "github.com/smartystreets/goconvey/convey"
)

func TestGetPostsByDiscussionID(t *testing.T) {
	//ctx := context.Background()
	//now := time.Now()
	//discussionID := "discussion1"
	//participantID := "participant1"
	//postID := "post1"
	//postObject := model.Post{
	//	ID:            "post1",
	//	CreatedAt:     now,
	//	UpdatedAt:     now,
	//	DiscussionID:  &discussionID,
	//	ParticipantID: &participantID,
	//	PostContentID: &postID,
	//}
	//
	//Convey("GetPostsByDiscussionID", t, func() {
	//	db, mock, err := sqlmock.New()
	//
	//	assert.Nil(t, err, "Failed setting up sqlmock db")
	//
	//	gormDB, _ := gorm.Open("postgres", db)
	//	mockDatastore := &delphisDB{
	//		dbConfig: config.TablesConfig{},
	//		sql:      gormDB,
	//		dynamo:   nil,
	//		encoder:  nil,
	//	}
	//	defer db.Close()
	//
	//	expectedQueryStr := `SELECT * FROM "posts" WHERE "posts"."deleted_at" IS NULL AND (("posts"."discussion_id" = $1))`
	//
	//	Convey("when the objects are found", func() {
	//		mock.ExpectQuery(regexp.QuoteMeta(expectedQueryStr)).WithArgs(discussionID).WillReturnRows(
	//			sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "deleted_reason_code", "discussion_id", "participant_id", "post_content_id"}).
	//				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DeletedAt, postObject.DeletedReasonCode, postObject.DiscussionID, postObject.ParticipantID, postObject.PostContentID))
	//
	//		resp, err := mockDatastore.GetPostsByDiscussionID(ctx, discussionID)
	//
	//		So(err, ShouldBeNil)
	//		So(resp, ShouldNotBeNil)
	//		So(resp, ShouldResemble, []model.Post{postObject})
	//		//So(mock.ExpectationsWereMet(), ShouldBeNil)
	//	})
	//})
}

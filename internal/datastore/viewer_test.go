package datastore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/config"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestDelphisDB_UpsertViewer(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := "userID"
	discussionID := "discussionID"
	viewerID := "viewerID"
	postID := "post1"

	viewerObj := model.Viewer{
		ID:               viewerID,
		CreatedAt:        now,
		UpdatedAt:        now,
		DiscussionID:     &discussionID,
		LastViewedPostID: &postID,
		UserID:           &userID,
	}

	Convey("UpsertViewer", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		assert.Nil(t, err, "Failed setting up sqlmock db")

		gormDB, _ := gorm.Open("postgres", db)
		mockDatastore := &delphisDB{
			dbConfig: config.TablesConfig{},
			sql:      gormDB,
			dynamo:   nil,
			encoder:  nil,
		}
		defer db.Close()

		expectedFindQueryStr := `SELECT * FROM "viewers" WHERE "viewers"."deleted_at" IS NULL AND (("viewers"."id" = $1)) ORDER BY "viewers"."id" ASC LIMIT 1`
		createQueryStr := `INSERT INTO "viewers" ("id","created_at","updated_at","deleted_at","discussion_id","last_viewed","last_viewed_post_id","user_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "viewers"."id"`
		expectedNewObjectRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "discussion_id", "last_viewed", "last_viewed_post_id", "user_id"}).
			AddRow(viewerObj.ID, viewerObj.CreatedAt, viewerObj.UpdatedAt, viewerObj.DeletedAt, viewerObj.DiscussionID, viewerObj.LastViewed, viewerObj.LastViewedPostID, viewerObj.UserID)
		expectedUpdateStr := `UPDATE "viewers" SET "last_viewed_post_id" = $1, "updated_at" = $2 WHERE "viewers"."deleted_at" IS NULL AND "viewers"."id" = $3`
		expectedPostUpdateSelectStr := `SELECT * FROM "viewers" WHERE "viewers"."deleted_at" IS NULL AND "viewers"."id" = $1 ORDER BY "viewers"."id" ASC LIMIT 1`

		Convey("when find query fails with a non-not-found-error the function should return the error", func() {
			expectedError := fmt.Errorf("Some fake error")
			mock.ExpectQuery(expectedFindQueryStr).WithArgs(viewerObj.ID).WillReturnError(expectedError)

			resp, err := mockDatastore.UpsertViewer(ctx, viewerObj)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when find query returns not-found-error it should call create", func() {
			Convey("when create returns an error it should return it", func() {
				expectedError := fmt.Errorf("Some fake error")

				mock.ExpectQuery(expectedFindQueryStr).WithArgs(viewerObj.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					viewerObj.ID, viewerObj.CreatedAt, viewerObj.UpdatedAt, viewerObj.DeletedAt,
					viewerObj.DiscussionID, viewerObj.LastViewed, viewerObj.LastViewedPostID, viewerObj.UserID,
				).WillReturnError(expectedError)

				resp, err := mockDatastore.UpsertViewer(ctx, viewerObj)

				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, expectedError)
				So(resp, ShouldBeNil)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})

			Convey("when create succeeds it should return the new object", func() {
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(viewerObj.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					viewerObj.ID, viewerObj.CreatedAt, viewerObj.UpdatedAt, viewerObj.DeletedAt,
					viewerObj.DiscussionID, viewerObj.LastViewed, viewerObj.LastViewedPostID, viewerObj.UserID,
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(viewerObj.ID))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(viewerObj.ID).WillReturnRows(expectedNewObjectRow)

				resp, err := mockDatastore.UpsertViewer(ctx, viewerObj)

				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(resp, ShouldResemble, &viewerObj)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})
		})

		Convey("when the object is found we should update it", func() {
			Convey("when updating if it fails then return the error", func() {
				expectedError := fmt.Errorf("something went wrong")
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(viewerObj.ID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(expectedUpdateStr).WithArgs(
					viewerObj.LastViewedPostID, sqlmock.AnyArg(), viewerObj.ID,
				).WillReturnError(expectedError)

				resp, err := mockDatastore.UpsertViewer(ctx, viewerObj)

				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, expectedError)
				So(resp, ShouldBeNil)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})

			Convey("when updating if it succeeds it should return the updated object", func() {

				mock.ExpectQuery(expectedFindQueryStr).WithArgs(viewerObj.ID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(expectedUpdateStr).WithArgs(
					viewerObj.LastViewedPostID, sqlmock.AnyArg(), viewerObj.ID,
				).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedPostUpdateSelectStr).WithArgs(viewerObj.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "discussion_id", "last_viewed", "last_viewed_post_id", "user_id"}).
						AddRow(viewerObj.ID, viewerObj.CreatedAt, viewerObj.UpdatedAt, viewerObj.DeletedAt, viewerObj.DiscussionID,
							viewerObj.LastViewed, viewerObj.LastViewedPostID, viewerObj.UserID))

				resp, err := mockDatastore.UpsertViewer(ctx, viewerObj)

				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})
		})
	})
}

func TestDelphisDB_GetViewerByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := "userID"
	discussionID := "discussionID"
	viewerID := "viewerID"
	postID := "post1"

	viewerObj := model.Viewer{
		ID:               viewerID,
		CreatedAt:        now,
		UpdatedAt:        now,
		DiscussionID:     &discussionID,
		LastViewedPostID: &postID,
		UserID:           &userID,
	}

	Convey("GetViewerByID", t, func() {
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

		expectedQueryString := `SELECT * FROM "viewers" WHERE "viewers"."deleted_at" IS NULL AND (("viewers"."id" = $1)) ORDER BY "viewers"."id" ASC LIMIT 1`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(viewerObj.ID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetViewerByID(ctx, viewerObj.ID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and does not find a record", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(viewerObj.ID).WillReturnError(gorm.ErrRecordNotFound)

			resp, err := mockDatastore.GetViewerByID(ctx, viewerObj.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "last_viewed", "last_viewed_post_id",
				"discussion_id", "user_id"}).AddRow(viewerObj.ID, viewerObj.CreatedAt, viewerObj.UpdatedAt, viewerObj.DeletedAt,
				viewerObj.LastViewed, viewerObj.LastViewedPostID, viewerObj.DiscussionID, viewerObj.UserID)

			mock.ExpectQuery(expectedQueryString).WithArgs(viewerObj.ID).WillReturnRows(rs)

			resp, err := mockDatastore.GetViewerByID(ctx, viewerObj.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &viewerObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetViewersByIDs(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := "userID"
	discussionID := "discussionID"
	viewerID := "viewerID"
	postID := "post1"

	viewerObj := model.Viewer{
		ID:               viewerID,
		CreatedAt:        now,
		UpdatedAt:        now,
		DiscussionID:     &discussionID,
		LastViewedPostID: &postID,
		UserID:           &userID,
	}

	viewerIDs := []string{viewerObj.ID, viewerObj.ID}

	Convey("GetViewersByIDs", t, func() {
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

		expectedQueryString := `SELECT * FROM "viewers" WHERE "viewers"."deleted_at" IS NULL AND (("viewers"."id" IN ($1,$2)))`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(viewerIDs[0], viewerIDs[1]).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetViewersByIDs(ctx, viewerIDs)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and does not find a record", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(viewerIDs[0], viewerIDs[1]).WillReturnError(gorm.ErrRecordNotFound)

			resp, err := mockDatastore.GetViewersByIDs(ctx, viewerIDs)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "last_viewed", "last_viewed_post_id",
				"discussion_id", "user_id"}).AddRow(viewerObj.ID, viewerObj.CreatedAt, viewerObj.UpdatedAt, viewerObj.DeletedAt,
				viewerObj.LastViewed, viewerObj.LastViewedPostID, viewerObj.DiscussionID, viewerObj.UserID).
				AddRow(viewerObj.ID, viewerObj.CreatedAt, viewerObj.UpdatedAt, viewerObj.DeletedAt,
					viewerObj.LastViewed, viewerObj.LastViewedPostID, viewerObj.DiscussionID, viewerObj.UserID)

			mock.ExpectQuery(expectedQueryString).WithArgs(viewerIDs[0], viewerIDs[1]).WillReturnRows(rs)

			resp, err := mockDatastore.GetViewersByIDs(ctx, viewerIDs)

			verifyMap := map[string]*model.Viewer{
				viewerIDs[0]: &viewerObj,
				viewerIDs[1]: &viewerObj,
			}

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, verifyMap)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func Test_MarshalViewer(t *testing.T) {
	// type args struct {
	// 	viewer model.Viewer
	// }

	// lastViewed := time.Now()

	// haveViewerObj := model.Viewer{
	// 	ID:        "11111",
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// 	NotificationPreferences: model.ViewerNotificationPreferences{
	// 		ID: "55555",
	// 	},
	// 	DiscussionID:     "22222",
	// 	Discussion:       &model.Discussion{},
	// 	LastViewed:       &lastViewed,
	// 	LastViewedPostID: aws.String("33333"),
	// 	LastViewedPost:   &model.Post{},
	// 	Bookmarks:        &model.PostsConnection{},
	// 	UserID:           "44444",
	// 	User:             &model.User{},
	// }

	// datastoreObj := NewDatastore(config.DBConfig{})

	// tests := []struct {
	// 	name string
	// 	args args
	// 	want map[string]*dynamodb.AttributeValue
	// }{
	// 	{
	// 		name: "fully filled object",
	// 		args: args{
	// 			viewer: haveViewerObj,
	// 		},
	// 		want: map[string]*dynamodb.AttributeValue{
	// 			"ViewerID": {
	// 				S: aws.String(haveViewerObj.ID),
	// 			},
	// 			"CreatedAt": {
	// 				S: aws.String(haveViewerObj.CreatedAt.Format(time.RFC3339Nano)),
	// 			},
	// 			"UpdatedAt": {
	// 				S: aws.String(haveViewerObj.UpdatedAt.Format(time.RFC3339Nano)),
	// 			},
	// 			"DeletedAt": {
	// 				NULL: aws.Bool(true),
	// 			},
	// 			"NotificationPreferences": {
	// 				M: map[string]*dynamodb.AttributeValue{
	// 					"ID": {
	// 						S: aws.String(haveViewerObj.NotificationPreferences.ID),
	// 					},
	// 				},
	// 			},
	// 			"DiscussionID": {
	// 				S: aws.String(haveViewerObj.DiscussionID),
	// 			},
	// 			"LastViewed": {
	// 				S: aws.String(haveViewerObj.LastViewed.Format(time.RFC3339Nano)),
	// 			},
	// 			"LastViewedPostID": {
	// 				S: haveViewerObj.LastViewedPostID,
	// 			},
	// 			"UserID": {
	// 				S: aws.String(haveViewerObj.UserID),
	// 			},
	// 		},
	// 	},
	// }

	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		marshaled, err := datastoreObj.marshalMap(tt.args.viewer)
	// 		if err != nil {
	// 			t.Errorf("Caught an error marshaling: %+v", err)
	// 			return
	// 		}
	// 		if !reflect.DeepEqual(marshaled, tt.want) {
	// 			t.Errorf("These objects did not match. Got: %+v\n\n Want: %+v", marshaled, tt.want)
	// 		}
	// 	})
	// }
}

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

func TestDelphisDB_UpsertUser(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := "userID"
	userObj := model.User{
		ID:        userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	Convey("UpsertUser", t, func() {
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

		expectedFindQueryStr := `SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND (("users"."id" = $1)) ORDER BY "users"."id" ASC LIMIT 1`
		createQueryStr := `INSERT INTO "users" ("id","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4) RETURNING "users"."id"`
		expectedNewObjectRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at"}).
			AddRow(userObj.ID, userObj.CreatedAt, userObj.UpdatedAt, userObj.DeletedAt)

		Convey("when find query fails with a non-not-found-error the function should return the error", func() {
			expectedError := fmt.Errorf("Some fake error")
			mock.ExpectQuery(expectedFindQueryStr).WithArgs(userObj.ID).WillReturnError(expectedError)

			resp, err := mockDatastore.UpsertUser(ctx, userObj)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when find query returns not-found-error it should call create", func() {
			Convey("when create returns an error it should return it", func() {
				expectedError := fmt.Errorf("Some fake error")

				mock.ExpectQuery(expectedFindQueryStr).WithArgs(userObj.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					userObj.ID, userObj.CreatedAt, userObj.UpdatedAt, userObj.DeletedAt,
				).WillReturnError(expectedError)

				resp, err := mockDatastore.UpsertUser(ctx, userObj)

				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, expectedError)
				So(resp, ShouldBeNil)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})

			Convey("when create succeeds it should return the new object", func() {
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(userObj.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					userObj.ID, userObj.CreatedAt, userObj.UpdatedAt, userObj.DeletedAt,
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userObj.ID))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(userObj.ID).WillReturnRows(expectedNewObjectRow)

				resp, err := mockDatastore.UpsertUser(ctx, userObj)

				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(resp, ShouldResemble, &userObj)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})
		})

		Convey("when the object is found we should update it", func() {
			Convey("when updating if it fails then return the error", func() {

				mock.ExpectQuery(expectedFindQueryStr).WithArgs(userObj.ID).WillReturnRows(expectedNewObjectRow)

				resp, err := mockDatastore.UpsertUser(ctx, userObj)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})
		})
	})
}

func TestDelphisDB_GetUserByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := "userID"
	userObj := model.User{
		ID:        userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	parID := "parID"
	discussionID := "discussionID"
	viewerID := "viewerID"
	flairID := "flairID"
	gradientColor := model.GradientColorAzalea
	postID := "post1"

	parObj := model.Participant{
		ID:            parID,
		ParticipantID: 0,
		DiscussionID:  &discussionID,
		ViewerID:      &viewerID,
		FlairID:       &flairID,
		GradientColor: &gradientColor,
		UserID:        &userID,
		HasJoined:     true,
		IsAnonymous:   false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	viewerObj := model.Viewer{
		ID:               viewerID,
		CreatedAt:        now,
		UpdatedAt:        now,
		DiscussionID:     &discussionID,
		LastViewedPostID: &postID,
		UserID:           &userID,
	}

	profileID := "profileID"
	profileObj := model.UserProfile{
		ID:            profileID,
		CreatedAt:     now,
		UpdatedAt:     now,
		DisplayName:   "name",
		UserID:        &userID,
		TwitterHandle: "handle",
	}

	Convey("GetUserByID", t, func() {
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

		expectedQueryString := `SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND (("users"."id" = $1)) ORDER BY "users"."id" ASC LIMIT 1`
		expectedParticipantString := `SELECT * FROM "participants"  WHERE "participants"."deleted_at" IS NULL AND (("user_id" IN ($1))) ORDER BY "participants"."id" ASC`
		expectedViewersString := `SELECT * FROM "viewers"  WHERE "viewers"."deleted_at" IS NULL AND (("user_id" IN ($1))) ORDER BY "viewers"."id" ASC`
		expectedProfileString := `SELECT * FROM "user_profiles"  WHERE "user_profiles"."deleted_at" IS NULL AND (("user_id" IN ($1))) ORDER BY "user_profiles"."id" ASC`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(userObj.ID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetUserByID(ctx, userObj.ID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and does not find a record", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(userObj.ID).WillReturnError(gorm.ErrRecordNotFound)

			resp, err := mockDatastore.GetUserByID(ctx, userObj.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and errors on participants query", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at"}).
				AddRow(userObj.ID, userObj.CreatedAt, userObj.UpdatedAt, userObj.DeletedAt)

			mock.ExpectQuery(expectedQueryString).WithArgs(userObj.ID).WillReturnRows(rs)
			mock.ExpectQuery(expectedParticipantString).WithArgs(userObj.ID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetUserByID(ctx, userObj.ID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and errors on viewers query", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at"}).
				AddRow(userObj.ID, userObj.CreatedAt, userObj.UpdatedAt, userObj.DeletedAt)

			participantRs := sqlmock.NewRows([]string{"id", "participant_id", "created_at", "updated_at", "deleted_at", "discussion_id",
				"viewer_id", "user_id", "flair_id", "is_anonymous", "gradient_color", "has_joined"}).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID, parObj.FlairID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined)

			mock.ExpectQuery(expectedQueryString).WithArgs(userObj.ID).WillReturnRows(rs)
			mock.ExpectQuery(expectedParticipantString).WithArgs(userObj.ID).WillReturnRows(participantRs)
			mock.ExpectQuery(expectedViewersString).WithArgs(userObj.ID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetUserByID(ctx, userObj.ID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and errors on user profiles query", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at"}).
				AddRow(userObj.ID, userObj.CreatedAt, userObj.UpdatedAt, userObj.DeletedAt)

			participantRs := sqlmock.NewRows([]string{"id", "participant_id", "created_at", "updated_at", "deleted_at", "discussion_id",
				"viewer_id", "user_id", "flair_id", "is_anonymous", "gradient_color", "has_joined"}).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID, parObj.FlairID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined)

			viewerRs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "last_viewed", "last_viewed_post_id",
				"discussion_id", "user_id"}).AddRow(viewerObj.ID, viewerObj.CreatedAt, viewerObj.UpdatedAt, viewerObj.DeletedAt,
				viewerObj.LastViewed, viewerObj.LastViewedPostID, viewerObj.DiscussionID, viewerObj.UserID)

			mock.ExpectQuery(expectedQueryString).WithArgs(userObj.ID).WillReturnRows(rs)
			mock.ExpectQuery(expectedParticipantString).WithArgs(userObj.ID).WillReturnRows(participantRs)
			mock.ExpectQuery(expectedViewersString).WithArgs(userObj.ID).WillReturnRows(viewerRs)
			mock.ExpectQuery(expectedProfileString).WithArgs(userObj.ID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetUserByID(ctx, userObj.ID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at"}).
				AddRow(userObj.ID, userObj.CreatedAt, userObj.UpdatedAt, userObj.DeletedAt)

			participantRs := sqlmock.NewRows([]string{"id", "participant_id", "created_at", "updated_at", "deleted_at", "discussion_id",
				"viewer_id", "user_id", "flair_id", "is_anonymous", "gradient_color", "has_joined"}).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID, parObj.FlairID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined)

			viewerRs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "last_viewed", "last_viewed_post_id",
				"discussion_id", "user_id"}).AddRow(viewerObj.ID, viewerObj.CreatedAt, viewerObj.UpdatedAt, viewerObj.DeletedAt,
				viewerObj.LastViewed, viewerObj.LastViewedPostID, viewerObj.DiscussionID, viewerObj.UserID)

			profileRs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "display_name",
				"user_id", "twitter_handle"}).AddRow(profileObj.ID, profileObj.CreatedAt, profileObj.UpdatedAt, profileObj.DeletedAt,
				profileObj.DisplayName, profileObj.UserID, profileObj.TwitterHandle)

			mock.ExpectQuery(expectedQueryString).WithArgs(userObj.ID).WillReturnRows(rs)
			mock.ExpectQuery(expectedParticipantString).WithArgs(userObj.ID).WillReturnRows(participantRs)
			mock.ExpectQuery(expectedViewersString).WithArgs(userObj.ID).WillReturnRows(viewerRs)
			mock.ExpectQuery(expectedProfileString).WithArgs(userObj.ID).WillReturnRows(profileRs)

			resp, err := mockDatastore.GetUserByID(ctx, userObj.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

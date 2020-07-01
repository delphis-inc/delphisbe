package datastore

import (
	"context"
	"database/sql"
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

func TestDelphisDB_GetModeratorByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	profileID := "profileID"
	modObj := model.Moderator{
		ID:            modID,
		CreatedAt:     now,
		UpdatedAt:     now,
		UserProfileID: &profileID,
		UserProfile: &model.UserProfile{
			ID: profileID,
		},
	}

	Convey("GetModeratorByID", t, func() {
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

		expectedQueryStr := `SELECT * FROM "moderators" WHERE "moderators"."deleted_at" IS NULL AND (("moderators"."id" = $1)) ORDER BY "moderators"."id" ASC LIMIT 1`
		expectedProfileQueryStr := `SELECT * FROM "user_profiles"  WHERE "user_profiles"."deleted_at" IS NULL AND (("id" IN ($1))) ORDER BY "user_profiles"."id" ASC`

		Convey("when record is not found", func() {
			expectedError := gorm.ErrRecordNotFound
			mock.ExpectQuery(expectedQueryStr).WithArgs(modObj.ID).WillReturnError(expectedError)

			resp, err := mockDatastore.GetModeratorByID(ctx, modID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when an error occurs", func() {
			expectedError := fmt.Errorf("Something went wrong")
			mock.ExpectQuery(expectedQueryStr).WithArgs(modObj.ID).WillReturnError(expectedError)

			resp, err := mockDatastore.GetModeratorByID(ctx, modID)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the objects are found", func() {
			mock.ExpectQuery(expectedQueryStr).WithArgs(modObj.ID).WillReturnRows(
				sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_profile_id"}).
					AddRow(modObj.ID, modObj.CreatedAt, modObj.UpdatedAt, modObj.DeletedAt, modObj.UserProfileID))
			mock.ExpectQuery(expectedProfileQueryStr).WithArgs(*modObj.UserProfileID).WillReturnRows(
				sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "display_name", "user_id", "twitter_handler"}).
					AddRow(modObj.UserProfile.ID, modObj.UserProfile.CreatedAt, modObj.UserProfile.UpdatedAt,
						modObj.UserProfile.DeletedAt, modObj.UserProfile.DisplayName, modObj.UserProfile.UserID,
						modObj.UserProfile.TwitterHandle))

			resp, err := mockDatastore.GetModeratorByID(ctx, modID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &modObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_CreateModerator(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	profileID := "profileID"
	userID := "userID"
	modObj := model.Moderator{
		ID:            modID,
		CreatedAt:     now,
		UpdatedAt:     now,
		UserProfileID: &profileID,
		UserProfile: &model.UserProfile{
			ID:            profileID,
			DisplayName:   "name",
			UserID:        &userID,
			TwitterHandle: "twitter",
		},
	}

	Convey("CreateModerator", t, func() {
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

		updateQueryStr := `UPDATE "user_profiles" SET "updated_at" = $1, "deleted_at" = $2, "display_name" = $3, "user_id" = $4, "twitter_handle" = $5  WHERE "user_profiles"."deleted_at" IS NULL AND "user_profiles"."id" = $6`
		insertModQueryStr := `INSERT INTO "moderators" ("id","created_at","updated_at","deleted_at","user_profile_id") VALUES ($1,$2,$3,$4,$5) RETURNING "moderators"."id"`
		expectedPostInsertSelectStr := `SELECT * FROM "moderators"  WHERE "moderators"."deleted_at" IS NULL AND (("moderators"."id" = $1)) ORDER BY "moderators"."id" ASC LIMIT 1`

		Convey("when create returns an error on updating the user_profile", func() {
			expectedError := fmt.Errorf("Some fake error")

			mock.ExpectBegin()
			mock.ExpectExec(updateQueryStr).WithArgs(
				sqlmock.AnyArg(), modObj.UserProfile.DeletedAt, modObj.UserProfile.DisplayName, *modObj.UserProfile.UserID, modObj.UserProfile.TwitterHandle, modObj.UserProfile.ID,
			).WillReturnError(expectedError)

			resp, err := mockDatastore.CreateModerator(ctx, modObj)

			assert.NotNil(t, err)
			assert.Nil(t, resp)
			assert.Nil(t, mock.ExpectationsWereMet())
		})

		Convey("when create returns error on inserting the moderator", func() {
			expectedError := fmt.Errorf("Some fake error")

			mock.ExpectBegin()
			mock.ExpectExec(updateQueryStr).WithArgs(
				sqlmock.AnyArg(), modObj.UserProfile.DeletedAt, modObj.UserProfile.DisplayName, *modObj.UserProfile.UserID, modObj.UserProfile.TwitterHandle, modObj.UserProfile.ID,
			).WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectQuery(insertModQueryStr).WithArgs(
				modObj.ID, modObj.CreatedAt, modObj.UpdatedAt, modObj.DeletedAt, modObj.UserProfileID,
			).WillReturnError(expectedError)

			resp, err := mockDatastore.CreateModerator(ctx, modObj)

			assert.NotNil(t, err)
			assert.Nil(t, resp)
			assert.Nil(t, mock.ExpectationsWereMet())
		})

		Convey("when create succeeds and errors on returning the new object", func() {
			expectedError := fmt.Errorf("Some fake error")
			mock.ExpectBegin()
			mock.ExpectExec(updateQueryStr).WithArgs(
				sqlmock.AnyArg(), modObj.UserProfile.DeletedAt, modObj.UserProfile.DisplayName, *modObj.UserProfile.UserID, modObj.UserProfile.TwitterHandle, modObj.UserProfile.ID,
			).WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectQuery(insertModQueryStr).WithArgs(
				modObj.ID, modObj.CreatedAt, modObj.UpdatedAt, modObj.DeletedAt, modObj.UserProfileID,
			).WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(modObj.ID))
			mock.ExpectCommit()
			mock.ExpectQuery(expectedPostInsertSelectStr).WithArgs(modID).WillReturnError(expectedError)

			resp, err := mockDatastore.CreateModerator(ctx, modObj)

			assert.NotNil(t, err)
			assert.Nil(t, resp)
			assert.Nil(t, mock.ExpectationsWereMet())
		})

		Convey("when create succeeds it should return the new object", func() {
			mock.ExpectBegin()
			mock.ExpectExec(updateQueryStr).WithArgs(
				sqlmock.AnyArg(), modObj.UserProfile.DeletedAt, modObj.UserProfile.DisplayName, *modObj.UserProfile.UserID, modObj.UserProfile.TwitterHandle, modObj.UserProfile.ID,
			).WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectQuery(insertModQueryStr).WithArgs(
				modObj.ID, modObj.CreatedAt, modObj.UpdatedAt, modObj.DeletedAt, modObj.UserProfileID,
			).WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(modObj.ID))
			mock.ExpectCommit()
			mock.ExpectQuery(expectedPostInsertSelectStr).WithArgs(modID).WillReturnRows(
				sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_profile_id"}).
					AddRow(modObj.ID, modObj.CreatedAt, modObj.UpdatedAt, modObj.DeletedAt, modObj.UserProfileID))

			resp, err := mockDatastore.CreateModerator(ctx, modObj)

			testModObj := modObj
			testModObj.UserProfile = nil

			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, *resp, testModObj)
			assert.Nil(t, mock.ExpectationsWereMet())
		})
	})
}

func TestDelphisDB_GetModeratorByUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	userID := "userID"
	profileID := "profileID"
	modObj := model.Moderator{
		ID:            modID,
		CreatedAt:     now,
		UpdatedAt:     now,
		UserProfileID: &profileID,
	}

	Convey("GetModeratorByUserID", t, func() {
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

			resp, err := mockDatastore.GetModeratorByUserID(ctx, userID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getModeratorByUserIDString).WithArgs(userID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetModeratorByUserID(ctx, userID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution does not find a record", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getModeratorByUserIDString).WithArgs(userID).WillReturnError(sql.ErrNoRows)

			resp, err := mockDatastore.GetModeratorByUserID(ctx, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a post", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"m.id", "m.created_at", "m.updated_at", "m.deleted_at", "m.user_profile_id"}).
				AddRow(modObj.ID, modObj.CreatedAt, modObj.UpdatedAt, modObj.DeletedAt, modObj.UserProfileID)

			mock.ExpectQuery(getModeratorByUserIDString).WithArgs(userID).WillReturnRows(rs)

			resp, err := mockDatastore.GetModeratorByUserID(ctx, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &modObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetModeratorByUserIDAndDiscussionID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	userID := "userID"
	discussionID := "discussionID"
	profileID := "profileID"
	modObj := model.Moderator{
		ID:            modID,
		CreatedAt:     now,
		UpdatedAt:     now,
		UserProfileID: &profileID,
	}

	Convey("GetModeratorByUserIDAndDiscussionID", t, func() {
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

			resp, err := mockDatastore.GetModeratorByUserIDAndDiscussionID(ctx, userID, discussionID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getModeratorByUserIDAndDiscussionIDString).WithArgs(userID, discussionID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetModeratorByUserIDAndDiscussionID(ctx, userID, discussionID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution does not find a record", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getModeratorByUserIDAndDiscussionIDString).WithArgs(userID, discussionID).WillReturnError(sql.ErrNoRows)

			resp, err := mockDatastore.GetModeratorByUserIDAndDiscussionID(ctx, userID, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a post", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"m.id", "m.created_at", "m.updated_at", "m.deleted_at", "m.user_profile_id"}).
				AddRow(modObj.ID, modObj.CreatedAt, modObj.UpdatedAt, modObj.DeletedAt, modObj.UserProfileID)

			mock.ExpectQuery(getModeratorByUserIDAndDiscussionIDString).WithArgs(userID, discussionID).WillReturnRows(rs)

			resp, err := mockDatastore.GetModeratorByUserIDAndDiscussionID(ctx, userID, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &modObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

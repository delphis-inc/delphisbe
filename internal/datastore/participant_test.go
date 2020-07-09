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

func TestDelphisDB_GetParticipantByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	parID := "parID"
	discussionID := "discussionID"
	viewerID := "viewerID"
	flairID := "flairID"
	gradientColor := model.GradientColorAzalea
	userID := "userID"
	inviterParticipantID := 0
	parObj := model.Participant{
		ID:                   parID,
		ParticipantID:        0,
		DiscussionID:         &discussionID,
		ViewerID:             &viewerID,
		FlairID:              &flairID,
		GradientColor:        &gradientColor,
		UserID:               &userID,
		HasJoined:            true,
		IsAnonymous:          false,
		CreatedAt:            now,
		UpdatedAt:            now,
		InviterParticipantID: &inviterParticipantID,
	}

	Convey("GetParticipantByID", t, func() {
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

		expectedQueryString := `SELECT * FROM "participants" WHERE "participants"."deleted_at" IS NULL AND (("participants"."id" IN ($1)))`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(parObj.ID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetParticipantByID(ctx, parID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and does not find a record", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(parObj.ID).WillReturnError(gorm.ErrRecordNotFound)

			resp, err := mockDatastore.GetParticipantByID(ctx, parID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a discussions", func() {
			rs := sqlmock.NewRows([]string{"id", "participant_id", "created_at", "updated_at", "deleted_at", "discussion_id",
				"viewer_id", "user_id", "flair_id", "is_anonymous", "gradient_color", "has_joined", "is_banned", "inviter_participant_id"}).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID, parObj.FlairID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned, parObj.InviterParticipantID)

			mock.ExpectQuery(expectedQueryString).WithArgs(parObj.ID).WillReturnRows(rs)

			resp, err := mockDatastore.GetParticipantByID(ctx, parID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &parObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetParticipantsByIDs(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	parID := "parID"
	discussionID := "discussionID"
	viewerID := "viewerID"
	flairID := "flairID"
	gradientColor := model.GradientColorAzalea
	userID := "userID"
	inviterParticipantID := 0
	parObj := model.Participant{
		ID:                   parID,
		ParticipantID:        0,
		DiscussionID:         &discussionID,
		ViewerID:             &viewerID,
		FlairID:              &flairID,
		GradientColor:        &gradientColor,
		UserID:               &userID,
		HasJoined:            true,
		IsAnonymous:          false,
		CreatedAt:            now,
		UpdatedAt:            now,
		IsBanned:             false,
		InviterParticipantID: &inviterParticipantID,
	}

	participants := []model.Participant{parObj, parObj}
	participantIDs := []string{parObj.ID, parObj.ID}

	Convey("GetParticipantByID", t, func() {
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

		expectedQueryString := `SELECT * FROM "participants" WHERE "participants"."deleted_at" IS NULL AND (("participants"."id" IN ($1,$2)))`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(participantIDs[0], participantIDs[1]).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetParticipantsByIDs(ctx, participantIDs)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and does not find a record", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(participantIDs[0], participantIDs[1]).WillReturnError(gorm.ErrRecordNotFound)

			resp, err := mockDatastore.GetParticipantsByIDs(ctx, participantIDs)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a discussions", func() {
			rs := sqlmock.NewRows([]string{"id", "participant_id", "created_at", "updated_at", "deleted_at", "discussion_id",
				"viewer_id", "user_id", "flair_id", "is_anonymous", "gradient_color", "has_joined", "is_banned", "inviter_participant_id"}).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID, parObj.FlairID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned, parObj.InviterParticipantID).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID, parObj.FlairID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned, parObj.InviterParticipantID)

			mock.ExpectQuery(expectedQueryString).WithArgs(participantIDs[0], participantIDs[1]).WillReturnRows(rs)

			resp, err := mockDatastore.GetParticipantsByIDs(ctx, participantIDs)

			verifyMap := map[string]*model.Participant{
				participants[0].ID: &parObj,
				participants[1].ID: &parObj,
			}

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, verifyMap)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetParticipantsByDiscussionID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	parID := "parID"
	discussionID := "discussionID"
	viewerID := "viewerID"
	flairID := "flairID"
	gradientColor := model.GradientColorAzalea
	userID := "userID"
	inviterParticipantID := 0
	parObj := model.Participant{
		ID:                   parID,
		ParticipantID:        0,
		DiscussionID:         &discussionID,
		ViewerID:             &viewerID,
		FlairID:              &flairID,
		GradientColor:        &gradientColor,
		UserID:               &userID,
		HasJoined:            true,
		IsAnonymous:          false,
		CreatedAt:            now,
		UpdatedAt:            now,
		IsBanned:             true,
		InviterParticipantID: &inviterParticipantID,
	}

	Convey("GetParticipantsByDiscussionID", t, func() {
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

		expectedQueryString := `SELECT * FROM "participants" WHERE "participants"."deleted_at" IS NULL AND (("participants"."discussion_id" = $1)) ORDER BY participant_id desc`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(parObj.DiscussionID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetParticipantsByDiscussionID(ctx, discussionID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and does not find a record", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(parObj.DiscussionID).WillReturnError(gorm.ErrRecordNotFound)

			resp, err := mockDatastore.GetParticipantsByDiscussionID(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a discussions", func() {
			rs := sqlmock.NewRows([]string{"id", "participant_id", "created_at", "updated_at", "deleted_at", "discussion_id",
				"viewer_id", "user_id", "flair_id", "is_anonymous", "gradient_color", "has_joined", "is_banned", "inviter_participant_id"}).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID, parObj.FlairID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned, parObj.InviterParticipantID).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID, parObj.FlairID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned, parObj.InviterParticipantID)

			mock.ExpectQuery(expectedQueryString).WithArgs(parObj.DiscussionID).WillReturnRows(rs)

			resp, err := mockDatastore.GetParticipantsByDiscussionID(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []model.Participant{parObj, parObj})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetParticipantsByDiscussionIDUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	parID := "parID"
	discussionID := "discussionID"
	viewerID := "viewerID"
	flairID := "flairID"
	gradientColor := model.GradientColorAzalea
	userID := "userID"
	inviterParticipantID := 0
	parObj := model.Participant{
		ID:                   parID,
		ParticipantID:        0,
		DiscussionID:         &discussionID,
		ViewerID:             &viewerID,
		FlairID:              &flairID,
		GradientColor:        &gradientColor,
		UserID:               &userID,
		HasJoined:            true,
		IsAnonymous:          false,
		CreatedAt:            now,
		UpdatedAt:            now,
		IsBanned:             false,
		InviterParticipantID: &inviterParticipantID,
	}

	Convey("GetParticipantsByDiscussionIDUserID", t, func() {
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

		expectedQueryString := `SELECT * FROM "participants" WHERE "participants"."deleted_at" IS NULL AND (("participants"."discussion_id" = $1) AND ("participants"."user_id" = $2)) ORDER BY participant_id desc LIMIT 2`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(parObj.DiscussionID, parObj.UserID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetParticipantsByDiscussionIDUserID(ctx, discussionID, userID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and does not find a record", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(parObj.DiscussionID, parObj.UserID).WillReturnError(gorm.ErrRecordNotFound)

			resp, err := mockDatastore.GetParticipantsByDiscussionIDUserID(ctx, discussionID, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a discussions", func() {
			rs := sqlmock.NewRows([]string{"id", "participant_id", "created_at", "updated_at", "deleted_at", "discussion_id",
				"viewer_id", "user_id", "flair_id", "is_anonymous", "gradient_color", "has_joined", "is_banned", "inviter_participant_id"}).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID, parObj.FlairID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned, parObj.InviterParticipantID).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID, parObj.FlairID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned, parObj.InviterParticipantID)

			mock.ExpectQuery(expectedQueryString).WithArgs(parObj.DiscussionID, parObj.UserID).WillReturnRows(rs)

			resp, err := mockDatastore.GetParticipantsByDiscussionIDUserID(ctx, discussionID, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []model.Participant{parObj, parObj})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetParticipantByDiscussionIDParticipantID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	parID := "parID"
	discussionID := "discussionID"
	viewerID := "viewerID"
	flairID := "flairID"
	gradientColor := model.GradientColorAzalea
	userID := "userID"
	participantID := 1
	inviterParticipantID := 0
	parObj := model.Participant{
		ID:                   parID,
		ParticipantID:        participantID,
		DiscussionID:         &discussionID,
		ViewerID:             &viewerID,
		FlairID:              &flairID,
		GradientColor:        &gradientColor,
		UserID:               &userID,
		HasJoined:            true,
		IsAnonymous:          false,
		CreatedAt:            now,
		UpdatedAt:            now,
		IsBanned:             false,
		InviterParticipantID: &inviterParticipantID,
	}

	Convey("GetParticipantByDiscussionIDParticipantID", t, func() {
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

		expectedQueryString := `SELECT * FROM "participants" WHERE "participants"."deleted_at" IS NULL AND (("participants"."discussion_id" = $1 AND "participants"."participant_id" = $2)) LIMIT 1`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(parObj.DiscussionID, parObj.ParticipantID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetParticipantByDiscussionIDParticipantID(ctx, discussionID, participantID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and does not find a record", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(parObj.DiscussionID, parObj.ParticipantID).WillReturnError(gorm.ErrRecordNotFound)

			resp, err := mockDatastore.GetParticipantByDiscussionIDParticipantID(ctx, discussionID, participantID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a discussions", func() {
			rs := sqlmock.NewRows([]string{"id", "participant_id", "created_at", "updated_at", "deleted_at", "discussion_id",
				"viewer_id", "user_id", "flair_id", "is_anonymous", "gradient_color", "has_joined", "is_banned", "inviter_participant_id"}).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID, parObj.FlairID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned, parObj.InviterParticipantID)

			mock.ExpectQuery(expectedQueryString).WithArgs(parObj.DiscussionID, parObj.ParticipantID).WillReturnRows(rs)

			resp, err := mockDatastore.GetParticipantByDiscussionIDParticipantID(ctx, discussionID, participantID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &parObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_UpsertParticipant(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	parID := "parID"
	discussionID := "discussionID"
	viewerID := "viewerID"
	flairID := "flairID"
	gradientColor := model.GradientColorAzalea
	userID := "userID"
	inviterParticipantID := 0
	parObj := model.Participant{
		ID:                   parID,
		ParticipantID:        0,
		DiscussionID:         &discussionID,
		ViewerID:             &viewerID,
		FlairID:              &flairID,
		GradientColor:        &gradientColor,
		UserID:               &userID,
		HasJoined:            true,
		IsAnonymous:          false,
		CreatedAt:            now,
		UpdatedAt:            now,
		IsBanned:             true,
		InviterParticipantID: &inviterParticipantID,
	}

	Convey("UpsertParticipant", t, func() {
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

		expectedFindQueryStr := `SELECT * FROM "participants" WHERE "participants"."deleted_at" IS NULL AND (("participants"."id" = $1)) ORDER BY "participants"."id" ASC LIMIT 1`
		createQueryStr := `INSERT INTO "participants" ("id","participant_id","created_at","updated_at","deleted_at","discussion_id","viewer_id","flair_id","gradient_color","user_id","is_banned","has_joined","is_anonymous","inviter_participant_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "participants"."id"`
		expectedNewObjectRow := sqlmock.NewRows([]string{"id", "participant_id", "created_at", "updated_at", "deleted_at", "discussion_id", "viewer_id", "flair_id", "gradient_color", "user_id", "is_banned", "has_joined", "is_anonymous", "inviter_participant_id"}).
			AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt, parObj.DiscussionID, parObj.ViewerID, parObj.FlairID, parObj.GradientColor, parObj.UserID, parObj.IsBanned, parObj.HasJoined, parObj.IsAnonymous, parObj.InviterParticipantID)
		expectedUpdateStr := `UPDATE "participants" SET "flair_id" = $1, "gradient_color" = $2, "has_joined" = $3, "inviter_participant_id" = $4, "is_banned" = $5, "updated_at" = $6 WHERE "participants"."deleted_at" IS NULL AND "participants"."id" = $7`
		expectedPostUpdateSelectStr := `SELECT * FROM "participants" WHERE "participants"."deleted_at" IS NULL AND "participants"."id" = $1 ORDER BY "participants"."id" ASC LIMIT 1`

		Convey("when find query fails with a non-not-found-error the function should return the error", func() {
			expectedError := fmt.Errorf("Some fake error")
			mock.ExpectQuery(expectedFindQueryStr).WithArgs(parObj.ID).WillReturnError(expectedError)

			resp, err := mockDatastore.UpsertParticipant(ctx, parObj)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when find query returns not-found-error it should call create", func() {
			Convey("when create returns an error it should return it", func() {
				expectedError := fmt.Errorf("Some fake error")

				mock.ExpectQuery(expectedFindQueryStr).WithArgs(parObj.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.FlairID, parObj.GradientColor, parObj.UserID, parObj.IsBanned,
					parObj.HasJoined, parObj.IsAnonymous, parObj.InviterParticipantID,
				).WillReturnError(expectedError)

				resp, err := mockDatastore.UpsertParticipant(ctx, parObj)

				assert.NotNil(t, err)
				assert.Nil(t, resp)
				assert.Nil(t, mock.ExpectationsWereMet())
			})

			Convey("when create succeeds it should return the new object", func() {
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(parObj.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.FlairID, parObj.GradientColor, parObj.UserID, parObj.IsBanned,
					parObj.HasJoined, parObj.IsAnonymous, parObj.InviterParticipantID,
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(parObj.ID))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(parObj.ID).WillReturnRows(expectedNewObjectRow)

				resp, err := mockDatastore.UpsertParticipant(ctx, parObj)

				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, *resp, parObj)
				assert.Nil(t, mock.ExpectationsWereMet())
			})
		})

		Convey("when the object is found we should update it", func() {
			Convey("when updating if it fails then return the error", func() {
				expectedError := fmt.Errorf("something went wrong")
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(parObj.ID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(expectedUpdateStr).WithArgs(
					parObj.FlairID, parObj.GradientColor, parObj.HasJoined, parObj.InviterParticipantID, parObj.IsBanned, sqlmock.AnyArg(), parObj.ID,
				).WillReturnError(expectedError)

				resp, err := mockDatastore.UpsertParticipant(ctx, parObj)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})

			Convey("when updating if it succeeds it should return the updated object", func() {
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(parObj.ID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(expectedUpdateStr).WithArgs(
					parObj.FlairID, parObj.GradientColor, parObj.HasJoined, parObj.InviterParticipantID, parObj.IsBanned, sqlmock.AnyArg(), parObj.ID,
				).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedPostUpdateSelectStr).WithArgs(parObj.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "participant_id", "created_at", "updated_at",
						"deleted_at", "discussion_id", "viewer_id", "flair_id", "gradient_color", "user_id",
						"has_joined", "is_anonymous", "is_banned", "inviter_participant_id"}).
						AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt,
							parObj.DeletedAt, parObj.DiscussionID, parObj.ViewerID, parObj.FlairID,
							parObj.GradientColor, parObj.UserID, parObj.HasJoined, parObj.IsAnonymous, parObj.IsBanned, parObj.InviterParticipantID))

				resp, err := mockDatastore.UpsertParticipant(ctx, parObj)

				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Nil(t, mock.ExpectationsWereMet())
			})
		})
	})
}

func TestDelphisDB_AssignFlair(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	parID := "parID"
	discussionID := "discussionID"
	viewerID := "viewerID"
	flairID := "flairID"
	gradientColor := model.GradientColorAzalea
	userID := "userID"
	inviterParticipantID := 0
	parObj := model.Participant{
		ID:                   parID,
		ParticipantID:        0,
		DiscussionID:         &discussionID,
		ViewerID:             &viewerID,
		FlairID:              &flairID,
		GradientColor:        &gradientColor,
		UserID:               &userID,
		HasJoined:            true,
		IsAnonymous:          false,
		CreatedAt:            now,
		UpdatedAt:            now,
		InviterParticipantID: &inviterParticipantID,
	}

	Convey("AssignFlair", t, func() {
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

		expectedUpdateStr := `UPDATE "participants" SET "flair_id" = $1 WHERE "participants"."deleted_at" IS NULL AND "participants"."id" = $2`

		Convey("when find query fails with a non-not-found-error the function should return the error", func() {
			expectedError := fmt.Errorf("Some fake error")
			mock.ExpectBegin()
			mock.ExpectExec(expectedUpdateStr).WithArgs(parObj.FlairID, parObj.ID).WillReturnError(expectedError)

			resp, err := mockDatastore.AssignFlair(ctx, parObj, &flairID)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldResemble, &parObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when find query updates successfully", func() {
			mock.ExpectBegin()
			mock.ExpectExec(expectedUpdateStr).WithArgs(parObj.FlairID, parObj.ID).WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()

			resp, err := mockDatastore.AssignFlair(ctx, parObj, &flairID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &parObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetTotalParticipantCountByDiscussionID(t *testing.T) {
	ctx := context.Background()

	discussionID := "discussionID"

	Convey("GetTotalParticipantCountByDiscussionID", t, func() {
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

		expectedCountQuery := `SELECT count(*) FROM "participants"  WHERE "participants"."deleted_at" IS NULL AND (("participants"."discussion_id" = $1))`

		Convey("when find query fails on the count", func() {
			expectedError := fmt.Errorf("Some fake error")
			mock.ExpectQuery(expectedCountQuery).WithArgs(discussionID).WillReturnError(expectedError)

			count := mockDatastore.GetTotalParticipantCountByDiscussionID(ctx, discussionID)

			So(count, ShouldEqual, 0)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when find query returns a count", func() {
			mock.ExpectQuery(expectedCountQuery).WithArgs(discussionID).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

			count := mockDatastore.GetTotalParticipantCountByDiscussionID(ctx, discussionID)

			So(count, ShouldEqual, 2)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

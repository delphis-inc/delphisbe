package datastore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestDelphisDB_GetParticipantByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	parID := "parID"
	discussionID := "discussionID"
	viewerID := "viewerID"
	gradientColor := model.GradientColorAzalea
	userID := "userID"
	parObj := model.Participant{
		ID:            parID,
		ParticipantID: 0,
		DiscussionID:  &discussionID,
		ViewerID:      &viewerID,
		GradientColor: &gradientColor,
		UserID:        &userID,
		HasJoined:     true,
		IsAnonymous:   false,
		CreatedAt:     now,
		UpdatedAt:     now,
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
				"viewer_id", "user_id", "is_anonymous", "gradient_color", "has_joined", "is_banned"}).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned)

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
	gradientColor := model.GradientColorAzalea
	userID := "userID"
	parObj := model.Participant{
		ID:            parID,
		ParticipantID: 0,
		DiscussionID:  &discussionID,
		ViewerID:      &viewerID,
		GradientColor: &gradientColor,
		UserID:        &userID,
		HasJoined:     true,
		IsAnonymous:   false,
		CreatedAt:     now,
		UpdatedAt:     now,
		IsBanned:      false,
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
				"viewer_id", "user_id", "is_anonymous", "gradient_color", "has_joined", "is_banned"}).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned)

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
	gradientColor := model.GradientColorAzalea
	userID := "userID"
	parObj := model.Participant{
		ID:            parID,
		ParticipantID: 0,
		DiscussionID:  &discussionID,
		ViewerID:      &viewerID,
		GradientColor: &gradientColor,
		UserID:        &userID,
		HasJoined:     true,
		IsAnonymous:   false,
		CreatedAt:     now,
		UpdatedAt:     now,
		IsBanned:      true,
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
				"viewer_id", "user_id", "is_anonymous", "gradient_color", "has_joined", "is_banned"}).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned)

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
	gradientColor := model.GradientColorAzalea
	userID := "userID"
	parObj := model.Participant{
		ID:            parID,
		ParticipantID: 0,
		DiscussionID:  &discussionID,
		ViewerID:      &viewerID,
		GradientColor: &gradientColor,
		UserID:        &userID,
		HasJoined:     true,
		IsAnonymous:   false,
		CreatedAt:     now,
		UpdatedAt:     now,
		IsBanned:      false,
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
				"viewer_id", "user_id", "is_anonymous", "gradient_color", "has_joined", "is_banned"}).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned)

			mock.ExpectQuery(expectedQueryString).WithArgs(parObj.DiscussionID, parObj.UserID).WillReturnRows(rs)

			resp, err := mockDatastore.GetParticipantsByDiscussionIDUserID(ctx, discussionID, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []model.Participant{parObj, parObj})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetModeratorParticipantsByDiscussionID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	parID := "parID"
	discussionID := "discussionID"
	viewerID := "viewerID"
	gradientColor := model.GradientColorAzalea
	userID := "userID"
	parObj := model.Participant{
		ID:            parID,
		ParticipantID: 0,
		DiscussionID:  &discussionID,
		ViewerID:      &viewerID,
		GradientColor: &gradientColor,
		UserID:        &userID,
		HasJoined:     true,
		IsAnonymous:   false,
		CreatedAt:     now,
		UpdatedAt:     now,
		IsBanned:      false,
	}

	Convey("GetModeratorParticipantsByDiscussionID", t, func() {
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
		joinUserProfiles := "JOIN user_profiles ON user_profiles.user_id = participants.user_id"
		joinModerators := "JOIN moderators ON user_profiles.id = moderators.user_profile_id"
		joinDiscussions := "JOIN discussions ON participants.discussion_id = discussions.id"
		expectedQueryString := `SELECT "participants".* FROM "participants" ` + joinUserProfiles + ` ` + joinModerators + ` ` + joinDiscussions + ` WHERE "participants"."deleted_at" IS NULL AND ((("discussions"."moderator_id" = "moderators"."id") AND "discussions"."id" = $1)) ORDER BY participant_id desc LIMIT 2`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(parObj.DiscussionID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetModeratorParticipantsByDiscussionID(ctx, discussionID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and does not find a record", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(parObj.DiscussionID).WillReturnError(gorm.ErrRecordNotFound)

			resp, err := mockDatastore.GetModeratorParticipantsByDiscussionID(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a discussions", func() {
			rs := sqlmock.NewRows([]string{"id", "participant_id", "created_at", "updated_at", "deleted_at", "discussion_id",
				"viewer_id", "user_id", "is_anonymous", "gradient_color", "has_joined", "is_banned"}).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt,
					parObj.DiscussionID, parObj.ViewerID, parObj.UserID,
					parObj.IsAnonymous, parObj.GradientColor, parObj.HasJoined, parObj.IsBanned)

			mock.ExpectQuery(expectedQueryString).WithArgs(parObj.DiscussionID).WillReturnRows(rs)

			resp, err := mockDatastore.GetModeratorParticipantsByDiscussionID(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []model.Participant{parObj, parObj})
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
	gradientColor := model.GradientColorAzalea
	userID := "userID"
	parObj := model.Participant{
		ID:            parID,
		ParticipantID: 0,
		DiscussionID:  &discussionID,
		ViewerID:      &viewerID,
		GradientColor: &gradientColor,
		UserID:        &userID,
		HasJoined:     true,
		IsAnonymous:   false,
		CreatedAt:     now,
		UpdatedAt:     now,
		IsBanned:      true,
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
		createQueryStr := `INSERT INTO "participants" ("id","participant_id","created_at","updated_at","deleted_at","discussion_id","viewer_id","gradient_color","user_id","is_banned","has_joined","is_anonymous","muted_until") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING "participants"."id"`
		expectedNewObjectRow := sqlmock.NewRows([]string{"id", "participant_id", "created_at", "updated_at", "deleted_at", "discussion_id", "viewer_id", "gradient_color", "user_id", "is_banned", "has_joined", "is_anonymous", "muted_until"}).
			AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt, parObj.DeletedAt, parObj.DiscussionID, parObj.ViewerID, parObj.GradientColor, parObj.UserID, parObj.IsBanned, parObj.HasJoined, parObj.IsAnonymous, parObj.MutedUntil)
		expectedUpdateStr := `UPDATE "participants" SET "gradient_color" = $1, "has_joined" = $2, "is_banned" = $3, "updated_at" = $4 WHERE "participants"."deleted_at" IS NULL AND "participants"."id" = $5`
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
					parObj.DiscussionID, parObj.ViewerID, parObj.GradientColor, parObj.UserID, parObj.IsBanned,
					parObj.HasJoined, parObj.IsAnonymous, parObj.MutedUntil,
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
					parObj.DiscussionID, parObj.ViewerID, parObj.GradientColor, parObj.UserID, parObj.IsBanned,
					parObj.HasJoined, parObj.IsAnonymous, parObj.MutedUntil,
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
					parObj.GradientColor, parObj.HasJoined, parObj.IsBanned, sqlmock.AnyArg(), parObj.ID,
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
					parObj.GradientColor, parObj.HasJoined, parObj.IsBanned, sqlmock.AnyArg(), parObj.ID,
				).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedPostUpdateSelectStr).WithArgs(parObj.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "participant_id", "created_at", "updated_at",
						"deleted_at", "discussion_id", "viewer_id", "gradient_color", "user_id",
						"has_joined", "is_anonymous", "is_banned"}).
						AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt,
							parObj.DeletedAt, parObj.DiscussionID, parObj.ViewerID,
							parObj.GradientColor, parObj.UserID, parObj.HasJoined, parObj.IsAnonymous, parObj.IsBanned))

				resp, err := mockDatastore.UpsertParticipant(ctx, parObj)

				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Nil(t, mock.ExpectationsWereMet())
			})
		})
	})
}

func TestDelphisDB_SetParticipantsMutedUntil(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	parID := "parID"
	discussionID := "discussionID"
	viewerID := "viewerID"
	gradientColor := model.GradientColorAzalea
	userID := "userID"
	parObj := model.Participant{
		ID:            parID,
		ParticipantID: 0,
		DiscussionID:  &discussionID,
		ViewerID:      &viewerID,
		GradientColor: &gradientColor,
		UserID:        &userID,
		HasJoined:     true,
		IsAnonymous:   false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	emptyListObj := []*model.Participant{}
	parListObj := []*model.Participant{&parObj}
	timeObj := time.Now().Add(time.Duration(60) * time.Second)

	Convey("SetParticipantsMutedUntil", t, func() {
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

		expectedUpdateStr := `UPDATE "participants" SET "muted_until" = $1 WHERE (id IN ($2))`
		expectedSelectStr := `SELECT * FROM "participants" WHERE "participants"."deleted_at" IS NULL AND ((id IN ($1)))`

		Convey("when update query errors out", func() {
			expectedError := fmt.Errorf("Some fake error")
			mock.ExpectBegin()
			mock.ExpectExec(expectedUpdateStr).WithArgs(&timeObj, parID).WillReturnError(expectedError)

			resp, err := mockDatastore.SetParticipantsMutedUntil(ctx, parListObj, &timeObj)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldResemble, emptyListObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when find query errors out", func() {
			expectedError := fmt.Errorf("Some fake error")
			mock.ExpectBegin()
			mock.ExpectExec(expectedUpdateStr).WithArgs(&timeObj, parID).WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()
			mock.ExpectQuery(expectedSelectStr).WithArgs(parID).WillReturnError(expectedError)

			resp, err := mockDatastore.SetParticipantsMutedUntil(ctx, parListObj, &timeObj)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldResemble, emptyListObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when find query succeeds", func() {
			expectedParObj := parObj
			expectedParObj.MutedUntil = &timeObj
			expectedParListObj := []*model.Participant{&expectedParObj}
			mock.ExpectBegin()
			mock.ExpectExec(expectedUpdateStr).WithArgs(&timeObj, parID).WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()
			mock.ExpectQuery(expectedSelectStr).WithArgs(parID).WillReturnRows(sqlmock.NewRows([]string{
				"id", "participant_id", "created_at", "updated_at",
				"deleted_at", "discussion_id", "viewer_id", "gradient_color", "user_id",
				"has_joined", "is_anonymous", "is_banned", "muted_until"}).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt,
					parObj.DeletedAt, parObj.DiscussionID, parObj.ViewerID,
					parObj.GradientColor, parObj.UserID, parObj.HasJoined, parObj.IsAnonymous, parObj.IsBanned, timeObj))

			resp, err := mockDatastore.SetParticipantsMutedUntil(ctx, parListObj, &timeObj)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, expectedParListObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when everything succeeds with nil time", func() {
			expectedParObj := parObj
			expectedParObj.MutedUntil = nil
			expectedParListObj := []*model.Participant{&expectedParObj}
			mock.ExpectBegin()
			mock.ExpectExec(expectedUpdateStr).WithArgs(nil, parID).WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()
			mock.ExpectQuery(expectedSelectStr).WithArgs(parID).WillReturnRows(sqlmock.NewRows([]string{
				"id", "participant_id", "created_at", "updated_at",
				"deleted_at", "discussion_id", "viewer_id", "gradient_color", "user_id",
				"has_joined", "is_anonymous", "is_banned", "muted_until"}).
				AddRow(parObj.ID, parObj.ParticipantID, parObj.CreatedAt, parObj.UpdatedAt,
					parObj.DeletedAt, parObj.DiscussionID, parObj.ViewerID,
					parObj.GradientColor, parObj.UserID, parObj.HasJoined, parObj.IsAnonymous, parObj.IsBanned, nil))

			resp, err := mockDatastore.SetParticipantsMutedUntil(ctx, parListObj, nil)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, expectedParListObj)
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

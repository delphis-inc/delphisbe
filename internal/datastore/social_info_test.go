package datastore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestDelphisDB_UpsertSocialInfo(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := "userID"
	profileID := "profileID"
	socialObj := model.SocialInfo{
		CreatedAt:         now,
		UpdatedAt:         now,
		AccessToken:       "access_token",
		AccessTokenSecret: "secret",
		UserID:            userID,
		ProfileImageURL:   "url",
		ScreenName:        "screen_nam",
		IsVerified:        false,
		Network:           util.SocialNetworkTwitter,
		UserProfileID:     profileID,
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

		expectedFindQueryStr := `SELECT * FROM "social_infos" WHERE "social_infos"."deleted_at" IS NULL AND (("social_infos"."network" = $1) AND ("social_infos"."user_profile_id" = $2)) ORDER BY "social_infos"."network" ASC LIMIT 1`
		createQueryStr := `INSERT INTO "social_infos" ("created_at","updated_at","deleted_at","access_token","access_token_secret","user_id","profile_image_url","screen_name","is_verified","network","user_profile_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "social_infos"."network"`
		expectedNewObjectRow := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "user_profile_id", "network",
			"access_token", "access_token_secret", "user_id", "profile_image_url", "screen_name", "is_verified"}).
			AddRow(socialObj.CreatedAt, socialObj.UpdatedAt, socialObj.DeletedAt, socialObj.UserProfileID,
				socialObj.Network, socialObj.AccessToken, socialObj.AccessTokenSecret, socialObj.UserID,
				socialObj.ProfileImageURL, socialObj.ScreenName, socialObj.IsVerified)
		expectedUpdateStr := `UPDATE "social_infos" SET "access_token" = $1, "access_token_secret" = $2, "profile_image_url" = $3, "screen_name" = $4, "updated_at" = $5, "user_id" = $6 WHERE "social_infos"."deleted_at" IS NULL AND "social_infos"."network" = $7 AND "social_infos"."user_profile_id" = $8`
		expectedPostUpdateSelectStr := `SELECT * FROM "social_infos" WHERE "social_infos"."deleted_at" IS NULL AND "social_infos"."network" = $1 AND "social_infos"."user_profile_id" = $2 ORDER BY "social_infos"."network" ASC LIMIT 1`

		Convey("when find query fails with a non-not-found-error the function should return the error", func() {
			expectedError := fmt.Errorf("Some fake error")
			mock.ExpectQuery(expectedFindQueryStr).WithArgs(socialObj.Network, socialObj.UserProfileID).WillReturnError(expectedError)

			resp, err := mockDatastore.UpsertSocialInfo(ctx, socialObj)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when find query returns not-found-error it should call create", func() {
			Convey("when create returns an error it should return it", func() {
				expectedError := fmt.Errorf("Some fake error")

				mock.ExpectQuery(expectedFindQueryStr).WithArgs(socialObj.Network, socialObj.UserProfileID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					socialObj.CreatedAt, socialObj.UpdatedAt, socialObj.DeletedAt, socialObj.AccessToken, socialObj.AccessTokenSecret,
					socialObj.UserID, socialObj.ProfileImageURL, socialObj.ScreenName, socialObj.IsVerified,
					socialObj.Network, socialObj.UserProfileID,
				).WillReturnError(expectedError)

				resp, err := mockDatastore.UpsertSocialInfo(ctx, socialObj)

				assert.NotNil(t, err)
				assert.Nil(t, resp)
				assert.Nil(t, mock.ExpectationsWereMet())
			})

			Convey("when create succeeds it should return the new object", func() {
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(socialObj.Network, socialObj.UserProfileID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					socialObj.CreatedAt, socialObj.UpdatedAt, socialObj.DeletedAt, socialObj.AccessToken, socialObj.AccessTokenSecret,
					socialObj.UserID, socialObj.ProfileImageURL, socialObj.ScreenName, socialObj.IsVerified,
					socialObj.Network, socialObj.UserProfileID,
				).WillReturnRows(sqlmock.NewRows([]string{"network"}).AddRow(socialObj.Network))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(socialObj.Network, socialObj.UserProfileID).WillReturnRows(expectedNewObjectRow)

				resp, err := mockDatastore.UpsertSocialInfo(ctx, socialObj)

				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, *resp, socialObj)
				assert.Nil(t, mock.ExpectationsWereMet())
			})
		})

		Convey("when the object is found we should update it", func() {
			Convey("when updating if it fails then return the error", func() {
				expectedError := fmt.Errorf("something went wrong")
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(socialObj.Network, socialObj.UserProfileID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(expectedUpdateStr).WithArgs(
					socialObj.AccessToken, socialObj.AccessTokenSecret, socialObj.ProfileImageURL, socialObj.ScreenName,
					sqlmock.AnyArg(), socialObj.UserID, socialObj.Network, socialObj.UserProfileID,
				).WillReturnError(expectedError)

				resp, err := mockDatastore.UpsertSocialInfo(ctx, socialObj)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})

			Convey("when updating if it succeeds it should return the updated object", func() {
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(socialObj.Network, socialObj.UserProfileID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(expectedUpdateStr).WithArgs(
					socialObj.AccessToken, socialObj.AccessTokenSecret, socialObj.ProfileImageURL, socialObj.ScreenName,
					sqlmock.AnyArg(), socialObj.UserID, socialObj.Network, socialObj.UserProfileID,
				).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedPostUpdateSelectStr).WithArgs(socialObj.Network, socialObj.UserProfileID).
					WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "user_profile_id", "network",
						"access_token", "access_token_secret", "user_id", "profile_image_url", "screen_name", "is_verified"}).
						AddRow(socialObj.CreatedAt, socialObj.UpdatedAt, socialObj.DeletedAt, socialObj.UserProfileID,
							socialObj.Network, socialObj.AccessToken, socialObj.AccessTokenSecret, socialObj.UserID,
							socialObj.ProfileImageURL, socialObj.ScreenName, socialObj.IsVerified))

				resp, err := mockDatastore.UpsertSocialInfo(ctx, socialObj)

				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Nil(t, mock.ExpectationsWereMet())
			})
		})
	})
}

func TestDelphisDB_GetSocialInfosByUserProfileID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := "userID"
	profileID := "profileID"
	socialObj := model.SocialInfo{
		CreatedAt:         now,
		UpdatedAt:         now,
		AccessToken:       "access_token",
		AccessTokenSecret: "secret",
		UserID:            userID,
		ProfileImageURL:   "url",
		ScreenName:        "screen_nam",
		IsVerified:        false,
		Network:           util.SocialNetworkTwitter,
		UserProfileID:     profileID,
	}

	Convey("GetSocialInfosByUserProfileID", t, func() {
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

		expectedQueryString := `SELECT * FROM "social_infos" WHERE "social_infos"."deleted_at" IS NULL AND (("social_infos"."user_profile_id" = $1))`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(socialObj.UserProfileID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetSocialInfosByUserProfileID(ctx, profileID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and does not find a record", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(socialObj.UserProfileID).WillReturnError(gorm.ErrRecordNotFound)

			resp, err := mockDatastore.GetSocialInfosByUserProfileID(ctx, profileID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a socials", func() {
			rs := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "user_profile_id", "network",
				"access_token", "access_token_secret", "user_id", "profile_image_url", "screen_name", "is_verified"}).
				AddRow(socialObj.CreatedAt, socialObj.UpdatedAt, socialObj.DeletedAt, socialObj.UserProfileID,
					socialObj.Network, socialObj.AccessToken, socialObj.AccessTokenSecret, socialObj.UserID,
					socialObj.ProfileImageURL, socialObj.ScreenName, socialObj.IsVerified).
				AddRow(socialObj.CreatedAt, socialObj.UpdatedAt, socialObj.DeletedAt, socialObj.UserProfileID,
					socialObj.Network, socialObj.AccessToken, socialObj.AccessTokenSecret, socialObj.UserID,
					socialObj.ProfileImageURL, socialObj.ScreenName, socialObj.IsVerified)

			mock.ExpectQuery(expectedQueryString).WithArgs(socialObj.UserProfileID).WillReturnRows(rs)

			resp, err := mockDatastore.GetSocialInfosByUserProfileID(ctx, profileID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []model.SocialInfo{socialObj, socialObj})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

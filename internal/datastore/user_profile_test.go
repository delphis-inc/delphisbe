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

func TestDelphisDB_GetUserProfileByUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := "userID"
	profileID := "profileID"
	profileObj := model.UserProfile{
		ID:            profileID,
		CreatedAt:     now,
		UpdatedAt:     now,
		DisplayName:   "name",
		UserID:        &userID,
		TwitterHandle: "handle",
	}

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

	userObj := model.User{
		ID:        userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	profileObj.SocialInfos = []model.SocialInfo{socialObj}

	Convey("GetUserProfileByUserID", t, func() {
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
		expectedProfileString := `SELECT * FROM "user_profiles" WHERE "user_profiles"."deleted_at" IS NULL AND (("user_id" = $1)) ORDER BY "user_profiles"."id" ASC`
		expectedSocialString := `SELECT * FROM "social_infos" WHERE "social_infos"."deleted_at" IS NULL AND (("user_profile_id" IN ($1))) ORDER BY "social_infos"."network" ASC`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(*profileObj.UserID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetUserProfileByUserID(ctx, *profileObj.UserID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and does not find a record", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(*profileObj.UserID).WillReturnError(gorm.ErrRecordNotFound)

			resp, err := mockDatastore.GetUserProfileByUserID(ctx, *profileObj.UserID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and errors on user profiles", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at"}).
				AddRow(userObj.ID, userObj.CreatedAt, userObj.UpdatedAt, userObj.DeletedAt)

			mock.ExpectQuery(expectedQueryString).WithArgs(*profileObj.UserID).WillReturnRows(rs)
			mock.ExpectQuery(expectedProfileString).WithArgs(*profileObj.UserID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetUserProfileByUserID(ctx, *profileObj.UserID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and errors on social info", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at"}).
				AddRow(userObj.ID, userObj.CreatedAt, userObj.UpdatedAt, userObj.DeletedAt)

			profileRs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "display_name",
				"user_id", "twitter_handle"}).AddRow(profileObj.ID, profileObj.CreatedAt, profileObj.UpdatedAt, profileObj.DeletedAt,
				profileObj.DisplayName, profileObj.UserID, profileObj.TwitterHandle)

			mock.ExpectQuery(expectedQueryString).WithArgs(*profileObj.UserID).WillReturnRows(rs)
			mock.ExpectQuery(expectedProfileString).WithArgs(*profileObj.UserID).WillReturnRows(profileRs)
			mock.ExpectQuery(expectedSocialString).WithArgs(profileObj.ID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetUserProfileByUserID(ctx, *profileObj.UserID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at"}).
				AddRow(userObj.ID, userObj.CreatedAt, userObj.UpdatedAt, userObj.DeletedAt)

			profileRs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "display_name",
				"user_id", "twitter_handle"}).AddRow(profileObj.ID, profileObj.CreatedAt, profileObj.UpdatedAt, profileObj.DeletedAt,
				profileObj.DisplayName, profileObj.UserID, profileObj.TwitterHandle)

			socialRs := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "user_profile_id", "network",
				"access_token", "access_token_secret", "user_id", "profile_image_url", "screen_name", "is_verified"}).
				AddRow(socialObj.CreatedAt, socialObj.UpdatedAt, socialObj.DeletedAt, socialObj.UserProfileID,
					socialObj.Network, socialObj.AccessToken, socialObj.AccessTokenSecret, socialObj.UserID,
					socialObj.ProfileImageURL, socialObj.ScreenName, socialObj.IsVerified)

			mock.ExpectQuery(expectedQueryString).WithArgs(*profileObj.UserID).WillReturnRows(rs)
			mock.ExpectQuery(expectedProfileString).WithArgs(*profileObj.UserID).WillReturnRows(profileRs)
			mock.ExpectQuery(expectedSocialString).WithArgs(profileObj.ID).WillReturnRows(socialRs)

			resp, err := mockDatastore.GetUserProfileByUserID(ctx, *profileObj.UserID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &profileObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetUserProfileByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := "userID"
	profileID := "profileID"
	profileObj := model.UserProfile{
		ID:            profileID,
		CreatedAt:     now,
		UpdatedAt:     now,
		DisplayName:   "name",
		UserID:        &userID,
		TwitterHandle: "handle",
	}

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

	profileObj.SocialInfos = []model.SocialInfo{socialObj}

	Convey("GetUserProfileByID", t, func() {
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

		expectedQueryString := `SELECT * FROM "user_profiles" WHERE "user_profiles"."deleted_at" IS NULL AND (("user_profiles"."id" = $1)) ORDER BY "user_profiles"."id" ASC LIMIT 1`
		expectedSocialString := `SELECT * FROM "social_infos"  WHERE "social_infos"."deleted_at" IS NULL AND (("user_profile_id" IN ($1))) ORDER BY "social_infos"."network" ASC`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(profileObj.ID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetUserProfileByID(ctx, profileObj.ID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and does not find a record", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(profileObj.ID).WillReturnError(gorm.ErrRecordNotFound)

			resp, err := mockDatastore.GetUserProfileByID(ctx, profileObj.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and errors on social infos", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "display_name",
				"user_id", "twitter_handle"}).AddRow(profileObj.ID, profileObj.CreatedAt, profileObj.UpdatedAt, profileObj.DeletedAt,
				profileObj.DisplayName, profileObj.UserID, profileObj.TwitterHandle)

			mock.ExpectQuery(expectedQueryString).WithArgs(profileObj.ID).WillReturnRows(rs)
			mock.ExpectQuery(expectedSocialString).WithArgs(profileObj.ID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetUserProfileByID(ctx, profileObj.ID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "display_name",
				"user_id", "twitter_handle"}).AddRow(profileObj.ID, profileObj.CreatedAt, profileObj.UpdatedAt, profileObj.DeletedAt,
				profileObj.DisplayName, profileObj.UserID, profileObj.TwitterHandle)

			socialRs := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "user_profile_id", "network",
				"access_token", "access_token_secret", "user_id", "profile_image_url", "screen_name", "is_verified"}).
				AddRow(socialObj.CreatedAt, socialObj.UpdatedAt, socialObj.DeletedAt, socialObj.UserProfileID,
					socialObj.Network, socialObj.AccessToken, socialObj.AccessTokenSecret, socialObj.UserID,
					socialObj.ProfileImageURL, socialObj.ScreenName, socialObj.IsVerified)

			mock.ExpectQuery(expectedQueryString).WithArgs(profileObj.ID).WillReturnRows(rs)
			mock.ExpectQuery(expectedSocialString).WithArgs(profileObj.ID).WillReturnRows(socialRs)

			resp, err := mockDatastore.GetUserProfileByID(ctx, profileObj.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &profileObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_CreateOrUpdateUserProfile(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := "userID"
	profileID := "profileID"
	profileObj := model.UserProfile{
		ID:            profileID,
		CreatedAt:     now,
		UpdatedAt:     now,
		DisplayName:   "name",
		UserID:        &userID,
		TwitterHandle: "handle",
	}

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

	profileObj.SocialInfos = []model.SocialInfo{socialObj}

	Convey("CreateOrUpdateUserProfile", t, func() {
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

		expectedFindQueryStr := `SELECT * FROM "user_profiles" WHERE "user_profiles"."deleted_at" IS NULL AND (("user_profiles"."id" = $1)) ORDER BY "user_profiles"."id" ASC LIMIT 1`
		expectedFindSocialStr := `SELECT * FROM "social_infos"  WHERE "social_infos"."deleted_at" IS NULL AND (("user_profile_id" IN ($1))) ORDER BY "social_infos"."network" ASC`
		createQueryStr := `INSERT INTO "user_profiles" ("id","created_at","updated_at","deleted_at","display_name","user_id","twitter_handle") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "user_profiles"."id"`
		expectedNewObjectRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "display_name", "user_id", "twitter_handle"}).
			AddRow(profileObj.ID, profileObj.CreatedAt, profileObj.UpdatedAt, profileObj.DeletedAt, profileObj.DisplayName, profileObj.UserID, profileObj.TwitterHandle)
		expectedNewSocialObjectRow := sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "user_profile_id", "network",
			"access_token", "access_token_secret", "user_id", "profile_image_url", "screen_name", "is_verified"}).
			AddRow(socialObj.CreatedAt, socialObj.UpdatedAt, socialObj.DeletedAt, socialObj.UserProfileID,
				socialObj.Network, socialObj.AccessToken, socialObj.AccessTokenSecret, socialObj.UserID,
				socialObj.ProfileImageURL, socialObj.ScreenName, socialObj.IsVerified)
		expectedUpdateSocialStr := `UPDATE "social_infos" SET "created_at" = $1, "updated_at" = $2, "deleted_at" = $3, "access_token" = $4, "access_token_secret" = $5, "user_id" = $6, "profile_image_url" = $7, "screen_name" = $8, "is_verified" = $9  WHERE "social_infos"."deleted_at" IS NULL AND "social_infos"."network" = $10 AND "social_infos"."user_profile_id" = $11`
		expectedUpdateStr := `UPDATE "user_profiles" SET "display_name" = $1, "twitter_handle" = $2, "updated_at" = $3 WHERE "user_profiles"."deleted_at" IS NULL AND "user_profiles"."id" = $4`
		expectedPostUpdateSelectStr := `SELECT * FROM "user_profiles" WHERE "user_profiles"."deleted_at" IS NULL AND "user_profiles"."id" = $1 ORDER BY "user_profiles"."id" ASC LIMIT 1`
		expectedPostUpdateSocialSelectStr := `SELECT * FROM "social_infos"  WHERE "social_infos"."deleted_at" IS NULL AND (("user_profile_id" IN ($1))) ORDER BY "social_infos"."network" ASC`

		Convey("when find query fails with a non-not-found-error the function should return the error", func() {
			expectedError := fmt.Errorf("Some fake error")
			mock.ExpectQuery(expectedFindQueryStr).WithArgs(profileObj.ID).WillReturnError(expectedError)

			resp, inserted, err := mockDatastore.CreateOrUpdateUserProfile(ctx, profileObj)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(inserted, ShouldBeFalse)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when find query returns not-found-error it should call create", func() {
			Convey("when create returns an error it should return it", func() {
				expectedError := fmt.Errorf("Some fake error")

				mock.ExpectQuery(expectedFindQueryStr).WithArgs(profileObj.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					profileObj.ID, profileObj.CreatedAt, profileObj.UpdatedAt, profileObj.DeletedAt,
					profileObj.DisplayName, profileObj.UserID, profileObj.TwitterHandle,
				).WillReturnError(expectedError)

				resp, inserted, err := mockDatastore.CreateOrUpdateUserProfile(ctx, profileObj)

				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, expectedError)
				So(inserted, ShouldBeFalse)
				So(resp, ShouldBeNil)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})

			Convey("when create succeeds but social update fails", func() {
				expectedError := fmt.Errorf("Some fake error")

				mock.ExpectQuery(expectedFindQueryStr).WithArgs(profileObj.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					profileObj.ID, profileObj.CreatedAt, profileObj.UpdatedAt, profileObj.DeletedAt,
					profileObj.DisplayName, profileObj.UserID, profileObj.TwitterHandle,
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(profileObj.ID))
				mock.ExpectExec(expectedUpdateSocialStr).WithArgs(
					socialObj.CreatedAt, sqlmock.AnyArg(), socialObj.DeletedAt, socialObj.AccessToken,
					socialObj.AccessTokenSecret, socialObj.UserID, socialObj.ProfileImageURL, socialObj.ScreenName,
					socialObj.IsVerified, socialObj.Network, socialObj.UserProfileID,
				).WillReturnError(expectedError)

				resp, inserted, err := mockDatastore.CreateOrUpdateUserProfile(ctx, profileObj)

				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, expectedError)
				So(inserted, ShouldBeFalse)
				So(resp, ShouldBeNil)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})

			Convey("when create succeeds it should return the new object", func() {
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(profileObj.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					profileObj.ID, profileObj.CreatedAt, profileObj.UpdatedAt, profileObj.DeletedAt,
					profileObj.DisplayName, profileObj.UserID, profileObj.TwitterHandle,
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(profileObj.ID))
				mock.ExpectExec(expectedUpdateSocialStr).WithArgs(
					socialObj.CreatedAt, sqlmock.AnyArg(), socialObj.DeletedAt, socialObj.AccessToken,
					socialObj.AccessTokenSecret, socialObj.UserID, socialObj.ProfileImageURL, socialObj.ScreenName,
					socialObj.IsVerified, socialObj.Network, socialObj.UserProfileID,
				).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(profileObj.ID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectQuery(expectedFindSocialStr).WithArgs(profileObj.ID).WillReturnRows(expectedNewSocialObjectRow)

				resp, inserted, err := mockDatastore.CreateOrUpdateUserProfile(ctx, profileObj)

				So(err, ShouldBeNil)
				So(inserted, ShouldBeTrue)
				So(resp, ShouldNotBeNil)
				So(resp, ShouldResemble, &profileObj)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})
		})

		Convey("when the object is found we should update it", func() {
			Convey("when updating if it fails then return the error", func() {
				expectedError := fmt.Errorf("something went wrong")
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(profileObj.ID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(expectedUpdateStr).WithArgs(
					profileObj.DisplayName, profileObj.TwitterHandle, sqlmock.AnyArg(), profileObj.ID,
				).WillReturnError(expectedError)

				resp, inserted, err := mockDatastore.CreateOrUpdateUserProfile(ctx, profileObj)

				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, expectedError)
				So(inserted, ShouldBeFalse)
				So(resp, ShouldBeNil)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})

			Convey("when updating and fails on updating social info", func() {
				expectedError := fmt.Errorf("Some fake error")

				mock.ExpectQuery(expectedFindQueryStr).WithArgs(profileObj.ID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(expectedUpdateStr).WithArgs(
					profileObj.DisplayName, profileObj.TwitterHandle, sqlmock.AnyArg(), profileObj.ID,
				).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(expectedUpdateSocialStr).WithArgs(
					socialObj.CreatedAt, sqlmock.AnyArg(), socialObj.DeletedAt, socialObj.AccessToken,
					socialObj.AccessTokenSecret, socialObj.UserID, socialObj.ProfileImageURL, socialObj.ScreenName,
					socialObj.IsVerified, socialObj.Network, socialObj.UserProfileID,
				).WillReturnError(expectedError)

				resp, inserted, err := mockDatastore.CreateOrUpdateUserProfile(ctx, profileObj)

				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, expectedError)
				So(inserted, ShouldBeFalse)
				So(resp, ShouldBeNil)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})

			Convey("when updating if it succeeds it should return the updated object", func() {
				nilUserIDRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "display_name", "user_id", "twitter_handle"}).
					AddRow(profileObj.ID, profileObj.CreatedAt, profileObj.UpdatedAt, profileObj.DeletedAt, profileObj.DisplayName, profileObj.UserID, profileObj.TwitterHandle)

				mock.ExpectQuery(expectedFindQueryStr).WithArgs(profileObj.ID).WillReturnRows(nilUserIDRow)
				mock.ExpectBegin()
				mock.ExpectExec(expectedUpdateStr).WithArgs(
					profileObj.DisplayName, profileObj.TwitterHandle, sqlmock.AnyArg(), profileObj.ID,
				).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(expectedUpdateSocialStr).WithArgs(
					socialObj.CreatedAt, sqlmock.AnyArg(), socialObj.DeletedAt, socialObj.AccessToken,
					socialObj.AccessTokenSecret, socialObj.UserID, socialObj.ProfileImageURL, socialObj.ScreenName,
					socialObj.IsVerified, socialObj.Network, socialObj.UserProfileID,
				).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedPostUpdateSelectStr).WithArgs(profileObj.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "display_name", "user_id", "twitter_handle"}).
						AddRow(profileObj.ID, profileObj.CreatedAt, profileObj.UpdatedAt, profileObj.DeletedAt,
							profileObj.DisplayName, profileObj.UserID, profileObj.TwitterHandle))
				mock.ExpectQuery(expectedPostUpdateSocialSelectStr).WithArgs(profileObj.ID).
					WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "user_profile_id", "network",
						"access_token", "access_token_secret", "user_id", "profile_image_url", "screen_name", "is_verified"}).
						AddRow(socialObj.CreatedAt, socialObj.UpdatedAt, socialObj.DeletedAt, socialObj.UserProfileID,
							socialObj.Network, socialObj.AccessToken, socialObj.AccessTokenSecret, socialObj.UserID,
							socialObj.ProfileImageURL, socialObj.ScreenName, socialObj.IsVerified))

				resp, inserted, err := mockDatastore.CreateOrUpdateUserProfile(ctx, profileObj)

				So(err, ShouldBeNil)
				So(inserted, ShouldBeFalse)
				So(resp, ShouldNotBeNil)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})
		})
	})
}

func Test_MarshalUserProfile(t *testing.T) {
	// type args struct {
	// 	userProfile model.UserProfile
	// }

	// haveUserProfileObj := model.UserProfile{
	// 	ID:                     "11111",
	// 	DisplayName:            "Delphis Hello",
	// 	UserID:                 "22222",
	// 	TwitterHandle:          "delphishq",
	// 	ModeratedDiscussionIDs: []string{"33333"},
	// 	ModeratedDiscussions:   []model.Discussion{model.Discussion{}},
	// 	TwitterInfo: model.SocialInfo{
	// 		AccessToken:       "44444",
	// 		AccessTokenSecret: "55555",
	// 		UserID:            "55555",
	// 		ProfileImageURL:   "https://a.b/c.png",
	// 		ScreenName:        "delphishq",
	// 		IsVerified:        true,
	// 	},
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
	// 			userProfile: haveUserProfileObj,
	// 		},
	// 		want: map[string]*dynamodb.AttributeValue{
	// 			"ID": {
	// 				S: aws.String(haveUserProfileObj.ID),
	// 			},
	// 			"DisplayName": {
	// 				S: aws.String(haveUserProfileObj.DisplayName),
	// 			},
	// 			"UserID": {
	// 				S: aws.String(haveUserProfileObj.UserID),
	// 			},
	// 			"TwitterHandle": {
	// 				S: aws.String(haveUserProfileObj.TwitterHandle),
	// 			},
	// 			"ModeratedDiscussionIDs": {
	// 				SS: []*string{aws.String(haveUserProfileObj.ModeratedDiscussionIDs[0])},
	// 			},
	// 			"TwitterInfo": {
	// 				M: map[string]*dynamodb.AttributeValue{
	// 					"AccessToken": {
	// 						S: aws.String(haveUserProfileObj.TwitterInfo.AccessToken),
	// 					},
	// 					"AccessTokenSecret": {
	// 						S: aws.String(haveUserProfileObj.TwitterInfo.AccessTokenSecret),
	// 					},
	// 					"SocialUserID": {
	// 						S: aws.String(haveUserProfileObj.TwitterInfo.UserID),
	// 					},
	// 					"ProfileImageURL": {
	// 						S: aws.String(haveUserProfileObj.TwitterInfo.ProfileImageURL),
	// 					},
	// 					"ScreenName": {
	// 						S: aws.String(haveUserProfileObj.TwitterInfo.ScreenName),
	// 					},
	// 					"IsVerified": {
	// 						BOOL: aws.Bool(haveUserProfileObj.TwitterInfo.IsVerified),
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		marshaled, err := datastoreObj.marshalMap(tt.args.userProfile)
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

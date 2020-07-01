package datastore

import (
	"context"
	"fmt"
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

func Test_GetUserDevicesByUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	userID := "54321"
	userObject := model.UserDevice{
		ID:        "12345",
		CreatedAt: now,
		Platform:  "ios",
		LastSeen:  now,
		Token:     nil,
		UserID:    &userID,
	}

	Convey("GetUserDevicesByUserID", t, func() {
		db, mock, err := sqlmock.New()

		assert.Nil(t, err, "Failed setting up sqlmock db")

		gormDB, _ := gorm.Open("postgres", db)
		mockDatastore := &delphisDB{
			dbConfig: config.TablesConfig{},
			sql:      gormDB,
			dynamo:   nil,
			encoder:  nil,
		}
		defer db.Close()

		expectedQueryStr := `SELECT * FROM "user_devices"  WHERE "user_devices"."deleted_at" IS NULL AND (("user_devices"."user_id" = $1)) ORDER BY last_seen desc`

		Convey("when record is not found", func() {
			expectedError := gorm.ErrRecordNotFound
			mock.ExpectQuery(regexp.QuoteMeta(expectedQueryStr)).WithArgs(userObject.ID).WillReturnError(expectedError)

			resp, err := mockDatastore.GetUserDevicesByUserID(ctx, userObject.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when an error occurs", func() {
			expectedError := fmt.Errorf("Something went wrong")
			mock.ExpectQuery(regexp.QuoteMeta(expectedQueryStr)).WithArgs(userObject.ID).WillReturnError(expectedError)

			resp, err := mockDatastore.GetUserDevicesByUserID(ctx, userObject.ID)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the objects are found", func() {
			mock.ExpectQuery(regexp.QuoteMeta(expectedQueryStr)).WithArgs(userObject.ID).WillReturnRows(
				sqlmock.NewRows([]string{"id", "created_at", "deleted_at", "platform", "last_seen", "token", "user_id"}).
					AddRow(userObject.ID, userObject.CreatedAt, userObject.DeletedAt, userObject.Platform, userObject.LastSeen, userObject.Token, userObject.UserID).
					AddRow(userObject.ID, userObject.CreatedAt, userObject.DeletedAt, userObject.Platform, userObject.LastSeen, userObject.Token, userObject.UserID))

			resp, err := mockDatastore.GetUserDevicesByUserID(ctx, userObject.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []model.UserDevice{userObject, userObject})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func Test_UpsertUserDevice(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	userID := "54321"
	newObject := model.UserDevice{
		ID:        "12345",
		CreatedAt: now,
		Platform:  "ios",
		LastSeen:  now,
		Token:     nil,
		UserID:    &userID,
	}

	Convey("UpsertUserDevice", t, func() {
		db, mock, err := sqlmock.New()

		assert.Nil(t, err, "Failed setting up sqlmock db")

		gormDB, _ := gorm.Open("postgres", db)
		mockDatastore := &delphisDB{
			dbConfig: config.TablesConfig{},
			sql:      gormDB,
			dynamo:   nil,
			encoder:  nil,
		}
		defer db.Close()

		expectedFindQueryStr := `SELECT * FROM "user_devices" WHERE "user_devices"."deleted_at" IS NULL AND (("user_devices"."id" = $1)) ORDER BY "user_devices"."id" ASC LIMIT 1`
		createQueryStr := `INSERT INTO "user_devices" ("id","created_at","deleted_at","platform","last_seen","token","user_id") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "user_devices"."id"`
		expectedNewObjectRow := sqlmock.NewRows([]string{"id", "created_at", "deleted_at", "platform", "last_seen", "token", "user_id"}).
			AddRow(newObject.ID, newObject.CreatedAt, newObject.DeletedAt, newObject.Platform, newObject.LastSeen, newObject.Token, newObject.UserID)
		expectedUpdateStr := `UPDATE "user_devices" SET "last_seen" = $1, "user_id" = $2  WHERE "user_devices"."deleted_at" IS NULL AND "user_devices"."id" = $3`
		expectedPostUpdateSelectStr := `SELECT * FROM "user_devices" WHERE "user_devices"."deleted_at" IS NULL AND "user_devices"."id" = $1 ORDER BY "user_devices"."id" ASC LIMIT 1`

		Convey("when find query fails with a non-not-found-error the function should return the error", func() {
			expectedError := fmt.Errorf("Some fake error")
			mock.ExpectQuery(regexp.QuoteMeta(expectedFindQueryStr)).WithArgs(newObject.ID).WillReturnError(expectedError)

			resp, err := mockDatastore.UpsertUserDevice(ctx, newObject)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when find query returns not-found-error it should call create", func() {
			Convey("when create returns an error it should return it", func() {
				expectedError := fmt.Errorf("Some fake error")

				mock.ExpectQuery(regexp.QuoteMeta(expectedFindQueryStr)).WithArgs(newObject.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(createQueryStr)).WithArgs(
					newObject.ID, newObject.CreatedAt, newObject.DeletedAt, newObject.Platform, newObject.LastSeen, newObject.Token, newObject.UserID,
				).WillReturnError(expectedError)

				resp, err := mockDatastore.UpsertUserDevice(ctx, newObject)

				assert.NotNil(t, err)
				assert.Nil(t, resp)
				assert.Nil(t, mock.ExpectationsWereMet())
			})

			Convey("when create succeeds it should return the new object", func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectedFindQueryStr)).WithArgs(newObject.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(createQueryStr)).WithArgs(
					newObject.ID, newObject.CreatedAt, newObject.DeletedAt, newObject.Platform, newObject.LastSeen, newObject.Token, newObject.UserID,
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newObject.ID))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta(expectedFindQueryStr)).WithArgs(newObject.ID).WillReturnRows(expectedNewObjectRow)

				resp, err := mockDatastore.UpsertUserDevice(ctx, newObject)

				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, *resp, newObject)
				assert.Nil(t, mock.ExpectationsWereMet())
			})
		})

		Convey("when the object is found we should update it", func() {
			Convey("when updating if it fails then return the error", func() {
				expectedError := fmt.Errorf("something went wrong")
				mock.ExpectQuery(regexp.QuoteMeta(expectedFindQueryStr)).WithArgs(newObject.ID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(expectedUpdateStr)).WithArgs(
					newObject.LastSeen, *newObject.UserID, newObject.ID,
				).WillReturnError(expectedError)

				resp, err := mockDatastore.UpsertUserDevice(ctx, newObject)

				assert.NotNil(t, err)
				assert.Nil(t, resp)
				assert.Nil(t, mock.ExpectationsWereMet())
			})

			Convey("when updating if it succeeds it should return the updated object", func() {
				newObject.LastSeen = time.Time{}
				mock.ExpectQuery(regexp.QuoteMeta(expectedFindQueryStr)).WithArgs(newObject.ID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(expectedUpdateStr)).WithArgs(
					sqlmock.AnyArg(), *newObject.UserID, newObject.ID,
				).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(regexp.QuoteMeta(expectedPostUpdateSelectStr)).WithArgs(newObject.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "deleted_at", "platform", "last_seen", "token", "user_id"}).
					AddRow(newObject.ID, newObject.CreatedAt, newObject.DeletedAt, newObject.Platform, newObject.LastSeen, newObject.Token, newObject.UserID))

				resp, err := mockDatastore.UpsertUserDevice(ctx, newObject)

				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Nil(t, mock.ExpectationsWereMet())
			})
		})
	})
}

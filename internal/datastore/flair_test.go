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

func TestDelphisDB_UpsertFlair(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	flairID := "flairID"
	templateID := "templateID"
	userID := "userID"
	flairObj := model.Flair{
		ID:         flairID,
		TemplateID: templateID,
		CreatedAt:  now,
		UpdatedAt:  now,
		UserID:     userID,
	}

	Convey("UpsertUserDevice", t, func() {
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

		expectedFindQueryStr := `SELECT * FROM "flairs" WHERE "flairs"."deleted_at" IS NULL AND (("flairs"."id" = $1)) ORDER BY "flairs"."id" ASC LIMIT 1`
		createQueryStr := `INSERT INTO "flairs" ("id","template_id","created_at","updated_at","deleted_at","user_id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "flairs"."id"`
		expectedNewObjectRow := sqlmock.NewRows([]string{"id", "template_id", "created_at", "updated_at", "deleted_at", "user_id"}).
			AddRow(flairObj.ID, flairObj.TemplateID, flairObj.CreatedAt, flairObj.UpdatedAt, flairObj.DeletedAt, flairObj.UserID)
		expectedUpdateStr := `UPDATE "flairs" SET "template_id" = $1, "updated_at" = $2, "user_id" = $3 WHERE "flairs"."deleted_at" IS NULL AND "flairs"."id" = $4`
		expectedPostUpdateSelectStr := `SELECT * FROM "flairs" WHERE "flairs"."deleted_at" IS NULL AND "flairs"."id" = $1 ORDER BY "flairs"."id" ASC LIMIT 1`

		Convey("when find query fails with a non-not-found-error the function should return the error", func() {
			expectedError := fmt.Errorf("Some fake error")
			mock.ExpectQuery(expectedFindQueryStr).WithArgs(flairObj.ID).WillReturnError(expectedError)

			resp, err := mockDatastore.UpsertFlair(ctx, flairObj)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when find query returns not-found-error it should call create", func() {
			Convey("when create returns an error it should return it", func() {
				expectedError := fmt.Errorf("Some fake error")

				mock.ExpectQuery(expectedFindQueryStr).WithArgs(flairObj.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					flairObj.ID, flairObj.TemplateID, flairObj.CreatedAt, flairObj.UpdatedAt, flairObj.DeletedAt, flairObj.UserID,
				).WillReturnError(expectedError)

				resp, err := mockDatastore.UpsertFlair(ctx, flairObj)

				assert.NotNil(t, err)
				assert.Nil(t, resp)
				assert.Nil(t, mock.ExpectationsWereMet())
			})

			Convey("when create succeeds it should return the new object", func() {
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(flairObj.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					flairObj.ID, flairObj.TemplateID, flairObj.CreatedAt, flairObj.UpdatedAt, flairObj.DeletedAt, flairObj.UserID,
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(flairObj.ID))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(flairObj.ID).WillReturnRows(expectedNewObjectRow)

				resp, err := mockDatastore.UpsertFlair(ctx, flairObj)

				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, *resp, flairObj)
				assert.Nil(t, mock.ExpectationsWereMet())
			})
		})

		Convey("when the object is found we should update it", func() {
			Convey("when updating if it fails then return the error", func() {
				expectedError := fmt.Errorf("something went wrong")
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(flairObj.ID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(expectedUpdateStr).WithArgs(
					flairObj.TemplateID, sqlmock.AnyArg(), flairObj.UserID, flairObj.ID,
				).WillReturnError(expectedError)

				resp, err := mockDatastore.UpsertFlair(ctx, flairObj)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})

			Convey("when updating if it succeeds it should return the updated object", func() {
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(flairObj.ID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(expectedUpdateStr).WithArgs(
					flairObj.TemplateID, sqlmock.AnyArg(), flairObj.UserID, flairObj.ID,
				).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedPostUpdateSelectStr).WithArgs(flairObj.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "template_id", "created_at", "updated_at", "deleted_at", "user_id"}).
						AddRow(flairObj.ID, flairObj.TemplateID, flairObj.CreatedAt, flairObj.UpdatedAt, flairObj.DeletedAt, flairObj.UserID))

				resp, err := mockDatastore.UpsertFlair(ctx, flairObj)

				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Nil(t, mock.ExpectationsWereMet())
			})
		})
	})
}

func TestDelphisDB_GetFlairByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	flairID := "flairID"
	templateID := "templateID"
	userID := "userID"
	flairObj := model.Flair{
		ID:         flairID,
		TemplateID: templateID,
		CreatedAt:  now,
		UpdatedAt:  now,
		UserID:     userID,
	}

	Convey("GetFlairByID", t, func() {
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

		expectedQueryStr := `SELECT * FROM "flairs" WHERE "flairs"."deleted_at" IS NULL AND (("flairs"."id" = $1)) ORDER BY "flairs"."id" ASC LIMIT 1`

		Convey("when record is not found", func() {
			expectedError := gorm.ErrRecordNotFound
			mock.ExpectQuery(expectedQueryStr).WithArgs(flairObj.ID).WillReturnError(expectedError)

			resp, err := mockDatastore.GetFlairByID(ctx, flairID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when an error occurs", func() {
			expectedError := fmt.Errorf("Something went wrong")
			mock.ExpectQuery(expectedQueryStr).WithArgs(flairObj.ID).WillReturnError(expectedError)

			resp, err := mockDatastore.GetFlairByID(ctx, flairID)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the objects are found", func() {
			mock.ExpectQuery(expectedQueryStr).WithArgs(flairObj.ID).WillReturnRows(
				sqlmock.NewRows([]string{"id", "user_id", "template_id", "created_at", "updated_at", "deleted_at"}).
					AddRow(flairObj.ID, flairObj.UserID, flairObj.TemplateID, flairObj.CreatedAt, flairObj.UpdatedAt, flairObj.DeletedAt))

			resp, err := mockDatastore.GetFlairByID(ctx, flairID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &flairObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetFlairsByUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	flairID := "flairID"
	templateID := "templateID"
	userID := "userID"
	flairObj := model.Flair{
		ID:         flairID,
		TemplateID: templateID,
		CreatedAt:  now,
		UpdatedAt:  now,
		UserID:     userID,
	}

	Convey("GetFlairsByUserID", t, func() {
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

		expectedQueryStr := `SELECT * FROM "flairs" WHERE "flairs"."deleted_at" IS NULL AND (("flairs"."user_id" = $1))`

		Convey("when record is not found", func() {
			expectedError := gorm.ErrRecordNotFound
			mock.ExpectQuery(expectedQueryStr).WithArgs(flairObj.UserID).WillReturnError(expectedError)

			resp, err := mockDatastore.GetFlairsByUserID(ctx, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.Flair{})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when an error occurs", func() {
			expectedError := fmt.Errorf("Something went wrong")
			mock.ExpectQuery(expectedQueryStr).WithArgs(flairObj.UserID).WillReturnError(expectedError)

			resp, err := mockDatastore.GetFlairsByUserID(ctx, userID)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldResemble, []*model.Flair{})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the objects are found", func() {
			mock.ExpectQuery(expectedQueryStr).WithArgs(flairObj.UserID).WillReturnRows(
				sqlmock.NewRows([]string{"id", "user_id", "template_id", "created_at", "updated_at", "deleted_at"}).
					AddRow(flairObj.ID, flairObj.UserID, flairObj.TemplateID, flairObj.CreatedAt, flairObj.UpdatedAt, flairObj.DeletedAt).
					AddRow(flairObj.ID, flairObj.UserID, flairObj.TemplateID, flairObj.CreatedAt, flairObj.UpdatedAt, flairObj.DeletedAt))

			resp, err := mockDatastore.GetFlairsByUserID(ctx, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []*model.Flair{&flairObj, &flairObj})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_RemoveFlair(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	flairID := "flairID"
	templateID := "templateID"
	userID := "userID"
	flairObj := model.Flair{
		ID:         flairID,
		TemplateID: templateID,
		CreatedAt:  now,
		UpdatedAt:  now,
		UserID:     userID,
	}

	Convey("RemoveFlair", t, func() {
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

		expectedUpdateStr := `UPDATE "participants" SET "flair_id" = $1 WHERE ("participants"."flair_id" = $2)`
		expectedDeleteStr := `UPDATE "flairs" SET "deleted_at"=$1 WHERE "flairs"."deleted_at" IS NULL AND "flairs"."id" = $2`

		Convey("when record does not have an ID", func() {

			resp, err := mockDatastore.RemoveFlair(ctx, model.Flair{})

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &model.Flair{})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when an error occurs on update", func() {
			expectedError := fmt.Errorf("Something went wrong")
			mock.ExpectBegin()
			mock.ExpectExec(expectedUpdateStr).WithArgs().WillReturnError(expectedError)
			mock.ExpectRollback()

			resp, err := mockDatastore.RemoveFlair(ctx, flairObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldResemble, &flairObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when an error occurs", func() {
			expectedError := fmt.Errorf("Something went wrong")
			mock.ExpectBegin()
			mock.ExpectExec(expectedUpdateStr).WithArgs().WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectExec(expectedDeleteStr).WithArgs(sqlmock.AnyArg(), flairObj.ID).WillReturnError(expectedError)
			mock.ExpectRollback()

			resp, err := mockDatastore.RemoveFlair(ctx, flairObj)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldResemble, &flairObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the objects are found", func() {
			mock.ExpectBegin()
			mock.ExpectExec(expectedUpdateStr).WithArgs().WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectExec(expectedDeleteStr).WithArgs(sqlmock.AnyArg(), flairObj.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			resp, err := mockDatastore.RemoveFlair(ctx, flairObj)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &flairObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

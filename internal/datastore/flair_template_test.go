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

func TestDelphisDB_ListFlairTemplates(t *testing.T) {
	ctx := context.Background()
	ftID := "flairTemplateID"
	displayName := "displayName"
	imageURL := "imageURL"
	source := "test"
	ftObj := model.FlairTemplate{
		ID:          ftID,
		DisplayName: &displayName,
		ImageURL:    &imageURL,
		Source:      source,
	}

	Convey("ListFlairTemplates", t, func() {
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

		expectedQueryStr := `SELECT * FROM "flair_templates" WHERE "flair_templates"."deleted_at" IS NULL`
		expectedStrWithQuery := `SELECT * FROM "flair_templates" WHERE "flair_templates"."deleted_at" IS NULL AND ((source ILIKE $1) OR (display_name ILIKE $2))`

		Convey("when an error occurs", func() {
			expectedError := fmt.Errorf("Something went wrong")
			mock.ExpectQuery(expectedQueryStr).WillReturnError(expectedError)

			resp, err := mockDatastore.ListFlairTemplates(ctx, nil)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when a query is passed in and returns results", func() {
			query := "test"
			likeQuery := "%test%"
			mock.ExpectQuery(expectedStrWithQuery).WithArgs(likeQuery, likeQuery).
				WillReturnRows(sqlmock.NewRows([]string{"id", "display_name", "image_url", "source", "created_at", "updated_at", "deleted_at"}).
					AddRow(ftObj.ID, ftObj.DisplayName, ftObj.ImageURL, ftObj.Source, ftObj.CreatedAt, ftObj.UpdatedAt, ftObj.DeletedAt))

			resp, err := mockDatastore.ListFlairTemplates(ctx, &query)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []*model.FlairTemplate{&ftObj})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_UpsertFlairTemplate(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	ftID := "flairTemplateID"
	displayName := "displayName"
	imageURL := "imageURL"
	source := "test"
	ftObj := model.FlairTemplate{
		ID:          ftID,
		DisplayName: &displayName,
		ImageURL:    &imageURL,
		Source:      source,
		CreatedAt:   now,
		UpdatedAt:   now,
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

		expectedFindQueryStr := `SELECT * FROM "flair_templates" WHERE "flair_templates"."deleted_at" IS NULL AND (("flair_templates"."id" = $1)) ORDER BY "flair_templates"."id" ASC LIMIT 1`
		createQueryStr := `INSERT INTO "flair_templates" ("id","display_name","image_url","source","created_at","updated_at","deleted_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "flair_templates"."id"`
		expectedNewObjectRow := sqlmock.NewRows([]string{"id", "display_name", "image_url", "source", "created_at", "updated_at", "deleted_at"}).
			AddRow(ftObj.ID, ftObj.DisplayName, ftObj.ImageURL, ftObj.Source, ftObj.CreatedAt, ftObj.UpdatedAt, ftObj.DeletedAt)
		expectedUpdateStr := `UPDATE "flair_templates" SET "display_name" = $1, "image_url" = $2, "source" = $3, "updated_at" = $4 WHERE "flair_templates"."deleted_at" IS NULL AND "flair_templates"."id" = $5`
		expectedPostUpdateSelectStr := `SELECT * FROM "flair_templates" WHERE "flair_templates"."deleted_at" IS NULL AND "flair_templates"."id" = $1 ORDER BY "flair_templates"."id" ASC LIMIT 1`

		Convey("when find query fails with a non-not-found-error the function should return the error", func() {
			expectedError := fmt.Errorf("Some fake error")
			mock.ExpectQuery(expectedFindQueryStr).WithArgs(ftObj.ID).WillReturnError(expectedError)

			resp, err := mockDatastore.UpsertFlairTemplate(ctx, ftObj)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when find query returns not-found-error it should call create", func() {
			Convey("when create returns an error it should return it", func() {
				expectedError := fmt.Errorf("Some fake error")

				mock.ExpectQuery(expectedFindQueryStr).WithArgs(ftObj.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					ftObj.ID, ftObj.DisplayName, ftObj.ImageURL, ftObj.Source, ftObj.CreatedAt, ftObj.UpdatedAt, ftObj.DeletedAt,
				).WillReturnError(expectedError)

				resp, err := mockDatastore.UpsertFlairTemplate(ctx, ftObj)

				assert.NotNil(t, err)
				assert.Nil(t, resp)
				assert.Nil(t, mock.ExpectationsWereMet())
			})

			Convey("when create succeeds it should return the new object", func() {
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(ftObj.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					ftObj.ID, ftObj.DisplayName, ftObj.ImageURL, ftObj.Source, ftObj.CreatedAt, ftObj.UpdatedAt, ftObj.DeletedAt,
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(ftObj.ID))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(ftObj.ID).WillReturnRows(expectedNewObjectRow)

				resp, err := mockDatastore.UpsertFlairTemplate(ctx, ftObj)

				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, *resp, ftObj)
				assert.Nil(t, mock.ExpectationsWereMet())
			})
		})

		Convey("when the object is found we should update it", func() {
			Convey("when updating if it fails then return the error", func() {
				expectedError := fmt.Errorf("something went wrong")
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(ftObj.ID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(expectedUpdateStr).WithArgs(
					ftObj.DisplayName, ftObj.ImageURL, ftObj.Source, sqlmock.AnyArg(), ftObj.ID,
				).WillReturnError(expectedError)

				resp, err := mockDatastore.UpsertFlairTemplate(ctx, ftObj)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})

			Convey("when updating if it succeeds it should return the updated object", func() {
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(ftObj.ID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(expectedUpdateStr).WithArgs(
					ftObj.DisplayName, ftObj.ImageURL, ftObj.Source, sqlmock.AnyArg(), ftObj.ID,
				).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedPostUpdateSelectStr).WithArgs(ftObj.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "display_name", "image_url", "source", "created_at", "updated_at", "deleted_at"}).
						AddRow(ftObj.ID, ftObj.DisplayName, ftObj.ImageURL, ftObj.Source, ftObj.CreatedAt, ftObj.UpdatedAt, ftObj.DeletedAt))

				resp, err := mockDatastore.UpsertFlairTemplate(ctx, ftObj)

				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Nil(t, mock.ExpectationsWereMet())
			})
		})
	})
}

func TestDelphisDB_GetFlairTemplateByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	ftID := "flairTemplateID"
	displayName := "displayName"
	imageURL := "imageURL"
	source := "test"
	ftObj := model.FlairTemplate{
		ID:          ftID,
		DisplayName: &displayName,
		ImageURL:    &imageURL,
		Source:      source,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	Convey("GetFlairTemplateByID", t, func() {
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

		expectedQueryStr := `SELECT * FROM "flair_templates" WHERE "flair_templates"."deleted_at" IS NULL AND (("flair_templates"."id" = $1)) ORDER BY "flair_templates"."id" ASC LIMIT 1`

		Convey("when record is not found", func() {
			expectedError := gorm.ErrRecordNotFound
			mock.ExpectQuery(expectedQueryStr).WithArgs(ftObj.ID).WillReturnError(expectedError)

			resp, err := mockDatastore.GetFlairTemplateByID(ctx, ftObj.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when an error occurs", func() {
			expectedError := fmt.Errorf("Something went wrong")
			mock.ExpectQuery(expectedQueryStr).WithArgs(ftObj.ID).WillReturnError(expectedError)

			resp, err := mockDatastore.GetFlairTemplateByID(ctx, ftObj.ID)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the objects are found", func() {
			mock.ExpectQuery(expectedQueryStr).WithArgs(ftObj.ID).WillReturnRows(
				sqlmock.NewRows([]string{"id", "display_name", "image_url", "source", "created_at", "updated_at", "deleted_at"}).
					AddRow(ftObj.ID, ftObj.DisplayName, ftObj.ImageURL, ftObj.Source, ftObj.CreatedAt, ftObj.UpdatedAt, ftObj.DeletedAt))

			resp, err := mockDatastore.GetFlairTemplateByID(ctx, ftObj.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &ftObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_RemoveFlairTemplate(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	ftID := "flairTemplateID"
	displayName := "displayName"
	imageURL := "imageURL"
	source := "test"
	ftObj := model.FlairTemplate{
		ID:          ftID,
		DisplayName: &displayName,
		ImageURL:    &imageURL,
		Source:      source,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	Convey("RemoveFlairTemplate", t, func() {
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

		expectedUpdateStr := `UPDATE "participants" SET "flair_id" = $1, "updated_at" = $2 FROM flairs WHERE (flair_id = flairs.id AND flairs.template_id = $3)`
		expectedDeleteFlairStr := `UPDATE "flairs" SET "deleted_at"=$1 WHERE "flairs"."deleted_at" IS NULL AND (("flairs"."template_id" = $2))`
		expectedDeleteFlairTemplateStr := `UPDATE "flair_templates" SET "deleted_at"=$1  WHERE "flair_templates"."deleted_at" IS NULL AND "flair_templates"."id" = $2`

		Convey("when record does not have an ID", func() {

			resp, err := mockDatastore.RemoveFlairTemplate(ctx, model.FlairTemplate{})

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &model.FlairTemplate{})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when an error occurs on update", func() {
			expectedError := fmt.Errorf("Something went wrong")
			mock.ExpectBegin()
			mock.ExpectExec(expectedUpdateStr).WithArgs().WillReturnError(expectedError)
			mock.ExpectRollback()

			resp, err := mockDatastore.RemoveFlairTemplate(ctx, ftObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldResemble, &ftObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when an error occurs when deleting flairs", func() {
			expectedError := fmt.Errorf("Something went wrong")
			mock.ExpectBegin()
			mock.ExpectExec(expectedUpdateStr).WithArgs().WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectExec(expectedDeleteFlairStr).WithArgs(sqlmock.AnyArg(), ftObj.ID).WillReturnError(expectedError)
			mock.ExpectRollback()

			resp, err := mockDatastore.RemoveFlairTemplate(ctx, ftObj)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldResemble, &ftObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the objects are found and an error occurs when deleting the template", func() {
			expectedError := fmt.Errorf("Something went wrong")
			mock.ExpectBegin()
			mock.ExpectExec(expectedUpdateStr).WithArgs().WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectExec(expectedDeleteFlairStr).WithArgs(sqlmock.AnyArg(), ftObj.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(expectedDeleteFlairTemplateStr).WithArgs(sqlmock.AnyArg(), ftObj.ID).WillReturnError(expectedError)
			mock.ExpectRollback()

			resp, err := mockDatastore.RemoveFlairTemplate(ctx, ftObj)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldResemble, &ftObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the objects are found and successfully deleted", func() {
			mock.ExpectBegin()
			mock.ExpectExec(expectedUpdateStr).WithArgs().WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectExec(expectedDeleteFlairStr).WithArgs(sqlmock.AnyArg(), ftObj.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(expectedDeleteFlairTemplateStr).WithArgs(sqlmock.AnyArg(), ftObj.ID).WillReturnResult(sqlmock.NewResult(2, 1))
			mock.ExpectCommit()

			resp, err := mockDatastore.RemoveFlairTemplate(ctx, ftObj)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &ftObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

	})
}

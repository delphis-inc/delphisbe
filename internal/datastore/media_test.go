package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestDelphisDB_GetMediaRecordByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	mediaID := "media"
	mediaType := "video"
	mediaObj := model.Media{
		ID:        mediaID,
		CreatedAt: now.Format(time.RFC3339),
		MediaType: &mediaType,
		MediaSize: &model.MediaSize{
			Height: 100,
			Width:  100,
			SizeKb: 100,
		},
	}

	sizeJson, _ := json.Marshal(mediaObj.MediaSize)

	Convey("GetMediaRecordByID", t, func() {
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

			resp, err := mockDatastore.GetMediaRecordByID(ctx, mediaID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getMediaRecordString).WithArgs(mediaID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetMediaRecordByID(ctx, mediaID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and the unmarshalling fails", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"id", "created_at", "deleted_at", "deleted_reason_code", "media_type", "media_size"}).
				AddRow(mediaObj.ID, mediaObj.CreatedAt, nil, mediaObj.DeletedReasonCode, mediaObj.MediaType, nil)

			mock.ExpectQuery(getMediaRecordString).WithArgs(mediaID).WillReturnRows(rs)

			resp, err := mockDatastore.GetMediaRecordByID(ctx, mediaID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and the post has been deleted", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"id", "created_at", "deleted_at", "deleted_reason_code", "media_type", "media_size"}).
				AddRow(mediaObj.ID, mediaObj.CreatedAt, &now, mediaObj.DeletedReasonCode, mediaObj.MediaType, sizeJson)

			testMediaObj := mediaObj
			testMediaObj.IsDeleted = true

			mock.ExpectQuery(getMediaRecordString).WithArgs(mediaID).WillReturnRows(rs)

			resp, err := mockDatastore.GetMediaRecordByID(ctx, mediaID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &testMediaObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a post", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"id", "created_at", "deleted_at", "deleted_reason_code", "media_type", "media_size"}).
				AddRow(mediaObj.ID, mediaObj.CreatedAt, nil, mediaObj.DeletedReasonCode, mediaObj.MediaType, sizeJson)

			mock.ExpectQuery(getMediaRecordString).WithArgs(mediaID).WillReturnRows(rs)

			resp, err := mockDatastore.GetMediaRecordByID(ctx, mediaID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &mediaObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_PutMediaRecord(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	mediaID := "media"
	mediaType := "video"
	assetLoc := "video/1"
	mediaObj := model.Media{
		ID:        mediaID,
		CreatedAt: now.Format(time.RFC3339),
		MediaType: &mediaType,
		MediaSize: &model.MediaSize{
			Height: 100,
			Width:  100,
			SizeKb: 100,
		},
		AssetLocation: &assetLoc,
	}

	sizeJson, _ := json.Marshal(mediaObj.MediaSize)

	Convey("PutMediaRecord", t, func() {
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
			mock.ExpectBegin()
			mockPreparedStatementsWithError(mock)

			tx, err := mockDatastore.BeginTx(ctx)
			err = mockDatastore.PutMediaRecord(ctx, tx, mediaObj)

			So(err, ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when json marshalling errors out", func() {
			testMediaObj := mediaObj
			testMediaSize := *testMediaObj.MediaSize
			testMediaSize.SizeKb = math.NaN()
			testMediaObj.MediaSize = &testMediaSize

			mock.ExpectBegin()
			mockPreparedStatements(mock)

			tx, err := mockDatastore.BeginTx(ctx)
			err = mockDatastore.PutMediaRecord(ctx, tx, testMediaObj)

			So(err, ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putMediaRecordString)
			mock.ExpectExec(putMediaRecordString).WithArgs(mediaObj.ID, mediaObj.MediaType,
				sizeJson).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			err = mockDatastore.PutMediaRecord(ctx, tx, mediaObj)

			So(err, ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when put post succeeds and returns an object", func() {

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putMediaRecordString)
			mock.ExpectExec(putMediaRecordString).WithArgs(mediaObj.ID, mediaObj.MediaType,
				sizeJson).WillReturnResult(sqlmock.NewResult(0, 0))

			tx, err := mockDatastore.BeginTx(ctx)
			err = mockDatastore.PutMediaRecord(ctx, tx, mediaObj)

			So(err, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

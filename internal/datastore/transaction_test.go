package datastore

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/internal/config"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestDelphisDB_BeginTx(t *testing.T) {
	ctx := context.Background()

	Convey("BeginTx", t, func() {
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

		Convey("when being transaction fails", func() {
			expectedError := fmt.Errorf("Some fake error")
			mock.ExpectBegin().WillReturnError(expectedError)

			tx, err := mockDatastore.BeginTx(ctx)

			So(err, ShouldNotBeNil)
			So(tx, ShouldBeNil)
		})

		Convey("when being transaction succeeds", func() {
			mock.ExpectBegin()

			tx, err := mockDatastore.BeginTx(ctx)

			So(err, ShouldBeNil)
			So(tx, ShouldNotBeNil)
		})
	})
}

func TestDelphisDB_RollbackTx_CommitTx(t *testing.T) {
	ctx := context.Background()

	Convey("RollbackTx", t, func() {
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

		Convey("when RollbackTx is called", func() {
			mock.ExpectBegin()
			mock.ExpectRollback()

			tx, err := mockDatastore.BeginTx(ctx)
			err = mockDatastore.RollbackTx(ctx, tx)

			So(err, ShouldBeNil)
			So(tx, ShouldNotBeNil)
		})
	})
}

func TestDelphisDB_CommitTx(t *testing.T) {
	ctx := context.Background()

	Convey("CommitTx", t, func() {
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

		Convey("when CommitTx is called", func() {
			mock.ExpectBegin()
			mock.ExpectCommit()

			tx, err := mockDatastore.BeginTx(ctx)
			err = mockDatastore.CommitTx(ctx, tx)

			So(err, ShouldBeNil)
			So(tx, ShouldNotBeNil)
		})
	})
}

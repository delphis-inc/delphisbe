package datastore

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestDelphisDB_InitializeStatements(t *testing.T) {
	ctx := context.Background()

	Convey("InitializeStatments", t, func() {
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

		Convey("when prepared statement fails to initialize", func() {
			mock.ExpectPrepare(getPostByIDString)
			mock.ExpectPrepare(getPostsByDiscussionIDFromCursorString).WillReturnError(fmt.Errorf("error"))
			err := mockDatastore.initializeStatements(ctx)

			So(err, ShouldNotBeNil)
		})

		tests := []struct {
			Name string
			Test string
		}{

			{
				Name: "getPostByIDString",
				Test: getPostByIDString,
			},
			{
				Name: "getPostsByDiscussionIDString",
				Test: getPostsByDiscussionIDString,
			},
			{
				Name: "getLastPostByDiscussionIDStmt",
				Test: getLastPostByDiscussionIDStmt,
			},
			{
				Name: "getPostsByDiscussionIDFromCursorString",
				Test: getPostsByDiscussionIDFromCursorString,
			},
			{
				Name: "putPostString",
				Test: putPostString,
			},
			{
				Name: "putPostContentsString",
				Test: putPostContentsString,
			},
			{
				Name: "deletePostByIDString",
				Test: deletePostByIDString,
			},
			{
				Name: "deletePostByParticipantIDDiscussionIDString",
				Test: deletePostByParticipantIDDiscussionIDString,
			},
			{
				Name: "putActivityString",
				Test: putActivityString,
			},
			{
				Name: "putMediaRecordString",
				Test: putMediaRecordString,
			},
			{
				Name: "getMediaRecordString",
				Test: getMediaRecordString,
			},
			{
				Name: "getDiscussionByLinkSlugString",

				Test: getDiscussionByLinkSlugString,
			},
			{
				Name: "getDiscussionArchiveByDiscussionIDString",

				Test: getDiscussionArchiveByDiscussionIDString,
			},
			{
				Name: "upsertDiscussionArchiveString",

				Test: upsertDiscussionArchiveString,
			},
			{
				Name: "getModeratorByUserIDString",
				Test: getModeratorByUserIDString,
			},
			{
				Name: "getModeratorByDiscussionIDString",
				Test: getModeratorByDiscussionIDString,
			},
			{
				Name: "getModeratorByUserIDAndDiscussionIDString",
				Test: getModeratorByUserIDAndDiscussionIDString,
			},
			{
				Name: "getModeratedDiscussionsByUserIDString",
				Test: getModeratedDiscussionsByUserIDString,
			},
			{
				Name: "getDiscussionAccessRequestByUserIDString",

				Test: getDiscussionAccessRequestByUserIDString,
			},
			{
				Name: "getDiscussionsByUserAccessString",

				Test: getDiscussionsByUserAccessString,
			},
			{
				Name: "upsertDiscussionUserAccessString",

				Test: upsertDiscussionUserAccessString,
			},
			{
				Name: "getDUAForEverythingNotificationsString",

				Test: getDUAForEverythingNotificationsString,
			},
			{
				Name: "getDUAForMentionNotificationsString",

				Test: getDUAForMentionNotificationsString,
			},
			{
				Name: "deleteDiscussionUserAccessString",

				Test: deleteDiscussionUserAccessString,
			},
			{
				Name: "getDiscussionRequestAccessByIDString",

				Test: getDiscussionRequestAccessByIDString,
			},
			{
				Name: "getDiscussionAccessRequestsString",

				Test: getDiscussionAccessRequestsString,
			},
			{
				Name: "getSentDiscussionAccessRequestsForUserString",

				Test: getSentDiscussionAccessRequestsForUserString,
			},
			{
				Name: "putDiscussionAccessRequestString",

				Test: putDiscussionAccessRequestString,
			},
			{
				Name: "updateDiscussionAccessRequestString",

				Test: updateDiscussionAccessRequestString,
			},
			{
				Name: "getAccessLinkBySlugString",

				Test: getAccessLinkBySlugString,
			},
			{
				Name: "getAccessLinkByDiscussionIDString",

				Test: getAccessLinkByDiscussionIDString,
			},
			{
				Name: "putAccessLinkForDiscussionString",

				Test: putAccessLinkForDiscussionString,
			},
			{
				Name: "getNextShuffleTimeForDiscussionIDString",

				Test: getNextShuffleTimeForDiscussionIDString,
			},
			{
				Name: "putNextShuffleTimeForDiscussionIDString",

				Test: putNextShuffleTimeForDiscussionIDString,
			},
			{
				Name: "getDiscussionsToShuffle",

				Test: getDiscussionsToShuffle,
			},
			{
				Name: "incrDiscussionShuffleCount",

				Test: incrDiscussionShuffleCount,
			},
			{
				Name: "getViewerForDiscussionIDUserID",

				Test: getViewerForDiscussionIDUserID,
			},
			{
				Name: "updateViewerLastViewed",

				Test: updateViewerLastViewed,
			},
		}

		for index, test := range tests {
			Convey(fmt.Sprintf("when prepared statement fails to initialize - %v", test.Name), func() {
				for i := 0; i < index; i++ {
					mock.ExpectPrepare(tests[i].Test)
				}

				mock.ExpectPrepare(test.Test).WillReturnError(fmt.Errorf("error"))
				err := mockDatastore.initializeStatements(ctx)

				So(err, ShouldNotBeNil)
			})
		}
	})
}

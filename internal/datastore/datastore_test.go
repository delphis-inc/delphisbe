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
				Name: "getDiscussionsForAutoPostString",

				Test: getDiscussionsForAutoPostString,
			},
			{
				Name: "getPublicDiscussionsString",
				Test: getPublicDiscussionsString,
			},
			{
				Name: "getModeratorByUserIDString",
				Test: getModeratorByUserIDString,
			},
			{
				Name: "getModeratorByUserIDAndDiscussionIDString",
				Test: getModeratorByUserIDAndDiscussionIDString,
			},
			{
				Name: "getImportedContentByIDString",
				Test: getImportedContentByIDString,
			},
			{
				Name: "getImportedContentForDiscussionString",
				Test: getImportedContentForDiscussionString,
			},
			{
				Name: "getScheduledImportedContentByDiscussionIDString",
				Test: getScheduledImportedContentByDiscussionIDString,
			},
			{
				Name: "putImportedContentString",
				Test: putImportedContentString,
			},
			{
				Name: "putImportedContentDiscussionQueueString",

				Test: putImportedContentDiscussionQueueString,
			},
			{
				Name: "updateImportedContentDiscussionQueueString",

				Test: updateImportedContentDiscussionQueueString,
			},
			{
				Name: "getImportedContentTagsString",
				Test: getImportedContentTagsString,
			},
			{
				Name: "getDiscussionTagsString",
				Test: getDiscussionTagsString,
			},
			{
				Name: "getMatchingTagsString",

				Test: getMatchingTagsString,
			},
			{
				Name: "putImportedContentTagsString",

				Test: putImportedContentTagsString,
			},
			{
				Name: "putDiscussionTagsString",

				Test: putDiscussionTagsString,
			},
			{
				Name: "deleteDiscussionTagsString",

				Test: deleteDiscussionTagsString,
			},
			{
				Name: "getDiscussionsByFlairTemplateForUserString",

				Test: getDiscussionsByFlairTemplateForUserString,
			},
			{
				Name: "getDiscussionsByUserAccessForUserString",

				Test: getDiscussionsByUserAccessForUserString,
			},
			{
				Name: "getDiscussionFlairAccessString",

				Test: getDiscussionFlairAccessString,
			},
			{
				Name: "upsertDiscussionFlairAccessString",

				Test: upsertDiscussionFlairAccessString,
			},
			{
				Name: "upsertDiscussionUserAccessString",

				Test: upsertDiscussionUserAccessString,
			},
			{
				Name: "deleteDiscussionFlairAccessString",

				Test: deleteDiscussionFlairAccessString,
			},
			{
				Name: "deleteDiscussionUserAccessString",

				Test: deleteDiscussionUserAccessString,
			},
			{
				Name: "getDiscussionInviteByIDString",

				Test: getDiscussionInviteByIDString,
			},
			{
				Name: "getDiscussionRequestAccessByIDString",

				Test: getDiscussionRequestAccessByIDString,
			},
			{
				Name: "getDiscussionInvitesForUserString",

				Test: getDiscussionInvitesForUserString,
			},
			{
				Name: "getSentDiscussionInvitesForUserString",

				Test: getSentDiscussionInvitesForUserString,
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
				Name: "getInviteLinksForDiscussion",

				Test: getInviteLinksForDiscussion,
			},
			{
				Name: "putDiscussionInviteRecordString",

				Test: putDiscussionInviteRecordString,
			},
			{
				Name: "putDiscussionAccessRequestString",

				Test: putDiscussionAccessRequestString,
			},
			{
				Name: "updateDiscussionInviteRecordString",

				Test: updateDiscussionInviteRecordString,
			},
			{
				Name: "updateDiscussionAccessRequestString",

				Test: updateDiscussionAccessRequestString,
			},
			{
				Name: "upsertInviteLinksForDiscussion",

				Test: upsertInviteLinksForDiscussion,
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

//func MakeDatastore(ctx context.Context, testData TestData) (Datastore, func() error, error) {
//	url := "postgres://chatham_local@localhost:5432/"
//
//	dbName, err := createTestDatabase(ctx, url+"?sslmode=disable")
//	if err != nil {
//		return nil, nil, errors.Wrap(err, "failed to create test database")
//	}
//
//	url = url + dbName + "?sslmode=disable"
//
//	db, err := connectTestDatabase(ctx, url)
//	if err != nil {
//		return nil, nil, errors.Wrap(err, "failed to create testing datastore")
//	}
//
//	// Create test tables with test data
//	closeFunc, err := db.CreateTestTables(ctx, testData)
//	if err != nil {
//		return nil, nil, errors.Wrap(err, "failed to create test tables")
//	}
//
//	return db, closeFunc, nil
//}
//
//func connectTestDatabase(ctx context.Context, url string) (Datastore, error) {
//	// Initialize gorm
//	testGorm, err := gorm.Open("postgres", url)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to open testing database - gorm")
//	}
//
//	// Initialize sql
//	testDB, err := sql.Open("postgres", url)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to open testing database - SQL")
//	}
//
//	// Create db object
//	return &delphisDB{
//		sql:       testGorm,
//		pg:        testDB,
//		prepStmts: &dbPrepStmts{},
//	}, nil
//}
//
//func createTestDatabase(ctx context.Context, url string) (string, error) {
//	db, err := gorm.Open("postgres", url)
//	if err != nil {
//		return "", err
//	}
//
//	rand.Seed(time.Now().UnixNano())
//	name := "tests_" + strconv.Itoa(rand.Int())
//
//	logrus.Infof("Database Name: %v\n", name)
//
//	db = db.Exec(fmt.Sprintf(`create database %s;`, name))
//	if db.Error != nil {
//		return "", err
//	}
//
//	return name, nil
//}

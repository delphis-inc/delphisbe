package datastore

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/lib/pq"

	"github.com/sirupsen/logrus"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/config"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestDelphisDB_PutPost(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	participantID := "participant1"
	postID := "post1"
	postObject := model.Post{
		ID:            "post1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		PostContent: &model.PostContent{
			ID:      postID,
			Content: "test",
		},
		QuotedPostID: &postID,
		PostType:     model.PostTypeStandard,
	}

	Convey("PutPost", t, func() {
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
			resp, err := mockDatastore.PutPost(ctx, tx, postObject)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putPostString)
			mock.ExpectQuery(putPostString).WithArgs(postObject.ID, postObject.DiscussionID, postObject.ParticipantID, postObject.PostContent.ID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.PutPost(ctx, tx, postObject)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when put post succeeds and returns an object", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "discussion_id", "participant_id", "post_content_id", "quoted_post_id", "media_id", "imported_content_id", "post_type"}).
				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DiscussionID, postObject.ParticipantID, postObject.PostContentID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType)

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putPostString)
			mock.ExpectQuery(putPostString).WithArgs(postObject.ID, postObject.DiscussionID, postObject.ParticipantID, postObject.PostContent.ID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType).WillReturnRows(rs)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.PutPost(ctx, tx, postObject)

			logrus.Infof("Resp: %+v\n", resp)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &postObject)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetPostsByDiscussionIDIter(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	participantID := "participant1"
	postID := "post1"
	postObject := model.Post{
		ID:            "post1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		PostContent: &model.PostContent{
			ID:      postID,
			Content: "test",
		},
		QuotedPostID: &postID,
		PostType:     model.PostTypeStandard,
	}

	emptyPost := model.Post{}

	Convey("GetPostsByDiscussionIDIter", t, func() {
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

			iter := mockDatastore.GetPostsByDiscussionIDIter(ctx, discussionID)

			So(iter.Next(&emptyPost), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getPostsByDiscussionIDString).WithArgs(discussionID).WillReturnError(fmt.Errorf("error"))

			iter := mockDatastore.GetPostsByDiscussionIDIter(ctx, discussionID)

			So(iter.Next(&emptyPost), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns posts", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"p.id", "p.created_at", "p.updated_at", "p.deleted_at", "p.deleted_reason_code", "p.discussion_id", "p.participant_id",
				"p.quoted_post_id", "p.media_id", "p.imported_content_id", "p.post_type", "pc.id", "pc.content", "pc.mentioned_entities"}).
				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DeletedAt, postObject.DeletedReasonCode, postObject.DiscussionID,
					postObject.ParticipantID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType, postObject.PostContent.ID, postObject.PostContent.Content, pq.Array(postObject.PostContent.MentionedEntities)).
				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DeletedAt, postObject.DeletedReasonCode, postObject.DiscussionID,
					postObject.ParticipantID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType, postObject.PostContent.ID, postObject.PostContent.Content, pq.Array(postObject.PostContent.MentionedEntities))

			mock.ExpectQuery(getPostsByDiscussionIDString).WithArgs(discussionID).WillReturnRows(rs)

			iter := mockDatastore.GetPostsByDiscussionIDIter(ctx, discussionID)

			So(iter.Next(&emptyPost), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetPostsByDiscussionIDFromCursorIter(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	participantID := "participant1"
	cursor := now.String()
	limit := 10
	postID := "post1"
	postObject := model.Post{
		ID:            "post1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		PostContent: &model.PostContent{
			ID:      postID,
			Content: "test",
		},
		QuotedPostID: &postID,
		PostType:     model.PostTypeStandard,
	}

	emptyPost := model.Post{}

	Convey("GetPostsByDiscussionIDFromCursorIter", t, func() {
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

			iter := mockDatastore.GetPostsByDiscussionIDFromCursorIter(ctx, discussionID, cursor, limit)

			So(iter.Next(&emptyPost), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getPostsByDiscussionIDFromCursorString).WithArgs(discussionID, cursor, limit).WillReturnError(fmt.Errorf("error"))

			iter := mockDatastore.GetPostsByDiscussionIDFromCursorIter(ctx, discussionID, cursor, limit)

			So(iter.Next(&emptyPost), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns posts", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"p.id", "p.created_at", "p.updated_at", "p.deleted_at", "p.deleted_reason_code", "p.discussion_id", "p.participant_id",
				"p.quoted_post_id", "p.media_id", "p.imported_content_id", "p.post_type", "pc.id", "pc.content", "pc.mentioned_entities"}).
				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DeletedAt, postObject.DeletedReasonCode, postObject.DiscussionID,
					postObject.ParticipantID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType, postObject.PostContent.ID, postObject.PostContent.Content, pq.Array(postObject.PostContent.MentionedEntities)).
				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DeletedAt, postObject.DeletedReasonCode, postObject.DiscussionID,
					postObject.ParticipantID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType, postObject.PostContent.ID, postObject.PostContent.Content, pq.Array(postObject.PostContent.MentionedEntities))

			mock.ExpectQuery(getPostsByDiscussionIDFromCursorString).WithArgs(discussionID, cursor, limit).WillReturnRows(rs)

			iter := mockDatastore.GetPostsByDiscussionIDFromCursorIter(ctx, discussionID, cursor, limit)

			So(iter.Next(&emptyPost), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func Test_GetPostsConnectionByDiscussionID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	participantID := "participant1"
	cursor := now.String()
	limit := 2
	postID := "post1"
	postObject := model.Post{
		ID:            "post1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		PostContent: &model.PostContent{
			ID:      postID,
			Content: "test",
		},
		PostType: model.PostTypeStandard,
	}

	Convey("GetPostsByDiscussionIDFromCursorIter", t, func() {
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

		Convey("when limit less than two is passed in", func() {
			postConns, err := mockDatastore.GetPostsConnectionByDiscussionID(ctx, discussionID, cursor, 1)

			So(err, ShouldNotBeNil)
			So(postConns, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when preparing statements returns an error", func() {
			mockPreparedStatementsWithError(mock)

			postConns, err := mockDatastore.GetPostsConnectionByDiscussionID(ctx, discussionID, cursor, limit)

			So(err, ShouldNotBeNil)
			So(postConns, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when postIterCollect returns an error", func() {
			mockPreparedStatements(mock)

			mock.ExpectQuery(getPostsByDiscussionIDFromCursorString).WithArgs(discussionID, cursor, limit+1).WillReturnError(fmt.Errorf("error"))

			postConns, err := mockDatastore.GetPostsConnectionByDiscussionID(ctx, discussionID, cursor, limit)

			So(err, ShouldNotBeNil)
			So(postConns, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("whenthere are no records for the query", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"p.id", "p.created_at", "p.updated_at", "p.deleted_at", "p.deleted_reason_code", "p.discussion_id", "p.participant_id",
				"p.quoted_post_id", "p.media_id", "p.imported_content_id", "p.post_type", "pc.id", "pc.content", "pc.mentioned_entities"})

			mock.ExpectQuery(getPostsByDiscussionIDFromCursorString).WithArgs(discussionID, cursor, limit+1).WillReturnRows(rs)

			postConns, err := mockDatastore.GetPostsConnectionByDiscussionID(ctx, discussionID, cursor, limit)

			verifyPostConns := &model.PostsConnection{
				Edges: []*model.PostsEdge{},
				PageInfo: model.PageInfo{
					StartCursor: &cursor,
					EndCursor:   &cursor,
					HasNextPage: false,
				},
			}

			So(err, ShouldBeNil)
			So(postConns, ShouldNotBeNil)
			So(postConns, ShouldResemble, verifyPostConns)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns postConnections", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"p.id", "p.created_at", "p.updated_at", "p.deleted_at", "p.deleted_reason_code", "p.discussion_id", "p.participant_id",
				"p.quoted_post_id", "p.media_id", "p.imported_content_id", "p.post_type", "pc.id", "pc.content", "pc.mentioned_entities"}).
				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DeletedAt, postObject.DeletedReasonCode, postObject.DiscussionID,
					postObject.ParticipantID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType, postObject.PostContent.ID, postObject.PostContent.Content, pq.Array(postObject.PostContent.MentionedEntities)).
				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DeletedAt, postObject.DeletedReasonCode, postObject.DiscussionID,
					postObject.ParticipantID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType, postObject.PostContent.ID, postObject.PostContent.Content, pq.Array(postObject.PostContent.MentionedEntities)).
				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DeletedAt, postObject.DeletedReasonCode, postObject.DiscussionID,
					postObject.ParticipantID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType, postObject.PostContent.ID, postObject.PostContent.Content, pq.Array(postObject.PostContent.MentionedEntities))

			mock.ExpectQuery(getPostsByDiscussionIDFromCursorString).WithArgs(discussionID, cursor, limit+1).WillReturnRows(rs)

			postConns, err := mockDatastore.GetPostsConnectionByDiscussionID(ctx, discussionID, cursor, limit)

			cursor := postObject.CreatedAt.Format(time.RFC3339Nano)
			verifyPostConns := &model.PostsConnection{
				Edges: []*model.PostsEdge{
					{
						Cursor: postObject.CreatedAt.Format(time.RFC3339Nano),
						Node:   &postObject,
					},
					{
						Cursor: postObject.CreatedAt.Format(time.RFC3339Nano),
						Node:   &postObject,
					},
				},
				PageInfo: model.PageInfo{
					StartCursor: &cursor,
					EndCursor:   &cursor,
					HasNextPage: true,
				},
			}

			So(err, ShouldBeNil)
			So(postConns, ShouldNotBeNil)
			So(postConns, ShouldResemble, verifyPostConns)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetLastPostByDiscussionID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	participantID := "participant1"
	postID := "post1"
	postObject := model.Post{
		ID:            "post1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		PostContent: &model.PostContent{
			ID:      postID,
			Content: "test",
		},
		QuotedPostID: &postID,
		PostType:     model.PostTypeStandard,
	}
	minuteRange := 120

	Convey("GetLastPostByDiscussionID", t, func() {
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

			resp, err := mockDatastore.GetLastPostByDiscussionID(ctx, discussionID, minuteRange)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getLastPostByDiscussionIDStmt).WithArgs(discussionID, minuteRange).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetLastPostByDiscussionID(ctx, discussionID, minuteRange)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns no rows", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getLastPostByDiscussionIDStmt).WithArgs(discussionID, minuteRange).WillReturnError(sql.ErrNoRows)

			resp, err := mockDatastore.GetLastPostByDiscussionID(ctx, discussionID, minuteRange)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a post", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"p.id", "p.created_at", "p.updated_at", "p.deleted_at", "p.deleted_reason_code", "p.discussion_id", "p.participant_id",
				"p.quoted_post_id", "p.media_id", "p.imported_content_id", "p.post_type", "pc.id", "pc.content", "pc.mentioned_entities"}).
				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DeletedAt, postObject.DeletedReasonCode, postObject.DiscussionID,
					postObject.ParticipantID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType, postObject.PostContent.ID, postObject.PostContent.Content, pq.Array(postObject.PostContent.MentionedEntities))

			mock.ExpectQuery(getLastPostByDiscussionIDStmt).WithArgs(discussionID, minuteRange).WillReturnRows(rs)

			resp, err := mockDatastore.GetLastPostByDiscussionID(ctx, discussionID, minuteRange)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &postObject)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetPostByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	participantID := "participant1"
	postID := "post1"
	postObject := model.Post{
		ID:            "post1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		PostContent: &model.PostContent{
			ID:      postID,
			Content: "test",
		},
		QuotedPostID: &postID,
		PostType:     model.PostTypeStandard,
	}

	Convey("GetPostByID", t, func() {
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

			resp, err := mockDatastore.GetPostByID(ctx, postID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getPostByIDString).WithArgs(postID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetPostByID(ctx, postID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns no rows", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getPostByIDString).WithArgs(postID).WillReturnError(sql.ErrNoRows)

			resp, err := mockDatastore.GetPostByID(ctx, postID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a post", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"p.id", "p.created_at", "p.updated_at", "p.deleted_at", "p.deleted_reason_code", "p.discussion_id", "p.participant_id",
				"p.quoted_post_id", "p.media_id", "p.imported_content_id", "p.post_type", "pc.id", "pc.content", "pc.mentioned_entities"}).
				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DeletedAt, postObject.DeletedReasonCode, postObject.DiscussionID,
					postObject.ParticipantID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType, postObject.PostContent.ID, postObject.PostContent.Content, pq.Array(postObject.PostContent.MentionedEntities))

			mock.ExpectQuery(getPostByIDString).WithArgs(postID).WillReturnRows(rs)

			resp, err := mockDatastore.GetPostByID(ctx, postID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &postObject)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestPostIter_Next(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	participantID := "participant1"
	postID := "post1"
	postObject := model.Post{
		ID:            "post1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		PostContent: &model.PostContent{
			ID:      postID,
			Content: "test",
		},
		QuotedPostID: &postID,
		PostType:     model.PostTypeStandard,
	}

	emptyPost := model.Post{}

	Convey("PostIter Next", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := postIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Next(&emptyPost), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has a context error passed in", func() {
			ctx1, cancelFunc := context.WithCancel(ctx)
			cancelFunc()
			iter := postIter{
				ctx: ctx1,
			}

			So(iter.Next(&emptyPost), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has no more rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"p.id", "p.created_at", "p.updated_at", "p.deleted_at", "p.deleted_reason_code", "p.discussion_id", "p.participant_id",
				"p.quoted_post_id", "p.media_id", "p.imported_content_id", "p.post_type", "pc.id", "pc.content", "pc.mentioned_entities"})

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := postIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyPost), ShouldBeFalse)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on scan", func() {
			rs := sqlmock.NewRows([]string{"p.id", "p.created_at", "p.updated_at", "p.deleted_at", "p.deleted_reason_code", "p.discussion_id", "p.participant_id",
				"p.quoted_post_id", "p.media_id", "p.imported_content_id", "p.post_type", "pc.id", "pc.content"}).
				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DeletedAt, postObject.DeletedReasonCode, postObject.DiscussionID,
					postObject.ParticipantID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType, postObject.PostContent.ID, postObject.PostContent.Content)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := postIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyPost), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"p.id", "p.created_at", "p.updated_at", "p.deleted_at", "p.deleted_reason_code", "p.discussion_id", "p.participant_id",
				"p.quoted_post_id", "p.media_id", "p.imported_content_id", "p.post_type", "pc.id", "pc.content", "pc.mentioned_entities"}).
				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DeletedAt, postObject.DeletedReasonCode, postObject.DiscussionID,
					postObject.ParticipantID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType, postObject.PostContent.ID, postObject.PostContent.Content, pq.Array(postObject.PostContent.MentionedEntities)).
				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DeletedAt, postObject.DeletedReasonCode, postObject.DiscussionID,
					postObject.ParticipantID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType, postObject.PostContent.ID, postObject.PostContent.Content, pq.Array(postObject.PostContent.MentionedEntities))

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := postIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyPost), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestPostIter_Close(t *testing.T) {
	ctx := context.Background()

	Convey("PostIter Close", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := postIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on rows.Close", func() {
			rs := mock.NewRows([]string{"p.id", "p.created_at", "p.updated_at", "p.deleted_at", "p.deleted_reason_code", "p.discussion_id", "p.participant_id",
				"p.quoted_post_id", "p.media_id", "p.imported_content_id", "p.post_type", "pc.id", "pc.content", "pc.mentioned_entities"}).CloseError(fmt.Errorf("error"))

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := postIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_PostIterCollect(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	participantID := "participant1"
	postID := "post1"
	quotePostID := "quotePost"
	postObject := model.Post{
		ID:            "post1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		PostContent: &model.PostContent{
			ID:      postID,
			Content: "test",
		},
		PostType: model.PostTypeStandard,
	}

	quotePostObject := model.Post{
		ID:            quotePostID,
		CreatedAt:     now,
		UpdatedAt:     now,
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		PostContent: &model.PostContent{
			ID:      postID,
			Content: "test",
		},
		PostType: model.PostTypeStandard,
	}

	Convey("PostIterCollect", t, func() {
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

		Convey("when the iterator fails to close", func() {
			iter := &postIter{
				err: fmt.Errorf("error"),
			}

			posts, err := mockDatastore.PostIterCollect(ctx, iter)

			So(err, ShouldNotBeNil)
			So(posts, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on fetching the quoted post", func() {
			// Setup up post with quotedPost
			basePost := postObject
			basePost.QuotedPostID = &postID
			basePost.QuotedPost = &quotePostObject

			rs := sqlmock.NewRows([]string{"p.id", "p.created_at", "p.updated_at", "p.deleted_at", "p.deleted_reason_code", "p.discussion_id", "p.participant_id",
				"p.quoted_post_id", "p.media_id", "p.imported_content_id", "p.post_type", "pc.id", "pc.content", "pc.mentioned_entities"}).
				AddRow(basePost.ID, basePost.CreatedAt, basePost.UpdatedAt, basePost.DeletedAt, basePost.DeletedReasonCode, basePost.DiscussionID,
					basePost.ParticipantID, basePost.QuotedPostID, basePost.MediaID, basePost.ImportedContentID, basePost.PostType, basePost.PostContent.ID, basePost.PostContent.Content, pq.Array(basePost.PostContent.MentionedEntities))

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			mockPreparedStatements(mock)
			mock.ExpectQuery(getPostByIDString).WithArgs(postID).WillReturnError(fmt.Errorf("error"))

			iter := &postIter{
				ctx:  ctx,
				rows: rs1,
			}

			posts, err := mockDatastore.PostIterCollect(ctx, iter)

			So(err, ShouldNotBeNil)
			So(posts, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has results and returns slice of Posts", func() {
			rs := sqlmock.NewRows([]string{"p.id", "p.created_at", "p.updated_at", "p.deleted_at", "p.deleted_reason_code", "p.discussion_id", "p.participant_id",
				"p.quoted_post_id", "p.media_id", "p.imported_content_id", "p.post_type", "pc.id", "pc.content", "pc.mentioned_entities"}).
				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DeletedAt, postObject.DeletedReasonCode, postObject.DiscussionID,
					postObject.ParticipantID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType, postObject.PostContent.ID, postObject.PostContent.Content, pq.Array(postObject.PostContent.MentionedEntities)).
				AddRow(postObject.ID, postObject.CreatedAt, postObject.UpdatedAt, postObject.DeletedAt, postObject.DeletedReasonCode, postObject.DiscussionID,
					postObject.ParticipantID, postObject.QuotedPostID, postObject.MediaID, postObject.ImportedContentID, postObject.PostType, postObject.PostContent.ID, postObject.PostContent.Content, pq.Array(postObject.PostContent.MentionedEntities))

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := &postIter{
				ctx:  ctx,
				rows: rs1,
			}

			posts, err := mockDatastore.PostIterCollect(ctx, iter)

			So(err, ShouldBeNil)
			So(posts, ShouldNotBeNil)
			So(posts, ShouldResemble, []*model.Post{&postObject, &postObject})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has results and one post is a quotedPost", func() {
			// Setup up post with quotedPost
			basePost := postObject
			basePost.QuotedPostID = &postID
			basePost.QuotedPost = &quotePostObject

			rs := sqlmock.NewRows([]string{"p.id", "p.created_at", "p.updated_at", "p.deleted_at", "p.deleted_reason_code", "p.discussion_id", "p.participant_id",
				"p.quoted_post_id", "p.media_id", "p.imported_content_id", "p.post_type", "pc.id", "pc.content", "pc.mentioned_entities"}).
				AddRow(basePost.ID, basePost.CreatedAt, basePost.UpdatedAt, basePost.DeletedAt, basePost.DeletedReasonCode, basePost.DiscussionID,
					basePost.ParticipantID, basePost.QuotedPostID, basePost.MediaID, basePost.ImportedContentID, basePost.PostType, basePost.PostContent.ID, basePost.PostContent.Content, pq.Array(basePost.PostContent.MentionedEntities))

			quoteRow := sqlmock.NewRows([]string{"p.id", "p.created_at", "p.updated_at", "p.deleted_at", "p.deleted_reason_code", "p.discussion_id", "p.participant_id",
				"p.quoted_post_id", "p.media_id", "p.imported_content_id", "p.post_type", "pc.id", "pc.content", "pc.mentioned_entities"}).
				AddRow(quotePostObject.ID, quotePostObject.CreatedAt, quotePostObject.UpdatedAt, quotePostObject.DeletedAt, quotePostObject.DeletedReasonCode, quotePostObject.DiscussionID,
					quotePostObject.ParticipantID, quotePostObject.QuotedPostID, quotePostObject.MediaID, quotePostObject.ImportedContentID, quotePostObject.PostType, quotePostObject.PostContent.ID, quotePostObject.PostContent.Content, pq.Array(quotePostObject.PostContent.MentionedEntities))

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			mockPreparedStatements(mock)
			mock.ExpectQuery(getPostByIDString).WithArgs(postID).WillReturnRows(quoteRow)

			iter := &postIter{
				ctx:  ctx,
				rows: rs1,
			}

			posts, err := mockDatastore.PostIterCollect(ctx, iter)

			So(err, ShouldBeNil)
			So(posts, ShouldNotBeNil)
			So(posts, ShouldResemble, []*model.Post{&basePost})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

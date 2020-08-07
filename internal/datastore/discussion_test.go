package datastore

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"

	"github.com/delphis-inc/delphisbe/graph/model"
)

var emptyString = ""

func TestDelphisDB_GetDiscussionByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	discObj := model.Discussion{
		ID:            "discussion1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
		Title:         "test",
		AnonymityType: "",
		ModeratorID:   &modID,
		AutoPost:      false,
		IdleMinutes:   120,
		IconURL:       &emptyString,
	}
	modObj := model.Moderator{
		ID:        modID,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}

	Convey("GetDiscussionByID", t, func() {
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

		expectedQueryString := `SELECT * FROM "discussions"  WHERE "discussions"."deleted_at" IS NULL AND (("discussions"."id" IN ($1)))`
		expectedModQueryString := `SELECT * FROM "moderators"  WHERE "moderators"."deleted_at" IS NULL AND (("id" IN ($1)))`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(discObj.ID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetDiscussionByID(ctx, discObj.ID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and does not find a record", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "anonymity_type", "moderator_id", "icon_url",
				"auto_post", "idle_minutes"})

			mock.ExpectQuery(expectedQueryString).WithArgs(discObj.ID).WillReturnRows(rs)

			resp, err := mockDatastore.GetDiscussionByID(ctx, discObj.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a discussions", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "anonymity_type", "moderator_id", "icon_url",
				"auto_post", "idle_minutes"}).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt, discObj.Title, discObj.AnonymityType,
					discObj.ModeratorID, discObj.IconURL, discObj.AutoPost, discObj.IdleMinutes)

			modRs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_profile_id"}).
				AddRow(modObj.ID, modObj.CreatedAt, modObj.UpdatedAt, modObj.DeletedAt, modObj.UserProfileID).
				AddRow(modObj.ID, modObj.CreatedAt, modObj.UpdatedAt, modObj.DeletedAt, modObj.UserProfileID)

			mock.ExpectQuery(expectedQueryString).WithArgs(discObj.ID).WillReturnRows(rs)
			mock.ExpectQuery(expectedModQueryString).WithArgs(discObj.ModeratorID).WillReturnRows(modRs)

			resp, err := mockDatastore.GetDiscussionByID(ctx, discObj.ID)

			discObj.Moderator = &modObj

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &discObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetDiscussionsByIDs(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	discObj1 := model.Discussion{
		ID:            "discussion1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
		Title:         "test",
		AnonymityType: "",
		ModeratorID:   &modID,
		AutoPost:      false,
		IdleMinutes:   120,
		IconURL:       &emptyString,
	}
	discObj2 := model.Discussion{
		ID:            "discussion2",
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
		Title:         "test",
		AnonymityType: "",
		ModeratorID:   &modID,
		AutoPost:      false,
		IdleMinutes:   120,
		IconURL:       &emptyString,
	}
	modObj := model.Moderator{
		ID:        modID,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}
	discs := []model.Discussion{discObj1, discObj2}
	discussionIDs := []string{discObj1.ID, discObj2.ID}

	Convey("GetDiscussionsByIDs", t, func() {
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

		expectedQueryString := `SELECT * FROM "discussions"  WHERE "discussions"."deleted_at" IS NULL AND (("discussions"."id" IN ($1,$2)))`
		expectedModQueryString := `SELECT * FROM "moderators"  WHERE "moderators"."deleted_at" IS NULL AND (("id" IN ($1,$2)))`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(discussionIDs[0], discussionIDs[1]).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetDiscussionsByIDs(ctx, discussionIDs)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and does not find a record", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "anonymity_type", "moderator_id", "icon_url",
				"auto_post", "idle_minutes"})

			mock.ExpectQuery(expectedQueryString).WithArgs(discussionIDs[0], discussionIDs[1]).WillReturnRows(rs)

			resp, err := mockDatastore.GetDiscussionsByIDs(ctx, discussionIDs)

			verifyMap := map[string]*model.Discussion{
				discs[0].ID: nil,
				discs[1].ID: nil,
			}

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, verifyMap)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a discussions", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "anonymity_type", "moderator_id", "icon_url",
				"auto_post", "idle_minutes"}).
				AddRow(discObj1.ID, discObj1.CreatedAt, discObj1.UpdatedAt, discObj1.DeletedAt, discObj1.Title, discObj1.AnonymityType,
					discObj1.ModeratorID, discObj1.IconURL, discObj1.AutoPost, discObj1.IdleMinutes).
				AddRow(discObj2.ID, discObj2.CreatedAt, discObj2.UpdatedAt, discObj2.DeletedAt, discObj2.Title, discObj2.AnonymityType,
					discObj2.ModeratorID, discObj2.IconURL, discObj2.AutoPost, discObj2.IdleMinutes)

			modRs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_profile_id"}).
				AddRow(modObj.ID, modObj.CreatedAt, modObj.UpdatedAt, modObj.DeletedAt, modObj.UserProfileID).
				AddRow(modObj.ID, modObj.CreatedAt, modObj.UpdatedAt, modObj.DeletedAt, modObj.UserProfileID)

			mock.ExpectQuery(expectedQueryString).WithArgs(discussionIDs[0], discussionIDs[1]).WillReturnRows(rs)
			mock.ExpectQuery(expectedModQueryString).WithArgs(discs[0].ModeratorID, discs[1].ModeratorID).WillReturnRows(modRs)

			resp, err := mockDatastore.GetDiscussionsByIDs(ctx, discussionIDs)

			discObj1.Moderator = &modObj
			discObj2.Moderator = &modObj

			verifyMap := map[string]*model.Discussion{
				discs[0].ID: &discObj1,
				discs[1].ID: &discObj2,
			}

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, verifyMap)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})

}

// GORM!!!
func TestDelphisDB_GetDiscussionByModeratorID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	discObj := model.Discussion{
		ID:            "discussion1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
		Title:         "test",
		AnonymityType: "",
		ModeratorID:   &modID,
		AutoPost:      false,
		IdleMinutes:   120,
		IconURL:       &emptyString,
	}

	Convey("GetDiscussionByModeratorID", t, func() {
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

		expectedQueryString := `SELECT * FROM "discussions" WHERE "discussions"."deleted_at" IS NULL AND (("discussions"."moderator_id" = $1)) ORDER BY "discussions"."id" ASC LIMIT 1`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(discObj.ModeratorID).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.GetDiscussionByModeratorID(ctx, *discObj.ModeratorID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns and does not find a record", func() {
			mock.ExpectQuery(expectedQueryString).WithArgs(discObj.ModeratorID).WillReturnError(gorm.ErrRecordNotFound)

			resp, err := mockDatastore.GetDiscussionByModeratorID(ctx, *discObj.ModeratorID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a discussions", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "anonymity_type", "moderator_id", "icon_url",
				"auto_post", "idle_minutes"}).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt, discObj.Title, discObj.AnonymityType,
					discObj.ModeratorID, discObj.IconURL, discObj.AutoPost, discObj.IdleMinutes)

			mock.ExpectQuery(expectedQueryString).WithArgs(discObj.ModeratorID).WillReturnRows(rs)

			resp, err := mockDatastore.GetDiscussionByModeratorID(ctx, *discObj.ModeratorID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &discObj)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_GetDiscussionsAutoPost(t *testing.T) {
	ctx := context.Background()
	discussionID := "discussion1"
	discAP := model.DiscussionAutoPost{
		ID:          discussionID,
		IdleMinutes: 120,
	}

	emptyAP := model.DiscussionAutoPost{}

	Convey("GetDiscussionsAutoPost", t, func() {
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

			iter := mockDatastore.GetDiscussionsAutoPost(ctx)

			So(iter.Next(&emptyAP), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getDiscussionsForAutoPostString).WillReturnError(fmt.Errorf("error"))

			iter := mockDatastore.GetDiscussionsAutoPost(ctx)

			So(iter.Next(&emptyAP), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns posts", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"id", "idle_minutes"}).
				AddRow(discAP.ID, discAP.IdleMinutes).
				AddRow(discAP.ID, discAP.IdleMinutes)

			mock.ExpectQuery(getDiscussionsForAutoPostString).WillReturnRows(rs)

			iter := mockDatastore.GetDiscussionsAutoPost(ctx)

			So(iter.Next(&emptyAP), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_ListDiscussions(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	discObj := model.Discussion{
		ID:            "discussion1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
		Title:         "test",
		AnonymityType: "",
		ModeratorID:   &modID,
		AutoPost:      false,
		IdleMinutes:   120,
		IconURL:       &emptyString,
	}
	modObj := model.Moderator{
		ID:        modID,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}

	Convey("ListDiscussions", t, func() {
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

		expectedQueryString := `SELECT * FROM "discussions" WHERE "discussions"."deleted_at" IS NULL`
		expectedModQueryString := `SELECT * FROM "moderators" WHERE "moderators"."deleted_at" IS NULL AND (("id" IN ($1,$2)))`

		Convey("when query execution returns an error", func() {
			mock.ExpectQuery(expectedQueryString).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.ListDiscussions(ctx)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a discussions", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "anonymity_type", "moderator_id", "icon_url",
				"auto_post", "idle_minutes"}).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt, discObj.Title, discObj.AnonymityType,
					discObj.ModeratorID, discObj.IconURL, discObj.AutoPost, discObj.IdleMinutes).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt, discObj.Title, discObj.AnonymityType,
					discObj.ModeratorID, discObj.IconURL, discObj.AutoPost, discObj.IdleMinutes)

			modRs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_profile_id"}).
				AddRow(modObj.ID, modObj.CreatedAt, modObj.UpdatedAt, modObj.DeletedAt, modObj.UserProfileID).
				AddRow(modObj.ID, modObj.CreatedAt, modObj.UpdatedAt, modObj.DeletedAt, modObj.UserProfileID)

			mock.ExpectQuery(expectedQueryString).WillReturnRows(rs)
			mock.ExpectQuery(expectedModQueryString).WithArgs(discObj.ModeratorID, discObj.ModeratorID).WillReturnRows(modRs)

			resp, err := mockDatastore.ListDiscussions(ctx)

			discObj.Moderator = &modObj

			verifyDC := model.DiscussionsConnection{
				Edges: []*model.DiscussionsEdge{
					{
						Node: &discObj,
					},
					{
						Node: &discObj,
					},
				},
				IDs: []string{discObj.ID, discObj.ID},
			}

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &verifyDC)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_ListDiscussionsByUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	userID := "userID"
	state := model.DiscussionUserAccessStateActive
	discObj := model.Discussion{
		ID:            "discussion1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
		Title:         "test",
		AnonymityType: "",
		ModeratorID:   &modID,
		AutoPost:      false,
		IdleMinutes:   120,
		IconURL:       &emptyString,
	}

	Convey("ListDiscussionsByUserID", t, func() {
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

			resp, err := mockDatastore.ListDiscussionsByUserID(ctx, userID, state)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getDiscussionsByUserAccessString).WithArgs(userID, state).WillReturnError(fmt.Errorf("error"))

			resp, err := mockDatastore.ListDiscussionsByUserID(ctx, userID, state)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns a discussions", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title",
				"anonymity_type", "moderator_id", "auto_post", "icon_url", "idle_minutes", "description",
				"title_history", "description_history", "discussion_joinability", "last_post_id", "last_post_created_at",
				"shuffle_count", "lock_status"}).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt,
					discObj.Title, discObj.AnonymityType, discObj.ModeratorID, discObj.AutoPost,
					discObj.IconURL, discObj.IdleMinutes, discObj.Description, discObj.TitleHistory,
					discObj.DescriptionHistory, discObj.DiscussionJoinability, discObj.LastPostID, discObj.LastPostCreatedAt,
					discObj.ShuffleCount, discObj.LockStatus).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt,
					discObj.Title, discObj.AnonymityType, discObj.ModeratorID, discObj.AutoPost,
					discObj.IconURL, discObj.IdleMinutes, discObj.Description, discObj.TitleHistory,
					discObj.DescriptionHistory, discObj.DiscussionJoinability, discObj.LastPostID, discObj.LastPostCreatedAt,
					discObj.ShuffleCount, discObj.LockStatus)

			mock.ExpectQuery(getDiscussionsByUserAccessString).WithArgs(userID, state).WillReturnRows(rs)

			resp, err := mockDatastore.ListDiscussionsByUserID(ctx, userID, state)

			verifyDC := model.DiscussionsConnection{
				Edges: []*model.DiscussionsEdge{
					{
						Node: &discObj,
					},
					{
						Node: &discObj,
					},
				},
				IDs: []string{discObj.ID, discObj.ID},
			}

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &verifyDC)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_UpsertDiscussion(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	iconURL := "http://"
	discObj := model.Discussion{
		ID:            "discussion1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
		Title:         "test",
		Description:   "description",
		AnonymityType: model.AnonymityTypeStrong,
		ModeratorID:   &modID,
		AutoPost:      false,
		IdleMinutes:   120,
		IconURL:       &iconURL,
	}
	modObj := model.Moderator{
		ID:        modID,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}

	Convey("UpsertDiscussion", t, func() {
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

		expectedFindQueryStr := `SELECT * FROM "discussions" WHERE "discussions"."deleted_at" IS NULL AND (("discussions"."id" = $1)) ORDER BY "discussions"."id" ASC LIMIT 1`
		createQueryStr := `INSERT INTO "discussions" ("id","created_at","updated_at","deleted_at","title","description","title_history","description_history","anonymity_type","moderator_id","auto_post","idle_minutes","icon_url","discussion_joinability","last_post_id","last_post_created_at","shuffle_count","lock_status") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) RETURNING "discussions"."id"`

		expectedNewObjectRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "title_history", "description_history", "anonymity_type", "moderator_id", "icon_url",
			"auto_post", "idle_minutes", "discussion_joinability"}).
			AddRow(discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt, discObj.Title, discObj.Description, discObj.TitleHistory, discObj.DescriptionHistory, discObj.AnonymityType,
				discObj.ModeratorID, discObj.IconURL, discObj.AutoPost, discObj.IdleMinutes, discObj.DiscussionJoinability)

		expectedUpdateStr := `UPDATE "discussions" SET "anonymity_type" = $1, "auto_post" = $2, "description" = $3, "description_history" = $4, "discussion_joinability" = $5, "icon_url" = $6, "idle_minutes" = $7, "last_post_created_at" = $8, "last_post_id" = $9, "lock_status" = $10, "title" = $11, "title_history" = $12, "updated_at" = $13 WHERE "discussions"."deleted_at" IS NULL AND "discussions"."id" = $14`
		expectedPostUpdateSelectStr := `SELECT * FROM "discussions" WHERE "discussions"."deleted_at" IS NULL AND "discussions"."id" = $1 ORDER BY "discussions"."id" ASC LIMIT 1`
		expectedPostUpdateModSelectStr := `SELECT * FROM "moderators"  WHERE "moderators"."deleted_at" IS NULL AND (("id" IN ($1))) ORDER BY "moderators"."id" ASC`

		Convey("when find query fails with a non-not-found-error the function should return the error", func() {
			expectedError := fmt.Errorf("Some fake error")
			mock.ExpectQuery(expectedFindQueryStr).WithArgs(discObj.ID).WillReturnError(expectedError)

			resp, err := mockDatastore.UpsertDiscussion(ctx, discObj)

			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when find query returns not-found-error it should call create", func() {
			Convey("when create returns an error it should return it", func() {
				expectedError := fmt.Errorf("Some fake error")

				mock.ExpectQuery(expectedFindQueryStr).WithArgs(discObj.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt, discObj.Title,
					discObj.Description, discObj.TitleHistory, discObj.DescriptionHistory, discObj.AnonymityType,
					discObj.ModeratorID, discObj.AutoPost, discObj.IdleMinutes, discObj.IconURL,
					discObj.DiscussionJoinability, discObj.LastPostID, discObj.LastPostCreatedAt,
					discObj.ShuffleCount, discObj.LockStatus,
				).WillReturnError(expectedError)

				resp, err := mockDatastore.UpsertDiscussion(ctx, discObj)

				assert.NotNil(t, err)
				assert.Nil(t, resp)
				assert.Nil(t, mock.ExpectationsWereMet())
			})

			Convey("when create succeeds it should return the new object", func() {
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(discObj.ID).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectQuery(createQueryStr).WithArgs(
					discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt, discObj.Title, discObj.Description,
					discObj.TitleHistory, discObj.DescriptionHistory, discObj.AnonymityType,
					discObj.ModeratorID, discObj.AutoPost, discObj.IdleMinutes, discObj.IconURL,
					discObj.DiscussionJoinability, discObj.LastPostID, discObj.LastPostCreatedAt,
					discObj.ShuffleCount, discObj.LockStatus,
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(discObj.ID))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(discObj.ID).WillReturnRows(expectedNewObjectRow)

				resp, err := mockDatastore.UpsertDiscussion(ctx, discObj)

				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(resp, ShouldResemble, &discObj)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})
		})

		Convey("when the object is found we should update it", func() {
			Convey("when updating if it fails then return the error", func() {
				expectedError := fmt.Errorf("something went wrong")
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(discObj.ID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(expectedUpdateStr).WithArgs(
					discObj.AnonymityType, discObj.AutoPost, discObj.Description, discObj.DescriptionHistory,
					discObj.DiscussionJoinability, discObj.IconURL, discObj.IdleMinutes,
					discObj.LastPostCreatedAt, discObj.LastPostID, discObj.LockStatus,
					discObj.Title, discObj.TitleHistory, sqlmock.AnyArg(), discObj.ID,
				).WillReturnError(expectedError)

				resp, err := mockDatastore.UpsertDiscussion(ctx, discObj)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})

			Convey("when updating if it fails on moderator query", func() {
				expectedError := fmt.Errorf("something went wrong")
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(discObj.ID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(expectedUpdateStr).WithArgs(
					discObj.AnonymityType, discObj.AutoPost, discObj.Description, discObj.DescriptionHistory,
					discObj.DiscussionJoinability, discObj.IconURL, discObj.IdleMinutes,
					discObj.LastPostCreatedAt, discObj.LastPostID, discObj.LockStatus,
					discObj.Title, discObj.TitleHistory, sqlmock.AnyArg(), discObj.ID,
				).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedPostUpdateSelectStr).WithArgs(discObj.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "title_history",
						"description_history", "anonymity_type", "moderator_id", "icon_url",
						"auto_post", "idle_minutes", "discussion_joinability", "last_post_id", "last_post_created_at"}).
						AddRow(discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt, discObj.Title,
							discObj.Description, discObj.TitleHistory, discObj.DescriptionHistory, discObj.AnonymityType,
							discObj.ModeratorID, discObj.IconURL, discObj.AutoPost, discObj.IdleMinutes,
							discObj.DiscussionJoinability, discObj.LastPostID, discObj.LastPostCreatedAt))
				mock.ExpectQuery(expectedPostUpdateModSelectStr).WithArgs(discObj.ModeratorID).WillReturnError(expectedError)

				resp, err := mockDatastore.UpsertDiscussion(ctx, discObj)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
				So(mock.ExpectationsWereMet(), ShouldBeNil)
			})

			Convey("when updating if it succeeds it should return the updated object", func() {
				mock.ExpectQuery(expectedFindQueryStr).WithArgs(discObj.ID).WillReturnRows(expectedNewObjectRow)
				mock.ExpectBegin()
				mock.ExpectExec(expectedUpdateStr).WithArgs(
					discObj.AnonymityType, discObj.AutoPost, discObj.Description, discObj.DescriptionHistory,
					discObj.DiscussionJoinability, discObj.IconURL, discObj.IdleMinutes,
					discObj.LastPostCreatedAt, discObj.LastPostID, discObj.LockStatus,
					discObj.Title, discObj.TitleHistory, sqlmock.AnyArg(), discObj.ID,
				).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				mock.ExpectQuery(expectedPostUpdateSelectStr).WithArgs(discObj.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "title_history",
						"description_history", "anonymity_type", "moderator_id", "icon_url",
						"auto_post", "idle_minutes", "discussion_joinability", "last_post_id", "last_post_created_at"}).
						AddRow(discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt, discObj.Title,
							discObj.Description, discObj.TitleHistory, discObj.DescriptionHistory, discObj.AnonymityType,
							discObj.ModeratorID, discObj.IconURL, discObj.AutoPost, discObj.IdleMinutes,
							discObj.DiscussionJoinability, discObj.LastPostID, discObj.LastPostCreatedAt))
				mock.ExpectQuery(expectedPostUpdateModSelectStr).WithArgs(discObj.ModeratorID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_profile_id"}).
						AddRow(modObj.ID, modObj.CreatedAt, modObj.UpdatedAt, modObj.DeletedAt, modObj.UserProfileID))

				resp, err := mockDatastore.UpsertDiscussion(ctx, discObj)

				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Nil(t, mock.ExpectationsWereMet())
			})
		})
	})
}

func TestDelphisDB_GetDiscussionTags(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	tagObject := model.Tag{
		ID:        discussionID,
		Tag:       "testTag",
		CreatedAt: now,
		DeletedAt: nil,
	}

	emptyTag := model.Tag{}

	Convey("GetDiscussionTags", t, func() {
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

			iter := mockDatastore.GetDiscussionTags(ctx, discussionID)

			So(iter.Next(&emptyTag), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mockPreparedStatements(mock)
			mock.ExpectQuery(getDiscussionTagsString).WithArgs(discussionID).WillReturnError(fmt.Errorf("error"))

			iter := mockDatastore.GetDiscussionTags(ctx, discussionID)

			So(iter.Next(&emptyTag), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution succeeds and returns posts", func() {
			mockPreparedStatements(mock)
			rs := sqlmock.NewRows([]string{"discussion_id", "tag", "created_at"}).
				AddRow(tagObject.ID, tagObject.Tag, tagObject.CreatedAt).
				AddRow(tagObject.ID, tagObject.Tag, tagObject.CreatedAt)

			mock.ExpectQuery(getDiscussionTagsString).WithArgs(discussionID).WillReturnRows(rs)

			iter := mockDatastore.GetDiscussionTags(ctx, discussionID)

			So(iter.Next(&emptyTag), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_IncrementDiscussionShuffleCount(t *testing.T) {
	ctx := context.Background()
	discussionID := "discussion1"

	Convey("IncrementDiscussionShuffleCount", t, func() {
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
			resp, err := mockDatastore.IncrementDiscussionShuffleCount(ctx, tx, discussionID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(incrDiscussionShuffleCount)
			mock.ExpectQuery(incrDiscussionShuffleCount).WithArgs(discussionID).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.IncrementDiscussionShuffleCount(ctx, tx, discussionID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
		Convey("when query execution has a conflilct and doesn't return a row", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(incrDiscussionShuffleCount)
			mock.ExpectQuery(incrDiscussionShuffleCount).WithArgs(discussionID).WillReturnError(sql.ErrNoRows)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.IncrementDiscussionShuffleCount(ctx, tx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when put post succeeds and returns an object", func() {
			newShuffleCount := 1
			rs := sqlmock.NewRows([]string{"shuffle_count"}).
				AddRow(newShuffleCount)

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(incrDiscussionShuffleCount)
			mock.ExpectQuery(incrDiscussionShuffleCount).WithArgs(discussionID).WillReturnRows(rs)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.IncrementDiscussionShuffleCount(ctx, tx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &newShuffleCount)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_PutDiscussionTags(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	tagObject := model.Tag{
		ID:        discussionID,
		Tag:       "testTag",
		CreatedAt: now,
		DeletedAt: nil,
	}

	Convey("PutDiscussionTags", t, func() {
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
			resp, err := mockDatastore.PutDiscussionTags(ctx, tx, tagObject)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putDiscussionTagsString)
			mock.ExpectQuery(putDiscussionTagsString).WithArgs(tagObject.ID, tagObject.Tag).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.PutDiscussionTags(ctx, tx, tagObject)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution has a conflilct and doesn't return a row", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putDiscussionTagsString)
			mock.ExpectQuery(putDiscussionTagsString).WithArgs(tagObject.ID, tagObject.Tag).WillReturnError(sql.ErrNoRows)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.PutDiscussionTags(ctx, tx, tagObject)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &model.Tag{})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when put post succeeds and returns an object", func() {
			rs := sqlmock.NewRows([]string{"id", "tag", "created_at"}).
				AddRow(tagObject.ID, tagObject.Tag, tagObject.CreatedAt)

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(putDiscussionTagsString)
			mock.ExpectQuery(putDiscussionTagsString).WithArgs(tagObject.ID, tagObject.Tag).WillReturnRows(rs)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.PutDiscussionTags(ctx, tx, tagObject)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &tagObject)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_DeleteDiscussionTags(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussion1"
	tagObject := model.Tag{
		ID:        discussionID,
		Tag:       "testTag",
		CreatedAt: now,
		DeletedAt: nil,
	}

	Convey("DeleteDiscussionTags", t, func() {
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
			resp, err := mockDatastore.DeleteDiscussionTags(ctx, tx, tagObject)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when query execution returns an error", func() {
			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(deleteDiscussionTagsString)
			mock.ExpectQuery(deleteDiscussionTagsString).WithArgs(tagObject.ID, tagObject.Tag).WillReturnError(fmt.Errorf("error"))

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.DeleteDiscussionTags(ctx, tx, tagObject)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when put post succeeds and returns an object", func() {
			rs := sqlmock.NewRows([]string{"id", "tag", "created_at", "deleted_at"}).
				AddRow(tagObject.ID, tagObject.Tag, tagObject.CreatedAt, tagObject.DeletedAt)

			mock.ExpectBegin()
			mockPreparedStatements(mock)
			mock.ExpectPrepare(deleteDiscussionTagsString)
			mock.ExpectQuery(deleteDiscussionTagsString).WithArgs(tagObject.ID, tagObject.Tag).WillReturnRows(rs)

			tx, err := mockDatastore.BeginTx(ctx)
			resp, err := mockDatastore.DeleteDiscussionTags(ctx, tx, tagObject)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, &tagObject)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDelphisDB_DiscussionAutoPostIterCollect(t *testing.T) {
	ctx := context.Background()
	discussionID := "discussion1"

	dapObj := model.DiscussionAutoPost{
		ID:          discussionID,
		IdleMinutes: 120,
	}

	Convey("DiscussionAutoPostIterCollect", t, func() {
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
			iter := &autoPostDiscussionIter{
				err: fmt.Errorf("error"),
			}

			resp, err := mockDatastore.DiscussionAutoPostIterCollect(ctx, iter)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has results and returns slice of DiscussionAutoPost", func() {
			rs := sqlmock.NewRows([]string{"id", "idle_minutes"}).
				AddRow(dapObj.ID, dapObj.IdleMinutes).
				AddRow(dapObj.ID, dapObj.IdleMinutes)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := &autoPostDiscussionIter{
				ctx:  ctx,
				rows: rs1,
			}

			resp, err := mockDatastore.DiscussionAutoPostIterCollect(ctx, iter)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []*model.DiscussionAutoPost{&dapObj, &dapObj})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

	})
}

func TestDelphisDB_DiscussionIterCollect(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	emptyString := ""
	discObj := model.Discussion{
		ID:            "discussion1",
		CreatedAt:     now,
		DeletedAt:     nil,
		Title:         "test",
		AnonymityType: "",
		ModeratorID:   &modID,
		AutoPost:      false,
		IdleMinutes:   120,
		IconURL:       &emptyString,
	}

	Convey("DiscussionIterCollect", t, func() {
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
			iter := &discussionIter{
				err: fmt.Errorf("error"),
			}

			resp, err := mockDatastore.DiscussionIterCollect(ctx, iter)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has results and returns slice of Discussions", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title",
				"anonymity_type", "moderator_id", "auto_post", "icon_url", "idle_minutes", "description",
				"title_history", "description_history", "discussion_joinability", "last_post_id", "last_post_created_at",
				"shuffle_count", "lock_status"}).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt,
					discObj.Title, discObj.AnonymityType, discObj.ModeratorID, discObj.AutoPost,
					discObj.IconURL, discObj.IdleMinutes, discObj.Description, discObj.TitleHistory,
					discObj.DescriptionHistory, discObj.DiscussionJoinability, discObj.LastPostID, discObj.LastPostCreatedAt,
					discObj.ShuffleCount, discObj.LockStatus).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt,
					discObj.Title, discObj.AnonymityType, discObj.ModeratorID, discObj.AutoPost,
					discObj.IconURL, discObj.IdleMinutes, discObj.Description, discObj.TitleHistory,
					discObj.DescriptionHistory, discObj.DiscussionJoinability, discObj.LastPostID, discObj.LastPostCreatedAt,
					discObj.ShuffleCount, discObj.LockStatus)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := &discussionIter{
				ctx:  ctx,
				rows: rs1,
			}

			resp, err := mockDatastore.DiscussionIterCollect(ctx, iter)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []*model.Discussion{&discObj, &discObj})
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestAutoPostDiscussionIter_Next(t *testing.T) {
	ctx := context.Background()
	discussionID := "discussion1"
	discAP := model.DiscussionAutoPost{
		ID:          discussionID,
		IdleMinutes: 120,
	}

	emptyAP := model.DiscussionAutoPost{}

	Convey("AutoPostDiscussionIter_Next", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := autoPostDiscussionIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Next(&emptyAP), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has a context error passed in", func() {
			ctx1, cancelFunc := context.WithCancel(ctx)
			cancelFunc()
			iter := autoPostDiscussionIter{
				ctx: ctx1,
			}

			So(iter.Next(&emptyAP), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has no more rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"id", "idle_minutes"})

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := autoPostDiscussionIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyAP), ShouldBeFalse)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on scan", func() {
			rs := sqlmock.NewRows([]string{"id"}).
				AddRow(discAP.ID)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := autoPostDiscussionIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyAP), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"id", "idle_minutes"}).
				AddRow(discAP.ID, discAP.IdleMinutes).
				AddRow(discAP.ID, discAP.IdleMinutes)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := autoPostDiscussionIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyAP), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestAutoPostDiscussionIter_Close(t *testing.T) {
	ctx := context.Background()

	Convey("AutoPostDiscussionIter_Close", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := autoPostDiscussionIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on rows.Close", func() {
			rs := sqlmock.NewRows([]string{"id", "idle_minutes"}).CloseError(fmt.Errorf("error"))

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := autoPostDiscussionIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDiscussionIter_Next(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	modID := "modID"
	emptyString := ""
	discObj := model.Discussion{
		ID:            "discussion1",
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
		Title:         "test",
		AnonymityType: "",
		ModeratorID:   &modID,
		AutoPost:      false,
		IdleMinutes:   120,
		IconURL:       &emptyString,
	}
	emptyDisc := model.Discussion{}

	Convey("DiscussionIter_Next", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := discussionIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Next(&emptyDisc), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has a context error passed in", func() {
			ctx1, cancelFunc := context.WithCancel(ctx)
			cancelFunc()
			iter := discussionIter{
				ctx: ctx1,
			}

			So(iter.Next(&emptyDisc), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has no more rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title",
				"anonymity_type", "moderator_id", "auto_post", "icon_url", "idle_minutes", "description",
				"title_history", "description_history", "discussion_joinability"})

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := discussionIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyDisc), ShouldBeFalse)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on scan", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title",
				"anonymity_type", "moderator_id", "auto_post", "icon_url", "idle_minutes", "description",
				"title_history", "description_history"}).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt,
					discObj.Title, discObj.AnonymityType, discObj.ModeratorID, discObj.AutoPost,
					discObj.IconURL, discObj.IdleMinutes, discObj.Description, discObj.TitleHistory,
					discObj.DescriptionHistory)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := discussionIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyDisc), ShouldBeFalse)
			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator has rows to iterate over", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title",
				"anonymity_type", "moderator_id", "auto_post", "icon_url", "idle_minutes", "description",
				"title_history", "description_history", "discussion_joinability", "last_post_id", "last_post_created_at",
				"shuffle_count", "lock_status"}).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt,
					discObj.Title, discObj.AnonymityType, discObj.ModeratorID, discObj.AutoPost,
					discObj.IconURL, discObj.IdleMinutes, discObj.Description, discObj.TitleHistory,
					discObj.DescriptionHistory, discObj.DiscussionJoinability, discObj.LastPostID, discObj.LastPostCreatedAt,
					discObj.ShuffleCount, discObj.LockStatus).
				AddRow(discObj.ID, discObj.CreatedAt, discObj.UpdatedAt, discObj.DeletedAt,
					discObj.Title, discObj.AnonymityType, discObj.ModeratorID, discObj.AutoPost,
					discObj.IconURL, discObj.IdleMinutes, discObj.Description, discObj.TitleHistory,
					discObj.DescriptionHistory, discObj.DiscussionJoinability, discObj.LastPostID, discObj.LastPostCreatedAt,
					discObj.ShuffleCount, discObj.LockStatus)

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := discussionIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Next(&emptyDisc), ShouldBeTrue)
			So(iter.Close(), ShouldBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

func TestDiscussionIter_Close(t *testing.T) {
	ctx := context.Background()

	Convey("DiscussionIter_Close", t, func() {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.Nil(t, err, "Failed setting up sqlmock db")

		defer db.Close()

		Convey("when the iterator has an error passed in", func() {
			iter := discussionIter{
				err: fmt.Errorf("error"),
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("when the iterator errors on rows.Close", func() {
			rs := sqlmock.NewRows([]string{"id", "created_at", "title", "anonymity_type",
				"moderator_id", "auto_post", "icon_url", "idle_minutes"}).CloseError(fmt.Errorf("error"))

			// Convert mocked rows to sql.Rows
			mock.ExpectQuery("SELECT").WillReturnRows(rs)
			rs1, _ := db.Query("SELECT")

			iter := discussionIter{
				ctx:  ctx,
				rows: rs1,
			}

			So(iter.Close(), ShouldNotBeNil)
			So(mock.ExpectationsWereMet(), ShouldBeNil)
		})
	})
}

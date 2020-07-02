package backend

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/auth"
	"github.com/nedrocks/delphisbe/internal/cache"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/nedrocks/delphisbe/internal/util"
	"github.com/nedrocks/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

type mockDiscAutoPostIter struct{}

func (m *mockDiscAutoPostIter) Next(discussion *model.DiscussionAutoPost) bool { return true }
func (m *mockDiscAutoPostIter) Close() error                                   { return fmt.Errorf("error") }

type mockTagIter struct{}

func (m *mockTagIter) Next(tag *model.Tag) bool { return true }
func (m *mockTagIter) Close() error             { return fmt.Errorf("error") }

func TestDelphisBackend_CreateNewDiscussion(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := "userID"
	anonymityType := model.AnonymityTypeStrong
	title := "test title"
	publicAccess := true
	modID := "modID"
	profileID := "profileID"

	userObj := model.User{
		ID:          userID,
		CreatedAt:   now,
		UpdatedAt:   now,
		UserProfile: &model.UserProfile{ID: profileID},
	}

	modObj := model.Moderator{
		ID: modID,
		UserProfile: &model.UserProfile{
			ID: profileID,
		},
	}

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
		PublicAccess:  false,
		IconURL:       "",
	}

	flairID := "flairID"
	templateID := "templateID"
	flairObj := model.Flair{
		ID:         flairID,
		TemplateID: templateID,
		CreatedAt:  now,
		UpdatedAt:  now,
		UserID:     userID,
	}

	discussionID := "discussionID"
	viewerID := "viewerID"
	postID := "post1"

	viewerObj := model.Viewer{
		ID:               viewerID,
		CreatedAt:        now,
		UpdatedAt:        now,
		DiscussionID:     &discussionID,
		LastViewedPostID: &postID,
		UserID:           &userID,
	}

	parObj := model.Participant{
		ID: "participantID",
	}

	tx := sql.Tx{}

	Convey("CreateNewDiscussion", t, func() {
		now := time.Now()
		cacheObj := cache.NewInMemoryCache()
		authObj := auth.NewDelphisAuth(nil)
		mockDB := &mocks.Datastore{}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            authObj,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when create moderator errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("CreateModerator", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.CreateNewDiscussion(ctx, &userObj, anonymityType, title, publicAccess)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when create discussion errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("CreateModerator", ctx, mock.Anything).Return(&modObj, nil)
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.CreateNewDiscussion(ctx, &userObj, anonymityType, title, publicAccess)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when create participant errors outs", func() {
			expectedError := fmt.Errorf("Some Error")

			mockDB.On("CreateModerator", ctx, mock.Anything).Return(&modObj, nil)
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(&discObj, nil)

			// Create participant functions
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.CreateNewDiscussion(ctx, &userObj, anonymityType, title, publicAccess)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when upsert links errors outs", func() {
			expectedError := fmt.Errorf("Some Error")

			mockDB.On("CreateModerator", ctx, mock.Anything).Return(&modObj, nil)
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(&discObj, nil)

			// Create participant functions
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(&userObj, nil)
			mockDB.On("GetTotalParticipantCountByDiscussionID", ctx, mock.Anything).Return(10)
			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return([]*model.Flair{&flairObj}, nil)
			mockDB.On("UpsertViewer", ctx, mock.Anything).Return(&viewerObj, nil)
			mockDB.On("UpsertParticipant", ctx, mock.Anything).Return(&parObj, nil)
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)

			// Create post functions
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutPostContent", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("PutPost", ctx, mock.Anything, mock.Anything).Return(&model.Post{ID: "post123"}, nil)
			mockDB.On("PutActivity", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)
			mockDB.On("GetDiscussionByID", ctx, mock.Anything).Return(&discObj, nil)
			mockDB.On("GetParticipantsByDiscussionID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)

			mockDB.On("BeginTx", ctx).Return(tx, nil)
			mockDB.On("UpsertInviteLinksByDiscussionID", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.CreateNewDiscussion(ctx, &userObj, anonymityType, title, publicAccess)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when the discussion is created successfully", func() {

			mockDB.On("CreateModerator", ctx, mock.Anything).Return(&modObj, nil)
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(&discObj, nil)

			// Create participant functions
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(&userObj, nil)
			mockDB.On("GetTotalParticipantCountByDiscussionID", ctx, mock.Anything).Return(10)
			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return([]*model.Flair{&flairObj}, nil)
			mockDB.On("UpsertViewer", ctx, mock.Anything).Return(&viewerObj, nil)
			mockDB.On("UpsertParticipant", ctx, mock.Anything).Return(&parObj, nil)
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)

			//// Create post functions
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutPostContent", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("PutPost", ctx, mock.Anything, mock.Anything).Return(&model.Post{ID: "post123"}, nil)
			mockDB.On("PutActivity", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)
			mockDB.On("GetDiscussionByID", ctx, mock.Anything).Return(&discObj, nil)
			mockDB.On("GetParticipantsByDiscussionID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)

			mockDB.On("BeginTx", ctx).Return(tx, nil)
			mockDB.On("UpsertInviteLinksByDiscussionID", ctx, mock.Anything, mock.Anything).Return(
				&model.DiscussionLinkAccess{DiscussionID: discussionID}, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.CreateNewDiscussion(ctx, &userObj, anonymityType, title, publicAccess)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_UpdateDiscussion(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussionID"

	modID := "modID"

	title := "test title"
	discInput := model.DiscussionInput{
		Title: &title,
	}

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
		PublicAccess:  false,
		IconURL:       "",
	}

	Convey("UpdateDiscussion", t, func() {
		now := time.Now()
		cacheObj := cache.NewInMemoryCache()
		authObj := auth.NewDelphisAuth(nil)
		mockDB := &mocks.Datastore{}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            authObj,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when get discussion by id errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetDiscussionByID", ctx, discussionID).Return(nil, expectedError)

			resp, err := backendObj.UpdateDiscussion(ctx, discussionID, discInput)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the discussion fails to update", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetDiscussionByID", ctx, discussionID).Return(&discObj, nil)
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.UpdateDiscussion(ctx, discussionID, discInput)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the discussion is updated successfully", func() {
			mockDB.On("GetDiscussionByID", ctx, discussionID).Return(&discObj, nil)
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(&discObj, nil)

			resp, err := backendObj.UpdateDiscussion(ctx, discussionID, discInput)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_GetDiscussionByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussionID"

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
		PublicAccess:  false,
		IconURL:       "",
	}

	Convey("GetDiscussionByID", t, func() {
		now := time.Now()
		cacheObj := cache.NewInMemoryCache()
		authObj := auth.NewDelphisAuth(nil)
		mockDB := &mocks.Datastore{}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            authObj,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when the query errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetDiscussionByID", ctx, discussionID).Return(nil, expectedError)

			resp, err := backendObj.GetDiscussionByID(ctx, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetDiscussionByID", ctx, discussionID).Return(&discObj, nil)

			resp, err := backendObj.GetDiscussionByID(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &discObj)
		})
	})
}

func TestDelphisBackend_GetDiscussionsByIDs(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussionID"

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
		PublicAccess:  false,
		IconURL:       "",
	}

	Convey("GetDiscussionsByIDs", t, func() {
		now := time.Now()
		cacheObj := cache.NewInMemoryCache()
		authObj := auth.NewDelphisAuth(nil)
		mockDB := &mocks.Datastore{}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            authObj,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when the query errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetDiscussionsByIDs", ctx, []string{discussionID}).Return(nil, expectedError)

			resp, err := backendObj.GetDiscussionsByIDs(ctx, []string{discussionID})

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			dbResp := map[string]*model.Discussion{
				discObj.ID: &discObj,
			}

			mockDB.On("GetDiscussionsByIDs", ctx, []string{discussionID}).Return(dbResp, nil)

			resp, err := backendObj.GetDiscussionsByIDs(ctx, []string{discussionID})

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, dbResp)
		})
	})
}

func TestDelphisBackend_GetDiscussionByModeratorID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := "discussionID"

	modID := "modID"

	discObj := model.Discussion{
		ID:            discussionID,
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
		Title:         "test",
		AnonymityType: "",
		ModeratorID:   &modID,
		AutoPost:      false,
		IdleMinutes:   120,
		PublicAccess:  false,
		IconURL:       "",
	}

	Convey("GetDiscussionByModeratorID", t, func() {
		now := time.Now()
		cacheObj := cache.NewInMemoryCache()
		authObj := auth.NewDelphisAuth(nil)
		mockDB := &mocks.Datastore{}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            authObj,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when the query errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetDiscussionByModeratorID", ctx, modID).Return(nil, expectedError)

			resp, err := backendObj.GetDiscussionByModeratorID(ctx, modID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetDiscussionByModeratorID", ctx, modID).Return(&discObj, nil)

			resp, err := backendObj.GetDiscussionByModeratorID(ctx, modID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &discObj)
		})
	})
}

func TestDelphisBackend_GetDiscussionsForAutoPost(t *testing.T) {
	ctx := context.Background()

	apObj := model.DiscussionAutoPost{
		ID:          "id",
		IdleMinutes: 120,
	}

	Convey("GetDiscussionsForAutoPost", t, func() {
		now := time.Now()
		cacheObj := cache.NewInMemoryCache()
		authObj := auth.NewDelphisAuth(nil)
		mockDB := &mocks.Datastore{}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            authObj,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when the query errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetDiscussionsAutoPost", ctx).Return(&mockDiscAutoPostIter{})
			mockDB.On("DiscussionAutoPostIterCollect", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetDiscussionsForAutoPost(ctx)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetDiscussionsAutoPost", ctx).Return(&mockDiscAutoPostIter{})
			mockDB.On("DiscussionAutoPostIterCollect", ctx, mock.Anything).Return([]*model.DiscussionAutoPost{&apObj}, nil)

			resp, err := backendObj.GetDiscussionsForAutoPost(ctx)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.DiscussionAutoPost{&apObj})
		})
	})
}

func TestDelphisBackend_ListDiscussions(t *testing.T) {
	ctx := context.Background()

	dcObj := model.DiscussionsConnection{
		IDs:   []string{"s1", "s2"},
		From:  0,
		To:    0,
		Edges: nil,
	}

	Convey("ListDiscussions", t, func() {
		now := time.Now()
		cacheObj := cache.NewInMemoryCache()
		authObj := auth.NewDelphisAuth(nil)
		mockDB := &mocks.Datastore{}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            authObj,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when the query errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("ListDiscussions", ctx).Return(nil, expectedError)

			resp, err := backendObj.ListDiscussions(ctx)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("ListDiscussions", ctx).Return(&dcObj, nil)

			resp, err := backendObj.ListDiscussions(ctx)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &dcObj)
		})
	})
}

func TestDelphisBackend_GetDiscussionTags(t *testing.T) {
	ctx := context.Background()

	discussionID := "discussionID"

	tagObj := model.Tag{
		ID: "tagID",
	}

	Convey("GetDiscussionTags", t, func() {
		now := time.Now()
		cacheObj := cache.NewInMemoryCache()
		authObj := auth.NewDelphisAuth(nil)
		mockDB := &mocks.Datastore{}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            authObj,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when the query errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetDiscussionTags", ctx, discussionID).Return(&mockTagIter{})
			mockDB.On("TagIterCollect", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetDiscussionTags(ctx, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetDiscussionTags", ctx, discussionID).Return(&mockTagIter{})
			mockDB.On("TagIterCollect", ctx, mock.Anything).Return([]*model.Tag{&tagObj}, nil)

			resp, err := backendObj.GetDiscussionTags(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.Tag{&tagObj})
		})
	})
}

func TestDelphisBackend_PutDiscussionTags(t *testing.T) {
	ctx := context.Background()

	discussionID := "discussionID"
	tagID := "tagID"
	tagObj := model.Tag{
		ID:  discussionID,
		Tag: tagID,
	}

	tags := []string{tagObj.Tag}
	tx := sql.Tx{}

	Convey("PutDiscussionTags", t, func() {
		now := time.Now()
		cacheObj := cache.NewInMemoryCache()
		authObj := auth.NewDelphisAuth(nil)
		mockDB := &mocks.Datastore{}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            authObj,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when no tags are passed in", func() {
			resp, err := backendObj.PutDiscussionTags(ctx, discussionID, nil)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when BeginTx errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(nil, expectedError)

			resp, err := backendObj.PutDiscussionTags(ctx, discussionID, tags)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when PutDiscussionTags errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutDiscussionTags", ctx, mock.Anything, tagObj).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.PutDiscussionTags(ctx, discussionID, tags)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when PutDiscussionTags and rollback errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutDiscussionTags", ctx, mock.Anything, tagObj).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.PutDiscussionTags(ctx, discussionID, tags)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when PutDiscussionTags succeeds and CommitTx errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutDiscussionTags", ctx, mock.Anything, tagObj).Return(&tagObj, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.PutDiscussionTags(ctx, discussionID, tags)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when PutDiscussionTags succeeds", func() {
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutDiscussionTags", ctx, mock.Anything, tagObj).Return(&tagObj, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.PutDiscussionTags(ctx, discussionID, tags)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.Tag{&tagObj})
		})
	})
}

func TestDelphisBackend_DeleteDiscussionTags(t *testing.T) {
	ctx := context.Background()

	discussionID := "discussionID"
	tagID := "tagID"
	tagObj := model.Tag{
		ID:  discussionID,
		Tag: tagID,
	}

	tags := []string{tagObj.Tag}
	tx := sql.Tx{}

	Convey("DeleteDiscussionTags", t, func() {
		now := time.Now()
		cacheObj := cache.NewInMemoryCache()
		authObj := auth.NewDelphisAuth(nil)
		mockDB := &mocks.Datastore{}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            authObj,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when no tags are passed in", func() {
			resp, err := backendObj.DeleteDiscussionTags(ctx, discussionID, nil)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when BeginTx errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(nil, expectedError)

			resp, err := backendObj.DeleteDiscussionTags(ctx, discussionID, tags)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when PutDiscussionTags errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("DeleteDiscussionTags", ctx, mock.Anything, tagObj).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.DeleteDiscussionTags(ctx, discussionID, tags)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when PutDiscussionTags and rollback errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("DeleteDiscussionTags", ctx, mock.Anything, tagObj).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.DeleteDiscussionTags(ctx, discussionID, tags)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when PutDiscussionTags succeeds and CommitTx errors outs", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("DeleteDiscussionTags", ctx, mock.Anything, tagObj).Return(&tagObj, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.DeleteDiscussionTags(ctx, discussionID, tags)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when PutDiscussionTags succeeds", func() {
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("DeleteDiscussionTags", ctx, mock.Anything, tagObj).Return(&tagObj, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.DeleteDiscussionTags(ctx, discussionID, tags)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.Tag{&tagObj})
		})
	})
}

func TestDelphisBackend_UpdateDiscussionObj(t *testing.T) {
	anonymityType := model.AnonymityTypeStrong
	title := "test title"
	autoPost := true
	idleMinutes := 120
	publicAccess := true
	iconUrl := "http://test.com"

	discInput := model.DiscussionInput{
		AnonymityType: &anonymityType,
		Title:         &title,
		AutoPost:      &autoPost,
		IdleMinutes:   &idleMinutes,
		PublicAccess:  &publicAccess,
		IconURL:       &iconUrl,
	}

	disc := model.Discussion{}

	Convey("UpdateDiscussionObj", t, func() {
		Convey("when it updates the discussion object successfully", func() {
			updateDiscussionObj(&disc, discInput)

			So(disc.AnonymityType, ShouldResemble, *discInput.AnonymityType)
			So(disc.Title, ShouldResemble, *discInput.Title)
			So(disc.AutoPost, ShouldResemble, *discInput.AutoPost)
			So(disc.IdleMinutes, ShouldResemble, *discInput.IdleMinutes)
			So(disc.PublicAccess, ShouldResemble, *discInput.PublicAccess)
			So(disc.IconURL, ShouldResemble, *discInput.IconURL)

		})
	})
}

func TestDelphisBackend_DedupeDiscussions(t *testing.T) {
	disc1 := model.Discussion{
		ID:    "id1",
		Title: "test1",
	}

	disc2 := model.Discussion{
		ID:    "id2",
		Title: "test2",
	}

	Convey("DedupeDiscussions", t, func() {
		Convey("when it dedupes the discussion objects successfully", func() {
			resp := dedupeDiscussions([]*model.Discussion{&disc1, &disc2, &disc1})

			So(resp, ShouldResemble, []*model.Discussion{&disc1, &disc2})

		})
	})
}

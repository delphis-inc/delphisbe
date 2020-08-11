package backend

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/delphis-inc/delphisbe/internal/backend/test_utils"

	"github.com/stretchr/testify/mock"

	"github.com/delphis-inc/delphisbe/graph/model"

	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/cache"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/delphis-inc/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

type mockPostIter struct{}

func (m *mockPostIter) Next(post *model.Post) bool { return true }
func (m *mockPostIter) Close() error               { return fmt.Errorf("error") }

func TestDelphisBackend_DeletePostByID(t *testing.T) {
	ctx := context.Background()

	discussionID := test_utils.DiscussionID
	participantID := test_utils.ParticipantID

	postObj := test_utils.TestPost()
	participantObj := test_utils.TestParticipant()
	participantUserID := "participant_user_id"
	participantObj.UserID = &participantUserID
	userObj := test_utils.TestUser()
	profile := test_utils.TestUserProfile()
	discObj := test_utils.TestDiscussion()
	moderatorObj := test_utils.TestModerator()
	userProfileObj := test_utils.TestUserProfile()

	userObj.UserProfile = &profile

	//tx := sql.Tx{}

	Convey("DeletePostByID", t, func() {
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

		Convey("when the discussion is not found", func() {
			Convey("because it returns nil", func() {
				mockDB.On("GetDiscussionByID", ctx, discussionID).Return(nil, nil)

				resp, err := backendObj.DeletePostByID(ctx, discussionID, postObj.ID, userObj.ID)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})
			Convey("because it reutrns an error", func() {
				mockDB.On("GetDiscussionByID", ctx, discussionID).Return(&discObj, fmt.Errorf("sth"))

				resp, err := backendObj.DeletePostByID(ctx, discussionID, postObj.ID, userObj.ID)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})
		})

		mockDB.On("GetDiscussionByID", ctx, discussionID).Return(&discObj, nil)

		Convey("when the moderator is not found", func() {
			Convey("because it returns nil", func() {
				mockDB.On("GetModeratorByID", ctx, *discObj.ModeratorID).Return(nil, nil)

				resp, err := backendObj.DeletePostByID(ctx, discussionID, postObj.ID, userObj.ID)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})
			Convey("because it returns an error", func() {
				mockDB.On("GetModeratorByID", ctx, *discObj.ModeratorID).Return(nil, fmt.Errorf("sth"))

				resp, err := backendObj.DeletePostByID(ctx, discussionID, postObj.ID, userObj.ID)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})
		})

		mockDB.On("GetModeratorByID", ctx, *discObj.ModeratorID).Return(&moderatorObj, nil)

		Convey("when user profile is not found", func() {
			Convey("because it returns nil", func() {
				mockDB.On("GetUserProfileByID", ctx, userProfileObj.ID).Return(nil, nil)

				resp, err := backendObj.DeletePostByID(ctx, discussionID, postObj.ID, userObj.ID)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})
			Convey("because it returns an error", func() {
				mockDB.On("GetUserProfileByID", ctx, userProfileObj.ID).Return(nil, fmt.Errorf("sth"))

				resp, err := backendObj.DeletePostByID(ctx, discussionID, postObj.ID, userObj.ID)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})
		})

		mockDB.On("GetUserProfileByID", ctx, userProfileObj.ID).Return(&userProfileObj, nil)

		Convey("when the post is not found", func() {
			Convey("because it returns nil", func() {
				mockDB.On("GetPostByID", ctx, postObj.ID).Return(nil, nil)

				resp, err := backendObj.DeletePostByID(ctx, discussionID, postObj.ID, userObj.ID)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})

			Convey("because it returns an error", func() {
				mockDB.On("GetPostByID", ctx, postObj.ID).Return(&postObj, fmt.Errorf("sth"))

				resp, err := backendObj.DeletePostByID(ctx, discussionID, postObj.ID, userObj.ID)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})
		})

		mockDB.On("GetPostByID", ctx, postObj.ID).Return(&postObj, nil)

		Convey("when the participant is not found", func() {
			Convey("because it returns nil", func() {
				mockDB.On("GetParticipantByID", ctx, participantID).Return(nil, nil)

				resp, err := backendObj.DeletePostByID(ctx, discussionID, postObj.ID, userObj.ID)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})
			Convey("because it returns an error", func() {
				mockDB.On("GetParticipantByID", ctx, participantID).Return(&participantObj, fmt.Errorf("sth"))

				resp, err := backendObj.DeletePostByID(ctx, discussionID, postObj.ID, userObj.ID)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})
		})

		mockDB.On("GetParticipantByID", ctx, participantID).Return(&participantObj, nil)

		Convey("when userID is not moderator or participant", func() {
			resp, err := backendObj.DeletePostByID(ctx, discussionID, postObj.ID, "baduserid")

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when post is already deleted", func() {
			postObj.DeletedAt = &now

			resp, err := backendObj.DeletePostByID(ctx, discussionID, postObj.ID, *userProfileObj.UserID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})

		postObj.DeletedAt = nil

		Convey("when user is moderator", func() {
			mockDB.On("DeletePostByID", ctx, postObj.ID, model.PostDeletedReasonModeratorRemoved).Return(&postObj, nil)

			resp, err := backendObj.DeletePostByID(ctx, discussionID, postObj.ID, *userProfileObj.UserID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})

		Convey("when user is participant", func() {
			mockDB.On("DeletePostByID", ctx, postObj.ID, model.PostDeletedReasonParticipantRemoved).Return(&postObj, nil)

			resp, err := backendObj.DeletePostByID(ctx, discussionID, postObj.ID, *participantObj.UserID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_CreatePost(t *testing.T) {
	ctx := context.Background()

	discussionID := test_utils.DiscussionID
	userID := test_utils.UserID
	participantID := test_utils.ParticipantID

	postObj := test_utils.TestPost()
	postInputObj := test_utils.TestPostContentInput()
	userObj := test_utils.TestUser()
	profile := test_utils.TestUserProfile()
	discObj := test_utils.TestDiscussion()

	userObj.UserProfile = &profile

	tx := sql.Tx{}

	Convey("CreatePost", t, func() {
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

		Convey("when the post type is not passed in", func() {
			tempPostInputObj := postInputObj
			tempPostInputObj.PostType = ""
			resp, err := backendObj.CreatePost(ctx, discussionID, userID, participantID, tempPostInputObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when BeginTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(nil, expectedError)

			resp, err := backendObj.CreatePost(ctx, discussionID, userID, participantID, postInputObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when PutPostContent errors out and Rollback fails", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutPostContent", ctx, mock.Anything, mock.Anything).Return(expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.CreatePost(ctx, discussionID, userID, participantID, postInputObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when PutPostContent errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutPostContent", ctx, mock.Anything, mock.Anything).Return(expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.CreatePost(ctx, discussionID, userID, participantID, postInputObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when PutPost errors out and Rollback fails", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutPostContent", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("PutPost", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.CreatePost(ctx, discussionID, userID, participantID, postInputObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when PutPost errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutPostContent", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("PutPost", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.CreatePost(ctx, discussionID, userID, participantID, postInputObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when PutActivity errors out and we fail to commit", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutPostContent", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("PutPost", ctx, mock.Anything, mock.Anything).Return(&postObj, nil)
			mockDB.On("PutActivity", ctx, mock.Anything, mock.Anything).Return(expectedError)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.CreatePost(ctx, discussionID, userID, participantID, postInputObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when GetDiscussionByID errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutPostContent", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("PutPost", ctx, mock.Anything, mock.Anything).Return(&postObj, nil)
			mockDB.On("PutActivity", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)
			mockDB.On("GetDiscussionByID", ctx, discussionID).Return(nil, expectedError)

			resp, err := backendObj.CreatePost(ctx, discussionID, userID, participantID, postInputObj)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})

		Convey("when GetUserDevicesByUserID errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutPostContent", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("PutPost", ctx, mock.Anything, mock.Anything).Return(&postObj, nil)
			mockDB.On("PutActivity", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)
			mockDB.On("GetDiscussionByID", ctx, discussionID).Return(&discObj, nil)
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(&discObj, nil)
			mockDB.On("GetDUAForEverythingNotifications", ctx, discussionID, userID).Return(nil)
			mockDB.On("DuaIterCollect", ctx, mock.Anything).Return(nil, nil)
			mockDB.On("GetUserDevicesByUserID", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.CreatePost(ctx, discussionID, userID, participantID, postInputObj)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})

		Convey("when post succeeds", func() {
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutPostContent", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("PutPost", ctx, mock.Anything, mock.Anything).Return(&postObj, nil)
			mockDB.On("PutActivity", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)
			mockDB.On("GetDiscussionByID", ctx, discussionID).Return(&discObj, nil)
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(&discObj, nil)
			mockDB.On("GetDUAForEverythingNotifications", ctx, discussionID, userID).Return(nil)
			mockDB.On("DuaIterCollect", ctx, mock.Anything).Return(nil, nil)
			mockDB.On("GetUserDevicesByUserID", ctx, mock.Anything).Return(nil, nil)

			resp, err := backendObj.CreatePost(ctx, discussionID, userID, participantID, postInputObj)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_CreateAlertPost(t *testing.T) {
	ctx := context.Background()

	discussionID := test_utils.DiscussionID

	postObj := test_utils.TestPost()
	userObj := test_utils.TestUser()
	modObj := test_utils.TestModerator()
	profile := test_utils.TestUserProfile()
	discObj := test_utils.TestDiscussion()
	parObj := test_utils.TestParticipant()

	userObj.UserProfile = &profile
	modObj.UserProfile = &profile

	tx := sql.Tx{}

	Convey("CreateAlertPost", t, func() {
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

		Convey("when the concierge post has been added to a discussion", func() {
			tempUserObj := userObj
			tempUserObj.ID = model.ConciergeUser
			resp, err := backendObj.CreateAlertPost(ctx, discussionID, &tempUserObj, false)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when GetParticipantsByDiscussionIDUserID errors out", func() {
			expectedError := fmt.Errorf("Some error")
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, discussionID, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.CreateAlertPost(ctx, discussionID, &userObj, false)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when GetParticipantsByDiscussionIDUserID returns nil", func() {
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, discussionID, mock.Anything).Return(nil, nil)

			resp, err := backendObj.CreateAlertPost(ctx, discussionID, &userObj, false)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when GetParticipantsByDiscussionIDUserID for an anonymous user returns nil", func() {
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, discussionID, mock.Anything).Return(nil, nil)

			resp, err := backendObj.CreateAlertPost(ctx, discussionID, &userObj, true)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when the alert post is created successfully", func() {
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, discussionID, mock.Anything).Return([]model.Participant{parObj}, nil)

			// Create post functions
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutPostContent", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("PutPost", ctx, mock.Anything, mock.Anything).Return(&postObj, nil)
			mockDB.On("PutActivity", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)
			mockDB.On("GetDiscussionByID", ctx, discussionID).Return(&discObj, nil)
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(&discObj, nil)
			mockDB.On("GetDUAForEverythingNotifications", ctx, discussionID, mock.Anything).Return(nil)
			mockDB.On("DuaIterCollect", ctx, mock.Anything).Return(nil, nil)
			mockDB.On("GetUserDevicesByUserID", ctx, mock.Anything).Return(nil, nil)

			resp, err := backendObj.CreateAlertPost(ctx, discussionID, &userObj, false)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_PostImportedContent(t *testing.T) {
	ctx := context.Background()

	discussionID := test_utils.DiscussionID
	userID := test_utils.UserID
	participantID := test_utils.ParticipantID
	contentID := test_utils.ContentID

	dripType := model.AutoDrip
	tags := []string{"tag1", "tag2"}

	postObj := test_utils.TestPost()
	userObj := test_utils.TestUser()
	profile := test_utils.TestUserProfile()
	discObj := test_utils.TestDiscussion()
	contentObj := test_utils.TestImportedContent()
	cqrObj := test_utils.TestContentQueueRecord()

	userObj.UserProfile = &profile

	tx := sql.Tx{}

	Convey("PostImportedContent", t, func() {
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

		Convey("when GetImportedContentByID errors out", func() {
			expectedError := fmt.Errorf("Some error")
			mockDB.On("GetImportedContentByID", ctx, contentID).Return(nil, expectedError)

			resp, err := backendObj.PostImportedContent(ctx, userID, participantID, discussionID, contentID, &now, tags, dripType)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when CreatePost errors out", func() {
			expectedError := fmt.Errorf("Some error")
			mockDB.On("GetImportedContentByID", ctx, contentID).Return(&contentObj, nil)
			mockDB.On("BeginTx", ctx).Return(nil, expectedError)

			resp, err := backendObj.PostImportedContent(ctx, userID, participantID, discussionID, contentID, &now, tags, dripType)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when PutImportedContentQueue errors out", func() {
			expectedError := fmt.Errorf("Some error")
			mockDB.On("GetImportedContentByID", ctx, contentID).Return(&contentObj, nil)

			// Create post functions
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutPostContent", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("PutPost", ctx, mock.Anything, mock.Anything).Return(&postObj, nil)
			mockDB.On("PutActivity", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)
			mockDB.On("GetDiscussionByID", ctx, discussionID).Return(&discObj, nil)
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(&discObj, nil)
			mockDB.On("GetDUAForEverythingNotifications", ctx, discussionID, userID).Return(nil)
			mockDB.On("DuaIterCollect", ctx, mock.Anything).Return(nil, nil)
			mockDB.On("GetUserDevicesByUserID", ctx, mock.Anything).Return(nil, nil)

			// PutImportedContentQueue
			mockDB.On("PutImportedContentDiscussionQueue", ctx, discussionID, contentID, mock.Anything, tags).Return(nil, expectedError)

			resp, err := backendObj.PostImportedContent(ctx, userID, participantID, discussionID, contentID, &now, tags, dripType)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when PostImportedContent completes successfully", func() {
			mockDB.On("GetImportedContentByID", ctx, contentID).Return(&contentObj, nil)

			// Create post functions
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutPostContent", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("PutPost", ctx, mock.Anything, mock.Anything).Return(&postObj, nil)
			mockDB.On("PutActivity", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)
			mockDB.On("GetDiscussionByID", ctx, discussionID).Return(&discObj, nil)
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(&discObj, nil)
			mockDB.On("GetDUAForEverythingNotifications", ctx, discussionID, userID).Return(nil)
			mockDB.On("DuaIterCollect", ctx, mock.Anything).Return(nil, nil)
			mockDB.On("GetUserDevicesByUserID", ctx, mock.Anything).Return(nil, nil)

			// PutImportedContentQueue
			mockDB.On("PutImportedContentDiscussionQueue", ctx, discussionID, contentID, mock.Anything, tags).Return(&cqrObj, nil)

			resp, err := backendObj.PostImportedContent(ctx, userID, participantID, discussionID, contentID, &now, tags, dripType)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_PutImportedContentQueue(t *testing.T) {
	ctx := context.Background()

	discussionID := test_utils.DiscussionID
	contentID := test_utils.ContentID

	tags := []string{"tag1", "tag2"}

	userObj := test_utils.TestUser()
	profile := test_utils.TestUserProfile()
	cqrObj := test_utils.TestContentQueueRecord()

	userObj.UserProfile = &profile

	Convey("PutImportedContentQueue", t, func() {
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

		Convey("when no tags have been passed in and GetMatchingTags errors out", func() {
			dripType := model.AutoDrip

			expectedError := fmt.Errorf("Some error")
			mockDB.On("GetMatchingTags", ctx, discussionID, contentID).Return(nil, expectedError)

			resp, err := backendObj.PutImportedContentQueue(ctx, discussionID, contentID, &now, nil, dripType)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when the drip type is Manual", func() {
			dripType := model.ManualDrip

			Convey("when UpdateImportedContentDiscussionQueue errors out", func() {
				expectedError := fmt.Errorf("Some error")
				timeZero := time.Time{}
				mockDB.On("UpdateImportedContentDiscussionQueue", ctx, discussionID, contentID, &timeZero).Return(nil, expectedError)

				resp, err := backendObj.PutImportedContentQueue(ctx, discussionID, contentID, &now, tags, dripType)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})

			Convey("when PutImportedContentDiscussionQueue errors out", func() {
				expectedError := fmt.Errorf("Some error")
				timeZero := time.Time{}
				mockDB.On("UpdateImportedContentDiscussionQueue", ctx, discussionID, contentID, &timeZero).Return(&cqrObj, nil)
				mockDB.On("PutImportedContentDiscussionQueue", ctx, discussionID, contentID, &now, tags).Return(nil, expectedError)

				resp, err := backendObj.PutImportedContentQueue(ctx, discussionID, contentID, &now, tags, dripType)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})

			Convey("when PutImportedContentQueue succeeds", func() {
				timeZero := time.Time{}
				mockDB.On("UpdateImportedContentDiscussionQueue", ctx, discussionID, contentID, &timeZero).Return(&cqrObj, nil)
				mockDB.On("PutImportedContentDiscussionQueue", ctx, discussionID, contentID, &now, tags).Return(&cqrObj, nil)

				resp, err := backendObj.PutImportedContentQueue(ctx, discussionID, contentID, &now, tags, dripType)

				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
			})
		})

		Convey("when the drip type is Scheduled", func() {
			dripType := model.ScheduledDrip

			Convey("when UpdateImportedContentDiscussionQueue errors out", func() {
				expectedError := fmt.Errorf("Some error")
				mockDB.On("UpdateImportedContentDiscussionQueue", ctx, discussionID, contentID, &now).Return(nil, expectedError)

				resp, err := backendObj.PutImportedContentQueue(ctx, discussionID, contentID, &now, tags, dripType)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})

			Convey("when PutImportedContentQueue succeeds", func() {
				mockDB.On("UpdateImportedContentDiscussionQueue", ctx, discussionID, contentID, &now).Return(&cqrObj, nil)

				resp, err := backendObj.PutImportedContentQueue(ctx, discussionID, contentID, &now, tags, dripType)

				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
			})
		})

		Convey("when the drip type is Auto", func() {
			dripType := model.AutoDrip

			Convey("when UpdateImportedContentDiscussionQueue errors out", func() {
				expectedError := fmt.Errorf("Some error")
				mockDB.On("PutImportedContentDiscussionQueue", ctx, discussionID, contentID, &now, tags).Return(nil, expectedError)

				resp, err := backendObj.PutImportedContentQueue(ctx, discussionID, contentID, &now, tags, dripType)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})

			Convey("when PutImportedContentQueue succeeds", func() {
				mockDB.On("PutImportedContentDiscussionQueue", ctx, discussionID, contentID, &now, tags).Return(&cqrObj, nil)

				resp, err := backendObj.PutImportedContentQueue(ctx, discussionID, contentID, &now, tags, dripType)

				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
			})
		})
	})
}

func TestDelphisBackend_GetPostsByDiscussionID(t *testing.T) {
	ctx := context.Background()
	userID := test_utils.UserID
	discussionID := test_utils.DiscussionID

	postObject := test_utils.TestPost()
	postContentObj := test_utils.TestPostContent()
	postObject.PostContent = &postContentObj

	Convey("GetPostsByDiscussionID", t, func() {
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

		Convey("when PostIterCollect errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetPostsByDiscussionIDIter", ctx, discussionID).Return(&mockPostIter{})
			mockDB.On("PostIterCollect", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetPostsByDiscussionID(ctx, userID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when GetPostsByDiscussionID succeeds", func() {
			mockDB.On("GetPostsByDiscussionIDIter", ctx, discussionID).Return(&mockPostIter{})
			mockDB.On("PostIterCollect", ctx, mock.Anything).Return([]*model.Post{&postObject}, nil)

			resp, err := backendObj.GetPostsByDiscussionID(ctx, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_GetPostByDiscussionPostID(t *testing.T) {
	ctx := context.Background()
	discussionID := test_utils.DiscussionID

	postObject := test_utils.TestPost()

	Convey("GetPostByDiscussionPostID", t, func() {
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

		Convey("when GetPostByID returns an error", func() {
			mockDB.On("GetPostByID", ctx, postObject.ID).Return(nil, fmt.Errorf("sth"))

			resp, err := backendObj.GetPostByDiscussionPostID(ctx, discussionID, postObject.ID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when response from GetPostID is nil", func() {
			mockDB.On("GetPostByID", ctx, postObject.ID).Return(nil, nil)

			resp, err := backendObj.GetPostByDiscussionPostID(ctx, discussionID, postObject.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when response from GetPostID has nil DiscussionID", func() {
			returnedPost := postObject
			returnedPost.DiscussionID = nil
			mockDB.On("GetPostByID", ctx, postObject.ID).Return(&returnedPost, nil)

			resp, err := backendObj.GetPostByDiscussionPostID(ctx, discussionID, postObject.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when response from GetPostID has different discussionID than requested", func() {
			returnedPost := postObject
			discID := "foo"
			returnedPost.DiscussionID = &discID
			mockDB.On("GetPostByID", ctx, postObject.ID).Return(&returnedPost, nil)

			resp, err := backendObj.GetPostByDiscussionPostID(ctx, discussionID, postObject.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when response from GetPostID is set and matches the discussionID", func() {
			returnedPost := postObject
			returnedPost.DiscussionID = &discussionID
			mockDB.On("GetPostByID", ctx, postObject.ID).Return(&returnedPost, nil)

			resp, err := backendObj.GetPostByDiscussionPostID(ctx, discussionID, postObject.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &returnedPost)
		})
	})

}

func TestDelphisBackend_GetPostsConnectionByDiscussionID(t *testing.T) {
	ctx := context.Background()
	discussionID := test_utils.DiscussionID
	limit := test_utils.Limit
	cursor := time.Now().Add(10 * time.Minute).Format(time.RFC3339Nano)

	postConnObj := test_utils.TestPostsConnection(cursor)
	postObject := test_utils.TestPost()
	postContentObj := test_utils.TestPostContent()
	profileObj := test_utils.TestUserProfile()
	postObject.PostContent = &postContentObj

	modObj := test_utils.TestModerator()
	modObj.UserProfile = &profileObj

	Convey("GetPostsConnectionByDiscussionID", t, func() {
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

		Convey("when limit is less than 2", func() {
			resp, err := backendObj.GetPostsConnectionByDiscussionID(ctx, discussionID, cursor, 0)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when GetPostsConnectionByDiscussionID errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetPostsConnectionByDiscussionID", ctx, discussionID, cursor, limit).Return(nil, expectedError)

			resp, err := backendObj.GetPostsConnectionByDiscussionID(ctx, discussionID, cursor, limit)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the cursor can not be parsed", func() {
			mockDB.On("GetPostsConnectionByDiscussionID", ctx, discussionID, mock.Anything, limit).Return(&postConnObj, nil)

			resp, err := backendObj.GetPostsConnectionByDiscussionID(ctx, discussionID, "bad cursor", limit)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_GetMentionedEntities(t *testing.T) {
	ctx := context.Background()
	discussionID := test_utils.DiscussionID

	participantID := test_utils.ParticipantID

	parID := strings.Join([]string{model.ParticipantPrefix, participantID}, ":")
	discID := strings.Join([]string{model.DiscussionPrefix, discussionID}, ":")
	unknownID := "chatham:abc"

	entityIDs := []string{parID, discID, unknownID}

	parObj := test_utils.TestParticipant()
	discObj := test_utils.TestDiscussion()

	parMap := map[string]*model.Participant{
		parObj.ID: &parObj,
	}

	discMap := map[string]*model.Discussion{
		discObj.ID: &discObj,
	}

	entityMap := map[string]model.Entity{
		parID:  &parObj,
		discID: &discObj,
	}

	Convey("GetMentionedEntities", t, func() {
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

		Convey("when a malformed entity is sent it", func() {
			resp, err := backendObj.GetMentionedEntities(ctx, []string{"a:b:c"})

			So(err, ShouldBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when GetParticipantsByIDs errors out", func() {
			expectedError := fmt.Errorf("Some error")
			mockDB.On("GetParticipantsByIDs", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetMentionedEntities(ctx, entityIDs)

			So(err, ShouldResemble, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when GetDiscussionsByIDs errors out", func() {
			expectedError := fmt.Errorf("Some error")
			mockDB.On("GetParticipantsByIDs", ctx, mock.Anything).Return(parMap, nil)
			mockDB.On("GetDiscussionsByIDs", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetMentionedEntities(ctx, entityIDs)

			So(err, ShouldResemble, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when GetMentionedEntities succeeds", func() {
			mockDB.On("GetParticipantsByIDs", ctx, mock.Anything).Return(parMap, nil)
			mockDB.On("GetDiscussionsByIDs", ctx, mock.Anything).Return(discMap, nil)

			resp, err := backendObj.GetMentionedEntities(ctx, entityIDs)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, entityMap)
		})
	})
}

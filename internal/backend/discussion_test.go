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
	"github.com/dghubble/go-twitter/twitter"

	"github.com/stretchr/testify/mock"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/cache"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/delphis-inc/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDelphisBackend_CreateNewDiscussion(t *testing.T) {
	ctx := context.Background()

	discussionID := test_utils.DiscussionID
	anonymityType := model.AnonymityTypeStrong
	title := "test title"
	description := "test description"
	publicAccess := true

	userObj := test_utils.TestUser()
	modObj := test_utils.TestModerator()
	profile := test_utils.TestUserProfile()
	discObj := test_utils.TestDiscussion()
	discussionSettings := test_utils.TestDiscussionCreationSettings()
	discussionUserAccess := test_utils.TestDiscussionUserAccess()

	userObj.UserProfile = &profile
	modObj.UserProfile = &profile

	viewerObj := test_utils.TestViewer()

	parObj := test_utils.TestParticipant()

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

		Convey("when create moderator errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("CreateModerator", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.CreateNewDiscussion(ctx, &userObj, anonymityType, title, description, publicAccess, discussionSettings)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when create discussion errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("CreateModerator", ctx, mock.Anything).Return(&modObj, nil)
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.CreateNewDiscussion(ctx, &userObj, anonymityType, title, description, publicAccess, discussionSettings)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when upsert discussion access errors out", func() {
			expectedError := fmt.Errorf("Some Error")

			mockDB.On("CreateModerator", ctx, mock.Anything).Return(&modObj, nil)
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(&discObj, nil)

			// Upsert discussion access functions
			mockDB.On("BeginTx", ctx).Return(nil, expectedError)

			resp, err := backendObj.CreateNewDiscussion(ctx, &userObj, anonymityType, title, description, publicAccess, discussionSettings)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when create participant errors out", func() {
			expectedError := fmt.Errorf("Some Error")

			mockDB.On("CreateModerator", ctx, mock.Anything).Return(&modObj, nil)
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(&discObj, nil)

			// Upsert discussion access
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("GetDiscussionUserAccess", ctx, mock.Anything, mock.Anything).Return(&discussionUserAccess, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, &tx, mock.Anything).Return(&discussionUserAccess, nil)
			mockDB.On("CommitTx", ctx, &tx).Return(nil)

			// Create participant functions
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.CreateNewDiscussion(ctx, &userObj, anonymityType, title, description, publicAccess, discussionSettings)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when upsert links errors out", func() {
			expectedError := fmt.Errorf("Some Error")

			mockDB.On("CreateModerator", ctx, mock.Anything).Return(&modObj, nil)
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(&discObj, nil)

			// Upsert discussion access
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("GetDiscussionUserAccess", ctx, mock.Anything, mock.Anything).Return(&discussionUserAccess, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, &tx, mock.Anything).Return(&discussionUserAccess, nil)
			mockDB.On("CommitTx", ctx, &tx).Return(nil)

			// Create participant functions
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(&userObj, nil)
			mockDB.On("GetTotalParticipantCountByDiscussionID", ctx, mock.Anything).Return(10)
			mockDB.On("GetViewerForDiscussion", ctx, mock.Anything, mock.Anything).Return(nil, nil)
			mockDB.On("UpsertViewer", ctx, mock.Anything).Return(&viewerObj, nil)
			mockDB.On("UpsertParticipant", ctx, mock.Anything).Return(&parObj, nil)
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)

			// Upsert discussion access
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&discussionUserAccess, nil)
			mockDB.On("CommitTx", ctx, &tx).Return(nil)

			// Create participant functions
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(&userObj, nil)
			mockDB.On("GetTotalParticipantCountByDiscussionID", ctx, mock.Anything).Return(10)
			mockDB.On("GetViewerForDiscussion", ctx, mock.Anything, mock.Anything).Return(nil, nil)
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
			mockDB.On("GetDUAForEverythingNotifications", ctx, discussionID, mock.Anything).Return(nil)
			mockDB.On("DuaIterCollect", ctx, mock.Anything).Return(nil, nil)
			mockDB.On("GetUserDevicesByUserID", ctx, mock.Anything).Return(nil, nil)

			mockDB.On("BeginTx", ctx).Return(tx, nil)
			mockDB.On("PutAccessLinkForDiscussion", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.CreateNewDiscussion(ctx, &userObj, anonymityType, title, description, publicAccess, discussionSettings)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when the discussion is created successfully", func() {

			mockDB.On("CreateModerator", ctx, mock.Anything).Return(&modObj, nil)
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(&discObj, nil)

			// Upsert discussion access
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("GetDiscussionUserAccess", ctx, mock.Anything, mock.Anything).Return(&discussionUserAccess, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, &tx, mock.Anything).Return(&discussionUserAccess, nil)
			mockDB.On("CommitTx", ctx, &tx).Return(nil)

			// Create participant functions
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(&userObj, nil)
			mockDB.On("GetTotalParticipantCountByDiscussionID", ctx, mock.Anything).Return(10)
			mockDB.On("GetViewerForDiscussion", ctx, mock.Anything, mock.Anything).Return(nil, nil)
			mockDB.On("UpsertViewer", ctx, mock.Anything).Return(&viewerObj, nil)
			mockDB.On("UpsertParticipant", ctx, mock.Anything).Return(&parObj, nil)
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)

			// Upsert discussion access
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&discussionUserAccess, nil)
			mockDB.On("CommitTx", ctx, &tx).Return(nil)

			// Create participant functions
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(&userObj, nil)
			mockDB.On("GetTotalParticipantCountByDiscussionID", ctx, mock.Anything).Return(10)
			mockDB.On("GetViewerForDiscussion", ctx, mock.Anything, mock.Anything).Return(nil, nil)
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
			mockDB.On("GetDUAForEverythingNotifications", ctx, discussionID, mock.Anything).Return(nil)
			mockDB.On("DuaIterCollect", ctx, mock.Anything).Return(nil, nil)
			mockDB.On("GetUserDevicesByUserID", ctx, mock.Anything).Return(nil, nil)

			mockDB.On("BeginTx", ctx).Return(tx, nil)
			mockDB.On("PutAccessLinkForDiscussion", ctx, mock.Anything, mock.Anything).Return(
				&model.DiscussionAccessLink{DiscussionID: discussionID}, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.CreateNewDiscussion(ctx, &userObj, anonymityType, title, description, publicAccess, discussionSettings)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_IncrementDiscussionShuffleCount(t *testing.T) {
	ctx := context.Background()

	discussionObj := test_utils.TestDiscussion()

	Convey("IncrementDiscussionShuffleCount", t, func() {
		var tx *sql.Tx
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
		Convey("when tx is not passed", func() {
			Convey("when beginning a tx fails it returns an error", func() {
				expectedError := fmt.Errorf("Some Error")
				mockDB.On("BeginTx", ctx).Return(nil, expectedError)

				resp, err := backendObj.IncrementDiscussionShuffleCount(ctx, nil, discussionObj.ID)

				So(err, ShouldEqual, expectedError)
				So(resp, ShouldBeNil)
			})
			Convey("when increment returns an error", func() {
				tx = &sql.Tx{}
				Convey("and rolling back tx fails", func() {
					mockDB.On("BeginTx", ctx).Return(tx, nil)
					mockDB.On("IncrementDiscussionShuffleCount", ctx, tx, discussionObj.ID).Return(nil, fmt.Errorf("sth"))
					mockDB.On("RollbackTx", ctx, tx).Return(fmt.Errorf("sth"))

					resp, err := backendObj.IncrementDiscussionShuffleCount(ctx, nil, discussionObj.ID)

					So(err, ShouldNotBeNil)
					So(resp, ShouldBeNil)
				})

				Convey("and rolling back tx succeeds", func() {
					mockDB.On("BeginTx", ctx).Return(tx, nil)
					mockDB.On("IncrementDiscussionShuffleCount", ctx, tx, discussionObj.ID).Return(nil, fmt.Errorf("sth"))
					mockDB.On("RollbackTx", ctx, tx).Return(nil)

					resp, err := backendObj.IncrementDiscussionShuffleCount(ctx, nil, discussionObj.ID)

					So(err, ShouldNotBeNil)
					So(resp, ShouldBeNil)
				})
			})
			Convey("when increment succeeds", func() {
				tx = &sql.Tx{}
				newShuffleCount := 1
				Convey("and commit fails and rollback fails", func() {
					mockDB.On("BeginTx", ctx).Return(tx, nil)
					mockDB.On("IncrementDiscussionShuffleCount", ctx, tx, discussionObj.ID).Return(&newShuffleCount, nil)
					mockDB.On("CommitTx", ctx, tx).Return(fmt.Errorf("sth"))
					mockDB.On("RollbackTx", ctx, tx).Return(fmt.Errorf("sth"))

					resp, err := backendObj.IncrementDiscussionShuffleCount(ctx, nil, discussionObj.ID)

					So(err, ShouldNotBeNil)
					So(resp, ShouldBeNil)
				})
				Convey("and commit fails and rollback succeeds", func() {
					mockDB.On("BeginTx", ctx).Return(tx, nil)
					mockDB.On("IncrementDiscussionShuffleCount", ctx, tx, discussionObj.ID).Return(&newShuffleCount, nil)
					mockDB.On("CommitTx", ctx, tx).Return(fmt.Errorf("sth"))
					mockDB.On("RollbackTx", ctx, tx).Return(nil)

					resp, err := backendObj.IncrementDiscussionShuffleCount(ctx, nil, discussionObj.ID)

					So(err, ShouldNotBeNil)
					So(resp, ShouldBeNil)
				})
				Convey("and commit succeeds", func() {
					mockDB.On("BeginTx", ctx).Return(tx, nil)
					mockDB.On("IncrementDiscussionShuffleCount", ctx, tx, discussionObj.ID).Return(&newShuffleCount, nil)
					mockDB.On("CommitTx", ctx, tx).Return(nil)

					resp, err := backendObj.IncrementDiscussionShuffleCount(ctx, nil, discussionObj.ID)

					So(err, ShouldBeNil)
					So(resp, ShouldEqual, &newShuffleCount)
				})
			})
		})
		Convey("when tx is passed", func() {
			tx := &sql.Tx{}
			Convey("and increment discussion returns an error", func() {
				mockDB.On("IncrementDiscussionShuffleCount", ctx, tx, discussionObj.ID).Return(nil, fmt.Errorf("sth"))

				resp, err := backendObj.IncrementDiscussionShuffleCount(ctx, tx, discussionObj.ID)

				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})
			Convey("and increment discussion succeeds", func() {
				newShuffleCount := 1
				mockDB.On("IncrementDiscussionShuffleCount", ctx, tx, discussionObj.ID).Return(&newShuffleCount, nil)

				resp, err := backendObj.IncrementDiscussionShuffleCount(ctx, tx, discussionObj.ID)

				So(err, ShouldBeNil)
				So(resp, ShouldEqual, &newShuffleCount)
			})
		})
	})
}

func TestDelphisBackend_GetDiscussionJoinabilityForUser(t *testing.T) {
	ctx := context.Background()

	discussionObj := test_utils.TestDiscussion()
	moderatorObj := test_utils.TestModerator()
	moderatorUserProfileID := "someOtherID"
	moderatorObj.UserProfileID = &moderatorUserProfileID
	discussionObj.Moderator = &moderatorObj
	userObj := test_utils.TestUser()
	userProfileObj := test_utils.TestUserProfile()
	userObj.UserProfile = &userProfileObj
	twitterSocialInfo := test_utils.TestSocialInfo()
	moderatorTwitterSocialInfo := twitterSocialInfo
	moderatorTwitterSocialInfo.ScreenName = "moderator"
	nonTwitterSocialInfo := test_utils.TestSocialInfo()
	nonTwitterSocialInfo.Network = "foo"
	meParticipant := test_utils.TestParticipant()

	Convey("GetDiscussionJoinabilityForUser", t, func() {
		now := time.Now()
		cacheObj := cache.NewInMemoryCache()
		authObj := auth.NewDelphisAuth(nil)
		mockDB := &mocks.Datastore{}
		mockTwitterBackend := &mocks.TwitterBackend{}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            authObj,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
			twitterBackend:  mockTwitterBackend,
		}
		mockTwitterClient := &mocks.TwitterClient{}

		Convey("when objects are nil", func() {
			Convey("when userObj is nil", func() {
				resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, nil, &discussionObj, nil)

				So(resp, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
			Convey("when userObj.UserProfile is nil", func() {
				testUserObj := userObj
				testUserObj.UserProfile = nil

				resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &testUserObj, &discussionObj, nil)

				So(resp, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
			Convey("when discussion is nil", func() {
				testUserObj := userObj
				testUserObj.UserProfile = nil

				resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, nil, nil)

				So(resp, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
		mockDB.On("GetDiscussionUserAccess", ctx, discussionObj.ID, userObj.ID).Return(nil, nil)
		Convey("when getting social infos returns an error", func() {
			mockDB.On("GetSocialInfosByUserProfileID", ctx, userObj.UserProfile.ID).Return(nil, fmt.Errorf("sth"))

			resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

			So(resp, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
		Convey("when social info returns an empty array", func() {
			mockDB.On("GetSocialInfosByUserProfileID", ctx, userObj.UserProfile.ID).Return([]model.SocialInfo{}, nil)

			resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

			So(resp, ShouldNotBeNil)
			So(resp.Response, ShouldEqual, model.DiscussionJoinabilityResponseDenied)
			So(err, ShouldBeNil)
		})
		Convey("when social info returns objects but no twitter auth", func() {
			mockDB.On("GetSocialInfosByUserProfileID", ctx, userObj.UserProfile.ID).Return([]model.SocialInfo{nonTwitterSocialInfo}, nil)

			resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

			So(resp, ShouldNotBeNil)
			So(resp.Response, ShouldEqual, model.DiscussionJoinabilityResponseDenied)
			So(err, ShouldBeNil)
		})

		mockDB.On("GetSocialInfosByUserProfileID", ctx, userObj.UserProfile.ID).Return([]model.SocialInfo{twitterSocialInfo}, nil)

		Convey("when meParticipant is non-null", func() {
			resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, &meParticipant)

			So(resp, ShouldNotBeNil)
			So(resp.Response, ShouldEqual, model.DiscussionJoinabilityResponseAlreadyJoined)
			So(err, ShouldBeNil)
		})

		Convey("when discussion joinability set to Twitter Friends", func() {
			discussionObj.DiscussionJoinability = model.DiscussionJoinabilitySettingAllowTwitterFriends
			Convey("when fetching moderator social info errors", func() {
				mockDB.On("GetSocialInfosByUserProfileID", ctx, *discussionObj.Moderator.UserProfileID).Return(nil, fmt.Errorf("sth"))

				resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

				So(resp, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
			Convey("when no social info is found for moderator", func() {
				mockDB.On("GetSocialInfosByUserProfileID", ctx, *discussionObj.Moderator.UserProfileID).Return([]model.SocialInfo{}, nil)

				resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

				So(resp, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
			Convey("when social info is found", func() {
				mockDB.On("GetSocialInfosByUserProfileID", ctx, *discussionObj.Moderator.UserProfileID).Return([]model.SocialInfo{moderatorTwitterSocialInfo}, nil)
				Convey("when getting twitter client fails", func() {
					mockTwitterBackend.On("GetTwitterClientWithAccessTokens", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("sth"))

					resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

					So(resp, ShouldBeNil)
					So(err, ShouldNotBeNil)
				})
				Convey("when checking whether user follows fails", func() {
					mockTwitterBackend.On("GetTwitterClientWithAccessTokens", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockTwitterClient, nil)
					mockTwitterClient.On("FriendshipLookup", twitterSocialInfo.ScreenName, moderatorTwitterSocialInfo.ScreenName).Return(nil, fmt.Errorf("sth"))

					resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

					So(resp, ShouldBeNil)
					So(err, ShouldNotBeNil)
				})
				Convey("when moderator is following user", func() {
					mockTwitterBackend.On("GetTwitterClientWithAccessTokens", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockTwitterClient, nil)
					mockTwitterClient.On("FriendshipLookup", twitterSocialInfo.ScreenName, moderatorTwitterSocialInfo.ScreenName).Return(&twitter.Relationship{Target: twitter.RelationshipTarget{Following: true}}, nil)

					resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

					So(resp, ShouldNotBeNil)
					So(resp.Response, ShouldEqual, model.DiscussionJoinabilityResponseApprovedNotJoined)
					So(err, ShouldBeNil)
				})
				Convey("when moderator is not following user", func() {
					mockTwitterBackend.On("GetTwitterClientWithAccessTokens", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockTwitterClient, nil)
					mockTwitterClient.On("FriendshipLookup", twitterSocialInfo.ScreenName, moderatorTwitterSocialInfo.ScreenName).Return(&twitter.Relationship{Target: twitter.RelationshipTarget{Following: false}}, nil)

					Convey("when getting discussion access errors", func() {
						mockDB.On("GetDiscussionAccessRequestByDiscussionIDUserID", ctx, discussionObj.ID, userObj.ID).Return(nil, fmt.Errorf("sth"))

						resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

						So(resp, ShouldBeNil)
						So(err, ShouldNotBeNil)
					})
					Convey("when nil is returned", func() {
						mockDB.On("GetDiscussionAccessRequestByDiscussionIDUserID", ctx, discussionObj.ID, userObj.ID).Return(nil, nil)

						resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

						So(resp, ShouldNotBeNil)
						So(resp.Response, ShouldEqual, model.DiscussionJoinabilityResponseApprovalRequired)
						So(err, ShouldBeNil)
					})
					Convey("when a status is returned", func() {
						Convey("when the status is accepted", func() {
							mockDB.On("GetDiscussionAccessRequestByDiscussionIDUserID", ctx, discussionObj.ID, userObj.ID).Return(&model.DiscussionAccessRequest{Status: model.InviteRequestStatusAccepted}, nil)

							resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

							So(resp, ShouldNotBeNil)
							So(resp.Response, ShouldEqual, model.DiscussionJoinabilityResponseApprovedNotJoined)
							So(err, ShouldBeNil)
						})
						Convey("when the status is pending", func() {
							mockDB.On("GetDiscussionAccessRequestByDiscussionIDUserID", ctx, discussionObj.ID, userObj.ID).Return(&model.DiscussionAccessRequest{Status: model.InviteRequestStatusPending}, nil)

							resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

							So(resp, ShouldNotBeNil)
							So(resp.Response, ShouldEqual, model.DiscussionJoinabilityResponseApprovalRequired)
							So(err, ShouldBeNil)
						})
						Convey("when the status is rejected", func() {
							mockDB.On("GetDiscussionAccessRequestByDiscussionIDUserID", ctx, discussionObj.ID, userObj.ID).Return(&model.DiscussionAccessRequest{Status: model.InviteRequestStatusRejected}, nil)

							resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

							So(resp, ShouldNotBeNil)
							So(resp.Response, ShouldEqual, model.DiscussionJoinabilityResponseDenied)
							So(err, ShouldBeNil)
						})
						Convey("when the status is cancelled", func() {
							mockDB.On("GetDiscussionAccessRequestByDiscussionIDUserID", ctx, discussionObj.ID, userObj.ID).Return(&model.DiscussionAccessRequest{Status: model.InviteRequestStatusCancelled}, nil)

							resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

							So(resp, ShouldNotBeNil)
							So(resp.Response, ShouldEqual, model.DiscussionJoinabilityResponseApprovalRequired)
							So(err, ShouldBeNil)
						})
						Convey("when the status is unknown", func() {
							mockDB.On("GetDiscussionAccessRequestByDiscussionIDUserID", ctx, discussionObj.ID, userObj.ID).Return(&model.DiscussionAccessRequest{Status: model.InviteRequestStatus("foo")}, nil)

							resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

							So(resp, ShouldNotBeNil)
							So(resp.Response, ShouldEqual, model.DiscussionJoinabilityResponseApprovalRequired)
							So(err, ShouldBeNil)
						})

					})
				})
			})
		})
		Convey("when discussion joinability set to require approval", func() {
			discussionObj.DiscussionJoinability = model.DiscussionJoinabilitySettingAllRequireApproval
			mockTwitterBackend.On("GetTwitterClientWithAccessTokens", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockTwitterClient, nil)
			mockTwitterClient.On("FriendshipLookup", twitterSocialInfo.ScreenName, moderatorTwitterSocialInfo.ScreenName).Return(&twitter.Relationship{Target: twitter.RelationshipTarget{Following: false}}, nil)

			Convey("when getting discussion access errors", func() {
				mockDB.On("GetDiscussionAccessRequestByDiscussionIDUserID", ctx, discussionObj.ID, userObj.ID).Return(nil, fmt.Errorf("sth"))

				resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

				So(resp, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
			Convey("when nil is returned", func() {
				mockDB.On("GetDiscussionAccessRequestByDiscussionIDUserID", ctx, discussionObj.ID, userObj.ID).Return(nil, nil)

				resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

				So(resp, ShouldNotBeNil)
				So(resp.Response, ShouldEqual, model.DiscussionJoinabilityResponseApprovalRequired)
				So(err, ShouldBeNil)
			})
			Convey("when a status is returned", func() {
				Convey("when the status is accepted", func() {
					mockDB.On("GetDiscussionAccessRequestByDiscussionIDUserID", ctx, discussionObj.ID, userObj.ID).Return(&model.DiscussionAccessRequest{Status: model.InviteRequestStatusAccepted}, nil)

					resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

					So(resp, ShouldNotBeNil)
					So(resp.Response, ShouldEqual, model.DiscussionJoinabilityResponseApprovedNotJoined)
					So(err, ShouldBeNil)
				})
				Convey("when the status is pending", func() {
					mockDB.On("GetDiscussionAccessRequestByDiscussionIDUserID", ctx, discussionObj.ID, userObj.ID).Return(&model.DiscussionAccessRequest{Status: model.InviteRequestStatusPending}, nil)

					resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

					So(resp, ShouldNotBeNil)
					So(resp.Response, ShouldEqual, model.DiscussionJoinabilityResponseApprovalRequired)
					So(err, ShouldBeNil)
				})
				Convey("when the status is rejected", func() {
					mockDB.On("GetDiscussionAccessRequestByDiscussionIDUserID", ctx, discussionObj.ID, userObj.ID).Return(&model.DiscussionAccessRequest{Status: model.InviteRequestStatusRejected}, nil)

					resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

					So(resp, ShouldNotBeNil)
					So(resp.Response, ShouldEqual, model.DiscussionJoinabilityResponseDenied)
					So(err, ShouldBeNil)
				})
				Convey("when the status is cancelled", func() {
					mockDB.On("GetDiscussionAccessRequestByDiscussionIDUserID", ctx, discussionObj.ID, userObj.ID).Return(&model.DiscussionAccessRequest{Status: model.InviteRequestStatusCancelled}, nil)

					resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

					So(resp, ShouldNotBeNil)
					So(resp.Response, ShouldEqual, model.DiscussionJoinabilityResponseApprovalRequired)
					So(err, ShouldBeNil)
				})
				Convey("when the status is unknown", func() {
					mockDB.On("GetDiscussionAccessRequestByDiscussionIDUserID", ctx, discussionObj.ID, userObj.ID).Return(&model.DiscussionAccessRequest{Status: model.InviteRequestStatus("foo")}, nil)

					resp, err := backendObj.GetDiscussionJoinabilityForUser(ctx, &userObj, &discussionObj, nil)

					So(resp, ShouldNotBeNil)
					So(resp.Response, ShouldEqual, model.DiscussionJoinabilityResponseApprovalRequired)
					So(err, ShouldBeNil)
				})
			})
		})
	})
}

func TestDelphisBackend_GetDiscussionByLinkSlug(t *testing.T) {
	ctx := context.Background()

	slug := "slug"

	discObj := test_utils.TestDiscussion()

	Convey("GetDiscussionByLinkSlug", t, func() {
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

		Convey("when the query errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetDiscussionByLinkSlug", ctx, slug).Return(nil, expectedError)

			resp, err := backendObj.GetDiscussionByLinkSlug(ctx, slug)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetDiscussionByLinkSlug", ctx, slug).Return(&discObj, nil)

			resp, err := backendObj.GetDiscussionByLinkSlug(ctx, slug)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &discObj)
		})
	})
}

func TestDelphisBackend_UpdateDiscussion(t *testing.T) {
	ctx := context.Background()

	discussionID := test_utils.DiscussionID
	newShuffleCount := 1

	discInput := test_utils.TestDiscussionInput()
	discObj := test_utils.TestDiscussion()
	postObj := test_utils.TestPost()
	contentObj := test_utils.TestPostContent()
	archiveObj := test_utils.TestDiscussionArchive()

	postObj.PostContent = &contentObj

	tx := sql.Tx{}

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

		Convey("when get discussion by id errors out", func() {
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

		Convey("when the title is changed, the history is updated", func() {
			mockDB.On("GetDiscussionByID", ctx, discussionID).Return(&discObj, nil)
			newTitle := "newTitle"
			updatedDiscussion := discObj
			updatedDiscussion.Title = newTitle
			updatedDiscussion.AddTitleToHistory(updatedDiscussion.Title)

			matcher := func(arg interface{}) bool {
				argAsDiscussion := arg.(model.Discussion)
				expectedTitleHistory, err := updatedDiscussion.TitleHistoryAsObject()
				actualTitleHistory, err2 := argAsDiscussion.TitleHistoryAsObject()
				return err == nil && err2 == nil && len(expectedTitleHistory) == len(actualTitleHistory) && expectedTitleHistory[0].Value == actualTitleHistory[0].Value
			}

			mockDB.On("UpsertDiscussion", ctx, mock.MatchedBy(matcher)).Return(&discObj, nil)
			updateInput := model.DiscussionInput{
				Title: &newTitle,
			}
			resp, err := backendObj.UpdateDiscussion(ctx, discussionID, updateInput)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})

		Convey("when the description is changed, the history is updated", func() {
			mockDB.On("GetDiscussionByID", ctx, discussionID).Return(&discObj, nil)
			newDescription := "newDescription"
			updatedDiscussion := discObj
			updatedDiscussion.Description = newDescription
			updatedDiscussion.AddDescriptionToHistory(updatedDiscussion.Description)

			matcher := func(arg interface{}) bool {
				argAsDiscussion := arg.(model.Discussion)
				expectedDescriptionHistory, err := updatedDiscussion.DescriptionHistoryAsObject()
				actualDescriptionHistory, err2 := argAsDiscussion.DescriptionHistoryAsObject()
				return err == nil && err2 == nil && len(expectedDescriptionHistory) == len(actualDescriptionHistory) && expectedDescriptionHistory[0].Value == actualDescriptionHistory[0].Value
			}

			mockDB.On("UpsertDiscussion", ctx, mock.MatchedBy(matcher)).Return(&discObj, nil)
			updateInput := model.DiscussionInput{
				Description: &newDescription,
			}
			resp, err := backendObj.UpdateDiscussion(ctx, discussionID, updateInput)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})

		Convey("when the discussion is locked, an entry in the archive table is upserted", func() {
			mockDB.On("GetDiscussionByID", ctx, discussionID).Return(&discObj, nil)
			mockDB.On("GetPostsByDiscussionIDIter", ctx, discussionID).Return(&mockPostIter{})
			mockDB.On("PostIterCollect", ctx, mock.Anything).Return([]*model.Post{&postObj}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpsertDiscussionArchive", ctx, mock.Anything, mock.Anything).Return(&archiveObj, nil)
			mockDB.On("IncrementDiscussionShuffleCount", ctx, mock.Anything, discObj.ID).Return(&newShuffleCount, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)
			updatedDiscussion := discObj
			updatedDiscussion.LockStatus = true
			trueVal := true

			mockDB.On("UpsertDiscussion", ctx, updatedDiscussion).Return(&updatedDiscussion, nil)
			updateInput := model.DiscussionInput{
				LockStatus: &trueVal,
			}
			resp, err := backendObj.UpdateDiscussion(ctx, discussionID, updateInput)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_GetDiscussionArchiveByDiscussionID(t *testing.T) {
	ctx := context.Background()

	discussionID := test_utils.DiscussionID

	archiveObj := test_utils.TestDiscussionArchive()

	Convey("GetDiscussionArchiveByDiscussionID", t, func() {
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

		Convey("when the query errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetDiscussionArchiveByDiscussionID", ctx, discussionID).Return(nil, expectedError)

			resp, err := backendObj.GetDiscussionArchiveByDiscussionID(ctx, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetDiscussionArchiveByDiscussionID", ctx, discussionID).Return(&archiveObj, nil)

			resp, err := backendObj.GetDiscussionArchiveByDiscussionID(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &archiveObj)
		})
	})
}

func TestDelphisBackend_CreateDiscussionArchive(t *testing.T) {
	ctx := context.Background()

	discussionID := test_utils.DiscussionID
	shuffleCount := 1
	newShuffleCount := shuffleCount + 1
	postObj := test_utils.TestPost()
	contentObj := test_utils.TestPostContent()
	discObj := test_utils.TestDiscussion()
	archiveObj := test_utils.TestDiscussionArchive()

	postObj.PostContent = &contentObj

	tx := sql.Tx{}

	Convey("CreateDiscussionArchive", t, func() {
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

		Convey("when get posts by discussion id errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetPostsByDiscussionIDIter", ctx, discussionID).Return(&mockPostIter{})
			mockDB.On("PostIterCollect", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.CreateDiscussionArchive(ctx, discussionID, shuffleCount)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when begin tx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetPostsByDiscussionIDIter", ctx, discussionID).Return(&mockPostIter{})
			mockDB.On("PostIterCollect", ctx, mock.Anything).Return([]*model.Post{&postObj}, nil)
			mockDB.On("BeginTx", ctx).Return(nil, expectedError)

			resp, err := backendObj.CreateDiscussionArchive(ctx, discussionID, shuffleCount)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when upsert discussion archive errors out and rollback fails", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetPostsByDiscussionIDIter", ctx, discussionID).Return(&mockPostIter{})
			mockDB.On("PostIterCollect", ctx, mock.Anything).Return([]*model.Post{&postObj}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpsertDiscussionArchive", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.CreateDiscussionArchive(ctx, discussionID, shuffleCount)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when upsert discussion archive errors out and rollback succeeds", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetPostsByDiscussionIDIter", ctx, discussionID).Return(&mockPostIter{})
			mockDB.On("PostIterCollect", ctx, mock.Anything).Return([]*model.Post{&postObj}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpsertDiscussionArchive", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.CreateDiscussionArchive(ctx, discussionID, shuffleCount)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when increment discussion shuffle count errors out and rollback fails", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetPostsByDiscussionIDIter", ctx, discussionID).Return(&mockPostIter{})
			mockDB.On("PostIterCollect", ctx, mock.Anything).Return([]*model.Post{&postObj}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpsertDiscussionArchive", ctx, mock.Anything, mock.Anything).Return(&archiveObj, nil)
			mockDB.On("IncrementDiscussionShuffleCount", ctx, mock.Anything, discObj.ID).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.CreateDiscussionArchive(ctx, discussionID, shuffleCount)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when increment discussion shuffle count errors out and rollback succeeds", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetPostsByDiscussionIDIter", ctx, discussionID).Return(&mockPostIter{})
			mockDB.On("PostIterCollect", ctx, mock.Anything).Return([]*model.Post{&postObj}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpsertDiscussionArchive", ctx, mock.Anything, mock.Anything).Return(&archiveObj, nil)
			mockDB.On("IncrementDiscussionShuffleCount", ctx, mock.Anything, discObj.ID).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.CreateDiscussionArchive(ctx, discussionID, shuffleCount)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when CreateDiscussion succeeds but commitTx fails", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetPostsByDiscussionIDIter", ctx, discussionID).Return(&mockPostIter{})
			mockDB.On("PostIterCollect", ctx, mock.Anything).Return([]*model.Post{&postObj}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpsertDiscussionArchive", ctx, mock.Anything, mock.Anything).Return(&archiveObj, nil)
			mockDB.On("IncrementDiscussionShuffleCount", ctx, mock.Anything, discObj.ID).Return(&newShuffleCount, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.CreateDiscussionArchive(ctx, discussionID, shuffleCount)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when CreateDiscussion succeeds", func() {
			mockDB.On("GetPostsByDiscussionIDIter", ctx, discussionID).Return(&mockPostIter{})
			mockDB.On("PostIterCollect", ctx, mock.Anything).Return([]*model.Post{&postObj}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpsertDiscussionArchive", ctx, mock.Anything, mock.Anything).Return(&archiveObj, nil)
			mockDB.On("IncrementDiscussionShuffleCount", ctx, mock.Anything, discObj.ID).Return(&newShuffleCount, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.CreateDiscussionArchive(ctx, discussionID, shuffleCount)

			So(err, ShouldBeNil)
			So(resp, ShouldEqual, &archiveObj)
		})

	})
}

func TestDelphisBackend_GetDiscussionByID(t *testing.T) {
	ctx := context.Background()

	discussionID := test_utils.DiscussionID

	discObj := test_utils.TestDiscussion()

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

		Convey("when the query errors out", func() {
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
	discussionID := test_utils.DiscussionID

	discObj := test_utils.TestDiscussion()

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

		Convey("when the query errors out", func() {
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

	modID := test_utils.ModeratorID

	discObj := test_utils.TestDiscussion()

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

		Convey("when the query errors out", func() {
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

func TestDelphisBackend_ListDiscussions(t *testing.T) {
	ctx := context.Background()

	dcObj := test_utils.TestDiscussionsConnection()

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

		Convey("when the query errors out", func() {
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

func TestDelphisBackend_ListDiscussionsByUserID(t *testing.T) {
	ctx := context.Background()

	dcObj := test_utils.TestDiscussionsConnection()
	state := model.DiscussionUserAccessStateActive

	Convey("ListDiscussionsByUserID", t, func() {
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
		userID := "userID"

		Convey("when the query errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("ListDiscussionsByUserID", ctx, userID, state).Return(nil, expectedError)

			resp, err := backendObj.ListDiscussionsByUserID(ctx, userID, state)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("ListDiscussionsByUserID", ctx, userID, state).Return(&dcObj, nil)

			resp, err := backendObj.ListDiscussionsByUserID(ctx, userID, state)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &dcObj)
		})
	})
}

func TestDelphisBackend_anonymizePostsForArchive(t *testing.T) {
	ctx := context.Background()

	shuffleCount := 1

	postObj := test_utils.TestPost()
	contentObj := test_utils.TestPostContent()
	postObj.PostContent = &contentObj

	participantHash := util.GenerateParticipantSeed(*postObj.DiscussionID, *postObj.ParticipantID, shuffleCount)
	participantName := util.GenerateFullDisplayName(participantHash)

	expectedResult := model.ArchivedPost{
		PostType:          postObj.PostType,
		ParticipantName:   participantName,
		Content:           postObj.PostContent.Content,
		MentionedEntities: []string{},
	}

	Convey("anonymizePostsForArchive", t, func() {
		now := time.Now()

		Convey("when there are no posts to archive", func() {
			resp, err := anonymizePostsForArchive(ctx, nil, shuffleCount)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.ArchivedPost{})
		})

		Convey("when there are only deleted posts to archive", func() {
			tempPost := postObj
			tempPost.DeletedAt = &now
			tempPosts := []*model.Post{&tempPost}
			resp, err := anonymizePostsForArchive(ctx, tempPosts, shuffleCount)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.ArchivedPost{})
		})

		Convey("when we successfully anonymize the post", func() {
			tempPosts := []*model.Post{&postObj}

			resp, err := anonymizePostsForArchive(ctx, tempPosts, shuffleCount)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []*model.ArchivedPost{&expectedResult})
		})

		Convey("when there are mentioned entities", func() {
			mentionedParticipantID := "12345"
			mentionedParticipant := strings.Join([]string{model.ParticipantPrefix, mentionedParticipantID}, ":")

			mentionedParticipantHash := util.GenerateParticipantSeed(*postObj.DiscussionID, mentionedParticipantID, shuffleCount)
			mentionedParticipantName := util.GenerateFullDisplayName(mentionedParticipantHash)

			tempContent := contentObj
			tempContent.MentionedEntities = []string{mentionedParticipant, "discussion:1234"}

			tempPost := postObj
			tempPost.PostContent = &tempContent
			tempPosts := []*model.Post{&tempPost}

			// Expected result
			testResult := expectedResult
			testResult.MentionedEntities = []string{mentionedParticipantName, "redacted_discussion"}

			resp, err := anonymizePostsForArchive(ctx, tempPosts, shuffleCount)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
			So(resp, ShouldResemble, []*model.ArchivedPost{&testResult})
		})

	})
}

func TestDelphisBackend_UpdateDiscussionObj(t *testing.T) {
	discInput := test_utils.TestDiscussionInput()

	disc := model.Discussion{}

	Convey("UpdateDiscussionObj", t, func() {
		Convey("when it updates the discussion object successfully", func() {
			updateDiscussionObj(&disc, discInput)

			So(disc.AnonymityType, ShouldResemble, *discInput.AnonymityType)
			So(disc.Title, ShouldResemble, *discInput.Title)
			So(*disc.IconURL, ShouldResemble, *discInput.IconURL)

		})
	})
}

func TestDelphisBackend_DedupeDiscussions(t *testing.T) {
	disc1 := test_utils.TestDiscussion()
	disc2 := test_utils.TestDiscussion()

	disc2.ID = "id2"

	Convey("DedupeDiscussions", t, func() {
		Convey("when it dedupes the discussion objects successfully", func() {
			resp := dedupeDiscussions([]*model.Discussion{&disc1, &disc2, &disc1})

			So(resp, ShouldResemble, []*model.Discussion{&disc1, &disc2})

		})
	})
}

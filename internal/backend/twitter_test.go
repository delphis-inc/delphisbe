package backend

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/backend/test_utils"
	"github.com/delphis-inc/delphisbe/internal/cache"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/delphis-inc/delphisbe/mocks"
	"github.com/dghubble/go-twitter/twitter"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

type mockEmptyDiscussionInviteIter struct{}

func (m *mockEmptyDiscussionInviteIter) Next(invite *model.DiscussionInvite) bool { return false }

func (m *mockEmptyDiscussionInviteIter) Close() error { return nil }

func TestDelphisBackend_GetTwitterAccessToken(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	Convey("GetTwitterAccessToken", t, func() {
		testAuthedUser := test_utils.TestDelphisAuthedUser()
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

		Convey("when user is not authed", func() {
			token, secret, err := backendObj.GetTwitterAccessToken(ctx)

			So(err, ShouldNotEqual, nil)
			So(token, ShouldEqual, "")
			So(secret, ShouldEqual, "")
		})
		ctx = auth.WithAuthedUser(ctx, &testAuthedUser)

		Convey("when user profile query errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetUserProfileByUserID", ctx, testAuthedUser.UserID).Return(nil, expectedError)

			token, secret, err := backendObj.GetTwitterAccessToken(ctx)

			So(err, ShouldEqual, expectedError)
			So(token, ShouldEqual, "")
			So(secret, ShouldEqual, "")
		})

		Convey("when social info query errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			userProfile := test_utils.TestUserProfile()
			mockDB.On("GetUserProfileByUserID", ctx, testAuthedUser.UserID).Return(&userProfile, nil)
			mockDB.On("GetSocialInfosByUserProfileID", ctx, userProfile.ID).Return(nil, expectedError)

			token, secret, err := backendObj.GetTwitterAccessToken(ctx)

			So(err, ShouldEqual, expectedError)
			So(token, ShouldEqual, "")
			So(secret, ShouldEqual, "")
		})

		Convey("when everything is ok", func() {
			userProfile := test_utils.TestUserProfile()
			socialInfo := []model.SocialInfo{test_utils.TestSocialInfo()}
			mockDB.On("GetUserProfileByUserID", ctx, testAuthedUser.UserID).Return(&userProfile, nil)
			mockDB.On("GetSocialInfosByUserProfileID", ctx, userProfile.ID).Return(socialInfo, nil)

			token, secret, err := backendObj.GetTwitterAccessToken(ctx)

			So(err, ShouldEqual, nil)
			So(token, ShouldEqual, socialInfo[0].AccessToken)
			So(secret, ShouldEqual, socialInfo[0].AccessTokenSecret)
		})

	})
}

func TestDelphisBackend_GetTwitterClientWithUserTokens(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	Convey("GetTwitterClientWithUserTokens", t, func() {
		testAuthedUser := test_utils.TestDelphisAuthedUser()
		cacheObj := cache.NewInMemoryCache()
		authObj := auth.NewDelphisAuth(nil)
		mockDB := &mocks.Datastore{}
		twitterConfig := config.TwitterConfig{
			ConsumerKey:    "ConsumerKey",
			ConsumerSecret: "ConsumerSecret",
		}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            authObj,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when user is not authed", func() {
			token, secret, err := backendObj.GetTwitterAccessToken(ctx)

			So(err, ShouldNotEqual, nil)
			So(token, ShouldEqual, "")
			So(secret, ShouldEqual, "")
		})
		ctx = auth.WithAuthedUser(ctx, &testAuthedUser)

		Convey("when keys are not setted", func() {
			userProfile := test_utils.TestUserProfile()
			socialInfo := []model.SocialInfo{test_utils.TestSocialInfo()}
			mockDB.On("GetUserProfileByUserID", ctx, testAuthedUser.UserID).Return(&userProfile, nil)
			mockDB.On("GetSocialInfosByUserProfileID", ctx, userProfile.ID).Return(socialInfo, nil)

			_, err := backendObj.GetTwitterClientWithUserTokens(ctx)

			So(err, ShouldNotEqual, nil)
		})

		backendObj.config.Twitter = twitterConfig
		Convey("when tokens and keys are setted", func() {
			userProfile := test_utils.TestUserProfile()
			socialInfo := []model.SocialInfo{test_utils.TestSocialInfo()}
			mockDB.On("GetUserProfileByUserID", ctx, testAuthedUser.UserID).Return(&userProfile, nil)
			mockDB.On("GetSocialInfosByUserProfileID", ctx, userProfile.ID).Return(socialInfo, nil)

			_, err := backendObj.GetTwitterClientWithUserTokens(ctx)

			So(err, ShouldEqual, nil)
		})

	})
}

func TestDelphisBackend_GetTwitterUserHandleAutocompletes(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	participantID := "participant1"
	discussionID := "discussion1"

	Convey("GetTwitterUserHandleAutocompletes", t, func() {
		cacheObj := cache.NewInMemoryCache()
		authObj := auth.NewDelphisAuth(nil)
		mockDB := &mocks.Datastore{}
		mockTwitter := mocks.TwitterClient{}
		mockQuery := "usernametest"
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            authObj,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when users search errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockTwitter.On("SearchUsers", mockQuery, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, expectedError)
			mockDB.On("GetInvitedTwitterHandlesByDiscussionIDAndInviterID", ctx, discussionID, participantID).Return([]*string{}, nil)
			results, err := backendObj.GetTwitterUserHandleAutocompletes(ctx, &mockTwitter, mockQuery, discussionID, participantID)

			So(err, ShouldEqual, expectedError)
			So(results, ShouldEqual, nil)
		})

		Convey("when existing invitations query errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockTwitter.On("SearchUsers", mockQuery, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return([]twitter.User{}, nil)
			mockDB.On("GetInvitedTwitterHandlesByDiscussionIDAndInviterID", ctx, discussionID, participantID).Return(nil, expectedError)
			results, err := backendObj.GetTwitterUserHandleAutocompletes(ctx, &mockTwitter, mockQuery, discussionID, participantID)

			So(err, ShouldEqual, expectedError)
			So(results, ShouldEqual, nil)
		})

		Convey("when users search gives few results", func() {
			var returnedResult []twitter.User
			var expectedResult []*model.TwitterUserInfo
			for i := 0; i < 10; i++ {
				returnedResult = append(returnedResult, twitter.User{
					ScreenName:           fmt.Sprintf("username%d", i),
					Name:                 fmt.Sprintf("User Name %d", i),
					Verified:             true,
					IDStr:                fmt.Sprintf("%08d", i),
					ProfileImageURLHttps: "https://example.com/image.png",
				})
				expectedResult = append(expectedResult, &model.TwitterUserInfo{
					Name:            fmt.Sprintf("username%d", i),
					DiplayName:      fmt.Sprintf("User Name %d", i),
					Verified:        true,
					ID:              fmt.Sprintf("%08d", i),
					ProfileImageURL: "https://example.com/image.png",
					Invited:         false,
				})
			}
			mockTwitter.On("SearchUsers", mockQuery, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(returnedResult, nil)
			mockDB.On("GetInvitedTwitterHandlesByDiscussionIDAndInviterID", ctx, discussionID, participantID).Return([]*string{}, nil)
			results, err := backendObj.GetTwitterUserHandleAutocompletes(ctx, &mockTwitter, mockQuery, discussionID, participantID)

			So(err, ShouldEqual, nil)
			So(len(results), ShouldEqual, len(expectedResult))
			for i := range results {
				So(reflect.DeepEqual(results[i], expectedResult[i]), ShouldEqual, true)
			}
		})

		Convey("when users search gives many results", func() {
			var returnedResult []twitter.User
			var expectedResult []*model.TwitterUserInfo
			for i := 0; i < 80; i++ {
				returnedResult = append(returnedResult, twitter.User{
					ScreenName:           fmt.Sprintf("username%d", i),
					Name:                 fmt.Sprintf("User Name %d", i),
					Verified:             true,
					IDStr:                fmt.Sprintf("%08d", i),
					ProfileImageURLHttps: "https://example.com/image.png",
				})
				expectedResult = append(expectedResult, &model.TwitterUserInfo{
					Name:            fmt.Sprintf("username%d", i),
					DiplayName:      fmt.Sprintf("User Name %d", i),
					Verified:        true,
					ID:              fmt.Sprintf("%08d", i),
					ProfileImageURL: "https://example.com/image.png",
					Invited:         false,
				})
			}
			mockTwitter.On("SearchUsers", mockQuery, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(returnedResult, nil)
			mockDB.On("GetInvitedTwitterHandlesByDiscussionIDAndInviterID", ctx, discussionID, participantID).Return([]*string{}, nil)
			results, err := backendObj.GetTwitterUserHandleAutocompletes(ctx, &mockTwitter, mockQuery, discussionID, participantID)

			So(err, ShouldEqual, nil)
			So(len(results), ShouldEqual, len(expectedResult))
			for i := range results {
				So(reflect.DeepEqual(results[i], expectedResult[i]), ShouldEqual, true)
			}
		})

		Convey("when users has already invited someone", func() {
			var returnedResult []twitter.User
			var expectedResult []*model.TwitterUserInfo
			var alreadyInvitedHandles []*string
			for i := 0; i < 10; i++ {
				twitterHandle := fmt.Sprintf("username%d", i)
				twitterUserName := fmt.Sprintf("User Name %d", i)
				if i < 3 {
					alreadyInvitedHandles = append(alreadyInvitedHandles, &twitterHandle)
				}
				returnedResult = append(returnedResult, twitter.User{
					ScreenName:           twitterHandle,
					Name:                 twitterUserName,
					Verified:             true,
					IDStr:                fmt.Sprintf("%08d", i),
					ProfileImageURLHttps: "https://example.com/image.png",
				})
				expectedResult = append(expectedResult, &model.TwitterUserInfo{
					Name:            twitterHandle,
					DiplayName:      twitterUserName,
					Verified:        true,
					ID:              fmt.Sprintf("%08d", i),
					ProfileImageURL: "https://example.com/image.png",
					Invited:         i < 3,
				})
			}
			mockTwitter.On("SearchUsers", mockQuery, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(returnedResult, nil)
			mockDB.On("GetInvitedTwitterHandlesByDiscussionIDAndInviterID", ctx, discussionID, participantID).Return(alreadyInvitedHandles, nil)
			results, err := backendObj.GetTwitterUserHandleAutocompletes(ctx, &mockTwitter, mockQuery, discussionID, participantID)

			So(err, ShouldEqual, nil)
			So(len(results), ShouldEqual, len(expectedResult))
			for i := range results {
				So(reflect.DeepEqual(results[i], expectedResult[i]), ShouldEqual, true)
			}
		})
	})

}

func TestDelphisBackend_InviteTwitterUsersToDiscussion(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	discussionID := test_utils.DiscussionID
	participantID := test_utils.InvitingParticipantID
	profileObj := test_utils.TestUserProfile()
	socialObj := test_utils.TestSocialInfo()
	userObj := test_utils.TestUser()
	inviteObj := test_utils.TestDiscussionInvite(model.InviteRequestStatusAccepted)
	tx := sql.Tx{}

	Convey("InviteTwitterUsersToDiscussion", t, func() {
		testAuthedUser := test_utils.TestDelphisAuthedUser()
		cacheObj := cache.NewInMemoryCache()
		authObj := auth.NewDelphisAuth(nil)
		mockDB := &mocks.Datastore{}
		mockTwitter := mocks.TwitterClient{}
		mockQuery := &model.TwitterUserInput{
			Name: "username",
		}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            authObj,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when user is not authed", func() {
			token, secret, err := backendObj.GetTwitterAccessToken(ctx)

			So(err, ShouldNotEqual, nil)
			So(token, ShouldEqual, "")
			So(secret, ShouldEqual, "")
		})
		ctx = auth.WithAuthedUser(ctx, &testAuthedUser)

		Convey("when users autoinvites", func() {
			var returnedResult []twitter.User
			for i := 0; i < 20; i++ {
				screenName := fmt.Sprintf("result%d", i)
				name := fmt.Sprintf("Result #%d", i)
				returnedResult = append(returnedResult, twitter.User{
					ScreenName: screenName,
					Name:       name,
					ID:         int64(i),
				})
			}

			mockHandles := []*model.TwitterUserInput{mockQuery}
			mockDB.On("CreateOrUpdateUserProfile", ctx, mock.Anything).Return(&profileObj, true, nil)
			mockDB.On("UpsertSocialInfo", ctx, mock.Anything).Return(&socialObj, nil)
			mockDB.On("UpsertUser", ctx, mock.Anything).Return(testAuthedUser.User, nil)
			mockTwitter.On("LookupUsers", mock.Anything).Return(returnedResult, nil)
			mockDB.On("GetInvitedTwitterHandlesByDiscussionIDAndInviterID", ctx, discussionID, participantID).Return([]*string{}, nil)

			results, err := backendObj.InviteTwitterUsersToDiscussion(ctx, &mockTwitter, mockHandles, discussionID, participantID)

			So(err, ShouldEqual, nil)
			So(len(results), ShouldEqual, 0)
		})

		Convey("when users lookup succeeds ", func() {
			var returnedResult []twitter.User
			var mockHandles []*model.TwitterUserInput
			for i := 0; i < 20; i++ {
				screenName := fmt.Sprintf("result%d", i)
				name := fmt.Sprintf("Result #%d", i)
				returnedResult = append(returnedResult, twitter.User{
					ScreenName: screenName,
					Name:       name,
					ID:         int64(i),
				})
				mockHandles = append(mockHandles, &model.TwitterUserInput{Name: screenName})
			}
			mockDB.On("CreateOrUpdateUserProfile", ctx, mock.Anything).Return(&profileObj, true, nil)
			mockDB.On("UpsertSocialInfo", ctx, mock.Anything).Return(&socialObj, nil)
			mockDB.On("UpsertUser", ctx, mock.Anything).Return(&userObj, nil)
			mockDB.On("GetDiscussionInvitesByUserIDAndStatus", ctx, mock.Anything, mock.Anything).Return(&mockEmptyDiscussionInviteIter{}, nil)
			mockDB.On("DiscussionInviteIterCollect", ctx, mock.Anything).Return([]*model.DiscussionInvite{}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(&inviteObj, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)
			mockTwitter.On("LookupUsers", mock.Anything).Return(returnedResult, nil)
			mockDB.On("GetInvitedTwitterHandlesByDiscussionIDAndInviterID", ctx, discussionID, participantID).Return([]*string{}, nil)

			results, err := backendObj.InviteTwitterUsersToDiscussion(ctx, &mockTwitter, mockHandles, discussionID, participantID)

			So(err, ShouldEqual, nil)
			So(len(results), ShouldEqual, len(returnedResult))
		})

		Convey("when some some users are already invited ", func() {
			var alreadyInvited []*string
			var returnedResult []twitter.User
			var mockHandles []*model.TwitterUserInput
			for i := 0; i < 20; i++ {
				screenName := fmt.Sprintf("result%d", i)
				name := fmt.Sprintf("Result #%d", i)

				if i < 10 {
					alreadyInvited = append(alreadyInvited, &screenName)
				}
				returnedResult = append(returnedResult, twitter.User{
					ScreenName: screenName,
					Name:       name,
					ID:         int64(i),
				})
				mockHandles = append(mockHandles, &model.TwitterUserInput{Name: screenName})
			}
			mockDB.On("CreateOrUpdateUserProfile", ctx, mock.Anything).Return(&profileObj, true, nil)
			mockDB.On("UpsertSocialInfo", ctx, mock.Anything).Return(&socialObj, nil)
			mockDB.On("UpsertUser", ctx, mock.Anything).Return(&userObj, nil)
			mockDB.On("GetDiscussionInvitesByUserIDAndStatus", ctx, mock.Anything, mock.Anything).Return(&mockEmptyDiscussionInviteIter{}, nil)
			mockDB.On("DiscussionInviteIterCollect", ctx, mock.Anything).Return([]*model.DiscussionInvite{}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(&inviteObj, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)
			mockTwitter.On("LookupUsers", mock.Anything).Return(returnedResult, nil)
			mockDB.On("GetInvitedTwitterHandlesByDiscussionIDAndInviterID", ctx, discussionID, participantID).Return(alreadyInvited, nil)

			results, err := backendObj.InviteTwitterUsersToDiscussion(ctx, &mockTwitter, mockHandles, discussionID, participantID)

			So(err, ShouldEqual, nil)
			So(len(results), ShouldEqual, len(returnedResult)-10)
		})

	})
}

package backend

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/auth"
	"github.com/nedrocks/delphisbe/internal/cache"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/nedrocks/delphisbe/internal/util"
	"github.com/nedrocks/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestDelphisBackend_GetConciergeParticipantID(t *testing.T) {
	ctx := context.Background()
	discussionID := "discussionID"

	parObj := model.Participant{
		ID: "participantID",
	}

	Convey("GetConciergeParticipantID", t, func() {
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

		Convey("when GetParticipantsByDiscussionIDUserID errors out", func() {
			expectedError := fmt.Errorf("Some Error")

			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetConciergeParticipantID(ctx, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldResemble, "")
		})

		Convey("when GetParticipantsByDiscussionIDUserID returns only an anonymouse participant", func() {
			tempParObj := parObj
			tempParObj.IsAnonymous = true
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{tempParObj}, nil)

			resp, err := backendObj.GetConciergeParticipantID(ctx, discussionID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldResemble, "")
		})

		Convey("when GetConciergeParticipantID succeeds", func() {
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)

			resp, err := backendObj.GetConciergeParticipantID(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, parObj.ID)
		})
	})
}

func TestDelphisBackend_HandleConciergeMutation(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	userID := "userID"
	discussionID := "discussion1"
	mutationID := string(model.MutationUpdateFlairAccessToDiscussion)
	selectedOptions := []string{"1"}

	parObj := model.Participant{
		ID: "participantID",
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

	displayName := "name"
	ftObj := model.FlairTemplate{
		ID:          templateID,
		DisplayName: &displayName,
	}

	Convey("HandleConciergeMutation", t, func() {
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

		Convey("when GetConciergeParticipantID errors out", func() {
			expectedError := fmt.Errorf("Some Error")

			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.HandleConciergeMutation(ctx, userID, discussionID, mutationID, selectedOptions)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when createFlairAccessConciergePost is successful", func() {
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)
			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return([]*model.Flair{&flairObj}, nil)
			mockDB.On("GetFlairTemplateByID", ctx, mock.Anything).Return(&ftObj, nil)

			resp, err := backendObj.HandleConciergeMutation(ctx, userID, discussionID, mutationID, selectedOptions)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})

		Convey("when createRenameChatAndEmojiConciergePost is successful", func() {
			tempMutationID := string(model.MutationUpdateDiscussionNameAndEmoji)
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)

			resp, err := backendObj.HandleConciergeMutation(ctx, userID, discussionID, tempMutationID, selectedOptions)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})

		Convey("when createViewerAccessConciergePost is successful", func() {
			tempMutationID := string(model.MutationUpdateViewerAccessibility)
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)
			mockDB.On("GetDiscussionByID", ctx, mock.Anything).Return(&model.Discussion{ID: "1234"}, nil)

			resp, err := backendObj.HandleConciergeMutation(ctx, userID, discussionID, tempMutationID, selectedOptions)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})

		Convey("when createInviteSettingConciergePost is successful", func() {
			tempMutationID := string(model.MutationUpdateInvitationApproval)
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)
			mockDB.On("GetDiscussionByID", ctx, mock.Anything).Return(&model.Discussion{ID: "1234"}, nil)

			resp, err := backendObj.HandleConciergeMutation(ctx, userID, discussionID, tempMutationID, selectedOptions)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})

		Convey("when an improper mutationID is passed in", func() {
			tempMutationID := "test"
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)
			mockDB.On("GetDiscussionByID", ctx, mock.Anything).Return(&model.Discussion{ID: "1234"}, nil)

			resp, err := backendObj.HandleConciergeMutation(ctx, userID, discussionID, tempMutationID, selectedOptions)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})
	})
}

func TestDelphisBackend_CreateFlairAccessConciergePost(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	userID := "userID"
	discussionID := "discussion1"

	parObj := model.Participant{
		ID: "participantID",
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

	displayName := "name"
	ftObj := model.FlairTemplate{
		ID:          templateID,
		DisplayName: &displayName,
	}

	Convey("CreateFlairAccessConciergePost", t, func() {
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

		Convey("when GetFlairsByUserID errors out", func() {
			expectedError := fmt.Errorf("Some Error")

			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.createFlairAccessConciergePost(ctx, userID, discussionID, parObj.ID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when GetFlairTemplateByID errors out", func() {
			expectedError := fmt.Errorf("Some Error")

			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return([]*model.Flair{&flairObj}, nil)
			mockDB.On("GetFlairTemplateByID", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.createFlairAccessConciergePost(ctx, userID, discussionID, parObj.ID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when GetFlairsByUserID is successful", func() {
			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return([]*model.Flair{&flairObj}, nil)
			mockDB.On("GetFlairTemplateByID", ctx, mock.Anything).Return(&ftObj, nil)

			resp, err := backendObj.createFlairAccessConciergePost(ctx, userID, discussionID, parObj.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_CreateInviteSettingConciergePost(t *testing.T) {
	ctx := context.Background()
	discussionID := "discussion1"

	parObj := model.Participant{
		ID: "participantID",
	}

	Convey("CreateInviteSettingConciergePost", t, func() {
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

		Convey("when GetDiscussionByID errors out", func() {
			expectedError := fmt.Errorf("Some Error")

			mockDB.On("GetDiscussionByID", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.createInviteSettingConciergePost(ctx, discussionID, parObj.ID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when query is successful", func() {
			mockDB.On("GetDiscussionByID", ctx, mock.Anything).Return(&model.Discussion{ID: "1234"}, nil)

			resp, err := backendObj.createInviteSettingConciergePost(ctx, discussionID, parObj.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_CreateViewerAccessConciergePost(t *testing.T) {
	ctx := context.Background()
	discussionID := "discussion1"

	parObj := model.Participant{
		ID: "participantID",
	}

	Convey("CreateViewerAccessConciergePost", t, func() {
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

		Convey("when GetDiscussionByID errors out", func() {
			expectedError := fmt.Errorf("Some Error")

			mockDB.On("GetDiscussionByID", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.createViewerAccessConciergePost(ctx, discussionID, parObj.ID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when query is successful", func() {
			mockDB.On("GetDiscussionByID", ctx, mock.Anything).Return(&model.Discussion{ID: "1234"}, nil)

			resp, err := backendObj.createViewerAccessConciergePost(ctx, discussionID, parObj.ID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

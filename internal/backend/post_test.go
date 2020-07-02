package backend

import (
	"context"
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

type mockPostIter struct{}

func (m *mockPostIter) Next(post *model.Post) bool { return true }
func (m *mockPostIter) Close() error               { return fmt.Errorf("error") }

func TestDelphisBackend_GetPostsByDiscussionID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	userID := "userID"
	discussionID := "discussion1"
	participantID := "participant1"
	postID := "post1"
	modID := "modID"
	profileID := "profileID"

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

	modObj := model.Moderator{
		ID: modID,
		UserProfile: &model.UserProfile{
			ID: profileID,
		},
	}

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

			resp, err := backendObj.GetPostsByDiscussionID(ctx, userID, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when CheckIfModeratorForDiscussion errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetPostsByDiscussionIDIter", ctx, discussionID).Return(&mockPostIter{})
			mockDB.On("PostIterCollect", ctx, mock.Anything).Return([]*model.Post{&postObject}, nil)
			mockDB.On("GetModeratorByUserIDAndDiscussionID", ctx, userID, discussionID).Return(nil, expectedError)

			resp, err := backendObj.GetPostsByDiscussionID(ctx, userID, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when GetNewDiscussionConciergePosts errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetPostsByDiscussionIDIter", ctx, discussionID).Return(&mockPostIter{})
			mockDB.On("PostIterCollect", ctx, mock.Anything).Return([]*model.Post{&postObject}, nil)
			mockDB.On("GetModeratorByUserIDAndDiscussionID", ctx, userID, discussionID).Return(&modObj, nil)
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetPostsByDiscussionID(ctx, userID, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when GetPostsByDiscussionID succeeds", func() {
			mockDB.On("GetPostsByDiscussionIDIter", ctx, discussionID).Return(&mockPostIter{})
			mockDB.On("PostIterCollect", ctx, mock.Anything).Return([]*model.Post{&postObject}, nil)
			mockDB.On("GetModeratorByUserIDAndDiscussionID", ctx, userID, discussionID).Return(&modObj, nil)
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)
			mockDB.On("GetFlairsByUserID", ctx, mock.Anything).Return([]*model.Flair{&flairObj}, nil)
			mockDB.On("GetFlairTemplateByID", ctx, mock.Anything).Return(&ftObj, nil)

			resp, err := backendObj.GetPostsByDiscussionID(ctx, userID, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

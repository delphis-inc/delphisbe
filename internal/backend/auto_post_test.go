package backend

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/delphis-inc/delphisbe/internal/backend/test_utils"

	"github.com/delphis-inc/delphisbe/graph/model"

	"github.com/delphis-inc/delphisbe/internal/cache"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/delphis-inc/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestDelphisBackend_AutoPostContent(t *testing.T) {
	ctx := context.Background()

	discussionID := test_utils.DiscussionID
	idleMinutes := test_utils.IdleMinutes
	limit := test_utils.Limit
	userID := model.ConciergeUser
	contentID := test_utils.ContentID

	apObj := test_utils.TestDiscussionAutoPost()
	icObj := test_utils.TestImportedContent()
	parObj := test_utils.TestParticipant()
	discObj := test_utils.TestDiscussion()

	tx := sql.Tx{}
	matchingTags := []string{"tag1"}

	Convey("AutoPostContent", t, func() {
		now := time.Now()
		cacheObj := cache.NewInMemoryCache()
		mockAuth := &mocks.DelphisAuth{}
		mockDB := &mocks.Datastore{}
		backendObj := &delphisBackend{
			db:              mockDB,
			auth:            mockAuth,
			cache:           cacheObj,
			discussionMutex: sync.Mutex{},
			config:          config.Config{},
			timeProvider:    &util.FrozenTime{NowTime: now},
		}

		Convey("when GetDiscussionsForAutoPost errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetDiscussionsAutoPost", ctx).Return(&mockDiscAutoPostIter{})
			mockDB.On("DiscussionAutoPostIterCollect", ctx, mock.Anything).Return(nil, expectedError)

			backendObj.AutoPostContent()
		})

		Convey("when checkIdleTime errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetDiscussionsAutoPost", ctx).Return(&mockDiscAutoPostIter{})
			mockDB.On("DiscussionAutoPostIterCollect", ctx, mock.Anything).Return([]*model.DiscussionAutoPost{&apObj}, nil)
			mockDB.On("GetLastPostByDiscussionID", ctx, discussionID, idleMinutes).Return(nil, expectedError)

			backendObj.AutoPostContent()
		})

		Convey("when postNextContent errors out", func() {
			Convey("when GetScheduledImportedContentByDiscussionID errors out", func() {
				expectedError := fmt.Errorf("Some Error")
				mockDB.On("GetDiscussionsAutoPost", ctx).Return(&mockDiscAutoPostIter{})
				mockDB.On("DiscussionAutoPostIterCollect", ctx, mock.Anything).Return([]*model.DiscussionAutoPost{&apObj}, nil)
				mockDB.On("GetLastPostByDiscussionID", ctx, discussionID, idleMinutes).Return(nil, nil)
				mockDB.On("GetScheduledImportedContentByDiscussionID", ctx, discussionID).Return(&mockImportedContentIter{})
				mockDB.On("ContentIterCollect", ctx, mock.Anything).Return(nil, expectedError)

				backendObj.AutoPostContent()
			})

			Convey("when GetImportedContentByDiscussionID errors out", func() {
				scheduleIter := mockImportedContentIter{}

				expectedError := fmt.Errorf("Some Error")
				mockDB.On("GetDiscussionsAutoPost", ctx).Return(&mockDiscAutoPostIter{})
				mockDB.On("DiscussionAutoPostIterCollect", ctx, mock.Anything).Return([]*model.DiscussionAutoPost{&apObj}, nil)
				mockDB.On("GetLastPostByDiscussionID", ctx, discussionID, idleMinutes).Return(nil, nil)
				mockDB.On("GetScheduledImportedContentByDiscussionID", ctx, discussionID).Return(&scheduleIter)
				mockDB.On("ContentIterCollect", ctx, scheduleIter).Return(nil, nil)
				mockDB.On("GetImportedContentByDiscussionID", ctx, discussionID, limit).Return(&mockImportedContentIter{})
				mockDB.On("ContentIterCollect", ctx, mock.Anything).Return(nil, expectedError)

				backendObj.AutoPostContent()
			})

			Convey("when GetImportedContentByDiscussionID returns 0 content", func() {
				mockDB.On("GetDiscussionsAutoPost", ctx).Return(&mockDiscAutoPostIter{})
				mockDB.On("DiscussionAutoPostIterCollect", ctx, mock.Anything).Return([]*model.DiscussionAutoPost{&apObj}, nil)
				mockDB.On("GetLastPostByDiscussionID", ctx, discussionID, idleMinutes).Return(nil, nil)
				mockDB.On("GetScheduledImportedContentByDiscussionID", ctx, discussionID).Return(&mockImportedContentIter{})
				mockDB.On("ContentIterCollect", ctx, mock.Anything).Return(nil, nil)
				mockDB.On("GetImportedContentByDiscussionID", ctx, discussionID, limit).Return(&mockImportedContentIter{})
				mockDB.On("ContentIterCollect", ctx, mock.Anything).Return(nil, nil)

				backendObj.AutoPostContent()
			})

			Convey("when GetParticipantsByDiscussionIDUserID errors out", func() {
				expectedError := fmt.Errorf("Some Error")
				mockDB.On("GetDiscussionsAutoPost", ctx).Return(&mockDiscAutoPostIter{})
				mockDB.On("DiscussionAutoPostIterCollect", ctx, mock.Anything).Return([]*model.DiscussionAutoPost{&apObj}, nil)
				mockDB.On("GetLastPostByDiscussionID", ctx, discussionID, idleMinutes).Return(nil, nil)
				mockDB.On("GetScheduledImportedContentByDiscussionID", ctx, discussionID).Return(&mockImportedContentIter{})
				mockDB.On("ContentIterCollect", ctx, mock.Anything).Return([]*model.ImportedContent{&icObj}, nil)
				mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, discussionID, userID).Return(nil, expectedError)

				backendObj.AutoPostContent()
			})

			Convey("when GetParticipantsByDiscussionIDUserID returns only an anonymous concierge", func() {
				tempParObj := parObj
				tempParObj.IsAnonymous = true
				mockDB.On("GetDiscussionsAutoPost", ctx).Return(&mockDiscAutoPostIter{})
				mockDB.On("DiscussionAutoPostIterCollect", ctx, mock.Anything).Return([]*model.DiscussionAutoPost{&apObj}, nil)
				mockDB.On("GetLastPostByDiscussionID", ctx, discussionID, idleMinutes).Return(nil, nil)
				mockDB.On("GetScheduledImportedContentByDiscussionID", ctx, discussionID).Return(&mockImportedContentIter{})
				mockDB.On("ContentIterCollect", ctx, mock.Anything).Return([]*model.ImportedContent{&icObj}, nil)
				mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, discussionID, userID).Return([]model.Participant{tempParObj}, nil)

				backendObj.AutoPostContent()
			})

			Convey("when PostImportedContent errors out", func() {
				expectedError := fmt.Errorf("Some Error")
				mockDB.On("GetDiscussionsAutoPost", ctx).Return(&mockDiscAutoPostIter{})
				mockDB.On("DiscussionAutoPostIterCollect", ctx, mock.Anything).Return([]*model.DiscussionAutoPost{&apObj}, nil)
				mockDB.On("GetLastPostByDiscussionID", ctx, discussionID, idleMinutes).Return(nil, nil)
				mockDB.On("GetScheduledImportedContentByDiscussionID", ctx, discussionID).Return(&mockImportedContentIter{})
				mockDB.On("ContentIterCollect", ctx, mock.Anything).Return([]*model.ImportedContent{&icObj}, nil)
				mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, discussionID, userID).Return([]model.Participant{parObj}, nil)
				mockDB.On("GetImportedContentByID", ctx, contentID).Return(nil, expectedError)

				backendObj.AutoPostContent()
			})

			Convey("when auto posting succeeds", func() {
				mockDB.On("GetDiscussionsAutoPost", ctx).Return(&mockDiscAutoPostIter{})
				mockDB.On("DiscussionAutoPostIterCollect", ctx, mock.Anything).Return([]*model.DiscussionAutoPost{&apObj}, nil)
				mockDB.On("GetLastPostByDiscussionID", ctx, discussionID, idleMinutes).Return(nil, nil)
				mockDB.On("GetScheduledImportedContentByDiscussionID", ctx, discussionID).Return(&mockImportedContentIter{})
				mockDB.On("ContentIterCollect", ctx, mock.Anything).Return([]*model.ImportedContent{&icObj}, nil)
				mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, discussionID, userID).Return([]model.Participant{parObj}, nil)

				mockDB.On("GetImportedContentByID", ctx, contentID).Return(&icObj, nil)
				// Create post functions
				mockDB.On("BeginTx", ctx).Return(&tx, nil)
				mockDB.On("PutPostContent", ctx, mock.Anything, mock.Anything).Return(nil)
				mockDB.On("PutPost", ctx, mock.Anything, mock.Anything).Return(&model.Post{ID: "post123"}, nil)
				mockDB.On("PutActivity", ctx, mock.Anything, mock.Anything).Return(nil)
				mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)
				mockDB.On("GetDiscussionByID", ctx, mock.Anything).Return(&discObj, nil)
				mockDB.On("GetParticipantsByDiscussionID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)

				// Put Imported Content Queue
				mockDB.On("GetMatchingTags", ctx, discussionID, contentID).Return(matchingTags, nil)
				mockDB.On("UpdateImportedContentDiscussionQueue", ctx, discussionID, contentID, mock.Anything).Return(
					&model.ContentQueueRecord{DiscussionID: discussionID}, nil)
				mockDB.On("PutImportedContentDiscussionQueue", ctx, discussionID, contentID, time.Now(), matchingTags).Return(&icObj, nil)
				mockDB.On("GetUserDevicesByUserID", ctx, mock.Anything).Return(nil, nil)

				backendObj.AutoPostContent()
			})
		})
	})
}

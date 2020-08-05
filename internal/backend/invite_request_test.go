package backend

import (
	"context"
	"database/sql"
	"fmt"
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

type mockDiscussionInviteIter struct{}

func (m *mockDiscussionInviteIter) Next(invite *model.DiscussionInvite) bool { return true }
func (m *mockDiscussionInviteIter) Close() error                             { return fmt.Errorf("error") }

type mockDiscussionAccessRequestIter struct{}

func (m *mockDiscussionAccessRequestIter) Next(equest *model.DiscussionAccessRequest) bool {
	return true
}
func (m *mockDiscussionAccessRequestIter) Close() error { return fmt.Errorf("error") }

func TestDelphisBackend_GetDiscussionInviteByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	inviteID := test_utils.InviteID

	inviteObj := test_utils.TestDiscussionInvite(model.InviteRequestStatusAccepted)

	Convey("GetDiscussionInviteByID", t, func() {
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
			mockDB.On("GetDiscussionInviteByID", ctx, inviteID).Return(nil, expectedError)

			resp, err := backendObj.GetDiscussionInviteByID(ctx, inviteID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetDiscussionInviteByID", ctx, inviteID).Return(&inviteObj, nil)

			resp, err := backendObj.GetDiscussionInviteByID(ctx, inviteID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &inviteObj)
		})
	})
}

func TestDelphisBackend_GetDiscussionRequestAccessByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	requestID := test_utils.RequestID

	requestObj := test_utils.TestDiscussionAccessRequest(model.InviteRequestStatusAccepted)

	Convey("GetDiscussionRequestAccessByID", t, func() {
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
			mockDB.On("GetDiscussionRequestAccessByID", ctx, requestID).Return(nil, expectedError)

			resp, err := backendObj.GetDiscussionRequestAccessByID(ctx, requestID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetDiscussionRequestAccessByID", ctx, requestID).Return(&requestObj, nil)

			resp, err := backendObj.GetDiscussionRequestAccessByID(ctx, requestID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &requestObj)
		})
	})
}

func TestDelphisBackend_GetDiscussionInvitesByUserIDAndStatus(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := test_utils.UserID
	status := model.InviteRequestStatusPending

	inviteObj := test_utils.TestDiscussionInvite(model.InviteRequestStatusPending)

	Convey("GetDiscussionInvitesByUserIDAndStatus", t, func() {
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
			mockDB.On("GetDiscussionInvitesByUserIDAndStatus", ctx, userID, status).Return(&mockDiscussionInviteIter{})
			mockDB.On("DiscussionInviteIterCollect", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetDiscussionInvitesByUserIDAndStatus(ctx, userID, status)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetDiscussionInvitesByUserIDAndStatus", ctx, userID, status).Return(&mockDiscussionInviteIter{})
			mockDB.On("DiscussionInviteIterCollect", ctx, mock.Anything).Return([]*model.DiscussionInvite{&inviteObj}, nil)

			resp, err := backendObj.GetDiscussionInvitesByUserIDAndStatus(ctx, userID, status)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.DiscussionInvite{&inviteObj})
		})
	})
}

func TestDelphisBackend_GetSentDiscussionInvitesByUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := test_utils.UserID

	inviteObj := test_utils.TestDiscussionInvite(model.InviteRequestStatusPending)

	Convey("GetSentDiscussionInvitesByUserID", t, func() {
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
			mockDB.On("GetSentDiscussionInvitesByUserID", ctx, userID).Return(&mockDiscussionInviteIter{})
			mockDB.On("DiscussionInviteIterCollect", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetSentDiscussionInvitesByUserID(ctx, userID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetSentDiscussionInvitesByUserID", ctx, userID).Return(&mockDiscussionInviteIter{})
			mockDB.On("DiscussionInviteIterCollect", ctx, mock.Anything).Return([]*model.DiscussionInvite{&inviteObj}, nil)

			resp, err := backendObj.GetSentDiscussionInvitesByUserID(ctx, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.DiscussionInvite{&inviteObj})
		})
	})
}

func TestDelphisBackend_GetDiscussionAccessRequestsByDiscussionID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	discussionID := test_utils.DiscussionID

	requestObj := test_utils.TestDiscussionAccessRequest(model.InviteRequestStatusAccepted)

	Convey("GetDiscussionAccessRequestsByDiscussionID", t, func() {
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
			mockDB.On("GetDiscussionAccessRequestsByDiscussionID", ctx, discussionID).Return(&mockDiscussionAccessRequestIter{})
			mockDB.On("AccessRequestIterCollect", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetDiscussionAccessRequestsByDiscussionID(ctx, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetDiscussionAccessRequestsByDiscussionID", ctx, discussionID).Return(&mockDiscussionAccessRequestIter{})
			mockDB.On("AccessRequestIterCollect", ctx, mock.Anything).Return([]*model.DiscussionAccessRequest{&requestObj}, nil)

			resp, err := backendObj.GetDiscussionAccessRequestsByDiscussionID(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.DiscussionAccessRequest{&requestObj})
		})
	})
}

func TestDelphisBackend_GetSentDiscussionAccessRequestsByUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := test_utils.UserID

	requestObj := test_utils.TestDiscussionAccessRequest(model.InviteRequestStatusAccepted)

	Convey("GetSentDiscussionAccessRequestsByUserID", t, func() {
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
			mockDB.On("GetSentDiscussionAccessRequestsByUserID", ctx, userID).Return(&mockDiscussionAccessRequestIter{})
			mockDB.On("AccessRequestIterCollect", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.GetSentDiscussionAccessRequestsByUserID(ctx, userID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetSentDiscussionAccessRequestsByUserID", ctx, userID).Return(&mockDiscussionAccessRequestIter{})
			mockDB.On("AccessRequestIterCollect", ctx, mock.Anything).Return([]*model.DiscussionAccessRequest{&requestObj}, nil)

			resp, err := backendObj.GetSentDiscussionAccessRequestsByUserID(ctx, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []*model.DiscussionAccessRequest{&requestObj})
		})
	})
}

func TestDelphisBackend_InviteUserToDiscussion(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := test_utils.UserID
	discussionID := test_utils.DiscussionID
	participantID := test_utils.InvitingParticipantID

	inviteObj := test_utils.TestDiscussionInvite(model.InviteRequestStatusAccepted)

	tx := sql.Tx{}

	Convey("InviteUserToDiscussion", t, func() {
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

		Convey("when BeginTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(nil, expectedError)

			resp, err := backendObj.InviteUserToDiscussion(ctx, userID, discussionID, participantID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when PutDiscussionInviteRecord errors out and RollbackFails", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.InviteUserToDiscussion(ctx, userID, discussionID, participantID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when PutDiscussionInviteRecord errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.InviteUserToDiscussion(ctx, userID, discussionID, participantID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when CommitTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(&inviteObj, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.InviteUserToDiscussion(ctx, userID, discussionID, participantID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when invite succeeds", func() {
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(&inviteObj, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.InviteUserToDiscussion(ctx, userID, discussionID, participantID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_RequestAccessToDiscussion(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	userID := test_utils.UserID
	discussionID := test_utils.DiscussionID

	requestObj := test_utils.TestDiscussionAccessRequest(model.InviteRequestStatusAccepted)

	tx := sql.Tx{}

	Convey("RequestAccessToDiscussion", t, func() {
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

		Convey("when BeginTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(nil, expectedError)

			resp, err := backendObj.RequestAccessToDiscussion(ctx, userID, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when PutDiscussionAccessRequestRecord errors out and RollbackFails", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutDiscussionAccessRequestRecord", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.RequestAccessToDiscussion(ctx, userID, discussionID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when PutDiscussionAccessRequestRecord errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutDiscussionAccessRequestRecord", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.RequestAccessToDiscussion(ctx, userID, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when CommitTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutDiscussionAccessRequestRecord", ctx, mock.Anything, mock.Anything).Return(&requestObj, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.RequestAccessToDiscussion(ctx, userID, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when invite succeeds", func() {
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutDiscussionAccessRequestRecord", ctx, mock.Anything, mock.Anything).Return(&requestObj, nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.RequestAccessToDiscussion(ctx, userID, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_RespondToInvitation(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	inviteID := test_utils.InviteID
	response := model.InviteRequestStatusAccepted

	inputObj := test_utils.TestAddDiscussionParticipantInput()
	inviteObj := test_utils.TestDiscussionInvite(model.InviteRequestStatusAccepted)
	duaObj := test_utils.TestDiscussionUserAccess()
	userObj := test_utils.TestUser()
	flairObj := test_utils.TestFlair()
	viewerObj := test_utils.TestViewer()
	parObj := test_utils.TestParticipant()
	discObj := test_utils.TestDiscussion()
	profileObj := test_utils.TestUserProfile()

	userObj.UserProfile = &profileObj

	tx := sql.Tx{}

	Convey("RequestAccessToDiscussion", t, func() {
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

		Convey("when BeginTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(nil, expectedError)

			resp, err := backendObj.RespondToInvitation(ctx, inviteID, response, inputObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when UpdateDiscussionInviteRecord errors out and Rollback fails", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpdateDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.RespondToInvitation(ctx, inviteID, response, inputObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when UpdateDiscussionInviteRecord errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpdateDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.RespondToInvitation(ctx, inviteID, response, inputObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when UpsertDiscussionUserAccess errors out and Rollback fails", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpdateDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(&inviteObj, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, mock.Anything, duaObj).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.RespondToInvitation(ctx, inviteID, response, inputObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when UpsertDiscussionUserAccess errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpdateDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(&inviteObj, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, mock.Anything, duaObj).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.RespondToInvitation(ctx, inviteID, response, inputObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when CreateParticipantForDiscussion errors out and Rollback fails", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpdateDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(&inviteObj, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, mock.Anything, duaObj).Return(&duaObj, nil)

			// Create participant functions
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(nil, expectedError)

			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.RespondToInvitation(ctx, inviteID, response, inputObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when CreateParticipantForDiscussion errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpdateDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(&inviteObj, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, mock.Anything, duaObj).Return(&duaObj, nil)

			// Create participant functions
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(nil, expectedError)

			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.RespondToInvitation(ctx, inviteID, response, inputObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when response succeeds", func() {
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpdateDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(&inviteObj, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, mock.Anything, duaObj).Return(&duaObj, nil)

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
			mockDB.On("UpsertDiscussion", ctx, mock.Anything).Return(&discObj, nil)
			mockDB.On("GetParticipantsByDiscussionID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)
			mockDB.On("GetUserDevicesByUserID", ctx, mock.Anything).Return(nil, nil)

			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.RespondToInvitation(ctx, inviteID, response, inputObj)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_RespondToRequestAccess(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	participantID := test_utils.InvitingParticipantID
	requestID := test_utils.RequestID
	response := model.InviteRequestStatusAccepted

	requestObj := test_utils.TestDiscussionAccessRequest(model.InviteRequestStatusAccepted)
	inviteObj := test_utils.TestDiscussionInvite(model.InviteRequestStatusAccepted)
	duaObj := test_utils.TestDiscussionUserAccess()

	duaObj.RequestID = &requestObj.ID

	tx := sql.Tx{}

	Convey("RespondToRequestAccess", t, func() {
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

		Convey("when BeginTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(nil, expectedError)

			resp, err := backendObj.RespondToRequestAccess(ctx, requestID, response, participantID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when UpdateDiscussionAccessRequestRecord errors out and Rollback fails", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpdateDiscussionAccessRequestRecord", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.RespondToRequestAccess(ctx, requestID, response, participantID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when UpdateDiscussionAccessRequestRecord errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpdateDiscussionAccessRequestRecord", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.RespondToRequestAccess(ctx, requestID, response, participantID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when UpsertDiscussionUserAccess errors out and Rollback fails", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpdateDiscussionAccessRequestRecord", ctx, mock.Anything, mock.Anything).Return(&requestObj, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, mock.Anything, duaObj).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.RespondToRequestAccess(ctx, requestID, response, participantID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when UpsertDiscussionUserAccess errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpdateDiscussionAccessRequestRecord", ctx, mock.Anything, mock.Anything).Return(&requestObj, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, mock.Anything, duaObj).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.RespondToRequestAccess(ctx, requestID, response, participantID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when PutDiscussionInviteRecord errors out and Rollback fails", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpdateDiscussionAccessRequestRecord", ctx, mock.Anything, mock.Anything).Return(&requestObj, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, mock.Anything, duaObj).Return(&duaObj, nil)
			mockDB.On("PutDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.RespondToRequestAccess(ctx, requestID, response, participantID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when PutDiscussionInviteRecord errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpdateDiscussionAccessRequestRecord", ctx, mock.Anything, mock.Anything).Return(&requestObj, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, mock.Anything, duaObj).Return(&duaObj, nil)
			mockDB.On("PutDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(nil, expectedError)
			mockDB.On("RollbackTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.RespondToRequestAccess(ctx, requestID, response, participantID)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when CommitTx errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpdateDiscussionAccessRequestRecord", ctx, mock.Anything, mock.Anything).Return(&requestObj, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, mock.Anything, duaObj).Return(&duaObj, nil)
			mockDB.On("PutDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(&inviteObj, nil)

			mockDB.On("CommitTx", ctx, mock.Anything).Return(expectedError)

			resp, err := backendObj.RespondToRequestAccess(ctx, requestID, response, participantID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when response succeeds", func() {
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("UpdateDiscussionAccessRequestRecord", ctx, mock.Anything, mock.Anything).Return(&requestObj, nil)
			mockDB.On("UpsertDiscussionUserAccess", ctx, mock.Anything, duaObj).Return(&duaObj, nil)
			mockDB.On("PutDiscussionInviteRecord", ctx, mock.Anything, mock.Anything).Return(&inviteObj, nil)

			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)

			resp, err := backendObj.RespondToRequestAccess(ctx, requestID, response, participantID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

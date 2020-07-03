package backend

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/auth"
	"github.com/nedrocks/delphisbe/internal/backend/test_utils"
	"github.com/nedrocks/delphisbe/internal/cache"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/nedrocks/delphisbe/internal/util"
	"github.com/nedrocks/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestDelphisBackend_CreateParticipantForDiscussion(t *testing.T) {
	ctx := context.Background()

	discussionID := test_utils.DiscussionID
	userID := test_utils.UserID
	count := 10

	inputObj := test_utils.TestAddDiscussionParticipantInput()
	userObj := test_utils.TestUser()
	modObj := test_utils.TestModerator()
	profile := test_utils.TestUserProfile()
	discObj := test_utils.TestDiscussion()
	flairObj := test_utils.TestFlair()
	viewerObj := test_utils.TestViewer()
	parObj := test_utils.TestParticipant()

	userObj.Flairs = []*model.Flair{&flairObj}
	userObj.UserProfile = &profile
	modObj.UserProfile = &profile

	tx := sql.Tx{}

	Convey("CreateParticipantForDiscussion", t, func() {
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

		Convey("when GetUserByID errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.CreateParticipantForDiscussion(ctx, discussionID, userID, inputObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when GetUserByID does not fnd a record for the user", func() {
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(nil, nil)

			resp, err := backendObj.CreateParticipantForDiscussion(ctx, discussionID, userID, inputObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when GetFlairsByUserID errors out", func() {
			tempUserObj := userObj
			tempUserObj.Flairs = nil
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(&tempUserObj, nil)
			mockDB.On("GetTotalParticipantCountByDiscussionID", ctx, discussionID).Return(count)
			mockDB.On("GetFlairsByUserID", ctx, userID).Return(nil, expectedError)

			resp, err := backendObj.CreateParticipantForDiscussion(ctx, discussionID, userID, inputObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when UpsertViewer errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(&userObj, nil)
			mockDB.On("GetTotalParticipantCountByDiscussionID", ctx, discussionID).Return(count)
			mockDB.On("UpsertViewer", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.CreateParticipantForDiscussion(ctx, discussionID, userID, inputObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when UpsertParticipant errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(&userObj, nil)
			mockDB.On("GetTotalParticipantCountByDiscussionID", ctx, discussionID).Return(count)
			mockDB.On("GetFlairsByUserID", ctx, userID).Return(&flairObj, nil)
			mockDB.On("UpsertViewer", ctx, mock.Anything).Return(&viewerObj, nil)
			mockDB.On("UpsertParticipant", ctx, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.CreateParticipantForDiscussion(ctx, discussionID, userID, inputObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when CreateAlertPost errors out", func() {
			expectedError := fmt.Errorf("Some Error")
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(&userObj, nil)
			mockDB.On("GetTotalParticipantCountByDiscussionID", ctx, discussionID).Return(count)
			mockDB.On("UpsertViewer", ctx, mock.Anything).Return(&viewerObj, nil)
			mockDB.On("UpsertParticipant", ctx, mock.Anything).Return(&parObj, nil)

			//// Create post functions
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)
			mockDB.On("BeginTx", ctx).Return(nil, expectedError)

			resp, err := backendObj.CreateParticipantForDiscussion(ctx, discussionID, userID, inputObj)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})

		Convey("when participant is created successfully", func() {
			mockDB.On("GetUserByID", ctx, mock.Anything).Return(&userObj, nil)
			mockDB.On("GetTotalParticipantCountByDiscussionID", ctx, discussionID).Return(count)
			mockDB.On("UpsertViewer", ctx, mock.Anything).Return(&viewerObj, nil)
			mockDB.On("UpsertParticipant", ctx, mock.Anything).Return(&parObj, nil)

			//// Create post functions
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)
			mockDB.On("BeginTx", ctx).Return(&tx, nil)
			mockDB.On("PutPostContent", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("PutPost", ctx, mock.Anything, mock.Anything).Return(&model.Post{ID: "post123"}, nil)
			mockDB.On("PutActivity", ctx, mock.Anything, mock.Anything).Return(nil)
			mockDB.On("CommitTx", ctx, mock.Anything).Return(nil)
			mockDB.On("GetDiscussionByID", ctx, mock.Anything).Return(&discObj, nil)
			mockDB.On("GetParticipantsByDiscussionID", ctx, mock.Anything, mock.Anything).Return([]model.Participant{parObj}, nil)
			mockDB.On("GetUserDevicesByUserID", ctx, mock.Anything).Return(nil, nil)

			resp, err := backendObj.CreateParticipantForDiscussion(ctx, discussionID, userID, inputObj)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_GetParticipantsByDiscussionID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	discussionID := test_utils.DiscussionID

	parObj := test_utils.TestParticipant()

	Convey("GetParticipantsByDiscussionID", t, func() {
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
			mockDB.On("GetParticipantsByDiscussionID", ctx, discussionID).Return(nil, expectedError)

			resp, err := backendObj.GetParticipantsByDiscussionID(ctx, discussionID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetParticipantsByDiscussionID", ctx, discussionID).Return([]model.Participant{parObj}, nil)

			resp, err := backendObj.GetParticipantsByDiscussionID(ctx, discussionID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, []model.Participant{parObj})
		})
	})
}

func TestDelphisBackend_GetParticipantsByDiscussionIDUserID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	discussionID := test_utils.DiscussionID
	userID := test_utils.UserID

	parObj := test_utils.TestParticipant()

	Convey("GetParticipantsByDiscussionIDUserID", t, func() {
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
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, discussionID, userID).Return(nil, expectedError)

			resp, err := backendObj.GetParticipantsByDiscussionIDUserID(ctx, discussionID, userID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetParticipantsByDiscussionIDUserID", ctx, discussionID, userID).Return([]model.Participant{parObj}, nil)

			resp, err := backendObj.GetParticipantsByDiscussionIDUserID(ctx, discussionID, userID)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func TestDelphisBackend_GetParticipantByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	participantID := test_utils.ParticipantID

	parObj := test_utils.TestParticipant()

	Convey("GetParticipantByID", t, func() {
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
			mockDB.On("GetParticipantByID", ctx, participantID).Return(nil, expectedError)

			resp, err := backendObj.GetParticipantByID(ctx, participantID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetParticipantByID", ctx, participantID).Return(&parObj, nil)

			resp, err := backendObj.GetParticipantByID(ctx, participantID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &parObj)
		})
	})
}

func TestDelphisBackend_GetParticipantsByIDs(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	participantID := test_utils.ParticipantID

	parObj := test_utils.TestParticipant()

	participantIDs := []string{participantID}

	verifyMap := map[string]*model.Participant{
		parObj.ID: &parObj,
	}

	Convey("GetParticipantsByIDs", t, func() {
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
			mockDB.On("GetParticipantsByIDs", ctx, participantIDs).Return(nil, expectedError)

			resp, err := backendObj.GetParticipantsByIDs(ctx, participantIDs)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("GetParticipantsByIDs", ctx, participantIDs).Return(verifyMap, nil)

			resp, err := backendObj.GetParticipantsByIDs(ctx, participantIDs)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, verifyMap)
		})
	})
}

func TestDelphisBackend_AssignFlair(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	flairID := test_utils.FlairID

	parObj := test_utils.TestParticipant()

	Convey("AssignFlair", t, func() {
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
			mockDB.On("AssignFlair", ctx, parObj, &flairID).Return(nil, expectedError)

			resp, err := backendObj.AssignFlair(ctx, parObj, flairID)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("AssignFlair", ctx, parObj, &flairID).Return(&parObj, nil)

			resp, err := backendObj.AssignFlair(ctx, parObj, flairID)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &parObj)
		})
	})
}

func TestDelphisBackend_UnassignFlair(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	parObj := test_utils.TestParticipant()

	Convey("UnassignFlair", t, func() {
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
			mockDB.On("AssignFlair", ctx, parObj, mock.Anything).Return(nil, expectedError)

			resp, err := backendObj.UnassignFlair(ctx, parObj)

			So(err, ShouldEqual, expectedError)
			So(resp, ShouldBeNil)
		})

		Convey("when the query returns successfully", func() {
			mockDB.On("AssignFlair", ctx, parObj, mock.Anything).Return(&parObj, nil)

			resp, err := backendObj.UnassignFlair(ctx, parObj)

			So(err, ShouldBeNil)
			So(resp, ShouldResemble, &parObj)
		})
	})
}

func TestDelphisBackend_UpdateParticipant(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	participantID := test_utils.ParticipantID
	anonParObj := test_utils.TestParticipant()
	nonAnonParObj := test_utils.TestParticipant()
	updateObj := test_utils.TestUpdateParticipantInput()

	anonParObj.IsAnonymous = true
	participants := UserDiscussionParticipants{
		Anon: &anonParObj,
	}

	Convey("UnassignFlair", t, func() {
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

		Convey("when the participantID doesn't match either participant passed in", func() {

			resp, err := backendObj.UpdateParticipant(ctx, participants, "badID", updateObj)

			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("when the user can create an non-anonymous participant", func() {
			mockDB.On("GetTotalParticipantCountByDiscussionID", ctx, mock.Anything).Return(10)
			mockDB.On("UpsertParticipant", ctx, mock.Anything).Return(&anonParObj, nil)

			resp, err := backendObj.UpdateParticipant(ctx, participants, participantID, updateObj)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})

		Convey("when the user can create an anonymous participant", func() {
			testParticipants := participants
			testParticipants.NonAnon = &nonAnonParObj
			testParticipants.Anon = nil

			mockDB.On("GetTotalParticipantCountByDiscussionID", ctx, mock.Anything).Return(10)
			mockDB.On("UpsertParticipant", ctx, mock.Anything).Return(&anonParObj, nil)

			resp, err := backendObj.UpdateParticipant(ctx, testParticipants, participantID, updateObj)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})

		Convey("when the user is switching to an existing participant", func() {
			testParticipants := participants
			testParticipants.NonAnon = &nonAnonParObj
			testParticipants.Anon.HasJoined = false

			mockDB.On("UpsertParticipant", ctx, mock.Anything).Return(&anonParObj, nil)

			resp, err := backendObj.UpdateParticipant(ctx, testParticipants, participantID, updateObj)

			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
		//
		//Convey("when the query returns successfully", func() {
		//	mockDB.On("AssignFlair", ctx, parObj, mock.Anything).Return(&parObj, nil)
		//
		//	resp, err := backendObj.UnassignFlair(ctx, parObj)
		//
		//	So(err, ShouldBeNil)
		//	So(resp, ShouldResemble, &parObj)
		//})
	})
}

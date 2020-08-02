package backend

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/delphis-inc/delphisbe/internal/backend/test_utils"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/cache"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/delphis-inc/delphisbe/mocks"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestDelphisBackend_GetConciergeParticipantID(t *testing.T) {
	ctx := context.Background()
	discussionID := test_utils.DiscussionID

	parObj := test_utils.TestParticipant()

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

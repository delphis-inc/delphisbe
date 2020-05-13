// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/nedrocks/delphisbe/graph/model"
)

// Datastore is an autogenerated mock type for the Datastore type
type Datastore struct {
	mock.Mock
}

// AssignFlair provides a mock function with given fields: ctx, participant, flairID
func (_m *Datastore) AssignFlair(ctx context.Context, participant model.Participant, flairID *string) (*model.Participant, error) {
	ret := _m.Called(ctx, participant, flairID)

	var r0 *model.Participant
	if rf, ok := ret.Get(0).(func(context.Context, model.Participant, *string) *model.Participant); ok {
		r0 = rf(ctx, participant, flairID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Participant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.Participant, *string) error); ok {
		r1 = rf(ctx, participant, flairID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateModerator provides a mock function with given fields: ctx, moderator
func (_m *Datastore) CreateModerator(ctx context.Context, moderator model.Moderator) (*model.Moderator, error) {
	ret := _m.Called(ctx, moderator)

	var r0 *model.Moderator
	if rf, ok := ret.Get(0).(func(context.Context, model.Moderator) *model.Moderator); ok {
		r0 = rf(ctx, moderator)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Moderator)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.Moderator) error); ok {
		r1 = rf(ctx, moderator)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateOrUpdateUserProfile provides a mock function with given fields: ctx, userProfile
func (_m *Datastore) CreateOrUpdateUserProfile(ctx context.Context, userProfile model.UserProfile) (*model.UserProfile, bool, error) {
	ret := _m.Called(ctx, userProfile)

	var r0 *model.UserProfile
	if rf, ok := ret.Get(0).(func(context.Context, model.UserProfile) *model.UserProfile); ok {
		r0 = rf(ctx, userProfile)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserProfile)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(context.Context, model.UserProfile) bool); ok {
		r1 = rf(ctx, userProfile)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, model.UserProfile) error); ok {
		r2 = rf(ctx, userProfile)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetDiscussionByID provides a mock function with given fields: ctx, id
func (_m *Datastore) GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.Discussion
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Discussion); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Discussion)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDiscussionByModeratorID provides a mock function with given fields: ctx, moderatorID
func (_m *Datastore) GetDiscussionByModeratorID(ctx context.Context, moderatorID string) (*model.Discussion, error) {
	ret := _m.Called(ctx, moderatorID)

	var r0 *model.Discussion
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Discussion); ok {
		r0 = rf(ctx, moderatorID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Discussion)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, moderatorID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDiscussionsByIDs provides a mock function with given fields: ctx, ids
func (_m *Datastore) GetDiscussionsByIDs(ctx context.Context, ids []string) (map[string]*model.Discussion, error) {
	ret := _m.Called(ctx, ids)

	var r0 map[string]*model.Discussion
	if rf, ok := ret.Get(0).(func(context.Context, []string) map[string]*model.Discussion); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*model.Discussion)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFlairByID provides a mock function with given fields: ctx, id
func (_m *Datastore) GetFlairByID(ctx context.Context, id string) (*model.Flair, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.Flair
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Flair); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Flair)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFlairTemplateByID provides a mock function with given fields: ctx, id
func (_m *Datastore) GetFlairTemplateByID(ctx context.Context, id string) (*model.FlairTemplate, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.FlairTemplate
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.FlairTemplate); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FlairTemplate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFlairsByUserID provides a mock function with given fields: ctx, userID
func (_m *Datastore) GetFlairsByUserID(ctx context.Context, userID string) ([]*model.Flair, error) {
	ret := _m.Called(ctx, userID)

	var r0 []*model.Flair
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.Flair); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Flair)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetModeratorByID provides a mock function with given fields: ctx, id
func (_m *Datastore) GetModeratorByID(ctx context.Context, id string) (*model.Moderator, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.Moderator
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Moderator); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Moderator)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetParticipantByDiscussionIDUserID provides a mock function with given fields: ctx, discussionID, userID
func (_m *Datastore) GetParticipantByDiscussionIDUserID(ctx context.Context, discussionID string, userID string) (*model.Participant, error) {
	ret := _m.Called(ctx, discussionID, userID)

	var r0 *model.Participant
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.Participant); ok {
		r0 = rf(ctx, discussionID, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Participant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, discussionID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetParticipantByID provides a mock function with given fields: ctx, participantID
func (_m *Datastore) GetParticipantByID(ctx context.Context, participantID string) (*model.Participant, error) {
	ret := _m.Called(ctx, participantID)

	var r0 *model.Participant
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Participant); ok {
		r0 = rf(ctx, participantID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Participant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, participantID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetParticipantsByDiscussionID provides a mock function with given fields: ctx, id
func (_m *Datastore) GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error) {
	ret := _m.Called(ctx, id)

	var r0 []model.Participant
	if rf, ok := ret.Get(0).(func(context.Context, string) []model.Participant); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Participant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPostContentByID provides a mock function with given fields: ctx, id
func (_m *Datastore) GetPostContentByID(ctx context.Context, id string) (*model.PostContent, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.PostContent
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.PostContent); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.PostContent)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPostsByDiscussionID provides a mock function with given fields: ctx, discussionID
func (_m *Datastore) GetPostsByDiscussionID(ctx context.Context, discussionID string) ([]*model.Post, error) {
	ret := _m.Called(ctx, discussionID)

	var r0 []*model.Post
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.Post); ok {
		r0 = rf(ctx, discussionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Post)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, discussionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSocialInfosByUserProfileID provides a mock function with given fields: ctx, userProfileID
func (_m *Datastore) GetSocialInfosByUserProfileID(ctx context.Context, userProfileID string) ([]model.SocialInfo, error) {
	ret := _m.Called(ctx, userProfileID)

	var r0 []model.SocialInfo
	if rf, ok := ret.Get(0).(func(context.Context, string) []model.SocialInfo); ok {
		r0 = rf(ctx, userProfileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.SocialInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userProfileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTotalParticipantCountByDiscussionID provides a mock function with given fields: ctx, discussionID
func (_m *Datastore) GetTotalParticipantCountByDiscussionID(ctx context.Context, discussionID string) int {
	ret := _m.Called(ctx, discussionID)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, string) int); ok {
		r0 = rf(ctx, discussionID)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetUserByID provides a mock function with given fields: ctx, userID
func (_m *Datastore) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	ret := _m.Called(ctx, userID)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.User); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserDevicesByUserID provides a mock function with given fields: ctx, userID
func (_m *Datastore) GetUserDevicesByUserID(ctx context.Context, userID string) ([]model.UserDevice, error) {
	ret := _m.Called(ctx, userID)

	var r0 []model.UserDevice
	if rf, ok := ret.Get(0).(func(context.Context, string) []model.UserDevice); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.UserDevice)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserProfileByID provides a mock function with given fields: ctx, id
func (_m *Datastore) GetUserProfileByID(ctx context.Context, id string) (*model.UserProfile, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.UserProfile
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.UserProfile); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserProfile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserProfileByUserID provides a mock function with given fields: ctx, userID
func (_m *Datastore) GetUserProfileByUserID(ctx context.Context, userID string) (*model.UserProfile, error) {
	ret := _m.Called(ctx, userID)

	var r0 *model.UserProfile
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.UserProfile); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserProfile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetViewersByIDs provides a mock function with given fields: ctx, viewerIDs
func (_m *Datastore) GetViewersByIDs(ctx context.Context, viewerIDs []string) (map[string]*model.Viewer, error) {
	ret := _m.Called(ctx, viewerIDs)

	var r0 map[string]*model.Viewer
	if rf, ok := ret.Get(0).(func(context.Context, []string) map[string]*model.Viewer); ok {
		r0 = rf(ctx, viewerIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*model.Viewer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, viewerIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDiscussions provides a mock function with given fields: ctx
func (_m *Datastore) ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error) {
	ret := _m.Called(ctx)

	var r0 *model.DiscussionsConnection
	if rf, ok := ret.Get(0).(func(context.Context) *model.DiscussionsConnection); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DiscussionsConnection)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListFlairTemplates provides a mock function with given fields: ctx, query
func (_m *Datastore) ListFlairTemplates(ctx context.Context, query *string) ([]*model.FlairTemplate, error) {
	ret := _m.Called(ctx, query)

	var r0 []*model.FlairTemplate
	if rf, ok := ret.Get(0).(func(context.Context, *string) []*model.FlairTemplate); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.FlairTemplate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *string) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutParticipant provides a mock function with given fields: ctx, participant
func (_m *Datastore) PutParticipant(ctx context.Context, participant model.Participant) (*model.Participant, error) {
	ret := _m.Called(ctx, participant)

	var r0 *model.Participant
	if rf, ok := ret.Get(0).(func(context.Context, model.Participant) *model.Participant); ok {
		r0 = rf(ctx, participant)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Participant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.Participant) error); ok {
		r1 = rf(ctx, participant)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutPost provides a mock function with given fields: ctx, post
func (_m *Datastore) PutPost(ctx context.Context, post model.Post) (*model.Post, error) {
	ret := _m.Called(ctx, post)

	var r0 *model.Post
	if rf, ok := ret.Get(0).(func(context.Context, model.Post) *model.Post); ok {
		r0 = rf(ctx, post)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Post)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.Post) error); ok {
		r1 = rf(ctx, post)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveFlair provides a mock function with given fields: ctx, flair
func (_m *Datastore) RemoveFlair(ctx context.Context, flair model.Flair) (*model.Flair, error) {
	ret := _m.Called(ctx, flair)

	var r0 *model.Flair
	if rf, ok := ret.Get(0).(func(context.Context, model.Flair) *model.Flair); ok {
		r0 = rf(ctx, flair)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Flair)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.Flair) error); ok {
		r1 = rf(ctx, flair)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveFlairTemplate provides a mock function with given fields: ctx, flairTemplate
func (_m *Datastore) RemoveFlairTemplate(ctx context.Context, flairTemplate model.FlairTemplate) (*model.FlairTemplate, error) {
	ret := _m.Called(ctx, flairTemplate)

	var r0 *model.FlairTemplate
	if rf, ok := ret.Get(0).(func(context.Context, model.FlairTemplate) *model.FlairTemplate); ok {
		r0 = rf(ctx, flairTemplate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FlairTemplate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.FlairTemplate) error); ok {
		r1 = rf(ctx, flairTemplate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertDiscussion provides a mock function with given fields: ctx, discussion
func (_m *Datastore) UpsertDiscussion(ctx context.Context, discussion model.Discussion) (*model.Discussion, error) {
	ret := _m.Called(ctx, discussion)

	var r0 *model.Discussion
	if rf, ok := ret.Get(0).(func(context.Context, model.Discussion) *model.Discussion); ok {
		r0 = rf(ctx, discussion)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Discussion)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.Discussion) error); ok {
		r1 = rf(ctx, discussion)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertFlair provides a mock function with given fields: ctx, flair
func (_m *Datastore) UpsertFlair(ctx context.Context, flair model.Flair) (*model.Flair, error) {
	ret := _m.Called(ctx, flair)

	var r0 *model.Flair
	if rf, ok := ret.Get(0).(func(context.Context, model.Flair) *model.Flair); ok {
		r0 = rf(ctx, flair)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Flair)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.Flair) error); ok {
		r1 = rf(ctx, flair)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertFlairTemplate provides a mock function with given fields: ctx, flairTemplate
func (_m *Datastore) UpsertFlairTemplate(ctx context.Context, flairTemplate model.FlairTemplate) (*model.FlairTemplate, error) {
	ret := _m.Called(ctx, flairTemplate)

	var r0 *model.FlairTemplate
	if rf, ok := ret.Get(0).(func(context.Context, model.FlairTemplate) *model.FlairTemplate); ok {
		r0 = rf(ctx, flairTemplate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FlairTemplate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.FlairTemplate) error); ok {
		r1 = rf(ctx, flairTemplate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertSocialInfo provides a mock function with given fields: ctx, obj
func (_m *Datastore) UpsertSocialInfo(ctx context.Context, obj model.SocialInfo) (*model.SocialInfo, error) {
	ret := _m.Called(ctx, obj)

	var r0 *model.SocialInfo
	if rf, ok := ret.Get(0).(func(context.Context, model.SocialInfo) *model.SocialInfo); ok {
		r0 = rf(ctx, obj)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.SocialInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.SocialInfo) error); ok {
		r1 = rf(ctx, obj)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertUser provides a mock function with given fields: ctx, user
func (_m *Datastore) UpsertUser(ctx context.Context, user model.User) (*model.User, error) {
	ret := _m.Called(ctx, user)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(context.Context, model.User) *model.User); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertUserDevice provides a mock function with given fields: ctx, userDevice
func (_m *Datastore) UpsertUserDevice(ctx context.Context, userDevice model.UserDevice) (*model.UserDevice, error) {
	ret := _m.Called(ctx, userDevice)

	var r0 *model.UserDevice
	if rf, ok := ret.Get(0).(func(context.Context, model.UserDevice) *model.UserDevice); ok {
		r0 = rf(ctx, userDevice)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserDevice)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.UserDevice) error); ok {
		r1 = rf(ctx, userDevice)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertViewer provides a mock function with given fields: ctx, viewer
func (_m *Datastore) UpsertViewer(ctx context.Context, viewer model.Viewer) (*model.Viewer, error) {
	ret := _m.Called(ctx, viewer)

	var r0 *model.Viewer
	if rf, ok := ret.Get(0).(func(context.Context, model.Viewer) *model.Viewer); ok {
		r0 = rf(ctx, viewer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Viewer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.Viewer) error); ok {
		r1 = rf(ctx, viewer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

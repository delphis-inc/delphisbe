package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"sort"

	"github.com/delphis-inc/delphisbe/graph/generated"
	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/sirupsen/logrus"
)

func (r *userResolver) Participants(ctx context.Context, obj *model.User) ([]*model.Participant, error) {
	participants := make([]*model.Participant, 0)
	for _, p := range obj.Participants {
		discussion := p.Discussion
		if discussion == nil && p.DiscussionID != nil {
			var err error
			discussion, err = r.DAOManager.GetDiscussionByID(ctx, *p.DiscussionID)
			if err != nil {
				logrus.Debugf("Caught an error getting discussion for participant")
				// Silently ignore for now.
			}
			p.Discussion = discussion
		}
		if discussion != nil && discussion.DeletedAt == nil {
			participants = append(participants, p)
		}
	}

	sort.Slice(participants, func(i, j int) bool {
		return participants[i].CreatedAt.Before(participants[j].CreatedAt)
	})

	return participants, nil
}

func (r *userResolver) Viewers(ctx context.Context, obj *model.User) ([]*model.Viewer, error) {
	viewers := make([]*model.Viewer, 0)
	for _, v := range obj.Viewers {
		discussion := v.Discussion
		if v != nil && discussion == nil && v.DiscussionID != nil {
			var err error
			discussion, err = r.DAOManager.GetDiscussionByID(ctx, *v.DiscussionID)
			if err != nil {
				logrus.Debugf("Caught an error getting discussion for participant")
				// Silently ignore for now.
			}
			v.Discussion = discussion
		}
		if v != nil && v.Discussion != nil && v.Discussion.DeletedAt == nil {
			viewers = append(viewers, v)
		}
	}
	sort.Slice(viewers, func(i, j int) bool {
		return viewers[i].CreatedAt.Before(viewers[j].CreatedAt)
	})
	return viewers, nil
}

func (r *userResolver) Bookmarks(ctx context.Context, obj *model.User) ([]*model.PostBookmark, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Profile(ctx context.Context, obj *model.User) (*model.UserProfile, error) {
	if obj.UserProfile == nil {
		userProfile, err := r.DAOManager.GetUserProfileByUserID(ctx, obj.ID)
		if err != nil {
			return nil, err
		}
		obj.UserProfile = userProfile
	}
	return obj.UserProfile, nil
}

func (r *userResolver) Flairs(ctx context.Context, obj *model.User) ([]*model.Flair, error) {
	return r.DAOManager.GetFlairsByUserID(ctx, obj.ID)
}

func (r *userResolver) Devices(ctx context.Context, obj *model.User) ([]*model.UserDevice, error) {
	devices, err := r.DAOManager.GetUserDevicesByUserID(ctx, obj.ID)
	if err != nil {
		return nil, err
	}

	resp := make([]*model.UserDevice, len(devices))
	for i := range devices {
		resp[i] = &devices[i]
	}
	return resp, nil
}

func (r *userResolver) Discussions(ctx context.Context, obj *model.User, state model.DiscussionUserAccessState) ([]*model.Discussion, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	// Do we want a mod override?
	if authedUser.UserID != obj.ID {
		return nil, fmt.Errorf("unauthorized")
	}

	return r.DAOManager.GetDiscussionAccessByUserID(ctx, authedUser.UserID, state)
}

func (r *userResolver) DiscussionInvites(ctx context.Context, obj *model.User, status model.InviteRequestStatus) ([]*model.DiscussionInvite, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	// Do we want a mod override?
	if authedUser.UserID != obj.ID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Do we want a different status? Made it generic for easy change
	return r.DAOManager.GetDiscussionInvitesByUserIDAndStatus(ctx, authedUser.UserID, status)
}

func (r *userResolver) SentDiscussionInvites(ctx context.Context, obj *model.User) ([]*model.DiscussionInvite, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	// Do we want a mod override?
	if authedUser.UserID != obj.ID {
		return nil, fmt.Errorf("unauthorized")
	}

	return r.DAOManager.GetSentDiscussionInvitesByUserID(ctx, authedUser.UserID)
}

func (r *userResolver) SentDiscussionAccessRequests(ctx context.Context, obj *model.User) ([]*model.DiscussionAccessRequest, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	// Do we want a mod override?
	if authedUser.UserID != obj.ID {
		return nil, fmt.Errorf("unauthorized")
	}

	return r.DAOManager.GetSentDiscussionAccessRequestsByUserID(ctx, authedUser.UserID)
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }

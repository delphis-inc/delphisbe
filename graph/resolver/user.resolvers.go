// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"
	"fmt"
	"sort"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *userResolver) Participants(ctx context.Context, obj *model.User) ([]*model.Participant, error) {
	if obj.Participants == nil || len(obj.Participants) == 0 {
		participantMap, err := r.DAOManager.GetParticipantsByIDs(ctx, obj.DiscussionParticipants.Keys)

		if err != nil {
			return nil, err
		}

		participants := make([]*model.Participant, 0)
		for _, p := range participantMap {
			if p != nil {
				participants = append(participants, p)
			}
		}

		sort.Slice(participants, func(i, j int) bool {
			return participants[i].CreatedAt.Before(participants[j].CreatedAt)
		})
		obj.Participants = participants
	}
	return obj.Participants, nil
}

func (r *userResolver) Viewers(ctx context.Context, obj *model.User) ([]*model.Viewer, error) {
	if obj.Viewers == nil || len(obj.Viewers) == 0 {
		viewerMap, err := r.DAOManager.GetViewersByIDs(ctx, obj.DiscussionViewers.Keys)

		if err != nil {
			return nil, err
		}

		viewers := make([]*model.Viewer, 0)
		for _, v := range viewerMap {
			if v != nil {
				viewers = append(viewers, v)
			}
		}

		sort.Slice(viewers, func(i, j int) bool {
			return viewers[i].CreatedAt.Before(viewers[j].CreatedAt)
		})
		obj.Viewers = viewers
	}
	return obj.Viewers, nil
}

func (r *userResolver) Bookmarks(ctx context.Context, obj *model.User) ([]*model.PostBookmark, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Profile(ctx context.Context, obj *model.User) (*model.UserProfile, error) {
	if obj.UserProfile == nil {
		userProfile, err := r.DAOManager.GetUserProfileByID(ctx, obj.UserProfileID)
		if err != nil {
			return nil, err
		}
		obj.UserProfile = userProfile
	}
	return obj.UserProfile, nil
}

func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }

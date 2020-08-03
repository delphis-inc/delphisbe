package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/delphis-inc/delphisbe/graph/generated"
	"github.com/delphis-inc/delphisbe/graph/model"
)

func (r *participantResolver) Discussion(ctx context.Context, obj *model.Participant) (*model.Discussion, error) {
	if obj.Discussion == nil && obj.DiscussionID != nil {
		disc, err := r.resolveDiscussionByID(ctx, *obj.DiscussionID)
		if err != nil {
			return nil, err
		}
		obj.Discussion = disc
	}

	return obj.Discussion, nil
}

func (r *participantResolver) Viewer(ctx context.Context, obj *model.Participant) (*model.Viewer, error) {
	if obj.IsBanned {
		return nil, nil
	}

	if obj.Viewer == nil && obj.ViewerID != nil {
		viewerObj, err := r.DAOManager.GetViewerByID(ctx, *obj.ViewerID)

		if err != nil {
			return nil, err
		}

		obj.Viewer = viewerObj
	}
	return obj.Viewer, nil
}

func (r *participantResolver) DiscussionNotificationPreferences(ctx context.Context, obj *model.Participant) (model.DiscussionNotificationPreferences, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *participantResolver) Posts(ctx context.Context, obj *model.Participant) ([]*model.Post, error) {
	if obj.IsBanned {
		return nil, nil
	}

	if obj.DiscussionID == nil {
		return nil, nil
	}
	posts, err := r.DAOManager.GetPostsByDiscussionID(ctx, *obj.UserID, *obj.DiscussionID)
	if err != nil {
		return nil, err
	}

	n := 0
	for _, post := range posts {
		if post.ParticipantID != nil {
			if *post.ParticipantID == obj.ID {
				posts[n] = post
				n++
			}
		}
	}
	posts = posts[:n]
	return posts, nil
}

func (r *participantResolver) Flair(ctx context.Context, obj *model.Participant) (*model.Flair, error) {
	if obj.IsBanned {
		return nil, nil
	}

	if obj.Flair == nil && obj.FlairID != nil {
		return r.DAOManager.GetFlairByID(ctx, *obj.FlairID)
	}
	return obj.Flair, nil
}

func (r *participantResolver) Inviter(ctx context.Context, obj *model.Participant) (*model.Participant, error) {
	if obj.InviterID == nil {
		// By default the inviter is the moderator
		inviter, err := r.DAOManager.GetModeratorParticipantsByDiscussionID(ctx, *obj.DiscussionID)
		if err != nil {
			return nil, err
		}
		if inviter == nil {
			return nil, fmt.Errorf("Could not retrieve discussion's moderator")
		}

		if inviter.NonAnon != nil {
			return inviter.NonAnon, nil
		} else if inviter.Anon != nil {
			return inviter.Anon, nil
		}

		// We have a problem here
		return nil, nil
	}

	inviter, err := r.DAOManager.GetParticipantByID(ctx, *obj.InviterID)
	if err != nil {
		return nil, err
	}
	return inviter, nil
}

func (r *participantResolver) UserProfile(ctx context.Context, obj *model.Participant) (*model.UserProfile, error) {
	if obj.IsBanned {
		return nil, nil
	}

	if obj.IsAnonymous || obj.UserID == nil {
		return nil, nil
	}

	userProfile, err := r.DAOManager.GetUserProfileByUserID(ctx, *obj.UserID)

	if err != nil {
		return nil, err
	}

	return userProfile, nil
}

func (r *participantResolver) AnonDisplayName(ctx context.Context, obj *model.Participant) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Participant returns generated.ParticipantResolver implementation.
func (r *Resolver) Participant() generated.ParticipantResolver { return &participantResolver{r} }

type participantResolver struct{ *Resolver }

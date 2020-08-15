package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/delphis-inc/delphisbe/graph/generated"
	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/util"
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

func (r *participantResolver) Posts(ctx context.Context, obj *model.Participant) ([]*model.Post, error) {
	if obj.IsBanned {
		return nil, nil
	}

	if obj.DiscussionID == nil {
		return nil, nil
	}
	posts, err := r.DAOManager.GetPostsByDiscussionID(ctx, *obj.DiscussionID)
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

func (r *participantResolver) GradientColor(ctx context.Context, obj *model.Participant) (*model.GradientColor, error) {
	var gradientColor model.GradientColor = model.GradientColorUnknown
	if obj.IsAnonymous {
		if obj.Discussion == nil && obj.DiscussionID != nil {
			disc, err := r.resolveDiscussionByID(ctx, *obj.DiscussionID)
			if err != nil {
				return nil, err
			}
			obj.Discussion = disc
		}
		if obj.Discussion == nil {
			return nil, fmt.Errorf("Could not find discussion for participant")
		}

		hashAsInt64 := util.GenerateParticipantSeed(obj.Discussion.ID, obj.ID, obj.Discussion.ShuffleCount)
		gradientColor = util.GenerateGradient(hashAsInt64)
	}
	return &gradientColor, nil
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
	if !obj.IsAnonymous {
		return nil, nil
	}

	if obj.Discussion == nil && obj.DiscussionID != nil {
		disc, err := r.resolveDiscussionByID(ctx, *obj.DiscussionID)
		if err != nil {
			return nil, err
		}
		obj.Discussion = disc
	}

	if obj.Discussion == nil {
		return nil, fmt.Errorf("Could not find associated discussion")
	}

	hashAsInt64 := util.GenerateParticipantSeed(obj.Discussion.ID, obj.ID, obj.Discussion.ShuffleCount)
	fullDisplayName := util.GenerateFullDisplayName(hashAsInt64)
	return &fullDisplayName, nil
}

func (r *participantResolver) MutedForSeconds(ctx context.Context, obj *model.Participant) (*int, error) {
	var result int

	result = 0
	if obj.MutedUntil != nil {
		result = int(obj.MutedUntil.Sub(time.Now()).Seconds())
	}
	return &result, nil
}

// Participant returns generated.ParticipantResolver implementation.
func (r *Resolver) Participant() generated.ParticipantResolver { return &participantResolver{r} }

type participantResolver struct{ *Resolver }

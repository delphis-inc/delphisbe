package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"crypto/sha1"
	"encoding/binary"
	"fmt"

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

		hashAsInt64 := generateParticipantSeed(*obj.Discussion, *obj)
		gradientColor = util.GenerateGradient(hashAsInt64)
	}
	return &gradientColor, nil
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

	hashAsInt64 := generateParticipantSeed(*obj.Discussion, *obj)
	fullDisplayName := util.GenerateFullDisplayName(hashAsInt64)
	return &fullDisplayName, nil
}

// Participant returns generated.ParticipantResolver implementation.
func (r *Resolver) Participant() generated.ParticipantResolver { return &participantResolver{r} }

type participantResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func generateParticipantSeed(discussion model.Discussion, participant model.Participant) uint64 {
	// We generate the display name by SHA-1(discussion_id, participant.id, shuffle_id) without
	// commas, just concatenated.
	h := sha1.Sum([]byte(fmt.Sprintf("%s%s%d", discussion.ID, participant.ID, discussion.ShuffleID)))
	return binary.BigEndian.Uint64(h[:])
}

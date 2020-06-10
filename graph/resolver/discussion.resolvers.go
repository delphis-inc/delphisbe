package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/auth"
)

func (r *discussionResolver) Moderator(ctx context.Context, obj *model.Discussion) (*model.Moderator, error) {
	if obj.Moderator == nil && obj.ModeratorID != nil {
		moderator, err := r.DAOManager.GetModeratorByID(ctx, *obj.ModeratorID)
		if err != nil {
			return nil, err
		}
		obj.Moderator = moderator
	}
	return obj.Moderator, nil
}

func (r *discussionResolver) Posts(ctx context.Context, obj *model.Discussion) ([]*model.Post, error) {
	posts, err := r.DAOManager.GetPostsByDiscussionID(ctx, obj.ID)

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *discussionResolver) Participants(ctx context.Context, obj *model.Discussion) ([]*model.Participant, error) {
	if obj.Participants == nil {
		participants, err := r.DAOManager.GetParticipantsByDiscussionID(ctx, obj.ID)

		if err != nil {
			return nil, err
		}

		particPointers := make([]*model.Participant, 0)
		for i, elem := range participants {
			elem.Discussion = obj
			particPointers = append(particPointers, &participants[i])
		}

		obj.Participants = particPointers
	}
	return obj.Participants, nil
}

func (r *discussionResolver) CreatedAt(ctx context.Context, obj *model.Discussion) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

func (r *discussionResolver) UpdatedAt(ctx context.Context, obj *model.Discussion) (string, error) {
	return obj.UpdatedAt.Format(time.RFC3339), nil
}

func (r *discussionResolver) MeParticipant(ctx context.Context, obj *model.Discussion) (*model.Participant, error) {
	currentUser := auth.GetAuthedUser(ctx)
	if currentUser == nil {
		// Only works for logged in. Won't throw an error here though.
		return nil, nil
	}
	// TODO We should return your most recent participant by post created
	participants, err := r.MeAvailableParticipants(ctx, obj)
	if err != nil {
		return nil, err
	}

	if len(participants) == 1 {
		return participants[0], nil
	} else if len(participants) == 2 {
		if participants[0].UpdatedAt.After(participants[1].UpdatedAt) {
			return participants[0], nil
		} else {
			return participants[1], nil
		}
	}

	// Not a participant -- Not returning an error yet.
	return nil, nil
}

func (r *discussionResolver) MeAvailableParticipants(ctx context.Context, obj *model.Discussion) ([]*model.Participant, error) {
	currentUser := auth.GetAuthedUser(ctx)
	if currentUser == nil {
		return nil, nil
	}

	participantResponse, err := r.DAOManager.GetParticipantsByDiscussionIDUserID(ctx, obj.ID, currentUser.UserID)
	if err != nil {
		return nil, err
	}

	participantArr := make([]*model.Participant, 0)
	if participantArr != nil {
		if participantResponse.NonAnon != nil {
			participantArr = append(participantArr, participantResponse.NonAnon)
		}
		if participantResponse.Anon != nil {
			participantArr = append(participantArr, participantResponse.Anon)
		}
	}

	return participantArr, nil
}

func (r *discussionResolver) Tags(ctx context.Context, obj *model.Discussion) ([]*model.Tag, error) {
	return r.DAOManager.GetDiscussionTags(ctx, obj.ID)
}

func (r *discussionResolver) UpcomingContent(ctx context.Context, obj *model.Discussion) ([]*model.ImportedContent, error) {
	// TODO: Determine UX. Should this be merged with the data that has already been scheduled?
	// At one point is data too stale and should no longer be shown?
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	// Only allow the mod to view possible imported content
	modCheck, err := r.DAOManager.CheckIfModeratorForDiscussion(ctx, authedUser.UserID, obj.ID)
	if err != nil || !modCheck {
		return nil, fmt.Errorf("unauthorized")
	}

	return r.DAOManager.GetUpcomingImportedContentByDiscussionID(ctx, obj.ID)
}

func (r *tagResolver) CreatedAt(ctx context.Context, obj *model.Tag) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

func (r *tagResolver) IsDeleted(ctx context.Context, obj *model.Tag) (bool, error) {
	return obj.DeletedAt != nil, nil
}

// Discussion returns generated.DiscussionResolver implementation.
func (r *Resolver) Discussion() generated.DiscussionResolver { return &discussionResolver{r} }

// Tag returns generated.TagResolver implementation.
func (r *Resolver) Tag() generated.TagResolver { return &tagResolver{r} }

type discussionResolver struct{ *Resolver }
type tagResolver struct{ *Resolver }

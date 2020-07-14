package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/delphis-inc/delphisbe/graph/generated"
	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/backend"
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
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	posts, err := r.DAOManager.GetPostsByDiscussionID(ctx, authedUser.UserID, obj.ID)

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *discussionResolver) PostsConnection(ctx context.Context, obj *model.Discussion, after *string) (*model.PostsConnection, error) {
	/* Hardcode the number of posts per page. This can be changed to be settable by the client in the query. */
	limit := backend.PostPerPageLimit

	/* Sanity check. If no "after" parameter is specified, we set it to a time far into the future, for which
	   no post can yet have been created (not even cosidering large clocks drift) */
	if after == nil {
		futureTime := time.Now().AddDate(1, 0, 0).Format(time.RFC3339Nano)
		after = &futureTime
	} else if _, err := time.Parse(time.RFC3339, *after); err != nil {
		return nil, errors.New("The 'After' parameter is badly formatted: " + *after)
	}

	return r.DAOManager.GetPostsConnectionByDiscussionID(ctx, obj.ID, *after, limit)
}

func (r *discussionResolver) Participants(ctx context.Context, obj *model.Discussion) ([]*model.Participant, error) {
	if obj.Participants == nil {
		participants, err := r.DAOManager.GetParticipantsByDiscussionID(ctx, obj.ID)

		if err != nil {
			return nil, err
		}

		particPointers := make([]*model.Participant, 0)
		for i, elem := range participants {
			if !elem.IsBanned {
				elem.Discussion = obj
				particPointers = append(particPointers, &participants[i])
			}
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

func (r *discussionResolver) FlairTemplates(ctx context.Context, obj *model.Discussion) ([]*model.FlairTemplate, error) {
	return r.DAOManager.GetDiscussionFlairTemplateAccessByDiscussionID(ctx, obj.ID)
}

func (r *discussionResolver) AccessRequests(ctx context.Context, obj *model.Discussion) ([]*model.DiscussionAccessRequest, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	// Only allow the mod to view access requests
	modCheck, err := r.DAOManager.CheckIfModeratorForDiscussion(ctx, authedUser.UserID, obj.ID)
	if err != nil || !modCheck {
		return nil, fmt.Errorf("unauthorized")
	}

	return r.DAOManager.GetDiscussionAccessRequestsByDiscussionID(ctx, obj.ID)
}

func (r *discussionResolver) DiscussionLinksAccess(ctx context.Context, obj *model.Discussion) (*model.DiscussionLinkAccess, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	// Only allow the mod to view invite links
	modCheck, err := r.DAOManager.CheckIfModeratorForDiscussion(ctx, authedUser.UserID, obj.ID)
	if err != nil || !modCheck {
		return nil, fmt.Errorf("unauthorized")
	}

	return r.DAOManager.GetInviteLinksByDiscussionID(ctx, obj.ID)
}

func (r *discussionAccessRequestResolver) User(ctx context.Context, obj *model.DiscussionAccessRequest) (*model.User, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	if authedUser.UserID != obj.UserID {
		// Only allow the mod to view access requests
		modCheck, err := r.DAOManager.CheckIfModeratorForDiscussion(ctx, authedUser.UserID, obj.ID)
		if err != nil || !modCheck {
			return nil, fmt.Errorf("unauthorized")
		}
	}

	return r.DAOManager.GetUserByID(ctx, obj.UserID)
}

func (r *discussionAccessRequestResolver) Discussion(ctx context.Context, obj *model.DiscussionAccessRequest) (*model.Discussion, error) {
	return r.DAOManager.GetDiscussionByID(ctx, obj.DiscussionID)
}

func (r *discussionFlairTemplateAccessResolver) ID(ctx context.Context, obj *model.DiscussionFlairTemplateAccess) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *discussionFlairTemplateAccessResolver) Discussion(ctx context.Context, obj *model.DiscussionFlairTemplateAccess) (*model.Discussion, error) {
	return r.DAOManager.GetDiscussionByID(ctx, obj.DiscussionID)
}

func (r *discussionFlairTemplateAccessResolver) FlairTemplate(ctx context.Context, obj *model.DiscussionFlairTemplateAccess) (*model.FlairTemplate, error) {
	return r.DAOManager.GetFlairTemplateByID(ctx, obj.FlairTemplateID)
}

func (r *discussionFlairTemplateAccessResolver) CreatedAt(ctx context.Context, obj *model.DiscussionFlairTemplateAccess) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

func (r *discussionFlairTemplateAccessResolver) UpdatedAt(ctx context.Context, obj *model.DiscussionFlairTemplateAccess) (string, error) {
	return obj.UpdatedAt.Format(time.RFC3339), nil
}

func (r *discussionFlairTemplateAccessResolver) IsDeleted(ctx context.Context, obj *model.DiscussionFlairTemplateAccess) (bool, error) {
	return obj.DeletedAt != nil, nil
}

func (r *discussionInviteResolver) Discussion(ctx context.Context, obj *model.DiscussionInvite) (*model.Discussion, error) {
	return r.DAOManager.GetDiscussionByID(ctx, obj.DiscussionID)
}

func (r *discussionInviteResolver) InvitingParticipant(ctx context.Context, obj *model.DiscussionInvite) (*model.Participant, error) {
	return r.DAOManager.GetParticipantByID(ctx, obj.InvitingParticipantID)
}

func (r *discussionLinkAccessResolver) InviteLinkURL(ctx context.Context, obj *model.DiscussionLinkAccess) (string, error) {
	return strings.Join([]string{model.InviteLinkHostname, obj.DiscussionID, obj.InviteLinkSlug}, "/"), nil
}

func (r *discussionLinkAccessResolver) VipInviteLinkURL(ctx context.Context, obj *model.DiscussionLinkAccess) (string, error) {
	return strings.Join([]string{model.InviteLinkHostname, obj.DiscussionID, obj.VipInviteLinkSlug}, "/"), nil
}

func (r *tagResolver) CreatedAt(ctx context.Context, obj *model.Tag) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

func (r *tagResolver) IsDeleted(ctx context.Context, obj *model.Tag) (bool, error) {
	return obj.DeletedAt != nil, nil
}

// Discussion returns generated.DiscussionResolver implementation.
func (r *Resolver) Discussion() generated.DiscussionResolver { return &discussionResolver{r} }

// DiscussionAccessRequest returns generated.DiscussionAccessRequestResolver implementation.
func (r *Resolver) DiscussionAccessRequest() generated.DiscussionAccessRequestResolver {
	return &discussionAccessRequestResolver{r}
}

// DiscussionFlairTemplateAccess returns generated.DiscussionFlairTemplateAccessResolver implementation.
func (r *Resolver) DiscussionFlairTemplateAccess() generated.DiscussionFlairTemplateAccessResolver {
	return &discussionFlairTemplateAccessResolver{r}
}

// DiscussionInvite returns generated.DiscussionInviteResolver implementation.
func (r *Resolver) DiscussionInvite() generated.DiscussionInviteResolver {
	return &discussionInviteResolver{r}
}

// DiscussionLinkAccess returns generated.DiscussionLinkAccessResolver implementation.
func (r *Resolver) DiscussionLinkAccess() generated.DiscussionLinkAccessResolver {
	return &discussionLinkAccessResolver{r}
}

// Tag returns generated.TagResolver implementation.
func (r *Resolver) Tag() generated.TagResolver { return &tagResolver{r} }

type discussionResolver struct{ *Resolver }
type discussionAccessRequestResolver struct{ *Resolver }
type discussionFlairTemplateAccessResolver struct{ *Resolver }
type discussionInviteResolver struct{ *Resolver }
type discussionLinkAccessResolver struct{ *Resolver }
type tagResolver struct{ *Resolver }

package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"math"
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

	posts, err := r.DAOManager.GetPostsByDiscussionID(ctx, obj.ID)

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

func (r *discussionResolver) TitleHistory(ctx context.Context, obj *model.Discussion) ([]*model.HistoricalString, error) {
	return obj.TitleHistoryAsObject()
}

func (r *discussionResolver) DescriptionHistory(ctx context.Context, obj *model.Discussion) ([]*model.HistoricalString, error) {
	return obj.DescriptionHistoryAsObject()
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

func (r *discussionResolver) MeCanJoinDiscussion(ctx context.Context, obj *model.Discussion) (*model.CanJoinDiscussionResponse, error) {
	// Things to check here:
	// 1. Is the user logged out? => no (DENIED)
	// 2. Is the user only logged in with Apple? => no (DENIED)
	// 3. Is the user already part of the discussion? => yes (AWAITING_APPROVAL)
	// 4. Is the discussion set to manual approval? AND has the user already requested to join?
	//   4a. If yes then => yes (with awaitingApproval <- true, requiresApproval <- true) (AWAITING_APPROVAL)
	//   4b. If no then => yes (with awaitingApproval <- false, requiresApproval <- true) (APPROVAL_REQUIRED)
	//   4c. If yes AND the user has been rejected => no (with awaitingApproval <- false, requiresApproval <- true) (DENIED)
	//   4d. If yes AND the user has been approved => yes (with awaitingApproval <- false, requiresApproval <- false) (APPROVED_NOT_JOINED)
	// 5. Is the discussion set to automatic approval if following on Twitter? AND is the user NOT followed by moderator?
	//   <SEE #4>
	// 6. Is the discussion set to automatic approval if following on Twitter? AND is the user followed by the moderator? => yes (APPROVED_NOT_JOINED)
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return &model.CanJoinDiscussionResponse{
			Response: model.DiscussionJoinabilityResponseDenied,
		}, nil
	}

	if authedUser.User == nil {
		var err error
		authedUser.User, err = r.DAOManager.GetUserByID(ctx, authedUser.UserID)
		if err != nil || authedUser.User == nil {
			return nil, fmt.Errorf("Error fetching user with ID (%s)", authedUser.UserID)
		}
	}

	meParticipant, err := r.MeParticipant(ctx, obj)
	if err != nil {
		return nil, fmt.Errorf("Error fetching user information")
	}

	return r.DAOManager.GetDiscussionJoinabilityForUser(ctx, authedUser.User, obj, meParticipant)
}

func (r *discussionResolver) MeViewer(ctx context.Context, obj *model.Discussion) (*model.Viewer, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("need auth")
	}

	// TODO: Do we need to check if the user has access?
	return r.DAOManager.GetViewerForDiscussion(ctx, obj.ID, authedUser.UserID, true)
}

func (r *discussionResolver) MeNotificationSettings(ctx context.Context, obj *model.Discussion) (*model.DiscussionUserNotificationSetting, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, nil
	}

	resp, err := r.DAOManager.GetDiscussionUserAccess(ctx, authedUser.UserID, obj.ID)
	if err != nil {
		return nil, fmt.Errorf("Error fetching user information")
	}
	if resp == nil {
		return nil, nil
	}

	return &resp.NotifSetting, nil
}

func (r *discussionResolver) MeDiscussionStatus(ctx context.Context, obj *model.Discussion) (*model.DiscussionUserAccessState, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, nil
	}

	resp, err := r.DAOManager.GetDiscussionUserAccess(ctx, authedUser.UserID, obj.ID)
	if err != nil {
		return nil, fmt.Errorf("Error fetching user information")
	}
	if resp == nil {
		return nil, nil
	}
	return &resp.State, nil
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

func (r *discussionResolver) DiscussionAccessLink(ctx context.Context, obj *model.Discussion) (*model.DiscussionAccessLink, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	// Only allow the mod to view invite links
	modCheck, err := r.DAOManager.CheckIfModeratorForDiscussion(ctx, authedUser.UserID, obj.ID)
	if err != nil || !modCheck {
		return nil, fmt.Errorf("unauthorized")
	}

	return r.DAOManager.GetAccessLinkByDiscussionID(ctx, obj.ID)
}

func (r *discussionResolver) DiscussionJoinability(ctx context.Context, obj *model.Discussion) (model.DiscussionJoinabilitySetting, error) {
	if string(obj.DiscussionJoinability) == "" {
		return model.DiscussionJoinabilitySettingAllRequireApproval, nil
	}

	return obj.DiscussionJoinability, nil
}

func (r *discussionResolver) SecondsUntilShuffle(ctx context.Context, obj *model.Discussion) (*int, error) {
	nextShuffle, err := r.DAOManager.GetNextDiscussionShuffleTime(ctx, obj.ID)
	if err != nil {
		return nil, err
	}

	if nextShuffle == nil || nextShuffle.ShuffleTime == nil {
		return nil, nil
	}

	seconds := int(math.Max(0, time.Until(*nextShuffle.ShuffleTime).Seconds()))

	return &seconds, nil
}

func (r *discussionResolver) Archive(ctx context.Context, obj *model.Discussion) (*model.DiscussionArchive, error) {
	return r.DAOManager.GetDiscussionArchiveByDiscussionID(ctx, obj.ID)
}

func (r *discussionAccessLinkResolver) Discussion(ctx context.Context, obj *model.DiscussionAccessLink) (*model.Discussion, error) {
	return r.DAOManager.GetDiscussionByID(ctx, obj.DiscussionID)
}

func (r *discussionAccessLinkResolver) URL(ctx context.Context, obj *model.DiscussionAccessLink) (string, error) {
	return strings.Join([]string{model.InviteLinkHostname, obj.LinkSlug}, "/"), nil
}

func (r *discussionAccessLinkResolver) IsDeleted(ctx context.Context, obj *model.DiscussionAccessLink) (bool, error) {
	return obj.DeletedAt != nil, nil
}

func (r *discussionAccessRequestResolver) UserProfile(ctx context.Context, obj *model.DiscussionAccessRequest) (*model.UserProfile, error) {
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth for acccess request")
	}

	if authedUser.UserID != obj.UserID {
		// Only allow the mod to view access requests
		modCheck, err := r.DAOManager.CheckIfModeratorForDiscussion(ctx, authedUser.UserID, obj.DiscussionID)
		if err != nil || !modCheck {
			return nil, fmt.Errorf("unauthorized")
		}
	}
	return r.DAOManager.GetUserProfileByUserID(ctx, obj.UserID)
}

func (r *discussionAccessRequestResolver) Discussion(ctx context.Context, obj *model.DiscussionAccessRequest) (*model.Discussion, error) {
	return r.DAOManager.GetDiscussionByID(ctx, obj.DiscussionID)
}

func (r *discussionArchiveResolver) Archive(ctx context.Context, obj *model.DiscussionArchive) (string, error) {
	return string(obj.Archive.RawMessage), nil
}

func (r *discussionUserAccessResolver) Discussion(ctx context.Context, obj *model.DiscussionUserAccess) (*model.Discussion, error) {
	return r.DAOManager.GetDiscussionByID(ctx, obj.DiscussionID)
}

func (r *discussionUserAccessResolver) User(ctx context.Context, obj *model.DiscussionUserAccess) (*model.User, error) {
	return r.DAOManager.GetUserByID(ctx, obj.UserID)
}

func (r *discussionUserAccessResolver) IsDeleted(ctx context.Context, obj *model.DiscussionUserAccess) (bool, error) {
	return obj.DeletedAt != nil, nil
}

func (r *discussionUserAccessResolver) Request(ctx context.Context, obj *model.DiscussionUserAccess) (*model.DiscussionAccessRequest, error) {
	if obj.RequestID == nil {
		return nil, nil
	}
	return r.DAOManager.GetDiscussionRequestAccessByID(ctx, *obj.RequestID)
}

// Discussion returns generated.DiscussionResolver implementation.
func (r *Resolver) Discussion() generated.DiscussionResolver { return &discussionResolver{r} }

// DiscussionAccessLink returns generated.DiscussionAccessLinkResolver implementation.
func (r *Resolver) DiscussionAccessLink() generated.DiscussionAccessLinkResolver {
	return &discussionAccessLinkResolver{r}
}

// DiscussionAccessRequest returns generated.DiscussionAccessRequestResolver implementation.
func (r *Resolver) DiscussionAccessRequest() generated.DiscussionAccessRequestResolver {
	return &discussionAccessRequestResolver{r}
}

// DiscussionArchive returns generated.DiscussionArchiveResolver implementation.
func (r *Resolver) DiscussionArchive() generated.DiscussionArchiveResolver {
	return &discussionArchiveResolver{r}
}

// DiscussionUserAccess returns generated.DiscussionUserAccessResolver implementation.
func (r *Resolver) DiscussionUserAccess() generated.DiscussionUserAccessResolver {
	return &discussionUserAccessResolver{r}
}

type discussionResolver struct{ *Resolver }
type discussionAccessLinkResolver struct{ *Resolver }
type discussionAccessRequestResolver struct{ *Resolver }
type discussionArchiveResolver struct{ *Resolver }
type discussionUserAccessResolver struct{ *Resolver }

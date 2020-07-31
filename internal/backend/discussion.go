package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"

	"go.uber.org/multierr"

	"github.com/sirupsen/logrus"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/util"
)

const discussionSubscriberKey = "discussion_subscribers-%s"
const discussionEventSubscriberKey = "discussion_event_subscribers-%s"

func (d *delphisBackend) CreateNewDiscussion(ctx context.Context, creatingUser *model.User, anonymityType model.AnonymityType, title string, description string, publicAccess bool, discussionSettings model.DiscussionCreationSettings) (*model.Discussion, error) {
	moderatorObj := model.Moderator{
		ID:            util.UUIDv4(),
		UserProfileID: &creatingUser.UserProfile.ID,
	}
	_, err := d.db.CreateModerator(ctx, moderatorObj)
	if err != nil {
		return nil, err
	}

	discussionID := util.UUIDv4()
	now := time.Now()
	titleHistory := []model.HistoricalString{
		{
			Value:     title,
			CreatedAt: now,
		},
	}
	descriptionHistory := []model.HistoricalString{
		{
			Value:     description,
			CreatedAt: now,
		},
	}
	titleHistoryBytes, err := json.Marshal(titleHistory)
	if err != nil {
		return nil, err
	}
	descriptionHistoryBytes, err := json.Marshal(descriptionHistory)
	if err != nil {
		return nil, err
	}
	discussionObj := model.Discussion{
		CreatedAt:     now,
		UpdatedAt:     now,
		ID:            discussionID,
		AnonymityType: anonymityType,
		Title:         title,
		Description:   description,
		TitleHistory: postgres.Jsonb{
			RawMessage: titleHistoryBytes,
		},
		DescriptionHistory: postgres.Jsonb{
			RawMessage: descriptionHistoryBytes,
		},
		ModeratorID:           &moderatorObj.ID,
		PublicAccess:          publicAccess,
		DiscussionJoinability: discussionSettings.DiscussionJoinability,
	}

	_, err = d.db.UpsertDiscussion(ctx, discussionObj)
	if err != nil {
		return nil, err
	}

	// Create concierge participant
	trueObj := true
	// TODO: We should probably remove the concierge, right?
	if _, err := d.CreateParticipantForDiscussion(ctx, discussionObj.ID, model.ConciergeUser, model.AddDiscussionParticipantInput{HasJoined: &trueObj}); err != nil {
		logrus.WithError(err).Error("failed to create concierge user")
		return nil, err
	}

	// Create invite links for discussion
	if _, err := d.UpsertInviteLinksByDiscussionID(ctx, discussionID); err != nil {
		logrus.WithError(err).Error("failed to create invite links")
		return nil, err
	}

	return &discussionObj, nil
}

func (d *delphisBackend) GetDiscussionJoinabilityForUser(ctx context.Context, userObj *model.User, discussionObj *model.Discussion, meParticipant *model.Participant) (*model.CanJoinDiscussionResponse, error) {
	if userObj == nil || discussionObj == nil || userObj.UserProfile == nil {
		return nil, fmt.Errorf("No user available")
	}
	socialInfos, err := d.GetSocialInfosByUserProfileID(ctx, userObj.UserProfile.ID)
	if err != nil {
		return nil, fmt.Errorf("Error fetching user with ID (%s)", userObj.ID)
	}
	isTwitterAuth := false
	var twitterSocialInfo *model.SocialInfo
	for idx, s := range socialInfos {
		if s.Network == "twitter" {
			isTwitterAuth = true
			twitterSocialInfo = &socialInfos[idx]
			break
		}
	}
	if !isTwitterAuth {
		return &model.CanJoinDiscussionResponse{
			Response: model.DiscussionJoinabilityResponseDenied,
		}, nil
	}

	if meParticipant != nil {
		return &model.CanJoinDiscussionResponse{
			Response: model.DiscussionJoinabilityResponseAlreadyJoined,
		}, nil
	}

	if discussionObj.DiscussionJoinability == model.DiscussionJoinabilitySettingAllowTwitterFriends {
		// Now we need to know if this moderator follows the user on Twitter.
		moderatorSocialInfos, err := d.GetSocialInfosByUserProfileID(ctx, *discussionObj.Moderator.UserProfileID)
		if err != nil {
			return nil, fmt.Errorf("Error fetching moderator information")
		}
		var modSocialInfo *model.SocialInfo
		for idx, s := range moderatorSocialInfos {
			if s.Network == "twitter" {
				isTwitterAuth = true
				modSocialInfo = &moderatorSocialInfos[idx]
				break
			}
		}
		if modSocialInfo == nil {
			return nil, fmt.Errorf("Error fetching moderator information")
		}
		doesModeratorFollow, err := d.DoesTwitterUserFollowUser(ctx, *modSocialInfo, *twitterSocialInfo)
		if err != nil {
			return nil, err
		}

		if doesModeratorFollow {
			return &model.CanJoinDiscussionResponse{
				Response: model.DiscussionJoinabilityResponseApprovedNotJoined,
			}, nil
		} else {
			return d.getJoinabilityFromInviteStatus(ctx, discussionObj, userObj)
		}
	} else {
		return d.getJoinabilityFromInviteStatus(ctx, discussionObj, userObj)
	}
}

func (d *delphisBackend) getJoinabilityFromInviteStatus(ctx context.Context, discussionObj *model.Discussion, userObj *model.User) (*model.CanJoinDiscussionResponse, error) {
	requestAccess, err := d.db.GetDiscussionAccessRequestByDiscussionIDUserID(ctx, discussionObj.ID, userObj.ID)
	if err != nil {
		return nil, err
	}
	if requestAccess == nil {
		// No access request has been made.
		return &model.CanJoinDiscussionResponse{
			Response: model.DiscussionJoinabilityResponseApprovalRequired,
		}, nil
	} else {
		switch requestAccess.Status {
		case model.InviteRequestStatusAccepted:
			return &model.CanJoinDiscussionResponse{
				Response: model.DiscussionJoinabilityResponseApprovedNotJoined,
			}, nil
		case model.InviteRequestStatusPending:
			return &model.CanJoinDiscussionResponse{
				Response: model.DiscussionJoinabilityResponseApprovalRequired,
			}, nil
		case model.InviteRequestStatusRejected:
			return &model.CanJoinDiscussionResponse{
				Response: model.DiscussionJoinabilityResponseDenied,
			}, nil
		case model.InviteRequestStatusCancelled:
			return &model.CanJoinDiscussionResponse{
				Response: model.DiscussionJoinabilityResponseApprovalRequired,
			}, nil
		default:
			return &model.CanJoinDiscussionResponse{
				Response: model.DiscussionJoinabilityResponseApprovalRequired,
			}, nil
		}
	}
}

func (d *delphisBackend) UpdateDiscussion(ctx context.Context, id string, input model.DiscussionInput) (*model.Discussion, error) {
	discObj, err := d.db.GetDiscussionByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("failed to get discussion by ID")
		return nil, err
	}

	updateDiscussionObj(discObj, input)

	return d.db.UpsertDiscussion(ctx, *discObj)
}

func (d *delphisBackend) GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error) {
	return d.db.GetDiscussionByID(ctx, id)
}

func (d *delphisBackend) GetDiscussionsByIDs(ctx context.Context, ids []string) (map[string]*model.Discussion, error) {
	return d.db.GetDiscussionsByIDs(ctx, ids)
}

func (d *delphisBackend) GetDiscussionByModeratorID(ctx context.Context, moderatorID string) (*model.Discussion, error) {
	return d.db.GetDiscussionByModeratorID(ctx, moderatorID)
}

func (d *delphisBackend) GetDiscussionsForAutoPost(ctx context.Context) ([]*model.DiscussionAutoPost, error) {
	iter := d.db.GetDiscussionsAutoPost(ctx)
	return d.db.DiscussionAutoPostIterCollect(ctx, iter)
}

func (d *delphisBackend) ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error) {
	return d.db.ListDiscussions(ctx)
}

func (d *delphisBackend) ListDiscussionsByUserID(ctx context.Context, userID string) (*model.DiscussionsConnection, error) {
	return d.db.ListDiscussionsByUserID(ctx, userID)
}

func (d *delphisBackend) SubscribeToDiscussion(ctx context.Context, subscriberUserID string, postChannel chan *model.Post, discussionID string) error {
	cacheKey := fmt.Sprintf(discussionSubscriberKey, discussionID)
	d.discussionMutex.Lock()
	defer d.discussionMutex.Unlock()
	currentSubsIface, found := d.cache.Get(cacheKey)
	if !found {
		currentSubsIface = map[string]chan *model.Post{}
	}
	var currentSubs map[string]chan *model.Post
	var ok bool
	if currentSubs, ok = currentSubsIface.(map[string]chan *model.Post); !ok {
		currentSubs = map[string]chan *model.Post{}
	}
	currentSubs[subscriberUserID] = postChannel
	d.cache.Set(cacheKey, currentSubs, time.Hour)
	return nil
}

func (d *delphisBackend) UnSubscribeFromDiscussion(ctx context.Context, subscriberUserID string, discussionID string) error {
	cacheKey := fmt.Sprintf(discussionSubscriberKey, discussionID)
	d.discussionMutex.Lock()
	defer d.discussionMutex.Unlock()
	currentSubsIface, found := d.cache.Get(cacheKey)
	if !found {
		return nil
	}
	var currentSubs map[string]chan *model.Post
	var ok bool
	if currentSubs, ok = currentSubsIface.(map[string]chan *model.Post); !ok {
		currentSubs = map[string]chan *model.Post{}
	}
	delete(currentSubs, subscriberUserID)
	d.cache.Set(cacheKey, currentSubs, time.Hour)
	return nil
}

func (d *delphisBackend) SubscribeToDiscussionEvent(ctx context.Context, subscriberUserID string, eventChannel chan *model.DiscussionSubscriptionEvent, discussionID string) error {
	cacheKey := fmt.Sprintf(discussionEventSubscriberKey, discussionID)
	d.discussionMutex.Lock()
	defer d.discussionMutex.Unlock()
	currentSubsIface, found := d.cache.Get(cacheKey)
	if !found {
		currentSubsIface = map[string]chan *model.DiscussionSubscriptionEvent{}
	}
	var currentSubs map[string]chan *model.DiscussionSubscriptionEvent
	var ok bool
	if currentSubs, ok = currentSubsIface.(map[string]chan *model.DiscussionSubscriptionEvent); !ok {
		currentSubs = map[string]chan *model.DiscussionSubscriptionEvent{}
	}
	currentSubs[subscriberUserID] = eventChannel
	d.cache.Set(cacheKey, currentSubs, time.Hour)
	return nil
}

func (d *delphisBackend) UnSubscribeFromDiscussionEvent(ctx context.Context, subscriberUserID string, discussionID string) error {
	cacheKey := fmt.Sprintf(discussionEventSubscriberKey, discussionID)
	d.discussionMutex.Lock()
	defer d.discussionMutex.Unlock()
	currentSubsIface, found := d.cache.Get(cacheKey)
	if !found {
		return nil
	}
	var currentSubs map[string]chan *model.DiscussionSubscriptionEvent
	var ok bool
	if currentSubs, ok = currentSubsIface.(map[string]chan *model.DiscussionSubscriptionEvent); !ok {
		currentSubs = map[string]chan *model.DiscussionSubscriptionEvent{}
	}
	delete(currentSubs, subscriberUserID)
	d.cache.Set(cacheKey, currentSubs, time.Hour)
	return nil
}

func (d *delphisBackend) GetDiscussionTags(ctx context.Context, id string) ([]*model.Tag, error) {
	iter := d.db.GetDiscussionTags(ctx, id)
	return d.db.TagIterCollect(ctx, iter)
}

func (d *delphisBackend) PutDiscussionTags(ctx context.Context, discussionID string, tags []string) ([]*model.Tag, error) {
	if len(tags) == 0 {
		return nil, fmt.Errorf("no tags to add")
	}

	var addedTags []*model.Tag

	// Begin tx
	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return nil, err
	}

	tagObj := model.Tag{
		ID: discussionID,
	}
	for _, tag := range tags {
		tagObj.Tag = tag
		tagResp, err := d.db.PutDiscussionTags(ctx, tx, tagObj)
		if err != nil {
			logrus.WithError(err).Error("failed to PutDiscussionTags")

			// Rollback on errors
			if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
				logrus.WithError(txErr).Error("failed to rollback tx")
				return nil, multierr.Append(err, txErr)
			}
			return nil, err
		}
		addedTags = append(addedTags, tagResp)
	}

	// Commit transaction
	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("failed to commit post tx")
		return nil, err
	}

	return addedTags, nil
}

func (d *delphisBackend) DeleteDiscussionTags(ctx context.Context, discussionID string, tags []string) ([]*model.Tag, error) {
	if len(tags) == 0 {
		return nil, fmt.Errorf("no tags to delete")
	}

	var deletedTags []*model.Tag

	// Begin tx
	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return nil, err
	}

	tagObj := model.Tag{
		ID: discussionID,
	}
	for _, tag := range tags {
		tagObj.Tag = tag
		tagResp, err := d.db.DeleteDiscussionTags(ctx, tx, tagObj)
		if err != nil {
			logrus.WithError(err).Error("failed to PutDiscussionTags")

			// Rollback on errors
			if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
				logrus.WithError(txErr).Error("failed to rollback tx")
				return nil, multierr.Append(err, txErr)
			}
			return nil, err
		}
		deletedTags = append(deletedTags, tagResp)
	}

	// Commit transaction
	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("failed to commit post tx")
		return nil, err
	}

	return deletedTags, nil
}

func updateDiscussionObj(disc *model.Discussion, input model.DiscussionInput) {
	if input.AnonymityType != nil {
		disc.AnonymityType = *input.AnonymityType
	}
	if input.Title != nil {
		disc.AddTitleToHistory(*input.Title)
		disc.Title = *input.Title
	}
	if input.Description != nil {
		disc.AddDescriptionToHistory(*input.Description)
		disc.Description = *input.Description
	}
	if input.DiscussionJoinability != nil {
		disc.DiscussionJoinability = *input.DiscussionJoinability
	}
	if input.AutoPost != nil {
		disc.AutoPost = *input.AutoPost
	}
	if input.IdleMinutes != nil {
		disc.IdleMinutes = *input.IdleMinutes
	}
	if input.PublicAccess != nil {
		disc.PublicAccess = *input.PublicAccess
	}
	if input.IconURL != nil {
		disc.IconURL = input.IconURL
	}
}

func dedupeDiscussions(discussions []*model.Discussion) []*model.Discussion {
	hashMap := make(map[string]int)

	var results []*model.Discussion
	for _, val := range discussions {
		if _, ok := hashMap[val.ID]; !ok {
			results = append(results, val)
		}
		hashMap[val.ID]++
	}
	return results
}

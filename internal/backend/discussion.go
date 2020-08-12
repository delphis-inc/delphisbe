package backend

import (
	"context"
	"database/sql"
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
		DiscussionJoinability: discussionSettings.DiscussionJoinability,
		LockStatus:            false,
	}

	_, err = d.db.UpsertDiscussion(ctx, discussionObj)
	if err != nil {
		return nil, err
	}

	if err := d.grantAccessAndCreateParticipants(ctx, discussionObj.ID, creatingUser.ID); err != nil {
		logrus.WithError(err).Error("failed to grant access and create participants for new discussion")
		return nil, err
	}

	// Create access links for discussion
	if _, err := d.PutAccessLinkForDiscussion(ctx, discussionID); err != nil {
		logrus.WithError(err).Error("failed to create invite links")
		return nil, err
	}

	return &discussionObj, nil
}

func (d *delphisBackend) IncrementDiscussionShuffleCount(ctx context.Context, tx *sql.Tx, id string) (*int, error) {
	isPassedTx := tx != nil
	if tx == nil {
		var err error
		tx, err = d.db.BeginTx(ctx)
		if err != nil {
			logrus.WithError(err).Error("failed to begin tx")
			return nil, err
		}
	}

	newShuffleCount, err := d.db.IncrementDiscussionShuffleCount(ctx, tx, id)
	if err != nil {
		if !isPassedTx {
			txErr := d.rollbackTx(ctx, tx)
			if txErr != nil {
				return nil, multierr.Append(err, txErr)
			}
		}
		return nil, err
	} else if !isPassedTx {
		err := d.db.CommitTx(ctx, tx)
		if err != nil {
			txErr := d.rollbackTx(ctx, tx)
			if txErr != nil {
				return nil, multierr.Append(err, txErr)
			}
			return nil, err
		}
	}
	return newShuffleCount, nil
}

func (d *delphisBackend) GetDiscussionJoinabilityForUser(ctx context.Context, userObj *model.User, discussionObj *model.Discussion, meParticipant *model.Participant) (*model.CanJoinDiscussionResponse, error) {
	if userObj == nil || discussionObj == nil || userObj.UserProfile == nil {
		return nil, fmt.Errorf("No user available")
	}

	discussionUserAccess, err := d.db.GetDiscussionUserAccess(ctx, discussionObj.ID, userObj.ID)
	if err != nil {
		return nil, err
	}
	if discussionUserAccess != nil && discussionUserAccess.DeletedAt == nil {
		response := model.DiscussionJoinabilityResponseApprovedNotJoined
		if meParticipant != nil {
			response = model.DiscussionJoinabilityResponseAlreadyJoined
		}
		return &model.CanJoinDiscussionResponse{
			Response: response,
		}, nil
	}

	socialInfos, err := d.GetSocialInfosByUserProfileID(ctx, userObj.UserProfile.ID)
	if err != nil {
		return nil, fmt.Errorf("Error fetching user with ID (%s)", userObj.ID)
	}
	isTwitterAuth := false
	var twitterSocialInfo *model.SocialInfo
	for idx, s := range socialInfos {
		if s.Network == util.SocialNetworkTwitter {
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
			if s.Network == util.SocialNetworkTwitter {
				isTwitterAuth = true
				modSocialInfo = &moderatorSocialInfos[idx]
				break
			}
		}
		if modSocialInfo == nil {
			return nil, fmt.Errorf("Error fetching moderator information")
		}
		twitterClient, err := d.GetTwitterClientWithAccessTokens(ctx, twitterSocialInfo.AccessToken, twitterSocialInfo.AccessTokenSecret)
		if err != nil {
			return nil, err
		}
		doesModeratorFollow, err := d.DoesTwitterUserFollowUser(ctx, twitterClient, *modSocialInfo, *twitterSocialInfo)
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

func (d *delphisBackend) GetDiscussionByLinkSlug(ctx context.Context, slug string) (*model.Discussion, error) {
	return d.db.GetDiscussionByLinkSlug(ctx, slug)
}

func (d *delphisBackend) UpdateDiscussion(ctx context.Context, id string, input model.DiscussionInput) (*model.Discussion, error) {
	discObj, err := d.db.GetDiscussionByID(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("failed to get discussion by ID")
		return nil, err
	}

	if discObj == nil {
		logrus.Infof("No discussion to update")
		return nil, nil
	}

	updateDiscussionObj(discObj, input)

	if input.LockStatus != nil && *input.LockStatus == true {
		if _, err := d.CreateDiscussionArchive(ctx, id, discObj.ShuffleCount); err != nil {
			logrus.WithError(err).Error("failed to make discussion archive")
			return nil, err
		}

	}

	return d.db.UpsertDiscussion(ctx, *discObj)
}

// May need to add paging
func (d *delphisBackend) GetDiscussionArchiveByDiscussionID(ctx context.Context, discussionID string) (*model.DiscussionArchive, error) {
	return d.db.GetDiscussionArchiveByDiscussionID(ctx, discussionID)
}

func (d *delphisBackend) CreateDiscussionArchive(ctx context.Context, discussionID string, shuffleCount int) (*model.DiscussionArchive, error) {
	// Get discssion's posts
	posts, err := d.GetPostsByDiscussionID(ctx, discussionID)
	if err != nil {
		logrus.WithError(err).Error("failed to get posts by discussionID")
		return nil, err
	}

	archivedPosts, err := d.anonymizePostsForArchive(ctx, posts, shuffleCount)
	if err != nil {
		logrus.WithError(err).Error("failed to anonymize posts for archive")
		return nil, err
	}

	postsBytes, err := json.Marshal(archivedPosts)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal posts")
		return nil, err
	}

	// Update archive table within transaction
	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to being tx")
		return nil, err
	}

	// Create archive object
	discArchive := model.DiscussionArchive{
		DiscussionID: discussionID,
		Archive:      postgres.Jsonb{postsBytes},
	}

	archiveObj, err := d.db.UpsertDiscussionArchive(ctx, tx, discArchive)
	if err != nil {
		logrus.WithError(err).Error("failed to upsert discussion archive")
		// Rollback on errors
		if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
			logrus.WithError(txErr).Error("failed to rollback tx")
			return nil, multierr.Append(err, txErr)
		}

		return nil, err
	}

	if _, err := d.IncrementDiscussionShuffleCount(ctx, tx, discussionID); err != nil {
		logrus.WithError(err).Error("failed to increment shuffle count")
		// Rollback on errors
		if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
			logrus.WithError(txErr).Error("failed to rollback tx")
			return nil, multierr.Append(err, txErr)
		}

		return nil, err
	}

	// Commit transaction
	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("failed to commit post tx")
		return nil, err
	}

	return archiveObj, nil
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

func (d *delphisBackend) ListDiscussionsByUserID(ctx context.Context, userID string, state model.DiscussionUserAccessState) (*model.DiscussionsConnection, error) {
	return d.db.ListDiscussionsByUserID(ctx, userID, state)
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

func (d *delphisBackend) grantAccessAndCreateParticipants(ctx context.Context, discussionID, userID string) error {
	userIDs := []string{model.ConciergeUser, userID}

	trueObj := true
	for _, id := range userIDs {
		state := model.DiscussionUserAccessStateActive
		notifSettings := model.DiscussionUserNotificationSettingEverything
		settings := model.DiscussionUserSettings{
			State:        &state,
			NotifSetting: &notifSettings,
		}

		// Grant access
		if _, err := d.UpsertUserDiscussionAccess(ctx, id, discussionID, settings); err != nil {
			logrus.WithError(err).Error("failed to grant access")
			return err
		}

		// Create participant
		if _, err := d.CreateParticipantForDiscussion(ctx, discussionID, id, model.AddDiscussionParticipantInput{HasJoined: &trueObj}); err != nil {
			logrus.WithError(err).Error("failed to create concierge participant")
			return err
		}
	}

	return nil
}

func (d *delphisBackend) anonymizePostsForArchive(ctx context.Context, posts []*model.Post, shuffleCount int) ([]*model.ArchivedPost, error) {
	archivedPosts := make([]*model.ArchivedPost, 0)

	for _, post := range posts {
		// Skip deleted posts for archive
		if post.DeletedAt != nil {
			continue
		}

		// Generate anonymous participant name
		if post.DiscussionID == nil || post.ParticipantID == nil {
			logrus.Warn("discussion, participant should not be null")
			continue
		}

		hashAsInt64 := util.GenerateParticipantSeed(*post.DiscussionID, *post.ParticipantID, shuffleCount)
		fullDisplayName := util.GenerateFullDisplayName(hashAsInt64)

		entities := make([]string, 0)
		if post.PostContent.MentionedEntities != nil {
			for _, entityID := range post.PostContent.MentionedEntities {
				entity, err := util.ReturnParsedEntityID(entityID)
				if err != nil {
					logrus.WithError(err).Error("failed to parse entity")
					return nil, err
				}
				if entity.Type == model.ParticipantPrefix {
					hashAsInt64 := util.GenerateParticipantSeed(*post.DiscussionID, entity.ID, shuffleCount)
					entities = append(entities, util.GenerateFullDisplayName(hashAsInt64))
				} else if entity.Type == model.DiscussionPrefix {
					entities = append(entities, "redacted_discussion")
				}
			}
		}

		archivePost := model.ArchivedPost{
			PostType:          post.PostType,
			CreatedAt:         post.CreatedAt,
			ParticipantName:   fullDisplayName,
			Content:           post.PostContent.Content,
			MentionedEntities: entities,
			MediaID:           post.MediaID,
		}

		archivedPosts = append(archivedPosts, &archivePost)
	}

	return archivedPosts, nil
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
	if input.IconURL != nil {
		disc.IconURL = input.IconURL
	}
	if input.LastPostID != nil {
		disc.LastPostID = input.LastPostID
	}
	if input.LastPostCreatedAt != nil {
		disc.LastPostCreatedAt = input.LastPostCreatedAt
	}
	if input.LockStatus != nil {
		disc.LockStatus = *input.LockStatus
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

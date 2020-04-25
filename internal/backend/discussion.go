package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/util"
)

const discussionSubscriberKey = "discussion_subscribers-%s"

func (d *delphisBackend) CreateNewDiscussion(ctx context.Context, creatingUser *model.User, anonymityType model.AnonymityType, title string) (*model.Discussion, error) {
	moderatorObj := model.Moderator{
		ID:            util.UUIDv4(),
		UserProfileID: &creatingUser.UserProfile.ID,
	}
	_, err := d.db.CreateModerator(ctx, moderatorObj)
	if err != nil {
		return nil, err
	}

	discussionID := util.UUIDv4()
	discussionObj := model.Discussion{
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ID:            discussionID,
		AnonymityType: anonymityType,
		Title:         title,
		ModeratorID:   &moderatorObj.ID,
	}

	_, err = d.db.UpsertDiscussion(ctx, discussionObj)

	if err != nil {
		return nil, err
	}

	return &discussionObj, nil
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

func (d *delphisBackend) ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error) {
	return d.db.ListDiscussions(ctx)
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

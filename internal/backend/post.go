package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/util"
	"github.com/sirupsen/logrus"
)

func (d *delphisBackend) CreatePost(ctx context.Context, discussionID string, participantID string, input model.PostContentInput) (*model.Post, error) {
	postContent := model.PostContent{
		ID:      util.UUIDv4(),
		Content: input.PostText,
	}

	post := model.Post{
		ID:            util.UUIDv4(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		PostContentID: &postContent.ID,
		PostContent:   &postContent,
	}

	// Weird logic to have quoted_post_id actually set to nil as opposed to a blank space
	if input.QuotedPostID != nil {
		post.QuotedPostID = input.QuotedPostID
	}

	postObj, err := d.db.PutPost(ctx, post)

	if err != nil {
		return nil, err
	}

	discussion, err := d.db.GetDiscussionByID(ctx, discussionID)
	if err != nil {
		logrus.WithError(err).Debugf("Skipping notification to subscribers because of an error")
	} else {
		_, err := d.SendNotificationsToSubscribers(ctx, discussion, &post)
		if err != nil {
			logrus.WithError(err).Warn("Failed to send push notifications on createPost")
		}
	}

	return postObj, nil
}

func (d *delphisBackend) NotifySubscribersOfCreatedPost(ctx context.Context, post *model.Post, discussionID string) error {
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
	for userID, channel := range currentSubs {
		if channel != nil {
			select {
			case channel <- post:
				logrus.Debugf("Sent message to channel for user ID: %s", userID)
			default:
				logrus.Debugf("No message was sent. Unsubscribing the user")
				delete(currentSubs, userID)
			}
		}
	}
	d.cache.Set(cacheKey, currentSubs, time.Hour)
	return nil
}

func (d *delphisBackend) GetPostsByDiscussionID(ctx context.Context, discussionID string) ([]*model.Post, error) {
	return d.db.GetPostsByDiscussionID(ctx, discussionID)
}

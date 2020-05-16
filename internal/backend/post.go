package backend

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/nedrocks/delphisbe/internal/datastore"

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
		QuotedPostID:  input.QuotedPostID,
	}

	// TODO: Wrap into a transaction here
	// Put post contents and post
	if err := d.db.PutPostContent(ctx, postContent); err != nil {
		logrus.WithError(err).Error("failed to PutPostContent")
		return nil, err
	}

	postObj, err := d.db.PutPost(ctx, post)
	if err != nil {
		logrus.WithError(err).Error("failed to PutPost")
		return nil, err
	}
	// End transaction

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
	iter := d.db.GetPostsByDiscussionIDIter(ctx, discussionID)
	return d.iterToPosts(ctx, iter)
	//return d.db.GetPostsByDiscussionID(ctx, discussionID)
}

// Testing function to keep functionality
func (d *delphisBackend) iterToPosts(ctx context.Context, iter datastore.PostIter) ([]*model.Post, error) {
	var posts []*model.Post
	post := model.Post{}

	defer iter.Close()

	for iter.Next(&post) {
		tempPost := post

		// Check if there is a quotedPostID. Fetch if so
		if tempPost.QuotedPostID != nil {
			var err error
			// TODO: potentially optimize into joins
			tempPost.QuotedPost, err = d.db.GetPostByID(ctx, *tempPost.QuotedPostID)
			if err != nil {
				// Do we want to fail the whole discussion if we can't get a quote?
				return nil, err
			}
		}

		posts = append(posts, &tempPost)
	}

	if err := iter.Close(); err != nil && err != io.EOF {
		logrus.WithError(err).Error("failed to close iter")
		return nil, err
	}

	return posts, nil
}

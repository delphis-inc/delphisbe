package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"

	"github.com/sirupsen/logrus"
)

// Can make concurrent and scale with the number of discussions
func (d *delphisBackend) AutoPostContent() {
	ctx := context.Background()
	// Get discussions that have autopost turned on
	discs, err := d.GetDiscussionsForAutoPost(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to get discussions for autopost")
		return
	}

	// Iterate over discussions that are set to autopost
	for _, disc := range discs {
		shouldPost, err := d.checkIdleTime(ctx, disc.ID, disc.IdleMinutes)
		if err != nil {
			logrus.WithError(err).Error("failed to check idle time")
			return
		}

		// If the discussion has been idle for long enough, post!
		if shouldPost {
			if err := d.postNextContent(ctx, disc.ID); err != nil {
				logrus.WithError(err).Error("failed to post next content")
				return
			}
		}
	}
}

func (d *delphisBackend) checkIdleTime(ctx context.Context, discussionID string, minutes int) (bool, error) {
	post, err := d.GetLastPostByDiscussionID(ctx, discussionID, minutes)
	if err != nil {
		logrus.WithError(err).Error("failed to get last post by discussion ID")
		return false, err
	}
	return post == nil, nil
}

func (d *delphisBackend) postNextContent(ctx context.Context, discussionID string) error {
	dripType := model.ScheduledDrip
	// Fetch next article
	iter := d.db.GetScheduledImportedContentByDiscussionID(ctx, discussionID)
	contents, err := d.iterToContent(ctx, iter)
	if err != nil {
		logrus.WithError(err).Error("failed to get scheduledContent")
		return err
	}
	if len(contents) == 0 {
		dripType = model.AutoDrip
		iter = d.db.GetImportedContentByDiscussionID(ctx, discussionID, 10)
		contents, err = d.iterToContent(ctx, iter)
		if err != nil {
			logrus.WithError(err).Error("failed to get imported content")
			return err
		}

		if len(contents) == 0 {
			return nil
		}
	}

	// Get concierge participant
	resp, err := d.GetParticipantsByDiscussionIDUserID(ctx, discussionID, model.ConciergeUser)
	if err != nil {
		logrus.WithError(err).Error("failed to fetch concierge participant")
		return err
	}
	if resp.NonAnon == nil {
		return fmt.Errorf("discussion is missing a concierge participant")
	}

	// Call post function
	content := contents[0]
	now := time.Now()
	if _, err := d.PostImportedContent(ctx, resp.NonAnon.ID, discussionID, content.ID, &now, content.Tags, dripType); err != nil {
		logrus.WithError(err).Error("failed to post imported content from autodrip")
		return err
	}
	return nil
}

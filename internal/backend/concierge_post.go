package backend

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/delphis-inc/delphisbe/graph/model"
)

func (d *delphisBackend) GetConciergeParticipantID(ctx context.Context, discussionID string) (string, error) {
	// Get concierge participant for posts
	participants, err := d.GetParticipantsByDiscussionIDUserID(ctx, discussionID, model.ConciergeUser)
	if err != nil {
		logrus.WithError(err).Error("failed to get concierge participant")
		return "", err
	}

	if participants.NonAnon == nil {
		return "", fmt.Errorf("no non-anonymous participant for the concierge")
	}

	return participants.NonAnon.ID, nil
}

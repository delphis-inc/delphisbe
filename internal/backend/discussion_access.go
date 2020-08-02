package backend

import (
	"context"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

func (d *delphisBackend) GetDiscussionAccessByUserID(ctx context.Context, userID string) ([]*model.Discussion, error) {
	// Get discussions the user was invited to
	userDiscIter := d.db.GetDiscussionsByUserAccess(ctx, userID)
	userDiscussions, err := d.db.DiscussionIterCollect(ctx, userDiscIter)
	if err != nil {
		logrus.WithError(err).Error("failed to get user access discussions")
		return nil, err
	}

	dedupedDiscs := dedupeDiscussions(userDiscussions)

	return dedupedDiscs, nil
}

func (d *delphisBackend) GrantUserDiscussionAccess(ctx context.Context, userID string, discussionID string) (*model.DiscussionUserAccess, error) {
	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return nil, err
	}
	access, err := d.db.UpsertDiscussionUserAccess(ctx, tx, discussionID, userID)
	if err != nil {
		// Rollback tx.
		if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
			logrus.WithError(txErr).Error("Failed to rollback tx")
			return nil, multierr.Append(err, txErr)
		}
		return nil, err
	}

	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("Failed to commit tx")
		return nil, err
	}

	return access, nil
}

func (d *delphisBackend) validateFlairTemplatesToAdd(ctx context.Context, userID string, templates []string) ([]string, error) {
	userFlairs, err := d.GetFlairsByUserID(ctx, userID)
	if err != nil {
		logrus.WithError(err).Error("failed to get flairs for user")
		return nil, err
	}
	var validatedTemplates []string

	flairMap := make(map[string]int)

	// Build map out of user's flairs
	for _, val := range userFlairs {
		flairMap[val.TemplateID]++
	}

	// Validate that the passed in flairs are owned by the user
	for _, val := range templates {
		if _, ok := flairMap[val]; ok {
			validatedTemplates = append(validatedTemplates, val)
		}
	}

	return validatedTemplates, nil
}

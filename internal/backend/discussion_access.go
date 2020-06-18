package backend

import (
	"context"
	"fmt"
	"io"

	"github.com/nedrocks/delphisbe/internal/datastore"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

func (d *delphisBackend) GetDiscussionAccessByUserID(ctx context.Context, userID string) ([]*model.Discussion, error) {
	var discussions []*model.Discussion

	// Get public discussions
	publicIter := d.db.GetPublicDiscussions(ctx)
	publicDiscussions, err := d.iterToDiscussions(ctx, publicIter)
	if err != nil {
		logrus.WithError(err).Error("failed to get public discussions")
		return nil, err
	}

	// Get discussions the user has access to by flair
	flairDiscIter := d.db.GetDiscussionsForFlairTemplateByUserID(ctx, userID)
	flairDiscussions, err := d.iterToDiscussions(ctx, flairDiscIter)
	if err != nil {
		logrus.WithError(err).Error("failed to get flair discussions")
		return nil, err
	}

	// Get discussions the user was invited to
	userDiscIter := d.db.GetDiscussionsForUserAccessByUserID(ctx, userID)
	userDiscussions, err := d.iterToDiscussions(ctx, userDiscIter)
	if err != nil {
		logrus.WithError(err).Error("failed to get user access discussions")
		return nil, err
	}

	discussions = append(publicDiscussions, flairDiscussions...)
	discussions = append(discussions, userDiscussions...)

	dedupedDiscs := dedupeDiscussions(discussions)

	return dedupedDiscs, nil
}

func (d *delphisBackend) GetDiscussionFlairTemplateAccessByDiscussionID(ctx context.Context, discussionID string) ([]*model.FlairTemplate, error) {
	iter := d.db.GetDiscussionFlairTemplatesAccessByDiscussionID(ctx, discussionID)
	return d.iterToFlairTemplates(ctx, iter)
}

func (d *delphisBackend) PutDiscussionFlairTemplatesAccess(ctx context.Context, userID string, discussionID string, flairTemplateIDs []string) ([]*model.DiscussionFlairTemplateAccess, error) {
	if len(flairTemplateIDs) == 0 {
		return nil, fmt.Errorf("no flair template IDs to add to discussion")
	}

	validatedTemplates, err := d.validateFlairTemplatesToAdd(ctx, userID, flairTemplateIDs)
	if err != nil {
		logrus.WithError(err).Error("failed to validate flair templates to add")
		return nil, err
	}

	var addedTemplates []*model.DiscussionFlairTemplateAccess

	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return nil, err
	}

	for _, id := range validatedTemplates {
		resp, err := d.db.UpsertDiscussionFlairTemplatesAccess(ctx, tx, discussionID, id)
		if err != nil {
			logrus.WithError(err).Error("failed to UpsertDiscussionFlairTemplatesAccess")

			// Rollback on errors
			if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
				logrus.WithError(txErr).Error("failed to rollback tx")
				return nil, multierr.Append(err, txErr)
			}
			return nil, err
		}
		addedTemplates = append(addedTemplates, resp)
	}

	// Commit transaction
	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("failed to commit post tx")
		return nil, err
	}

	return addedTemplates, nil
}

func (d *delphisBackend) DeleteDiscussionFlairTemplatesAccess(ctx context.Context, discussionID string, flairTemplateIDs []string) ([]*model.DiscussionFlairTemplateAccess, error) {
	if len(flairTemplateIDs) == 0 {
		return nil, fmt.Errorf("no flair template IDs to delete from discussion")
	}

	var deletedTemplates []*model.DiscussionFlairTemplateAccess

	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return nil, err
	}

	for _, id := range flairTemplateIDs {
		resp, err := d.db.DeleteDiscussionFlairTemplatesAccess(ctx, tx, discussionID, id)
		if err != nil {
			logrus.WithError(err).Error("failed to DeleteDiscussionFlairTemplatesAccess")

			// Rollback on errors
			if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
				logrus.WithError(txErr).Error("failed to rollback tx")
				return nil, multierr.Append(err, txErr)
			}
			return nil, err
		}
		deletedTemplates = append(deletedTemplates, resp)
	}

	// Commit transaction
	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("failed to commit post tx")
		return nil, err
	}

	return deletedTemplates, nil
}

func (d *delphisBackend) iterToDiscussions(ctx context.Context, iter datastore.DiscussionIter) ([]*model.Discussion, error) {
	var discussions []*model.Discussion
	disc := model.Discussion{}

	defer iter.Close()

	for iter.Next(&disc) {
		tempDisc := disc

		discussions = append(discussions, &tempDisc)
	}

	if err := iter.Close(); err != nil && err != io.EOF {
		logrus.WithError(err).Error("failed to close iter")
		return nil, err
	}

	return discussions, nil
}

func (d *delphisBackend) iterToFlairTemplates(ctx context.Context, iter datastore.DFAIter) ([]*model.FlairTemplate, error) {
	var templates []*model.FlairTemplate
	dfa := model.DiscussionFlairTemplateAccess{}

	defer iter.Close()

	for iter.Next(&dfa) {
		template, err := d.db.GetFlairTemplateByID(ctx, dfa.FlairTemplateID)
		if err != nil {
			logrus.WithError(err).Error("failed to get flair template by id")
			return nil, err
		}

		templates = append(templates, template)
	}

	if err := iter.Close(); err != nil && err != io.EOF {
		logrus.WithError(err).Error("failed to close iter")
		return nil, err
	}

	return templates, nil
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

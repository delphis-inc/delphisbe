package backend

import (
	"context"
	"io"
	"strings"

	"github.com/nedrocks/delphisbe/internal/datastore"

	"go.uber.org/multierr"

	"github.com/nedrocks/delphisbe/internal/util"
	"github.com/sirupsen/logrus"

	"github.com/nedrocks/delphisbe/graph/model"
)

func (d *delphisBackend) GetUpcomingImportedContentByDiscussionID(ctx context.Context, discussionID string) ([]*model.ImportedContent, error) {
	iter := d.db.GetScheduledImportedContentByDiscussionID(ctx, discussionID)
	schedContents, err := d.iterToContent(ctx, iter)
	if err != nil {
		logrus.WithError(err).Error("failed to get scheduledContent")
		return nil, err
	}

	// Do we want the client to pass in the limit?
	iter = d.db.GetImportedContentByDiscussionID(ctx, discussionID, 10)
	importedContents, err := d.iterToContent(ctx, iter)
	if err != nil {
		logrus.WithError(err).Error("failed to get imported content")
		return nil, err
	}

	return append(schedContents, importedContents...), nil
}

func (d *delphisBackend) GetImportedContentByID(ctx context.Context, id string) (*model.ImportedContent, error) {
	return d.db.GetImportedContentByID(ctx, id)
}

func (d *delphisBackend) GetMatchingsTags(ctx context.Context, discussionID, contentID string) ([]string, error) {
	return d.db.GetMatchingTags(ctx, discussionID, contentID)
}

func (d *delphisBackend) PutImportedContentAndTags(ctx context.Context, input model.ImportedContentInput) (*model.ImportedContent, error) {
	icID := util.UUIDv4()
	importedContentObj := model.ImportedContent{
		ID:          icID,
		ContentName: input.ContentName,
		ContentType: input.ContentType,
		Link:        input.Link,
		Overview:    input.Overview,
		Source:      input.Source,
	}

	// Begin tx
	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return nil, err
	}

	importedContent, err := d.db.PutImportedContent(ctx, tx, importedContentObj)
	if err != nil {
		logrus.WithError(err).Error("failed to PutImportedContent")

		// Rollback on errors
		if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
			logrus.WithError(txErr).Error("failed to rollback tx")
			return nil, multierr.Append(err, txErr)
		}
		return nil, err
	}

	if len(input.Tags) != 0 {
		// Split tag list
		tagList := strings.Split(input.Tags, ",")

		// Create tag object
		tagObj := model.Tag{
			ID: icID,
		}
		for _, tag := range tagList {
			tagObj.Tag = tag
			_, err := d.db.PutImportedContentTags(ctx, tx, tagObj)
			if err != nil {
				logrus.WithError(err).Error("failed to PutImportedContentTag")

				// Rollback on errors
				if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
					logrus.WithError(txErr).Error("failed to rollback tx")
					return nil, multierr.Append(err, txErr)
				}
				return nil, err
			}
		}
	}
	// Commit transaction
	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("failed to commit post tx")
		return nil, err
	}

	return importedContent, nil
}

func (d *delphisBackend) iterToContent(ctx context.Context, iter datastore.ContentIter) ([]*model.ImportedContent, error) {
	var contents []*model.ImportedContent
	content := model.ImportedContent{}

	defer iter.Close()

	for iter.Next(&content) {
		tempContent := content

		contents = append(contents, &tempContent)
	}

	if err := iter.Close(); err != nil && err != io.EOF {
		logrus.WithError(err).Error("failed to close iter")
		return nil, err
	}

	return contents, nil
}

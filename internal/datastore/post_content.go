package datastore

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *delphisDB) PutPostContent(ctx context.Context, tx *sql.Tx, postContent model.PostContent) error {
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("PutPost::failed to initialize statements")
		return err
	}

	_, err := tx.StmtContext(ctx, d.prepStmts.putPostContentsStmt).ExecContext(
		ctx,
		postContent.ID,
		postContent.Content,
		pq.Array(postContent.MentionedEntities),
	)
	if err != nil {
		logrus.WithError(err).Error("failed to execute putPostContentsStmt")
		return errors.Wrap(err, "failed to put postContents")
	}

	return nil
}

func (d *delphisDB) GetPostContentByID(ctx context.Context, id string) (*model.PostContent, error) {
	logrus.Debug("GetPostContentByID::SQL Query")
	found := model.PostContent{}
	if err := d.sql.First(&found, &model.PostContent{ID: id}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("Failed to get PostContent by ID")
		return nil, err
	}
	return &found, nil
}

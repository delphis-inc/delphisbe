package datastore

import (
	"context"
	"database/sql"
	"io"
	"time"

	"github.com/lib/pq"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *delphisDB) GetImportedContentByID(ctx context.Context, id string) (*model.ImportedContent, error) {
	logrus.Debug("GetImportedContentByID::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetImportedContentByID::failed to initialize statements")
		return nil, err
	}

	ic := model.ImportedContent{}
	if err := d.prepStmts.getImportedContentByIDStmt.QueryRowContext(
		ctx,
		id,
	).Scan(
		&ic.ID,
		&ic.CreatedAt,
		&ic.ContentName,
		&ic.ContentType,
		&ic.Link,
		&ic.Overview,
		&ic.Source,
	); err != nil {
		logrus.WithError(err).Error("failed to execute getImportedContentByIDStmt")
		return nil, err
	}

	return &ic, nil
}

func (d *delphisDB) GetImportedContentTags(ctx context.Context, id string) TagIter {
	logrus.Debug("GetImportedContentTags::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetImportedContentTags::failed to initialize statements")
		return &tagIter{err: err}
	}

	rows, err := d.prepStmts.getImportedContentTagsStmt.QueryContext(
		ctx,
		id,
	)
	if err != nil {
		logrus.WithError(err).Error("failed to query GetImportedContentTags")
		return &tagIter{err: err}
	}

	return &tagIter{
		ctx:  ctx,
		rows: rows,
	}
}

func (d *delphisDB) GetImportedContentByDiscussionID(ctx context.Context, discussionID string, limit int) ContentIter {
	logrus.Debug("GetImportedContentByDiscussionID::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetImportedContentByDiscussionID::failed to initialize statements")
		return &contentIter{err: err}
	}

	rows, err := d.prepStmts.getImportedContentForDiscussionStmt.QueryContext(
		ctx,
		discussionID,
		limit,
	)
	if err != nil {
		logrus.WithError(err).Error("failed to query GetImportedContentByDiscussionID")
		return &contentIter{err: err}
	}

	return &contentIter{
		ctx:  ctx,
		rows: rows,
	}
}

func (d *delphisDB) GetScheduledImportedContentByDiscussionID(ctx context.Context, discussionID string) ContentIter {
	logrus.Debug("GetScheduledImportedContentByDiscussionID::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetScheduledImportedContentByDiscussionID::failed to initialize statements")
		return &contentIter{err: err}
	}

	rows, err := d.prepStmts.getScheduledImportedContentByDiscussionIDStmt.QueryContext(
		ctx,
		discussionID,
	)
	if err != nil {
		logrus.WithError(err).Error("failed to query GetScheduledImportedContentByDiscussionID")
		return &contentIter{err: err}
	}

	return &contentIter{
		ctx:  ctx,
		rows: rows,
	}
}

func (d *delphisDB) PutImportedContent(ctx context.Context, tx *sql.Tx, ic model.ImportedContent) (*model.ImportedContent, error) {
	logrus.Debug("PutImportedContent::SQL Create")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("PutImportedContent::failed to initialize statements")
		return nil, err
	}

	if err := tx.StmtContext(ctx, d.prepStmts.putImportedContentStmt).QueryRowContext(
		ctx,
		ic.ID,
		ic.ContentName,
		ic.ContentType,
		ic.Link,
		ic.Overview,
		ic.Source,
	).Scan(
		&ic.ID,
		&ic.CreatedAt,
		&ic.ContentName,
		&ic.ContentType,
		&ic.Link,
		&ic.Overview,
		&ic.Source,
	); err != nil {
		logrus.WithError(err).Error("failed to execute putImportedContentStmt")
		return nil, err
	}

	return &ic, nil
}

func (d *delphisDB) PutImportedContentTags(ctx context.Context, tx *sql.Tx, tag model.Tag) (*model.Tag, error) {
	logrus.Debug("PutImportedContentTags::SQL Create")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("PutImportedContentTags::failed to initialize statements")
		return nil, err
	}

	if err := tx.StmtContext(ctx, d.prepStmts.putImportedContentTagsStmt).QueryRowContext(
		ctx,
		tag.ID,
		tag.Tag,
	).Scan(
		&tag.ID,
		&tag.Tag,
		&tag.CreatedAt,
	); err != nil {
		logrus.WithError(err).Error("failed to execute putImportedContentTagsStmt")
		return nil, err
	}

	return &tag, nil
}

func (d *delphisDB) GetMatchingTags(ctx context.Context, discussionID, importedContentID string) ([]string, error) {
	logrus.Debug("GetMatchingTags::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetMatchingTags::failed to initialize statements")
		return nil, err
	}

	tags := make([]string, 0)
	if err := d.prepStmts.getMatchingTagsStmt.QueryRowContext(
		ctx,
		discussionID,
		importedContentID,
	).Scan(
		pq.Array(&tags),
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("failed to GetMatchingTags")
		return nil, err
	}

	return tags, nil
}

func (d *delphisDB) PutImportedContentDiscussionQueue(ctx context.Context, discussionID, contentID string, postedAt *time.Time, matchingTags []string) (*model.ContentQueueRecord, error) {
	logrus.Debug("PutImportedContentDiscussionQueue::SQL Create")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("PutImportedContentDiscussionQueue::failed to initialize statements")
		return nil, err
	}

	contentQ := model.ContentQueueRecord{}
	if err := d.prepStmts.putImportedContentDiscussionQueueStmt.QueryRowContext(
		ctx,
		discussionID,
		contentID,
		postedAt,
		pq.Array(matchingTags),
	).Scan(
		&contentQ.DiscussionID,
		&contentQ.ImportedContentID,
		&contentQ.CreatedAt,
		&contentQ.UpdatedAt,
		&contentQ.DeletedAt,
		&contentQ.PostedAt,
		pq.Array(&contentQ.MatchingTags),
	); err != nil {
		logrus.WithError(err).Error("failed to execute putImportedContentDiscussionQueueStmt")
		return nil, err
	}

	return &contentQ, nil
}

func (d *delphisDB) UpdateImportedContentDiscussionQueue(ctx context.Context, discussionID, contentID string, postedAt *time.Time) (*model.ContentQueueRecord, error) {
	logrus.Debug("UpdateImportedContentDiscussionQueue::SQL Create")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("UpdateImportedContentDiscussionQueue::failed to initialize statements")
		return nil, err
	}

	contentQ := model.ContentQueueRecord{}
	if err := d.prepStmts.updateImportedContentDiscussionQueueStmt.QueryRowContext(
		ctx,
		discussionID,
		contentID,
		postedAt,
	).Scan(
		&contentQ.DiscussionID,
		&contentQ.ImportedContentID,
		&contentQ.CreatedAt,
		&contentQ.UpdatedAt,
		&contentQ.DeletedAt,
		&contentQ.PostedAt,
		pq.Array(&contentQ.MatchingTags),
	); err != nil {
		logrus.WithError(err).Error("failed to execute updateImportedContentDiscussionQueueStmt")
		return nil, err
	}

	return &contentQ, nil
}

type tagIter struct {
	err  error
	ctx  context.Context
	rows *sql.Rows
}

func (iter *tagIter) Next(tag *model.Tag) bool {
	if iter.err != nil {
		logrus.WithError(iter.err).Error("iterator error")
		return false
	}

	if iter.err = iter.ctx.Err(); iter.err != nil {
		logrus.WithError(iter.err).Error("iterator context error")
		return false
	}

	if !iter.rows.Next() {
		return false
	}

	if iter.err = iter.rows.Scan(
		&tag.ID,
		&tag.Tag,
		&tag.CreatedAt,
	); iter.err != nil {
		logrus.WithError(iter.err).Error("iterator failed to scan row")
		return false
	}

	return true
}

func (iter *tagIter) Close() error {
	if err := iter.err; err != nil {
		logrus.WithError(err).Error("iter error on close")
		return err
	}
	if err := iter.rows.Close(); err != nil {
		logrus.WithError(err).Error("iter rows close on close")
		return err
	}

	return nil
}

type contentIter struct {
	err  error
	ctx  context.Context
	rows *sql.Rows
}

func (iter *contentIter) Next(content *model.ImportedContent) bool {
	if iter.err != nil {
		logrus.WithError(iter.err).Error("iterator error")
		return false
	}

	if iter.err = iter.ctx.Err(); iter.err != nil {
		logrus.WithError(iter.err).Error("iterator context error")
		return false
	}

	if !iter.rows.Next() {
		return false
	}

	if iter.err = iter.rows.Scan(
		&content.ID,
		&content.CreatedAt,
		&content.ContentName,
		&content.ContentType,
		&content.Link,
		&content.Overview,
		&content.Source,
		pq.Array(&content.Tags),
	); iter.err != nil {
		logrus.WithError(iter.err).Error("iterator failed to scan row")
		return false
	}

	return true
}

func (iter *contentIter) Close() error {
	if err := iter.err; err != nil {
		logrus.WithError(err).Error("iter error on close")
		return err
	}
	if err := iter.rows.Close(); err != nil {
		logrus.WithError(err).Error("iter rows close on close")
		return err
	}

	return nil
}

func (d *delphisDB) ContentIterCollect(ctx context.Context, iter ContentIter) ([]*model.ImportedContent, error) {
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

func (d *delphisDB) TagIterCollect(ctx context.Context, iter TagIter) ([]*model.Tag, error) {
	var tags []*model.Tag
	tag := model.Tag{}

	defer iter.Close()

	for iter.Next(&tag) {
		tempTag := tag

		tags = append(tags, &tempTag)
	}

	if err := iter.Close(); err != nil && err != io.EOF {
		logrus.WithError(err).Error("failed to close iter")
		return nil, err
	}

	return tags, nil
}

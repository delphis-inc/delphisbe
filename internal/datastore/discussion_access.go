package datastore

import (
	"context"
	"database/sql"
	"io"

	"github.com/sirupsen/logrus"

	"github.com/nedrocks/delphisbe/graph/model"
)

func (d *delphisDB) GetPublicDiscussions(ctx context.Context) DiscussionIter {
	logrus.Debug("GetPublicDiscussions::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetPublicDiscussions::failed to initialize statements")
		return &discussionIter{err: err}
	}

	rows, err := d.prepStmts.getPublicDiscussionsStmt.QueryContext(
		ctx,
	)
	if err != nil {
		logrus.WithError(err).Error("failed to query GetPublicDiscussions")
		return &discussionIter{err: err}
	}

	return &discussionIter{
		ctx:  ctx,
		rows: rows,
	}
}

func (d *delphisDB) GetDiscussionsForFlairTemplateByUserID(ctx context.Context, userID string) DiscussionIter {
	logrus.Debug("GetDiscussionsForFlairTemplateByUserID::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetDiscussionsForFlairTemplateByUserID::failed to initialize statements")
		return &discussionIter{err: err}
	}

	rows, err := d.prepStmts.getDiscussionsByFlairTemplateForUserStmt.QueryContext(
		ctx,
		userID,
	)
	if err != nil {
		logrus.WithError(err).Error("failed to query GetDiscussionsForFlairTemplateByUserID")
		return &discussionIter{err: err}
	}

	return &discussionIter{
		ctx:  ctx,
		rows: rows,
	}
}

func (d *delphisDB) GetDiscussionsForUserAccessByUserID(ctx context.Context, userID string) DiscussionIter {
	logrus.Debug("GetDiscussionsForUserAccessByUserID::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetDiscussionsForUserAccessByUserID::failed to initialize statements")
		return &discussionIter{err: err}
	}

	rows, err := d.prepStmts.getDiscussionsByUserAccessForUserStmt.QueryContext(
		ctx,
		userID,
	)
	if err != nil {
		logrus.WithError(err).Error("failed to query GetDiscussionsForUserAccessByUserID")
		return &discussionIter{err: err}
	}

	return &discussionIter{
		ctx:  ctx,
		rows: rows,
	}
}

func (d *delphisDB) GetDiscussionFlairTemplatesAccessByDiscussionID(ctx context.Context, discussionID string) DFAIter {
	logrus.Debug("GetDiscussionFlairTemplatesAccessByDiscussionID::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetDiscussionFlairTemplatesAccessByDiscussionID::failed to initialize statements")
		return &dfaIter{err: err}
	}

	rows, err := d.prepStmts.getDiscussionFlairAccessStmt.QueryContext(
		ctx,
		discussionID,
	)
	if err != nil {
		logrus.WithError(err).Error("failed to query etDiscussionFlairTemplatesAccessStmt")
		return &dfaIter{err: err}
	}

	return &dfaIter{
		ctx:  ctx,
		rows: rows,
	}
}

func (d *delphisDB) UpsertDiscussionFlairTemplatesAccess(ctx context.Context, tx *sql.Tx, discussionID, flairTemplateID string) (*model.DiscussionFlairTemplateAccess, error) {
	logrus.Debug("UpsertDiscussionFlairTemplatesAccess::SQL Create")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("UpsertDiscussionFlairTemplatesAccess::failed to initialize statements")
		return nil, err
	}

	dfa := model.DiscussionFlairTemplateAccess{}
	if err := tx.StmtContext(ctx, d.prepStmts.upsertDiscussionFlairAccessStmt).QueryRowContext(
		ctx,
		discussionID,
		flairTemplateID,
	).Scan(
		&dfa.DiscussionID,
		&dfa.FlairTemplateID,
		&dfa.CreatedAt,
		&dfa.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return &model.DiscussionFlairTemplateAccess{}, nil
		}
		logrus.WithError(err).Error("failed to execute upsertDiscussionFlairTemplatesStmt")
		return nil, err
	}

	return &dfa, nil
}

func (d *delphisDB) UpsertDiscussionUserAccess(ctx context.Context, tx *sql.Tx, discussionID, userID string) (*model.DiscussionUserAccess, error) {
	logrus.Debug("UpsertDiscussionUserAccess::SQL Create")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("UpsertDiscussionUserAccess::failed to initialize statements")
		return nil, err
	}

	dua := model.DiscussionUserAccess{}
	if err := tx.StmtContext(ctx, d.prepStmts.upsertDiscussionUserAccessStmt).QueryRowContext(
		ctx,
		discussionID,
		userID,
	).Scan(
		&dua.DiscussionID,
		&dua.UserID,
		&dua.CreatedAt,
		&dua.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return &model.DiscussionUserAccess{}, nil
		}
		logrus.WithError(err).Error("failed to execute upsertDiscussionUserAccess")
		return nil, err
	}

	return &dua, nil
}

func (d *delphisDB) DeleteDiscussionFlairTemplatesAccess(ctx context.Context, tx *sql.Tx, discussionID, flairTemplateID string) (*model.DiscussionFlairTemplateAccess, error) {
	logrus.Debug("DeleteDiscussionFlairTemplatesAccess::SQL Delete")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("DeleteDiscussionFlairTemplatesAccess::failed to initialize statements")
		return nil, err
	}

	dfa := model.DiscussionFlairTemplateAccess{}
	if err := tx.StmtContext(ctx, d.prepStmts.deleteDiscussionFlairAccessStmt).QueryRowContext(
		ctx,
		discussionID,
		flairTemplateID,
	).Scan(
		&dfa.DiscussionID,
		&dfa.FlairTemplateID,
		&dfa.CreatedAt,
		&dfa.UpdatedAt,
		&dfa.DeletedAt,
	); err != nil {
		logrus.WithError(err).Error("failed to execute deleteDiscussionFlairTemplatesStmt")
		return nil, err
	}

	return &dfa, nil
}

func (d *delphisDB) DeleteDiscussionUserAccess(ctx context.Context, tx *sql.Tx, discussionID, userID string) (*model.DiscussionUserAccess, error) {
	logrus.Debug("DeleteDiscussionUserAccess::SQL Delete")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("DeleteDiscussionUserAccess::failed to initialize statements")
		return nil, err
	}

	dua := model.DiscussionUserAccess{}
	if err := tx.StmtContext(ctx, d.prepStmts.deleteDiscussionUserAccessStmt).QueryRowContext(
		ctx,
		discussionID,
		userID,
	).Scan(
		&dua.DiscussionID,
		&dua.UserID,
		&dua.CreatedAt,
		&dua.UpdatedAt,
		&dua.DeletedAt,
	); err != nil {
		logrus.WithError(err).Error("failed to execute deleteDiscussionUserAccessStmt")
		return nil, err
	}

	return &dua, nil
}

type discussionIter struct {
	err  error
	ctx  context.Context
	rows *sql.Rows
}

func (iter *discussionIter) Next(discussion *model.Discussion) bool {
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
		&discussion.ID,
		&discussion.CreatedAt,
		&discussion.Title,
		&discussion.AnonymityType,
		&discussion.ModeratorID,
		&discussion.AutoPost,
		&discussion.IconURL,
		&discussion.IdleMinutes,
		&discussion.PublicAccess,
	); iter.err != nil {
		logrus.WithError(iter.err).Error("iterator failed to scan row")
		return false
	}

	return true
}

func (iter *discussionIter) Close() error {
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

type dfaIter struct {
	err  error
	ctx  context.Context
	rows *sql.Rows
}

func (iter *dfaIter) Next(dfa *model.DiscussionFlairTemplateAccess) bool {
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
		&dfa.DiscussionID,
		&dfa.FlairTemplateID,
		&dfa.CreatedAt,
		&dfa.UpdatedAt,
	); iter.err != nil {
		logrus.WithError(iter.err).Error("iterator failed to scan row")
		return false
	}

	return true
}

func (iter *dfaIter) Close() error {
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

func (d *delphisDB) DiscussionIterCollect(ctx context.Context, iter DiscussionIter) ([]*model.Discussion, error) {
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

func (d *delphisDB) FlairTemplatesIterCollect(ctx context.Context, iter DFAIter) ([]*model.FlairTemplate, error) {
	var templates []*model.FlairTemplate
	dfa := model.DiscussionFlairTemplateAccess{}

	defer iter.Close()

	for iter.Next(&dfa) {
		template, err := d.GetFlairTemplateByID(ctx, dfa.FlairTemplateID)
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

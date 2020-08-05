package datastore

import (
	"context"
	"database/sql"
	"io"

	"github.com/sirupsen/logrus"

	"github.com/delphis-inc/delphisbe/graph/model"
)

func (d *delphisDB) GetDiscussionsByUserAccess(ctx context.Context, userID string, state model.DiscussionUserAccessState) DiscussionIter {
	logrus.Debug("GetDiscussionsByUserAccess::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetDiscussionsByUserAccess::failed to initialize statements")
		return &discussionIter{err: err}
	}

	rows, err := d.prepStmts.getDiscussionsByUserAccessStmt.QueryContext(
		ctx,
		userID,
		state,
	)
	if err != nil {
		logrus.WithError(err).Error("failed to query GetDiscussionsByUserAccess")
		return &discussionIter{err: err}
	}

	return &discussionIter{
		ctx:  ctx,
		rows: rows,
	}
}

func (d *delphisDB) GetDiscussionUserAccess(ctx context.Context, discussionID, userID string) (*model.DiscussionUserAccess, error) {
	logrus.Debug("GetDiscussionUserAccess::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetDiscussionUserAccess::failed to initialize statements")
		return nil, err
	}

	dua := model.DiscussionUserAccess{}
	if err := d.prepStmts.getDiscussionUserAccessStmt.QueryRowContext(
		ctx,
		discussionID,
		userID,
	).Scan(
		&dua.DiscussionID,
		&dua.UserID,
		&dua.State,
		&dua.RequestID,
		&dua.CreatedAt,
		&dua.UpdatedAt,
		&dua.DeletedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("failed to query GetDiscussionsByUserAccess")
		return nil, err
	}

	return &dua, nil
}

func (d *delphisDB) UpsertDiscussionUserAccess(ctx context.Context, tx *sql.Tx, dua model.DiscussionUserAccess) (*model.DiscussionUserAccess, error) {
	logrus.Debug("UpsertDiscussionUserAccess::SQL Create")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("UpsertDiscussionUserAccess::failed to initialize statements")
		return nil, err
	}

	if err := tx.StmtContext(ctx, d.prepStmts.upsertDiscussionUserAccessStmt).QueryRowContext(
		ctx,
		dua.DiscussionID,
		dua.UserID,
		dua.State,
		dua.RequestID,
	).Scan(
		&dua.DiscussionID,
		&dua.UserID,
		&dua.State,
		&dua.RequestID,
		&dua.CreatedAt,
		&dua.UpdatedAt,
		&dua.DeletedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return &model.DiscussionUserAccess{}, nil
		}
		logrus.WithError(err).Error("failed to execute upsertDiscussionUserAccess")
		return nil, err
	}

	return &dua, nil
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

	titleHistory := make([]byte, 0)
	descriptionHistory := make([]byte, 0)

	if iter.err = iter.rows.Scan(
		&discussion.ID,
		&discussion.CreatedAt,
		&discussion.UpdatedAt,
		&discussion.DeletedAt,
		&discussion.Title,
		&discussion.AnonymityType,
		&discussion.ModeratorID,
		&discussion.AutoPost,
		&discussion.IconURL,
		&discussion.IdleMinutes,
		&discussion.Description,
		&titleHistory,
		&descriptionHistory,
		&discussion.DiscussionJoinability,
		&discussion.LastPostID,
		&discussion.LastPostCreatedAt,
	); iter.err != nil {
		logrus.WithError(iter.err).Error("iterator failed to scan row")
		return false
	}

	discussion.TitleHistory.RawMessage = titleHistory
	discussion.DescriptionHistory.RawMessage = descriptionHistory

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

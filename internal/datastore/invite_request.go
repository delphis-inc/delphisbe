package datastore

import (
	"context"
	"database/sql"
	"io"

	"github.com/sirupsen/logrus"

	"github.com/delphis-inc/delphisbe/graph/model"
)

func (d *delphisDB) GetDiscussionRequestAccessByID(ctx context.Context, id string) (*model.DiscussionAccessRequest, error) {
	logrus.Debug("GetDiscussionRequestAccessByID::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetDiscussionRequestAccessByID::failed to initialize statements")
		return nil, err
	}

	request := model.DiscussionAccessRequest{}
	if err := d.prepStmts.getDiscussionRequestAccessByIDStmt.QueryRowContext(
		ctx,
		id,
	).Scan(
		&request.ID,
		&request.UserID,
		&request.DiscussionID,
		&request.CreatedAt,
		&request.UpdatedAt,
		&request.Status,
	); err != nil {
		logrus.WithError(err).Error("failed to execute GetDiscussionRequestAccessByID")
		return nil, err
	}

	return &request, nil
}

func (d *delphisDB) GetDiscussionAccessRequestsByDiscussionID(ctx context.Context, discussionID string) DiscussionAccessRequestIter {
	logrus.Debug("GetDiscussionAccessRequestsByDiscussionID::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetDiscussionAccessRequestsByDiscussionID::failed to initialize statements")
		return &discussionAccessRequestIter{err: err}
	}

	rows, err := d.prepStmts.getDiscussionAccessRequestsStmt.QueryContext(
		ctx,
		discussionID,
	)
	if err != nil {
		logrus.WithError(err).Error("failed to query GetDiscussionAccessRequestsByDiscussionID")
		return &discussionAccessRequestIter{err: err}
	}

	return &discussionAccessRequestIter{
		ctx:  ctx,
		rows: rows,
	}

}

func (d *delphisDB) GetDiscussionAccessRequestByDiscussionIDUserID(ctx context.Context, discussionID string, userID string) (*model.DiscussionAccessRequest, error) {
	logrus.Debug("GetDiscussionAccessRequestsByDiscussionIDUserID::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetDiscussionAccessRequestsByDiscussionIDUserID::failed to initialize statements")
		return nil, err
	}

	accessRequest := model.DiscussionAccessRequest{}
	if err := d.prepStmts.getDiscussionAccessRequestByUserIDStmt.QueryRowContext(
		ctx,
		discussionID,
		userID,
	).Scan(
		&accessRequest.ID,
		&accessRequest.UserID,
		&accessRequest.DiscussionID,
		&accessRequest.CreatedAt,
		&accessRequest.UpdatedAt,
		&accessRequest.Status,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("Failed to get discussion access request")
		return nil, err
	}

	return &accessRequest, nil
}

func (d *delphisDB) GetSentDiscussionAccessRequestsByUserID(ctx context.Context, userID string) DiscussionAccessRequestIter {
	logrus.Debug("GetSentDiscussionAccessRequestsByUserID::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetSentDiscussionAccessRequestsByUserID::failed to initialize statements")
		return &discussionAccessRequestIter{err: err}
	}

	rows, err := d.prepStmts.getSentDiscussionAccessRequestsForUserStmt.QueryContext(
		ctx,
		userID,
	)
	if err != nil {
		logrus.WithError(err).Error("failed to query GetSentDiscussionAccessRequestsByUserID")
		return &discussionAccessRequestIter{err: err}
	}

	return &discussionAccessRequestIter{
		ctx:  ctx,
		rows: rows,
	}

}

func (d *delphisDB) PutDiscussionAccessRequestRecord(ctx context.Context, tx *sql.Tx, request model.DiscussionAccessRequest) (*model.DiscussionAccessRequest, error) {
	logrus.Debug("PutDiscussionAccessRequestRecord::SQL Create")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("PutDiscussionAccessRequestRecord::failed to initialize statements")
		return nil, err
	}

	if err := tx.StmtContext(ctx, d.prepStmts.putDiscussionAccessRequestStmt).QueryRowContext(
		ctx,
		request.ID,
		request.UserID,
		request.DiscussionID,
		request.Status,
	).Scan(
		&request.ID,
		&request.UserID,
		&request.DiscussionID,
		&request.CreatedAt,
		&request.UpdatedAt,
		&request.Status,
	); err != nil {
		logrus.WithError(err).Error("failed to execute PutDiscussionAccessRequestRecord")
		return nil, err
	}

	return &request, nil
}

func (d *delphisDB) UpdateDiscussionAccessRequestRecord(ctx context.Context, tx *sql.Tx, request model.DiscussionAccessRequest) (*model.DiscussionAccessRequest, error) {
	logrus.Debug("UpdateDiscussionAccessRequestRecord::SQL Create")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("UpdateDiscussionAccessRequestRecord::failed to initialize statements")
		return nil, err
	}

	if err := tx.StmtContext(ctx, d.prepStmts.updateDiscussionAccessRequestStmt).QueryRowContext(
		ctx,
		request.ID,
		request.Status,
	).Scan(
		&request.ID,
		&request.UserID,
		&request.DiscussionID,
		&request.CreatedAt,
		&request.UpdatedAt,
		&request.Status,
	); err != nil {
		logrus.WithError(err).Error("failed to execute UpdateDiscussionAccessRequestRecord")
		return nil, err
	}

	return &request, nil
}

type discussionAccessRequestIter struct {
	err  error
	ctx  context.Context
	rows *sql.Rows
}

func (iter *discussionAccessRequestIter) Next(request *model.DiscussionAccessRequest) bool {
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
		&request.ID,
		&request.UserID,
		&request.DiscussionID,
		&request.CreatedAt,
		&request.UpdatedAt,
		&request.Status,
	); iter.err != nil {
		logrus.WithError(iter.err).Error("iterator failed to scan row")
		return false
	}

	return true
}

func (iter *discussionAccessRequestIter) Close() error {
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

func (d *delphisDB) AccessRequestIterCollect(ctx context.Context, iter DiscussionAccessRequestIter) ([]*model.DiscussionAccessRequest, error) {
	var requests []*model.DiscussionAccessRequest
	request := model.DiscussionAccessRequest{}

	defer iter.Close()

	for iter.Next(&request) {
		tempRequest := request

		requests = append(requests, &tempRequest)
	}

	if err := iter.Close(); err != nil && err != io.EOF {
		logrus.WithError(err).Error("failed to close iter")
		return nil, err
	}

	return requests, nil
}

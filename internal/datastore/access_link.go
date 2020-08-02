package datastore

import (
	"context"
	"database/sql"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *delphisDB) GetAccessLinkBySlug(ctx context.Context, slug string) (*model.DiscussionAccessLink, error) {
	logrus.Debug("GetAccessLinkBySlug::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetAccessLinkBySlug::failed to initialize statements")
		return nil, err
	}

	dal := model.DiscussionAccessLink{}
	if err := d.prepStmts.getAccessLinkBySlugStmt.QueryRowContext(
		ctx,
		slug,
	).Scan(
		&dal.DiscussionID,
		&dal.LinkSlug,
		&dal.CreatedAt,
		&dal.UpdatedAt,
		&dal.DeletedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return &model.DiscussionAccessLink{}, nil
		}
		logrus.WithError(err).Error("failed to execute getAccessLinkBySlugStmt")
		return nil, err
	}

	return &dal, nil
}

func (d *delphisDB) GetAccessLinkByDiscussionID(ctx context.Context, discussionID string) (*model.DiscussionAccessLink, error) {
	logrus.Debug("GetAccessLinkByDiscussionID::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetAccessLinkByDiscussionID::failed to initialize statements")
		return nil, err
	}

	dal := model.DiscussionAccessLink{}
	if err := d.prepStmts.getAccessLinkByDiscussionIDString.QueryRowContext(
		ctx,
		discussionID,
	).Scan(
		&dal.DiscussionID,
		&dal.LinkSlug,
		&dal.CreatedAt,
		&dal.UpdatedAt,
		&dal.DeletedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return &model.DiscussionAccessLink{}, nil
		}
		logrus.WithError(err).Error("failed to execute getAccessLinkByDiscussionIDString")
		return nil, err
	}

	return &dal, nil
}

func (d *delphisDB) PutAccessLinkForDiscussion(ctx context.Context, tx *sql.Tx, input model.DiscussionAccessLink) (*model.DiscussionAccessLink, error) {
	logrus.Debug("PutAccessLinkForDiscussion::SQL Create")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("PutAccessLinkForDiscussion::failed to initialize statements")
		return nil, err
	}

	if err := tx.StmtContext(ctx, d.prepStmts.putAccessLinkForDiscussionString).QueryRowContext(
		ctx,
		input.DiscussionID,
		input.LinkSlug,
	).Scan(
		&input.DiscussionID,
		&input.LinkSlug,
		&input.CreatedAt,
		&input.UpdatedAt,
		&input.DeletedAt,
	); err != nil {
		logrus.WithError(err).Error("failed to execute putAccessLinkForDiscussionString")
		return nil, err
	}

	return &input, nil
}

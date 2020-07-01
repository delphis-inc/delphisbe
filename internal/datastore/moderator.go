package datastore

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *delphisDB) GetModeratorByID(ctx context.Context, id string) (*model.Moderator, error) {
	logrus.Debug("GetModeratorByID::SQL Query")
	moderator := model.Moderator{}
	if err := d.sql.Preload("UserProfile").First(&moderator, model.Moderator{ID: id}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetModeratorByID::Failed to get moderator")
		return nil, err
	}
	return &moderator, nil
}

func (d *delphisDB) CreateModerator(ctx context.Context, moderator model.Moderator) (*model.Moderator, error) {
	logrus.Debugf("CreateModerator::SQL Insert")
	found := model.Moderator{}
	if err := d.sql.Create(&moderator).First(&found, model.Moderator{ID: moderator.ID}).Error; err != nil {
		logrus.WithError(err).Errorf("Failed to create moderator")
		return nil, err
	}
	return &found, nil
}

func (d *delphisDB) GetModeratorByUserID(ctx context.Context, id string) (*model.Moderator, error) {
	logrus.Debug("GetModeratorByUserID::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetModeratorByUserID failed to initialize statements")
		return nil, err
	}

	moderator := model.Moderator{}
	if err := d.prepStmts.getModeratorByUserIDStmt.QueryRowContext(
		ctx,
		id,
	).Scan(
		&moderator.ID,
		&moderator.CreatedAt,
		&moderator.UpdatedAt,
		&moderator.DeletedAt,
		&moderator.UserProfileID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("failed to execute getModeratorByUserIDStmt")
		return nil, err
	}

	return &moderator, nil
}

func (d *delphisDB) GetModeratorByUserIDAndDiscussionID(ctx context.Context, userID, discussionID string) (*model.Moderator, error) {
	logrus.Debug("GetModeratorByUserIDAndDiscussionID::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetModeratorByUserIDAndDiscussionID failed to initialize statements")
		return nil, err
	}

	moderator := model.Moderator{}
	if err := d.prepStmts.getModeratorByUserIDAndDiscussionIDStmt.QueryRowContext(
		ctx,
		userID,
		discussionID,
	).Scan(
		&moderator.ID,
		&moderator.CreatedAt,
		&moderator.UpdatedAt,
		&moderator.DeletedAt,
		&moderator.UserProfileID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("failed to execute getModeratorByUserIDAndDiscussionIDStmt")
		return nil, err
	}

	return &moderator, nil
}

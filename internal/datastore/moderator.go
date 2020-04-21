package datastore

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *db) GetModeratorByID(ctx context.Context, id string) (*model.Moderator, error) {
	logrus.Debug("GetModeratorByID::SQL Query")
	moderator := model.Moderator{}
	if err := d.sql.Preload("UserProfile").First(&moderator, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetModeratorByID::Failed to get moderator")
		return nil, err
	}
	return &moderator, nil
}

func (d *db) CreateModerator(ctx context.Context, moderator model.Moderator) (*model.Moderator, error) {
	logrus.Debugf("CreateModerator::SQL Insert")
	if err := d.sql.Save(&moderator).Error; err != nil {
		logrus.WithError(err).Errorf("Failed to create moderator")
		return nil, err
	}
	return &moderator, nil
}

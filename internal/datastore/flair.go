package datastore

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *db) GetFlairByID(ctx context.Context, id string) (*model.Flair, error) {
	logrus.Debug("GetFlairByID::SQL Query")
	flair := model.Flair{}
	if err := d.sql.First(&flair, model.Flair{ID: id}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetFlairByID::Failed to get flair")
		return nil, err
	}
	return &flair, nil
}

func (d *db) GetFlairByUserIDFlairID(ctx context.Context, userID string, flairID string) (*model.Flair, error) {
	logrus.Debugf("GetFlairByUserIDFlairID::SQL Query")
	flair := model.Flair{}
	if err := d.sql.Joins("JOIN user_flairs ON user_id = ? AND flair_id = ?", userID, flairID).First(&flair).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetFlairByUserIDFlairID::Failed to get flair by user ID")
		return nil, err
	}
	return &flair, nil
}



func (d *db) GetFlairsByUserID(ctx context.Context, userID string) ([]model.Flair, error) {
	logrus.Debugf("GetFlairsByUserID::SQL Query")
	flairs := []model.Flair{}
	if err := d.sql.Joins("JOIN user_flairs ON user_id = ?", userID).Find(&flairs).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetFlairsByUserID::Failed to get flairs by user ID")
		return nil, err
	}
	return flairs, nil
}


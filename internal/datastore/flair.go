package datastore

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *db) UpsertFlair(ctx context.Context, data model.Flair) (*model.Flair, error) {
	logrus.Debug("UpsertFlair::SQL Create or Update")
	flair := model.Flair{}
	if err := d.sql.First(&flair, model.Flair{ID: data.ID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := d.sql.Create(&data).First(&flair, model.Flair{ID: data.ID}).Error; err != nil {
				logrus.WithError(err).Errorf("UpsertFlair::Failed to create new object")
				return nil, err
			}
		} else {
			logrus.WithError(err).Errorf("UpsertFlair::Failed checking for Flair object")
			return nil, err
		}
	} else {
		if err := d.sql.Model(&data).Updates(model.Flair{
			DisplayName: data.DisplayName,
			ImageURL:    data.ImageURL,
			Source:      data.Source,
		}).First(&flair).Error; err != nil {
			logrus.WithError(err).Errorf("UpsertFlair::Failed updating flair object")
			return nil, err
		}
	}
	return &flair, nil
}

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


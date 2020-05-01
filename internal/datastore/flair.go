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
			UserID:     data.UserID,
			TemplateID: data.TemplateID,
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

func (d *db) GetFlairsByUserID(ctx context.Context, userID string) ([]*model.Flair, error) {
	logrus.Debug("GetFlairsByUserID::SQL Query")
	flairs := []*model.Flair{}
	if err := d.sql.Where(model.Flair{UserID: userID}).Find(&flairs).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []*model.Flair{}, nil
		}
		logrus.WithError(err).Errorf("GetFlairsByUserID::Failed to get flairs by user ID")
		return []*model.Flair{}, err
	}
	return flairs, nil
}

func (d *db) RemoveFlair(ctx context.Context, flair model.Flair) (*model.Flair, error) {
	logrus.Debug("RemoveFlair::SQL Query")
	// Ensure that flair.ID is set, otherwise GORM could delete all flair
	if &flair.ID == nil {
		logrus.Errorf("Attempted to delete flair with no ID")
	} else if err := d.sql.Delete(&flair).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &flair, nil
		}
		logrus.WithError(err).Errorf("RemoveFlair::Failed to delete flair")
		return &flair, err
	}
	return &flair, nil
}

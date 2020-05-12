package datastore

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *db) UpsertUserDevice(ctx context.Context, userDevice model.UserDevice) (*model.UserDevice, error) {
	logrus.Debugf("UpsertUserDevice::SQL Insert/Update")
	found := model.UserDevice{}
	if err := d.sql.First(&found, model.UserDevice{ID: userDevice.ID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := d.sql.Create(&userDevice).First(&found, model.UserDevice{ID: userDevice.ID}).Error; err != nil {
				logrus.WithError(err).Errorf("UpserUserDevice::Failed to create new object")
				return nil, err
			}
		} else {
			logrus.WithError(err).Errorf("UpsertUserDevice::Failed checking for UserDeviceObject")
			return nil, err
		}
	} else {
		lastSeen := userDevice.LastSeen
		uninitializedTime := time.Time{}
		if lastSeen == uninitializedTime {
			lastSeen = time.Now()
		}
		if err := d.sql.Model(&userDevice).Updates(model.UserDevice{
			Token:    userDevice.Token,
			LastSeen: lastSeen,
			UserID:   userDevice.UserID,
		}).First(&found).Error; err != nil {
			logrus.WithError(err).Errorf("UpsertUserDevice::Failed updating UserDevice object")
			return nil, err
		}
	}
	return &found, nil
}

func (d *db) GetUserDevicesByUserID(ctx context.Context, userID string) ([]model.UserDevice, error) {
	logrus.Debugf("GetUserDevicesByUserID::SQL Query")
	userDevices := []model.UserDevice{}
	if err := d.sql.Where(&model.UserDevice{UserID: &userID}).Order("last_seen desc").Find(&userDevices).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetUserDevicesByUserID::Failed to get user devices by userID")
		return nil, err
	}
	return userDevices, nil
}

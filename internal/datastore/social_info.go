package datastore

import (
	"context"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func (d *delphisDB) UpsertSocialInfo(ctx context.Context, obj model.SocialInfo) (*model.SocialInfo, error) {
	logrus.Debugf("UpsertSocialInfo::SQL Create/Update")
	found := model.SocialInfo{}
	if err := d.sql.First(&found, model.SocialInfo{Network: obj.Network, UserProfileID: obj.UserProfileID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := d.sql.Create(&obj).First(&found, model.SocialInfo{Network: obj.Network, UserProfileID: obj.UserProfileID}).Error; err != nil {
				logrus.WithError(err).Errorf("UpsertSocialInfo::Failed to create new object")
				return nil, err
			}
		} else {
			logrus.WithError(err).Errorf("UpsertSocialInfo::Failed checking for SocialInfo object")
			return nil, err
		}
	} else {
		if err := d.sql.Model(&obj).Updates(model.SocialInfo{
			AccessToken:       obj.AccessToken,
			AccessTokenSecret: obj.AccessTokenSecret,
			UserID:            obj.UserID,
			ProfileImageURL:   obj.ProfileImageURL,
			ScreenName:        obj.ScreenName,
			IsVerified:        obj.IsVerified,
		}).First(&found).Error; err != nil {
			logrus.WithError(err).Errorf("UpsertSocialInfo::Failed updating SocialInfo object")
			return nil, err
		}
	}
	return &found, nil
}

func (d *delphisDB) GetSocialInfosByUserProfileID(ctx context.Context, userProfileID string) ([]model.SocialInfo, error) {
	logrus.Debugf("GetSocialInfosByUserProfileID::SQL Query")
	found := []model.SocialInfo{}
	if err := d.sql.Where(&model.SocialInfo{UserProfileID: userProfileID}).Find(&found).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			logrus.WithError(err).Errorf("GetSocialInfosByUserProfileID::Failed getting social infos")
			return nil, err
		}
	}
	return found, nil
}

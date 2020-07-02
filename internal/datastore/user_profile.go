package datastore

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *delphisDB) GetUserProfileByUserID(ctx context.Context, userID string) (*model.UserProfile, error) {
	logrus.Debugf("GetUserProfileByUserID::SQL Query, userID: %s", userID)
	user := model.User{}
	userProfile := model.UserProfile{}
	if err := d.sql.First(&user, &model.User{ID: userID}).Preload("SocialInfos").Related(&userProfile).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			logrus.Debugf("Value was not found")
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetUserProfileByUserID::Failed to get user profile by user ID")
		return nil, err
	}
	return &userProfile, nil
}

func (d *delphisDB) GetUserProfileByID(ctx context.Context, id string) (*model.UserProfile, error) {
	logrus.Debug("GetUserProfileByID::SQL Query")
	userProfile := model.UserProfile{}
	if err := d.sql.Preload("SocialInfos").First(&userProfile, &model.UserProfile{ID: id}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetUserProfileByID::Failed to get user profile by ID")
		return nil, err
	}

	return &userProfile, nil
}

func (d *delphisDB) CreateOrUpdateUserProfile(ctx context.Context, userProfile model.UserProfile) (*model.UserProfile, bool, error) {
	logrus.Debugf("CreateOrUpdateUserProfile::SQL Insert/Update: %+v", userProfile)
	found := model.UserProfile{}
	if err := d.sql.First(&found, model.UserProfile{ID: userProfile.ID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := d.sql.Preload("SocialInfos").Create(&userProfile).First(&found, model.UserProfile{ID: userProfile.ID}).Error; err != nil {
				logrus.WithError(err).Errorf("CreateOrUpdateUserProfile::Failed creating new object")
				return nil, false, err
			}
			return &userProfile, true, nil
		} else {
			logrus.WithError(err).Errorf("CreateOrUpdateUserProfile::Failed checking for UserProfile")
			return nil, false, err
		}
	} else {
		// Found so this is an update.
		toUpdate := model.UserProfile{
			DisplayName:   userProfile.DisplayName,
			TwitterHandle: userProfile.TwitterHandle,
		}

		// Can't mock this
		if *found.UserID == "" && userProfile.UserID != nil {
			toUpdate.UserID = userProfile.UserID
		}
		if err := d.sql.Preload("SocialInfos").Model(&userProfile).Updates(toUpdate).First(&found).Error; err != nil {
			logrus.WithError(err).Errorf("CreateOrUpdateUserProfile::Failed updating user profile")
			return nil, false, err
		}
		return &found, false, nil
	}
}

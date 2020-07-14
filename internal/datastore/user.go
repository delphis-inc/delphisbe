package datastore

import (
	"context"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func (d *delphisDB) UpsertUser(ctx context.Context, user model.User) (*model.User, error) {
	logrus.Debugf("UpsertUser::SQL Insert/Update")
	found := model.User{}
	if err := d.sql.First(&found, model.User{ID: user.ID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := d.sql.Create(&user).First(&found, model.User{ID: user.ID}).Error; err != nil {
				logrus.WithError(err).Errorf("UpsertUser::Failed to create new object")
				return nil, err
			}
		} else {
			logrus.WithError(err).Errorf("UpsertUser::Failed checking for User object")
			return nil, err
		}
	} else {
		if err := d.sql.Model(&user).Updates(model.User{
			// Nothing should actually update here rn.
		}).First(&found).Error; err != nil {
			logrus.WithError(err).Errorf("UpsertUser::Failed updating user object")
			return nil, err
		}
	}
	return &found, nil
}

func (d *delphisDB) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	logrus.Debug("GetUserByID::SQL Query")
	user := model.User{}
	if err := d.sql.Preload("Participants").Preload("Viewers").Preload("UserProfile").First(&user, &model.User{ID: userID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetUserByID::Failed to get user")
		return nil, err
	}

	return &user, nil
}

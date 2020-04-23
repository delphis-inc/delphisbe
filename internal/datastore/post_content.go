package datastore

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *db) GetPostContentByID(ctx context.Context, id string) (*model.PostContent, error) {
	logrus.Debug("GetPostContentByID::SQL Query")
	found := model.PostContent{}
	if err := d.sql.First(&found, &model.PostContent{ID: id}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("Failed to get PostContent by ID")
		return nil, err
	}
	return &found, nil
}

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

package datastore

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *delphisDB) ListFlairTemplates(ctx context.Context, query *string) ([]*model.FlairTemplate, error) {
	logrus.Debug("ListFlairTemplates::SQL Query")
	var flairTemplates []*model.FlairTemplate
	sql := d.sql
	if query != nil {
		likeQuery := "%" + *query + "%"
		sql = sql.Where("source ILIKE ?", likeQuery).Or("display_name ILIKE ?", likeQuery)
	}
	if err := sql.Find(&flairTemplates).Error; err != nil {
		logrus.WithError(err).Errorf("ListFlairTemplates::Failed to fetch FlairTemplates")
		return nil, err
	}
	return flairTemplates, nil
}

func (d *delphisDB) UpsertFlairTemplate(ctx context.Context, data model.FlairTemplate) (*model.FlairTemplate, error) {
	logrus.Debug("UpsertFlairTemplate::SQL Create or Update")
	flairTemplate := model.FlairTemplate{}
	if err := d.sql.First(&flairTemplate, model.FlairTemplate{ID: data.ID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := d.sql.Create(&data).First(&flairTemplate, model.FlairTemplate{ID: data.ID}).Error; err != nil {
				logrus.WithError(err).Errorf("UpsertFlairTemplate::Failed to create new object")
				return nil, err
			}
		} else {
			logrus.WithError(err).Errorf("UpsertFlairTemplate::Failed checking for FlairTemplate object")
			return nil, err
		}
	} else {
		if err := d.sql.Model(&data).Updates(model.FlairTemplate{
			DisplayName: data.DisplayName,
			ImageURL:    data.ImageURL,
			Source:      data.Source,
		}).First(&flairTemplate).Error; err != nil {
			logrus.WithError(err).Errorf("UpsertFlairTemplate::Failed updating flairTemplate object")
			return nil, err
		}
	}
	return &flairTemplate, nil
}

func (d *delphisDB) GetFlairTemplateByID(ctx context.Context, id string) (*model.FlairTemplate, error) {
	logrus.Debug("GetFlairTemplateByID::SQL Query")
	flairTemplate := model.FlairTemplate{}
	if err := d.sql.First(&flairTemplate, model.FlairTemplate{ID: id}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetFlairTemplateByID::Failed to get flair template")
		return nil, err
	}
	return &flairTemplate, nil
}

func (d *delphisDB) RemoveFlairTemplate(ctx context.Context, flairTemplate model.FlairTemplate) (*model.FlairTemplate, error) {
	logrus.Debug("RemoveFlairTemplate::SQL Query")
	// Ensure that flairTemplate.ID is set, otherwise GORM could delete all flairTemplate
	if &flairTemplate.ID == nil {
		logrus.Errorf("Attempted to delete flair template with no ID")
		return &flairTemplate, nil
	}
	err := d.sql.Transaction(func(tx *gorm.DB) error {
		// Set Null all participants referencing the flairs we just deleted.
		if err := tx.Unscoped().Model(model.Participant{}).
			Joins("FROM flairs").
			Where("flair_id = flairs.id AND flairs.template_id = ?", flairTemplate.ID).
			Update("flair_id", nil).Error; err != nil {
			logrus.WithError(err).Errorf("RemoveFlairTemplate::Failed to unassign related flairs")
			return err
		}

		// Delete all the flairs for this template
		if err := tx.Where(model.Flair{TemplateID: flairTemplate.ID}).
			Delete([]model.Flair{}).Error; err != nil {
			logrus.WithError(err).Errorf("RemoveFlairTemplate::Failed to delete related flairs")
			return err
		}

		// Delete the template itself
		if err := tx.Delete(&flairTemplate).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logrus.WithError(err).Errorf("RemoveFlairTemplate::Failed to delete flair template")
		return &flairTemplate, err
	}
	return &flairTemplate, nil
}

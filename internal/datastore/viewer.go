package datastore

import (
	"context"
	"database/sql"
	"time"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func (d *delphisDB) UpsertViewer(ctx context.Context, viewer model.Viewer) (*model.Viewer, error) {
	logrus.Debug("UpsertViewer::SQL Create/Update")
	found := model.Viewer{}
	if err := d.sql.First(&found, model.Viewer{ID: viewer.ID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := d.sql.Create(&viewer).First(&found, model.Viewer{ID: viewer.ID}).Error; err != nil {
				logrus.WithError(err).Errorf("UpsertViewer::Failed to create new object")
				return nil, err
			}
		} else {
			logrus.WithError(err).Errorf("UpsertViewer::Failed checking for Viewer Object")
			return nil, err
		}
	} else {
		if err := d.sql.Model(&viewer).Updates(model.Viewer{
			LastViewed:       viewer.LastViewed,
			LastViewedPostID: viewer.LastViewedPostID,
		}).First(&found).Error; err != nil {
			logrus.WithError(err).Errorf("UpsertViewer::Failed updating viewer object")
			return nil, err
		}
	}
	return &found, nil
}

func (d *delphisDB) SetViewerLastPostViewed(ctx context.Context, viewerID, postID string, viewedTime time.Time) (*model.Viewer, error) {
	logrus.Debug("SetViewerLastPostViewed::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("SetViewerLastPostViewed::failed to initialize statements")
		return nil, err
	}

	viewer := model.Viewer{}
	if err := d.prepStmts.updateViewerLastViewed.QueryRowContext(
		ctx,
		viewerID,
		viewedTime,
		postID,
	).Scan(
		&viewer.ID,
		&viewer.CreatedAt,
		&viewer.UpdatedAt,
		&viewer.LastViewed,
		&viewer.LastViewedPostID,
		&viewer.DiscussionID,
		&viewer.UserID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("failed to execute SetViewerLastPostViewed")
		return nil, err
	}

	return &viewer, nil
}

func (d *delphisDB) GetViewerForDiscussion(ctx context.Context, discussionID, userID string) (*model.Viewer, error) {
	logrus.Debug("GetViewerForDiscussion::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetViewerForDiscussion::failed to initialize statements")
		return nil, err
	}

	viewer := model.Viewer{}
	if err := d.prepStmts.getViewerForDiscussionIDUserID.QueryRowContext(
		ctx,
		discussionID,
		userID,
	).Scan(
		&viewer.ID,
		&viewer.CreatedAt,
		&viewer.UpdatedAt,
		&viewer.LastViewed,
		&viewer.LastViewedPostID,
		&viewer.DiscussionID,
		&viewer.UserID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("failed to execute GetViewerForDiscussion")
		return nil, err
	}

	return &viewer, nil
}

func (d *delphisDB) GetViewerByID(ctx context.Context, viewerID string) (*model.Viewer, error) {
	logrus.Debug("GetViewerByID::SQL Query")
	found := model.Viewer{}
	if err := d.sql.First(&found, model.Viewer{ID: viewerID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Error("GetViewerByID::Failed to get viewer")
		return nil, err
	}
	return &found, nil
}

func (d *delphisDB) GetViewersByIDs(ctx context.Context, viewerIDs []string) (map[string]*model.Viewer, error) {
	logrus.Debug("GetViewersByIDs::SQL Query")
	viewers := []model.Viewer{}
	if err := d.sql.Where(viewerIDs).Find(&viewers).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// This is a not found situation with multiple ids and I don't know what to do.
			return nil, nil
		} else {
			logrus.WithError(err).Errorf("GetViewersByIDs::Failed to get viewers by IDs")
			return nil, err
		}
	}
	retVal := map[string]*model.Viewer{}
	for _, id := range viewerIDs {
		retVal[id] = nil
	}
	for _, viewer := range viewers {
		retVal[viewer.ID] = &viewer
	}
	return retVal, nil
}

package datastore

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *db) UpsertViewer(ctx context.Context, viewer model.Viewer) (*model.Viewer, error) {
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

func (d *db) GetViewerByID(ctx context.Context, viewerID string) (*model.Viewer, error) {
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

func (d *db) GetViewersByIDs(ctx context.Context, viewerIDs []string) (map[string]*model.Viewer, error) {
	logrus.Debug("GetViewersByIDs::SQL Query")
	viewers := []model.Viewer{}
	if err := d.sql.Where(viewerIDs).Find(&viewers).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// This is a not found situation with multiple ids and I don't know what to do.
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

// func (d *db) GetViewersByIDsDynamo(ctx context.Context, discussionViewerKeys []model.DiscussionViewerKey) (map[model.DiscussionViewerKey]*model.Viewer, error) {
// 	if len(discussionViewerKeys) == 0 {
// 		return map[model.DiscussionViewerKey]*model.Viewer{}, nil
// 	}
// 	logrus.Debug("GetViewersByIDs: DynamoBatchGetItem")
// 	keys := make([]map[string]*dynamodb.AttributeValue, 0)
// 	for _, dv := range discussionViewerKeys {
// 		keys = append(keys, map[string]*dynamodb.AttributeValue{
// 			"DiscussionID": {
// 				S: aws.String(dv.DiscussionID),
// 			},
// 			"ViewerID": {
// 				S: aws.String(dv.ViewerID),
// 			},
// 		})
// 	}
// 	// NOTE: Unless we are fetching from the same discussion we need to use BatchGetItem instead
// 	// of Query here.
// 	res, err := d.dynamo.BatchGetItem(&dynamodb.BatchGetItemInput{
// 		RequestItems: map[string]*dynamodb.KeysAndAttributes{
// 			d.dbConfig.Viewers.TableName: {
// 				Keys: keys,
// 			},
// 		},
// 	})

// 	if err != nil {
// 		logrus.WithError(err).Errorf("GetViewersByIDs: Failed to retrieve viewers with ids: %+v", keys)
// 		return nil, err
// 	}

// 	viewerMap := map[model.DiscussionViewerKey]*model.Viewer{}
// 	for _, dv := range discussionViewerKeys {
// 		viewerMap[dv] = nil
// 	}
// 	elems := res.Responses[d.dbConfig.Viewers.TableName]
// 	for _, elem := range elems {
// 		viewerObj := model.Viewer{}
// 		err := dynamodbattribute.UnmarshalMap(elem, &viewerObj)
// 		if err != nil {
// 			logrus.WithError(err).Warnf("Failed to unmarshal viewer object: %+v", elem)
// 			continue
// 		}

// 		viewerMap[viewerObj.DiscussionViewerKey()] = &viewerObj
// 	}

// 	return viewerMap, nil
// }

// func (d *db) GetViewerByIDDynamo(ctx context.Context, discussionViewerKey model.DiscussionViewerKey) (*model.Viewer, error) {
// 	viewers, err := d.GetViewersByIDs(ctx, []model.DiscussionViewerKey{discussionViewerKey})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return viewers[discussionViewerKey], nil
// }

func (d *db) PutViewerDynamo(ctx context.Context, viewer model.Viewer) (*model.Viewer, error) {
	logrus.Debug("PutViewer::Dynamo PutItem")
	av, err := d.marshalMap(viewer)
	if err != nil {
		logrus.WithError(err).Errorf("PutViewer: Failed to marshal viewer object: %+v", viewer)
		return nil, err
	}

	_, err = d.dynamo.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.dbConfig.Viewers.TableName),
		Item:      av,
	})

	if err != nil {
		logrus.WithError(err).Errorf("PutViewer: Failed to put viewer object: %+v", av)
		return nil, err
	}

	return &viewer, nil
}

package datastore

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *db) PutViewer(ctx context.Context, viewer model.Viewer) (*model.Viewer, error) {
	logrus.Debug("PutViewer::Dynamo PutItem")
	av, err := dynamodbattribute.MarshalMap(viewer)
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

func (d *db) GetViewersByIDs(ctx context.Context, discussionViewers []model.DiscussionViewer) (map[model.DiscussionViewer]*model.Viewer, error) {
	logrus.Debug("GetViewersByIDs: DynamoBatchGetItem")
	keys := make([]map[string]*dynamodb.AttributeValue, 0)
	for _, dv := range discussionViewers {
		keys = append(keys, map[string]*dynamodb.AttributeValue{
			"DiscussionID": {
				S: aws.String(dv.DiscussionID),
			},
			"ViewerID": {
				S: aws.String(dv.ViewerID),
			},
		})
	}
	res, err := d.dynamo.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			d.dbConfig.Viewers.TableName: {
				Keys: keys,
			},
		},
	})

	if err != nil {
		logrus.WithError(err).Errorf("GetViewersByIDs: Failed to retrieve viewers with ids: %+v", keys)
		return nil, err
	}

	viewerMap := map[model.DiscussionViewer]*model.Viewer{}
	for _, dv := range discussionViewers {
		viewerMap[dv] = nil
	}
	elems := res.Responses[d.dbConfig.Viewers.TableName]
	for _, elem := range elems {
		viewerObj := model.Viewer{}
		err := dynamodbattribute.UnmarshalMap(elem, &viewerObj)
		if err != nil {
			logrus.WithError(err).Warnf("Failed to unmarshal viewer object: %+v", elem)
			continue
		}

		viewerMap[viewerObj.DiscussionViewerKey()] = &viewerObj
	}

	return viewerMap, nil
}

func (d *db) GetViewerByID(ctx context.Context, discussionViewer model.DiscussionViewer) (*model.Viewer, error) {
	viewers, err := d.GetViewersByIDs(ctx, []model.DiscussionViewer{discussionViewer})

	if err != nil {
		return nil, err
	}

	return viewers[discussionViewer], nil
}

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

package datastore

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *db) GetUserProfileByID(ctx context.Context, id string) (*model.UserProfile, error) {
	logrus.Debug("GetUserProfileByID: Dynamo GetItem")
	res, err := d.dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.dbConfig.UserProfiles.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		logrus.WithError(err).Errorf("GetUserProfileByID: Error getting user profile with id %s", id)
		return nil, err
	}

	userProfile := model.UserProfile{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &userProfile)

	if err != nil {
		logrus.WithError(err).Errorf("GetUserProfileByID: Failed unmarshaling user profile object %+v", res.Item)
		return nil, err
	}

	return &userProfile, nil
}

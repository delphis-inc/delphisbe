package datastore

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *db) AddParticipantToUser(ctx context.Context, userID string, discussionParticipantKey model.DiscussionParticipantKey) (*model.User, error) {
	logrus.Debug("AddParticipantToUser: Dynamo Update")
	av, err := dynamodbattribute.Marshal(discussionParticipantKey)
	if err != nil {
		logrus.WithError(err).Errorf("AddParticipantToUser: Failed to marshal value: %+v", discussionParticipantKey)
		return nil, err
	}
	res, err := d.dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		TableName:        aws.String(d.dbConfig.Users.TableName),
		UpdateExpression: aws.String("ADD DiscussionParticipants :pids"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userID),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pids": {
				SS: []*string{av.S},
			},
		},
		ReturnValues: aws.String("UPDATED_VALUES"),
	})

	if err != nil {
		logrus.WithError(err).Errorf("Failed updating user object")
		return nil, err
	}

	userObj := model.User{}

	err = dynamodbattribute.UnmarshalMap(res.Attributes, &userObj)

	if err != nil {
		logrus.WithError(err).Errorf("Failed unmarshaling returned value: %+v", res.Attributes)
		return nil, err
	}

	return &userObj, err
}

func (d *db) AddViewerToUser(ctx context.Context, userID string, discussionViewerKey model.DiscussionViewerKey) (*model.User, error) {
	logrus.Debug("AddViewerToUser: Dynamo Update")
	av, err := dynamodbattribute.Marshal(discussionViewerKey)
	if err != nil {
		logrus.WithError(err).Errorf("AddViewerToUser: Failed to marshal value: %+v", discussionViewerKey)
		return nil, err
	}
	res, err := d.dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		TableName:        aws.String(d.dbConfig.Users.TableName),
		UpdateExpression: aws.String("ADD DiscussionViewers :vids"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userID),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":vids": {
				SS: []*string{av.S},
			},
		},
		ReturnValues: aws.String("UPDATED_VALUES"),
	})

	if err != nil {
		logrus.WithError(err).Errorf("Failed updating user object")
		return nil, err
	}

	userObj := model.User{}
	err = dynamodbattribute.UnmarshalMap(res.Attributes, &userObj)

	if err != nil {
		logrus.WithError(err).Errorf("Failed to unmarshal response value: %+v", res.Attributes)
		return nil, err
	}

	return &userObj, err
}

func (d *db) PutUser(ctx context.Context, user model.User) (*model.User, error) {
	logrus.Debug("PutUser::Dynamo PutItem")
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		logrus.WithError(err).Errorf("PutUser: Failed to marshal user object: %+v", user)
		return nil, err
	}

	_, err = d.dynamo.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.dbConfig.Users.TableName),
		Item:      av,
	})

	if err != nil {
		logrus.WithError(err).Errorf("PutUser: Failed to put user object: %+v", av)
		return nil, err
	}

	return &user, nil
}

func (d *db) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	logrus.Debug("GetUserByID: Dynamo GetItem")
	res, err := d.dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.dbConfig.Users.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userID),
			},
		},
	})

	if err != nil {
		logrus.WithError(err).Errorf("GetUserByID: Failed to get user with ID: %s", userID)
		return nil, err
	}

	user := model.User{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &user)

	if err != nil {
		logrus.WithError(err).Errorf("GetUserByID: Failed to unmarshal user object: %+v", res.Item)
		return nil, err
	}

	return &user, nil
}

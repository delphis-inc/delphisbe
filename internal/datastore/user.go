package datastore

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *db) AddParticipantToUser(ctx context.Context, userID, participantID string) error {
	logrus.Debug("AddParticipantToUser: Dynamo Update")
	_, err := d.dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		TableName:        aws.String(d.dbConfig.Users.TableName),
		UpdateExpression: aws.String("ADD ParticipantIDs :pids"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userID),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pids": {
				SS: []*string{aws.String(participantID)},
			},
		},
	})

	return err
}

func (d *db) AddViewerToUser(ctx context.Context, userID, viewerID string) error {
	logrus.Debug("AddParticipantToUser: Dynamo Update")
	_, err := d.dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		TableName:        aws.String(d.dbConfig.Users.TableName),
		UpdateExpression: aws.String("ADD ViewerIDs :vids"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userID),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":vids": {
				SS: []*string{aws.String(viewerID)},
			},
		},
	})

	return err
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

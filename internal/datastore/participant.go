package datastore

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *db) GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error) {
	logrus.Debug("GetParticipantsByDiscussionID::Dynamo Query")
	res, err := d.dynamo.Query(&dynamodb.QueryInput{
		TableName: aws.String(d.dbConfig.Participants.TableName),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":did": {
				S: aws.String(id),
			},
		},
		KeyConditionExpression: aws.String("DiscussionID = :did"),
	})

	if err != nil {
		logrus.WithError(err).Errorf("GetParticipantsByDiscussionID: Failed to query participants for discussionID: %s", id)
		return nil, err
	}

	participants := make([]model.Participant, 0)
	if res != nil {
		for _, elem := range res.Items {
			participantObj := model.Participant{}
			err := dynamodbattribute.UnmarshalMap(elem, &participantObj)
			if err != nil {
				logrus.WithError(err).Warnf("GetParticipantsByDiscussionID: Failed unmarshaling participant object: %+v", elem)
				continue
			}
			participants = append(participants, participantObj)
		}
	}
	return participants, nil
}

func (d *db) PutParticipant(ctx context.Context, participant model.Participant) (*model.Participant, error) {
	logrus.Debug("PutParticipant::Dynamo PutItem")
	av, err := dynamodbattribute.MarshalMap(participant)
	if err != nil {
		logrus.WithError(err).Errorf("PutParticipant: Failed to marshal participant object: %+v", participant)
		return nil, err
	}

	_, err = d.dynamo.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.dbConfig.Participants.TableName),
		Item:      av,
	})

	if err != nil {
		logrus.WithError(err).Errorf("PutParticipant: Failed to put participant object: %+v", av)
		return nil, err
	}

	return &participant, nil
}

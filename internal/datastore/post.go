package datastore

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *db) PutPost(ctx context.Context, post model.Post) (*model.Post, error) {
	logrus.Debug("PutPost: DynamoDB PutItem")
	av, err := dynamodbattribute.MarshalMap(post)
	if err != nil {
		logrus.WithError(err).Errorf("PutPost: Failed to marshal post object: %+v", post)
		return nil, err
	}
	_, err = d.dynamo.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.dbConfig.Posts.TableName),
		Item:      av,
	})

	if err != nil {
		logrus.WithError(err).Errorf("PutPost: Failed to put post object: %+v", av)
		return nil, err
	}

	return &post, nil
}

func (d *db) GetPostsByDiscussionID(ctx context.Context, discussionID string) ([]*model.Post, error) {
	logrus.Debug("GetPostsByDiscussionID: DynamoDB Query")
	res, err := d.dynamo.Query(&dynamodb.QueryInput{
		TableName: aws.String(d.dbConfig.Posts.TableName),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":did": {
				S: aws.String(discussionID),
			},
		},
		KeyConditionExpression: aws.String("DiscussionID = :did"),
	})

	if err != nil {
		logrus.WithError(err).Errorf("GetPostsByDiscussionID: Failed to query dynamo for discussionID: %s", discussionID)
		return nil, err
	}

	postObjs := []*model.Post{}
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &postObjs)

	if err != nil {
		logrus.WithError(err).Errorf("Failed to unmarshal response values: %+v", res.Items)
		return nil, err
	}

	return postObjs, nil
}

package datastore

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *db) PutPost(ctx context.Context, post model.Post) (*model.Post, error) {
	logrus.Debug("PutPost::SQL Create")
	if err := d.sql.Create(&post).Error; err != nil {
		logrus.WithError(err).Errorf("Failed to create a post")
		return nil, err
	}

	return &post, nil
}

func (d *db) GetPostsByDiscussionID(ctx context.Context, discussionID string) ([]*model.Post, error) {
	logrus.Debug("GetPostsByDiscussionID::SQL Query")
	posts := []model.Post{}
	if err := d.sql.Where(model.Post{DiscussionID: &discussionID}).Find(&posts).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// Not sure if this will return not found error... If the discussion is empty maybe?
			// Should this be nil, nil?
			return []*model.Post{}, nil
		}
		logrus.WithError(err).Errorf("Failed to get posts by discussionID")
		return nil, err
	}

	returnedPosts := []*model.Post{}
	for _, p := range posts {
		returnedPosts = append(returnedPosts, &p)
	}
	return returnedPosts, nil
}

///////////////
// Dynamo functions
///////////////

func (d *db) PutPostDynamo(ctx context.Context, post model.Post) (*model.Post, error) {
	logrus.Debug("PutPost: DynamoDB PutItem")
	av, err := d.marshalMap(post)
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

func (d *db) GetPostsByDiscussionIDDynamo(ctx context.Context, discussionID string) ([]*model.Post, error) {
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

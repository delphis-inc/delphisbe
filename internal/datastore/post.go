package datastore

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *delphisDB) PutPost(ctx context.Context, post model.Post) (*model.Post, error) {
	logrus.Debug("PutPost::SQL Create")
	logrus.Infof("PutPost: QuotePost: %v\n", post.QuotedPostID)
	found := model.Post{}
	if err := d.sql.Create(&post).First(&found, model.Post{ID: post.ID}).Error; err != nil {
		logrus.WithError(err).Errorf("Failed to create a post")
		return nil, err
	}

	return &found, nil
}

func (d *delphisDB) GetPostsByDiscussionID(ctx context.Context, discussionID string) ([]*model.Post, error) {
	logrus.Debug("GetPostsByDiscussionID::SQL Query")
	posts := []model.Post{}
	if err := d.sql.Where(model.Post{DiscussionID: &discussionID}).Preload("PostContent").Find(&posts).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// Not sure if this will return not found error... If the discussion is empty maybe?
			// Should this be nil, nil?
			return []*model.Post{}, nil
		}
		logrus.WithError(err).Errorf("Failed to get posts by discussionID")
		return nil, err
	}

	logrus.Debugf("Found posts: %+v", posts)

	returnedPosts := []*model.Post{}
	for i := range posts {
		if posts[i].QuotedPostID != nil {
			var err error
			posts[i].QuotedPost, err = d.getPostByID(ctx, *posts[i].QuotedPostID)
			if err != nil {
				// Do we want to fail the whole discussion if we can't get a quote?
				return nil, err
			}
		}
		returnedPosts = append(returnedPosts, &posts[i])

	}

	logrus.Infof("Posts: %+v\n", returnedPosts)
	return returnedPosts, nil
}

func (d *delphisDB) getPostByID(ctx context.Context, postID string) (*model.Post, error) {
	logrus.Debug("GetPostByID::SQL Query")
	post := model.Post{}
	if err := d.sql.Where([]string{postID}).Preload("PostContent").Find(&post).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("Failed to get Post by ID")
		return nil, err
	}

	return &post, nil
}

///////////////
// Dynamo functions
///////////////

// func (d *db) PutPostDynamo(ctx context.Context, post model.Post) (*model.Post, error) {
// 	logrus.Debug("PutPost: DynamoDB PutItem")
// 	av, err := d.marshalMap(post)
// 	if err != nil {
// 		logrus.WithError(err).Errorf("PutPost: Failed to marshal post object: %+v", post)
// 		return nil, err
// 	}
// 	_, err = d.dynamo.PutItem(&dynamodb.PutItemInput{
// 		TableName: aws.String(d.dbConfig.Posts.TableName),
// 		Item:      av,
// 	})

// 	if err != nil {
// 		logrus.WithError(err).Errorf("PutPost: Failed to put post object: %+v", av)
// 		return nil, err
// 	}

// 	return &post, nil
// }

// func (d *db) GetPostsByDiscussionIDDynamo(ctx context.Context, discussionID string) ([]*model.Post, error) {
// 	logrus.Debug("GetPostsByDiscussionID: DynamoDB Query")
// 	res, err := d.dynamo.Query(&dynamodb.QueryInput{
// 		TableName: aws.String(d.dbConfig.Posts.TableName),
// 		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
// 			":did": {
// 				S: aws.String(discussionID),
// 			},
// 		},
// 		KeyConditionExpression: aws.String("DiscussionID = :did"),
// 	})

// 	if err != nil {
// 		logrus.WithError(err).Errorf("GetPostsByDiscussionID: Failed to query dynamo for discussionID: %s", discussionID)
// 		return nil, err
// 	}

// 	postObjs := []*model.Post{}
// 	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &postObjs)

// 	if err != nil {
// 		logrus.WithError(err).Errorf("Failed to unmarshal response values: %+v", res.Items)
// 		return nil, err
// 	}

// 	return postObjs, nil
// }

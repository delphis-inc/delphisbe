package datastore

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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

func (d *db) UpdateUserProfileTwitterInfo(ctx context.Context, userProfile model.UserProfile, twitterInfo model.SocialInfo) (*model.UserProfile, error) {
	logrus.Debug("UpdateUserProfileTwitterInfo: Dynamo UpdateItem")
	av, err := d.marshalMap(twitterInfo)
	if err != nil {
		logrus.WithError(err).Errorf("UpdateUserProfileTwitterInfo: Failed to marshal twitterInfo: %+v", twitterInfo)
		return nil, err
	}
	res, err := d.dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#TI": aws.String("TwitterInfo"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				M: av,
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userProfile.ID),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		TableName:        aws.String(d.dbConfig.UserProfiles.TableName),
		UpdateExpression: aws.String("SET #TI = :t"),
	})

	if err != nil {
		logrus.WithError(err).Errorf("Failed to update userProfile with ID (%s) with new values <redacted>", userProfile.ID)
		return nil, err
	}

	returnedProfile := model.UserProfile{}
	err = dynamodbattribute.UnmarshalMap(res.Attributes, &returnedProfile)

	if err != nil {
		logrus.WithError(err).Errorf("Failed unmarshaling return value <redacted>")
		return nil, err
	}

	return &returnedProfile, nil
}

func (d *db) UpdateUserProfileUserID(ctx context.Context, userProfileID string, userID string) (*model.UserProfile, error) {
	logrus.Debug("UpdateUserProfileUserID: Dynamo UpdateItem")
	res, err := d.dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#UI": aws.String("UserID"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":ui": {
				S: aws.String(userID),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userProfileID),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		TableName:        aws.String(d.dbConfig.UserProfiles.TableName),
		UpdateExpression: aws.String("SET #UI = :ui"),
	})

	if err != nil {
		logrus.WithError(err).Errorf("UpdateUserProfileUserID: Failed to update userProfileID (%s) with userID (%s)", userProfileID, userID)
		return nil, err
	}

	userProfile := model.UserProfile{}
	err = dynamodbattribute.UnmarshalMap(res.Attributes, &userProfile)
	if err != nil {
		logrus.WithError(err).Errorf("UpdateUserProfileUserID: Failed to unmarshal return value: %+v", res.Attributes)
		return nil, err
	}
	return &userProfile, nil
}

func (d *db) CreateOrUpdateUserProfile(ctx context.Context, userProfile model.UserProfile) (*model.UserProfile, bool, error) {
	logrus.Debug("PutUserProfileIfNotExists: Dynamo PutItem")
	av, err := d.marshalMap(userProfile)
	if err != nil {
		logrus.WithError(err).Errorf("PutUserProfileIfNotExists: Failed to marshal UserProfile object: %+v", userProfile)
		return nil, false, err
	}
	_, err = d.dynamo.PutItem(&dynamodb.PutItemInput{
		TableName:           aws.String(d.dbConfig.UserProfiles.TableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(ID)"),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
				// Need to just return the existing object
				existing, err := d.GetUserProfileByID(ctx, userProfile.ID)
				if err != nil {
					logrus.Errorf("PutUserProfileIfNotExists: Failed to get user profile to update with ID (%s)", userProfile.ID)
					return nil, false, err
				}
				existing.DisplayName, existing.TwitterHandle, existing.TwitterInfo = userProfile.DisplayName, userProfile.TwitterHandle, userProfile.TwitterInfo
				updated, err := d.UpdateUserProfileTwitterInfo(ctx, *existing, userProfile.TwitterInfo)
				return updated, false, err
			}
		}
		logrus.WithError(err).Errorf("PutUserProfileIfNotExists: Failed putting UserProfile: %+v", userProfile)
		return nil, false, err
	}

	return &userProfile, true, nil
}

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

func (d *db) GetUserProfileByUserID(ctx context.Context, userID string) (*model.UserProfile, error) {
	logrus.Debugf("GetUserProfileByUserID::SQL Query, userID: %s", userID)
	user := model.User{}
	userProfile := model.UserProfile{}
	if err := d.sql.First(&user, &model.User{ID: userID}).Preload("SocialInfos").Related(&userProfile).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			logrus.Debugf("Value was not found")
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetUserProfileByUserID::Failed to get user profile by user ID")
		return nil, err
	}
	return &userProfile, nil
}

func (d *db) GetUserProfileByID(ctx context.Context, id string) (*model.UserProfile, error) {
	logrus.Debug("GetUserProfileByID::SQL Query")
	userProfile := model.UserProfile{}
	if err := d.sql.Preload("SocialInfos").First(&userProfile, &model.UserProfile{ID: id}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetUserProfileByID::Failed to get user profile by ID")
		return nil, err
	}

	return &userProfile, nil
}

func (d *db) CreateOrUpdateUserProfile(ctx context.Context, userProfile model.UserProfile) (*model.UserProfile, bool, error) {
	logrus.Debugf("CreateOrUpdateUserProfile::SQL Insert/Update: %+v", userProfile)
	found := model.UserProfile{}
	if err := d.sql.First(&found, model.UserProfile{ID: userProfile.ID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := d.sql.Preload("SocialInfos").Create(&userProfile).First(&found, model.UserProfile{ID: userProfile.ID}).Error; err != nil {
				logrus.WithError(err).Errorf("CreateOrUpdateUserProfile::Failed creating new object")
				return nil, false, err
			}
			return &userProfile, true, nil
		} else {
			logrus.WithError(err).Errorf("CreateOrUpdateUserProfile::Failed checking for UserProfile")
			return nil, false, err
		}
	} else {
		// Found so this is an update.
		toUpdate := model.UserProfile{
			DisplayName:   userProfile.DisplayName,
			TwitterHandle: userProfile.TwitterHandle,
		}
		if found.UserID == nil && userProfile.UserID != nil {
			toUpdate.UserID = userProfile.UserID
		}
		if err := d.sql.Preload("SocialInfos").Model(&userProfile).Updates(toUpdate).First(&found).Error; err != nil {
			logrus.WithError(err).Errorf("CreateOrUpdateUserProfile::Failed updating user profile")
			return nil, false, err
		}
		return &found, false, nil
	}
}

func (d *db) AddModeratedDiscussionToUserProfileDynamo(ctx context.Context, userProfileID string, discussionID string) (*model.UserProfile, error) {
	logrus.Debug("AddModeratedDiscussionToUserProfile: Dynamo UpdateItem")
	res, err := d.dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#MD": aws.String("ModeratedDiscussionIDs"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":dids": {
				SS: []*string{aws.String(discussionID)},
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userProfileID),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		TableName:        aws.String(d.dbConfig.UserProfiles.TableName),
		UpdateExpression: aws.String("SET #MD = :dids"),
	})

	if err != nil {
		logrus.WithError(err).Errorf("Failed updating user profile object")
		return nil, err
	}

	userProfileObj := model.UserProfile{}

	err = dynamodbattribute.UnmarshalMap(res.Attributes, &userProfileObj)

	if err != nil {
		logrus.WithError(err).Errorf("Failed unmarshaling returned value: %+v", res.Attributes)
		return nil, err
	}

	return &userProfileObj, nil
}

func (d *db) GetUserProfileByIDDynamo(ctx context.Context, id string) (*model.UserProfile, error) {
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

// func (d *db) CreateOrUpdateUserProfileDynamo(ctx context.Context, userProfile model.UserProfile) (*model.UserProfile, bool, error) {
// 	logrus.Debug("PutUserProfileIfNotExists: Dynamo PutItem")
// 	av, err := d.marshalMap(userProfile)
// 	if err != nil {
// 		logrus.WithError(err).Errorf("PutUserProfileIfNotExists: Failed to marshal UserProfile object: %+v", userProfile)
// 		return nil, false, err
// 	}
// 	fmt.Printf("%+v\n", av)
// 	_, err = d.dynamo.PutItem(&dynamodb.PutItemInput{
// 		TableName:           aws.String(d.dbConfig.UserProfiles.TableName),
// 		Item:                av,
// 		ConditionExpression: aws.String("attribute_not_exists(ID)"),
// 	})

// 	if err != nil {
// 		if aerr, ok := err.(awserr.Error); ok {
// 			if aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
// 				// Need to just return the existing object
// 				existing, err := d.GetUserProfileByID(ctx, userProfile.ID)
// 				if err != nil {
// 					logrus.Errorf("PutUserProfileIfNotExists: Failed to get user profile to update with ID (%s)", userProfile.ID)
// 					return nil, false, err
// 				}
// 				existing.DisplayName, existing.TwitterHandle, existing.TwitterInfo = userProfile.DisplayName, userProfile.TwitterHandle, userProfile.TwitterInfo
// 				//updated, err := d.UpdateUserProfileTwitterInfo(ctx, *existing, userProfile.TwitterInfo)
// 				return nil, false, err
// 			}
// 		}
// 		logrus.WithError(err).Errorf("PutUserProfileIfNotExists: Failed putting UserProfile: %+v", userProfile)
// 		return nil, false, err
// 	}

// 	return &userProfile, true, nil
// }

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

func (d *db) UpsertUser(ctx context.Context, user model.User) (*model.User, error) {
	logrus.Debugf("UpsertUser::SQL Insert/Update")
	found := model.User{}
	if err := d.sql.First(&found, model.User{ID: user.ID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := d.sql.Create(&user).First(&found, model.User{ID: user.ID}).Error; err != nil {
				logrus.WithError(err).Errorf("UpsertUser::Failed to create new object")
				return nil, err
			}
		} else {
			logrus.WithError(err).Errorf("UpsertUser::Failed checking for User object")
			return nil, err
		}
	} else {
		if err := d.sql.Model(&user).Updates(model.User{
			// Nothing should actually update here rn.
		}).First(&found).Error; err != nil {
			logrus.WithError(err).Errorf("UpsertUser::Failed updating user object")
			return nil, err
		}
	}
	return &found, nil
}

func (d *db) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	logrus.Debug("GetUserByID::SQL Query")
	user := model.User{}
	if err := d.sql.Preload("Participants").Preload("Viewers").Preload("UserProfile").First(&user, &model.User{ID: userID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetUserByID::Failed to get user")
		return nil, err
	}
	logrus.Debugf("Found: %+v", user)

	return &user, nil
}

func (d *db) GetUserByIDDynamo(ctx context.Context, userID string) (*model.User, error) {
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

func (d *db) PutUserDynamo(ctx context.Context, user model.User) (*model.User, error) {
	logrus.Debug("PutUser::Dynamo PutItem")
	av, err := d.marshalMap(user)
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

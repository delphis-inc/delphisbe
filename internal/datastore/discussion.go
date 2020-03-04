package datastore

import (
	"context"

	"github.com/nedrocks/delphisbe/graph/model"
)

func (d *db) GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error) {
	// res, err := d.dynamo.GetItem(&dynamodb.GetItemInput{
	// 	TableName: aws.String(d.dbConfig.Discussions.TableName),
	// 	Key: map[string]*dynamodb.AttributeValue{
	// 		"ID": {
	// 			S: aws.String(id),
	// 		},
	// 	},
	// })

	// if err != nil {
	// 	logrus.WithError(err).Infof("Failed getting discussion by ID (%s)", id)
	// 	return nil, err
	// }

	// discussionObj := model.Discussion{}
	// err = dynamodbattribute.UnmarshalMap(res.Item, &discussionObj)

	// if err != nil {
	// 	logrus.WithError(err).Infof("Failed unmarshalling discussion by ID (%s)", id)
	// 	return nil, err
	// }

	// return &discussionObj, nil
	return &model.Discussion{
		ID: "12345",
	}, nil
}

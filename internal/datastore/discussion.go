package datastore

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nedrocks/delphisbe/graph/model"
)

const (
	tableName = "Discussions"
)

func (d *db) GetDiscussionByID(id string) (*model.Discussion, error) {
	d.dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: tableName,
	})
}

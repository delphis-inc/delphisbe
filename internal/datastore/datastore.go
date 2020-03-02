package datastore

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nedrocks/delphisbe/graph/model"
)

type Datastore interface {
	GetDiscussionByID(id string) (*model.Discussion, error)
}

type db struct {
	dynamo *dynamodb.DynamoDB
}

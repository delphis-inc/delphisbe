package datastore

import (
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/config"
)

type Datastore interface {
	GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error)
}

type db struct {
	dynamo   *dynamodb.DynamoDB
	dbConfig config.TablesConfig
}

func NewDatastore(dbConfig config.TablesConfig) Datastore {
	return &db{
		dbConfig: dbConfig,
	}
}

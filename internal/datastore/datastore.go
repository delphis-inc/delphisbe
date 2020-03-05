package datastore

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/sirupsen/logrus"
)

type Datastore interface {
	GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error)
	ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error)
	PutDiscussion(ctx context.Context, discussion model.Discussion) (*model.Discussion, error)
}

type db struct {
	dynamo   *dynamodb.DynamoDB
	dbConfig config.TablesConfig
}

func NewDatastore(dbConfig config.DBConfig) Datastore {
	creds := credentials.NewStaticCredentials("fakeMyKeyId", "fakeSecretAccessKey", "")
	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(dbConfig.Region),
		Endpoint:    aws.String(fmt.Sprintf("%s:%d", dbConfig.Host, dbConfig.Port)),
	})

	if err != nil {
		logrus.Println(err)
	}
	dbSvc := dynamodb.New(sess)
	return &db{
		dbConfig: dbConfig.TablesConfig,
		dynamo:   dbSvc,
	}
}

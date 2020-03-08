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
	GetDiscussionsByIDs(ctx context.Context, ids []string) (map[string]*model.Discussion, error)
	ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error)
	PutDiscussion(ctx context.Context, discussion model.Discussion) (*model.Discussion, error)
	GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error)
	GetParticipantsByIDs(ctx context.Context, discussionParticipantKeys []model.DiscussionParticipantKey) (map[model.DiscussionParticipantKey]*model.Participant, error)
	PutParticipant(ctx context.Context, participant model.Participant) (*model.Participant, error)
	AddParticipantToUser(ctx context.Context, userID string, discussionParticipant model.DiscussionParticipantKey) (*model.User, error)
	GetPostsByDiscussionID(ctx context.Context, discussionID string) ([]*model.Post, error)
	PutPost(ctx context.Context, post model.Post) (*model.Post, error)
	AddViewerToUser(ctx context.Context, userID string, discussionViewerKey model.DiscussionViewerKey) (*model.User, error)
	GetUserProfileByID(ctx context.Context, id string) (*model.UserProfile, error)
	PutUser(ctx context.Context, user model.User) (*model.User, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetViewersByIDs(ctx context.Context, discussionViewerKeys []model.DiscussionViewerKey) (map[model.DiscussionViewerKey]*model.Viewer, error)
	PutViewer(ctx context.Context, viewer model.Viewer) (*model.Viewer, error)
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

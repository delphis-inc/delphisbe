package datastore

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Datastore interface {
	GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error)
	GetDiscussionsByIDs(ctx context.Context, ids []string) (map[string]*model.Discussion, error)
	GetDiscussionByModeratorID(ctx context.Context, moderatorID string) (*model.Discussion, error)
	CreateModerator(ctx context.Context, moderator model.Moderator) (*model.Moderator, error)
	GetModeratorByID(ctx context.Context, id string) (*model.Moderator, error)
	ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error)
	UpsertDiscussion(ctx context.Context, discussion model.Discussion) (*model.Discussion, error)
	GetParticipantByID(ctx context.Context, participantID string) (*model.Participant, error)
	GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error)
	PutParticipant(ctx context.Context, participant model.Participant) (*model.Participant, error)
	GetPostsByDiscussionID(ctx context.Context, discussionID string) ([]*model.Post, error)
	GetPostContentByID(ctx context.Context, id string) (*model.PostContent, error)
	PutPost(ctx context.Context, post model.Post) (*model.Post, error)
	GetUserProfileByID(ctx context.Context, id string) (*model.UserProfile, error)
	GetUserProfileByUserID(ctx context.Context, userID string) (*model.UserProfile, error)
	GetSocialInfosByUserProfileID(ctx context.Context, userProfileID string) ([]model.SocialInfo, error)
	UpsertSocialInfo(ctx context.Context, obj model.SocialInfo) (*model.SocialInfo, error)
	CreateOrUpdateUserProfile(ctx context.Context, userProfile model.UserProfile) (*model.UserProfile, bool, error)
	UpsertUser(ctx context.Context, user model.User) (*model.User, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetViewersByIDs(ctx context.Context, viewerIDs []string) (map[string]*model.Viewer, error)
	UpsertViewer(ctx context.Context, viewer model.Viewer) (*model.Viewer, error)

	marshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error)
}

type db struct {
	dynamo   dynamodbiface.DynamoDBAPI
	sql      *gorm.DB
	dbConfig config.TablesConfig
	encoder  *dynamodbattribute.Encoder
}

//MarshalMap wraps the dynamodbattribute.MarshalMap with a defined encoder.
func (d *db) marshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	av, err := d.encoder.Encode(in)
	if err != nil || av == nil || av.M == nil {
		return map[string]*dynamodb.AttributeValue{}, err
	}

	return av.M, nil
}

func NewDatastore(config config.Config, awsSession *session.Session) Datastore {
	mySession := awsSession
	dbConfig := config.DBConfig
	if dbConfig.Host != "" && dbConfig.Port != 0 {
		mySession = mySession.Copy(awsSession.Config.WithEndpoint(fmt.Sprintf("%s:%d", dbConfig.Host, dbConfig.Port)))
		logrus.Debugf("endpoint: %v", mySession.Config.Endpoint)
	}
	dbSvc := dynamodb.New(mySession)
	return &db{
		dbConfig: dbConfig.TablesConfig,
		sql:      NewSQLDatastore(config.SQLDBConfig, awsSession),
		dynamo:   dbSvc,
		encoder: &dynamodbattribute.Encoder{
			MarshalOptions: dynamodbattribute.MarshalOptions{
				SupportJSONTags: false,
			},
			NullEmptyString: true,
		},
	}
}

func NewSQLDatastore(sqlDbConfig config.SQLDBConfig, awsSession *session.Session) *gorm.DB {
	dbURI := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", sqlDbConfig.Host, sqlDbConfig.Port, sqlDbConfig.Username, sqlDbConfig.DBName, sqlDbConfig.Password)
	logrus.Debugf("About to open connection to DB")
	db, err := gorm.Open("postgres", dbURI)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to open db")
		return nil
	}
	logrus.Debugf("Opened connection to DB!")

	// Set autoload
	//db = db.Set("gorm:auto_preload", true)
	//db = db.LogMode(true)
	// need to defer closing the db.
	//db.AutoMigrate(model.DatabaseModels...)

	return db
}

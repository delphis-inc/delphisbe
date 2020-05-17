package datastore

import (
	"context"
	sql2 "database/sql"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/pkg/errors"
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
	AssignFlair(ctx context.Context, participant model.Participant, flairID *string) (*model.Participant, error)
	GetFlairByID(ctx context.Context, id string) (*model.Flair, error)
	GetFlairsByUserID(ctx context.Context, userID string) ([]*model.Flair, error)
	RemoveFlair(ctx context.Context, flair model.Flair) (*model.Flair, error)
	UpsertFlair(ctx context.Context, flair model.Flair) (*model.Flair, error)
	ListFlairTemplates(ctx context.Context, query *string) ([]*model.FlairTemplate, error)
	GetFlairTemplateByID(ctx context.Context, id string) (*model.FlairTemplate, error)
	UpsertFlairTemplate(ctx context.Context, flairTemplate model.FlairTemplate) (*model.FlairTemplate, error)
	RemoveFlairTemplate(ctx context.Context, flairTemplate model.FlairTemplate) (*model.FlairTemplate, error)
	GetTotalParticipantCountByDiscussionID(ctx context.Context, discussionID string) int
	GetParticipantByID(ctx context.Context, participantID string) (*model.Participant, error)
	GetParticipantsByIDs(ctx context.Context, ids []string) ([]*model.Participant, error)
	GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error)
	GetParticipantsByDiscussionIDUserID(ctx context.Context, discussionID string, userID string) ([]model.Participant, error)
	UpsertParticipant(ctx context.Context, participant model.Participant) (*model.Participant, error)
	GetPostsByDiscussionID(ctx context.Context, discussionID string) ([]*model.Post, error)
	GetPostsByDiscussionIDIter(ctx context.Context, discussionID string) PostIter
	GetPostContentByID(ctx context.Context, id string) (*model.PostContent, error)
	PutPost(ctx context.Context, tx *sql2.Tx, post model.Post) (*model.Post, error)
	PutPostContent(ctx context.Context, tx *sql2.Tx, postContent model.PostContent) error
	GetUserProfileByID(ctx context.Context, id string) (*model.UserProfile, error)
	GetUserProfileByUserID(ctx context.Context, userID string) (*model.UserProfile, error)
	GetSocialInfosByUserProfileID(ctx context.Context, userProfileID string) ([]model.SocialInfo, error)
	UpsertSocialInfo(ctx context.Context, obj model.SocialInfo) (*model.SocialInfo, error)
	CreateOrUpdateUserProfile(ctx context.Context, userProfile model.UserProfile) (*model.UserProfile, bool, error)
	UpsertUser(ctx context.Context, user model.User) (*model.User, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetUserDevicesByUserID(ctx context.Context, userID string) ([]model.UserDevice, error)
	UpsertUserDevice(ctx context.Context, userDevice model.UserDevice) (*model.UserDevice, error)
	GetViewersByIDs(ctx context.Context, viewerIDs []string) (map[string]*model.Viewer, error)
	UpsertViewer(ctx context.Context, viewer model.Viewer) (*model.Viewer, error)
	GetPostByID(ctx context.Context, postID string) (*model.Post, error)
	PutActivity(ctx context.Context, tx *sql2.Tx, post *model.Post) error
	CreateTestTables(ctx context.Context, data TestData) (func() error, error)

	// TXN
	BeginTx(ctx context.Context) (*sql2.Tx, error)
	RollbackTx(ctx context.Context, tx *sql2.Tx) error
	CommitTx(ctx context.Context, tx *sql2.Tx) error
}

type delphisDB struct {
	dynamo   dynamodbiface.DynamoDBAPI
	sql      *gorm.DB
	pg       *sql2.DB
	dbConfig config.TablesConfig
	encoder  *dynamodbattribute.Encoder

	prepStmts *dbPrepStmts

	// Check if prepared statements are initialized
	ready   bool
	readyMu sync.RWMutex
}

type PostIter interface {
	Next(post *model.Post) bool
	Close() error
}

//MarshalMap wraps the dynamodbattribute.MarshalMap with a defined encoder.
func (d *delphisDB) marshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
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
		logrus.Debugf("endpoint: %v", *mySession.Config.Endpoint)
	}
	dbSvc := dynamodb.New(mySession)
	gormDB, db := NewSQLDatastore(config.SQLDBConfig, awsSession)
	return &delphisDB{
		dbConfig:  dbConfig.TablesConfig,
		sql:       gormDB,
		pg:        db,
		dynamo:    dbSvc,
		prepStmts: &dbPrepStmts{},
		encoder: &dynamodbattribute.Encoder{
			MarshalOptions: dynamodbattribute.MarshalOptions{
				SupportJSONTags: false,
			},
			NullEmptyString: true,
		},
	}
}

func NewSQLDatastore(sqlDbConfig config.SQLDBConfig, awsSession *session.Session) (*gorm.DB, *sql2.DB) {
	dbURI := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", sqlDbConfig.Host, sqlDbConfig.Port, sqlDbConfig.Username, sqlDbConfig.DBName, sqlDbConfig.Password)
	logrus.Debugf("About to open connection to DB - gorm")
	gormDB, err := gorm.Open("postgres", dbURI)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to open gormDB - gorm")
		return nil, nil
	}
	logrus.Debugf("Opened connection to DB! - gorm")

	// Open DB Connection for raw sql queries. We will currently support both as we slowly migrate
	logrus.Debugf("About to open connection to DB - pg")
	db, err := sql2.Open("postgres", dbURI)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to open gormDB - pg")
		return nil, nil
	}
	logrus.Debugf("Opened connection to DB! - pg")

	// Set autoload
	//gormDB = gormDB.Set("gorm:auto_preload", true)
	//gormDB = gormDB.LogMode(true)
	// need to defer closing the gormDB.
	//gormDB.AutoMigrate(model.DatabaseModels...)

	return gormDB, db
}

func (d *delphisDB) initializeStatements(ctx context.Context) (err error) {
	d.readyMu.RLock()
	ready := d.ready
	d.readyMu.RUnlock()

	if ready {
		return nil
	}

	d.readyMu.Lock()
	defer d.readyMu.Unlock()
	if d.ready {
		return nil
	}

	// POSTS
	if d.prepStmts.putPostStmt, err = d.pg.PrepareContext(ctx, putPostString); err != nil {
		logrus.WithError(err).Error("failed to prepare putPostStmt")
		return errors.Wrap(err, "failed to prepare putPostStmt")
	}

	if d.prepStmts.getPostsByDiscussionIDStmt, err = d.pg.PrepareContext(ctx, getPostsByDiscussionIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare getPostsByDiscussionIDStmt")
		return errors.Wrap(err, "failed to prepare getPostsByDiscussionIDStmt")
	}

	// POST CONTENTS
	if d.prepStmts.putPostContentsStmt, err = d.pg.PrepareContext(ctx, putPostContentsString); err != nil {
		logrus.WithError(err).Error("failed to prepare putPostContentsStmt")
		return errors.Wrap(err, "failed to prepare putPostContentsStmt")
	}

	// MENTIONS
	if d.prepStmts.putActivityStmt, err = d.pg.PrepareContext(ctx, putActivityString); err != nil {
		logrus.WithError(err).Error("failed to prepare putActivityStmt")
		return errors.Wrap(err, "failed to prepare putActivityStmt")
	}

	d.ready = true
	return
}

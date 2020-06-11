package datastore

import (
	"context"
	sql2 "database/sql"
	"fmt"
	"sync"
	"time"

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
	GetDiscussionsAutoPost(ctx context.Context) DiscussionIter
	GetDiscussionByModeratorID(ctx context.Context, moderatorID string) (*model.Discussion, error)
	CreateModerator(ctx context.Context, moderator model.Moderator) (*model.Moderator, error)
	GetModeratorByID(ctx context.Context, id string) (*model.Moderator, error)
	GetModeratorByUserID(ctx context.Context, id string) (*model.Moderator, error)
	GetModeratorByUserIDAndDiscussionID(ctx context.Context, userID, discussionID string) (*model.Moderator, error)
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
	GetParticipantsByIDs(ctx context.Context, ids []string) (map[string]*model.Participant, error)
	GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error)
	GetParticipantsByDiscussionIDUserID(ctx context.Context, discussionID string, userID string) ([]model.Participant, error)
	UpsertParticipant(ctx context.Context, participant model.Participant) (*model.Participant, error)
	GetPostsByDiscussionID(ctx context.Context, discussionID string) ([]*model.Post, error)
	GetPostsByDiscussionIDIter(ctx context.Context, discussionID string) PostIter
	GetLastPostByDiscussionID(ctx context.Context, discussionID string, minutes int) (*model.Post, error)
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
	PutMediaRecord(ctx context.Context, tx *sql2.Tx, media model.Media) error
	GetMediaRecordByID(ctx context.Context, mediaID string) (*model.Media, error)
	GetImportedContentByID(ctx context.Context, id string) (*model.ImportedContent, error)
	GetImportedContentTags(ctx context.Context, id string) TagIter
	GetDiscussionTags(ctx context.Context, id string) TagIter
	GetMatchingTags(ctx context.Context, discussionID, importedContentID string) ([]string, error)
	PutImportedContent(ctx context.Context, tx *sql2.Tx, ic model.ImportedContent) (*model.ImportedContent, error)
	PutImportedContentTags(ctx context.Context, tx *sql2.Tx, tag model.Tag) (*model.Tag, error)
	PutDiscussionTags(ctx context.Context, tx *sql2.Tx, tag model.Tag) (*model.Tag, error)
	DeleteDiscussionTags(ctx context.Context, tx *sql2.Tx, tag model.Tag) (*model.Tag, error)
	GetImportedContentByDiscussionID(ctx context.Context, discussionID string, limit int) ContentIter
	GetScheduledImportedContentByDiscussionID(ctx context.Context, discussionID string) ContentIter
	PutImportedContentDiscussionQueue(ctx context.Context, discussionID, contentID string, postedAt *time.Time, matchingTags []string) (*model.ContentQueueRecord, error)
	UpdateImportedContentDiscussionQueue(ctx context.Context, discussionID, contentID string, postedAt *time.Time) (*model.ContentQueueRecord, error)

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

type TagIter interface {
	Next(tag *model.Tag) bool
	Close() error
}

type ContentIter interface {
	Next(content *model.ImportedContent) bool
	Close() error
}

type DiscussionIter interface {
	Next(discussion *model.DiscussionAutoPost) bool
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
	if d.prepStmts.getPostsByDiscussionIDStmt, err = d.pg.PrepareContext(ctx, getPostsByDiscussionIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare getPostsByDiscussionIDStmt")
		return errors.Wrap(err, "failed to prepare getPostsByDiscussionIDStmt")
	}
	if d.prepStmts.getLastPostByDiscussionIDStmt, err = d.pg.PrepareContext(ctx, getLastPostByDiscussionIDStmt); err != nil {
		logrus.WithError(err).Error("failed to prepare getLastPostByDiscussionIDStmt")
		return errors.Wrap(err, "failed to prepare getLastPostByDiscussionIDStmt")
	}
	if d.prepStmts.putPostStmt, err = d.pg.PrepareContext(ctx, putPostString); err != nil {
		logrus.WithError(err).Error("failed to prepare putPostStmt")
		return errors.Wrap(err, "failed to prepare putPostStmt")
	}

	// POST CONTENTS
	if d.prepStmts.putPostContentsStmt, err = d.pg.PrepareContext(ctx, putPostContentsString); err != nil {
		logrus.WithError(err).Error("failed to prepare putPostContentsStmt")
		return errors.Wrap(err, "failed to prepare putPostContentsStmt")
	}

	// ACTIVITY
	if d.prepStmts.putActivityStmt, err = d.pg.PrepareContext(ctx, putActivityString); err != nil {
		logrus.WithError(err).Error("failed to prepare putActivityStmt")
		return errors.Wrap(err, "failed to prepare putActivityStmt")
	}

	// MEDIA
	if d.prepStmts.putMediaRecordStmt, err = d.pg.PrepareContext(ctx, putMediaRecordString); err != nil {
		logrus.WithError(err).Error("failed to prepare putMediaRecordStmt")
		return errors.Wrap(err, "failed to prepare putMediaRecordStmt")
	}
	if d.prepStmts.getMediaRecordStmt, err = d.pg.PrepareContext(ctx, getMediaRecordString); err != nil {
		logrus.WithError(err).Error("failed to prepare getMediaRecordStmt")
		return errors.Wrap(err, "failed to prepare getMediaRecordStmt")
	}

	// DISCUSSION
	if d.prepStmts.getDiscussionsForAutoPostStmt, err = d.pg.PrepareContext(ctx, getDiscussionsForAutoPostString); err != nil {
		logrus.WithError(err).Error("failed to prepare getDiscussionsForAutoPostStmt")
		return errors.Wrap(err, "failed to prepare getDiscussionsForAutoPostStmt")
	}

	// MODERATOR
	if d.prepStmts.getModeratorByUserIDStmt, err = d.pg.PrepareContext(ctx, getModeratorByUserIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare getModeratorByUserProfileID")
		return errors.Wrap(err, "failed to prepare getModeratorByUserProfileID")
	}
	if d.prepStmts.getModeratorByUserIDAndDiscussionIDStmt, err = d.pg.PrepareContext(ctx, getModeratorByUserIDAndDiscussionIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare getModeratorByUserIDAndDiscussionIDStmt")
		return errors.Wrap(err, "failed to prepare getModeratorByUserIDAndDiscussionIDStmt")
	}

	// IMPORTED CONTENT
	if d.prepStmts.getImportedContentByIDStmt, err = d.pg.PrepareContext(ctx, getImportedContentByIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare getImportedContentByIDStmt")
		return errors.Wrap(err, "failed to prepare getImportedContentByIDStmt")
	}
	if d.prepStmts.getImportedContentForDiscussionStmt, err = d.pg.PrepareContext(ctx, getImportedContentForDiscussionString); err != nil {
		logrus.WithError(err).Error("failed to prepare getImportedContentForDiscussionStmt")
		return errors.Wrap(err, "failed to prepare getImportedContentForDiscussionStmt")
	}
	if d.prepStmts.getScheduledImportedContentByDiscussionIDStmt, err = d.pg.PrepareContext(ctx, getScheduledImportedContentByDiscussionIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare getScheduledImportedContentByDiscussionIDStmt")
		return errors.Wrap(err, "failed to prepare getScheduledImportedContentByDiscussionIDStmt")
	}
	if d.prepStmts.putImportedContentStmt, err = d.pg.PrepareContext(ctx, putImportedContentString); err != nil {
		logrus.WithError(err).Error("failed to prepare putImportedContentStmt")
		return errors.Wrap(err, "failed to prepare putImportedContentStmt")
	}
	if d.prepStmts.putImportedContentDiscussionQueueStmt, err = d.pg.PrepareContext(ctx, putImportedContentDiscussionQueueString); err != nil {
		logrus.WithError(err).Error("failed to prepare putImportedContentDiscussionQueueStmt")
		return errors.Wrap(err, "failed to prepare putImportedContentDiscussionQueueStmt")
	}
	if d.prepStmts.updateImportedContentDiscussionQueueStmt, err = d.pg.PrepareContext(ctx, updateImportedContentDiscussionQueueString); err != nil {
		logrus.WithError(err).Error("failed to prepare updateImportedContentDiscussionQueueStmt")
		return errors.Wrap(err, "failed to prepare updateImportedContentDiscussionQueueStmt")
	}

	// TAGS
	if d.prepStmts.getImportedContentTagsStmt, err = d.pg.PrepareContext(ctx, getImportedContentTagsString); err != nil {
		logrus.WithError(err).Error("failed to prepare getImportedContentTagsStmt")
		return errors.Wrap(err, "failed to prepare getImportedContentTagsStmt")
	}
	if d.prepStmts.getDiscussionTagsStmt, err = d.pg.PrepareContext(ctx, getDiscussionTagsString); err != nil {
		logrus.WithError(err).Error("failed to prepare getDiscussionTagsStmt")
		return errors.Wrap(err, "failed to prepare getDiscussionTagsStmt")
	}
	if d.prepStmts.getMatchingTagsStmt, err = d.pg.PrepareContext(ctx, getMatchingTagsString); err != nil {
		logrus.WithError(err).Error("failed to prepare getMatchingTagsStmt")
		return errors.Wrap(err, "failed to prepare getMatchingTagsStmt")
	}
	if d.prepStmts.putImportedContentTagsStmt, err = d.pg.PrepareContext(ctx, putImportedContentTagsString); err != nil {
		logrus.WithError(err).Error("failed to prepare putImportedContentTagsStmt")
		return errors.Wrap(err, "failed to prepare putImportedContentTagsStmt")
	}
	if d.prepStmts.putDiscussionTagsStmt, err = d.pg.PrepareContext(ctx, putDiscussionTagsString); err != nil {
		logrus.WithError(err).Error("failed to prepare putDiscussionTagsStmt")
		return errors.Wrap(err, "failed to prepare putDiscussionTagsStmt")
	}
	if d.prepStmts.deleteDiscussionTagsStmt, err = d.pg.PrepareContext(ctx, deleteDiscussionTagsString); err != nil {
		logrus.WithError(err).Error("failed to prepare deleteDiscussionTagsStmt")
		return errors.Wrap(err, "failed to prepare deleteDiscussionTagsStmt")
	}

	d.ready = true
	return
}

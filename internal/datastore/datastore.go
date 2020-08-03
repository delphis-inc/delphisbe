package datastore

import (
	"context"
	"database/sql"
	sql2 "database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Datastore interface {
	GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error)
	GetDiscussionsByIDs(ctx context.Context, ids []string) (map[string]*model.Discussion, error)
	GetDiscussionsAutoPost(ctx context.Context) AutoPostDiscussionIter
	GetDiscussionByModeratorID(ctx context.Context, moderatorID string) (*model.Discussion, error)
	CreateModerator(ctx context.Context, moderator model.Moderator) (*model.Moderator, error)
	GetModeratorByID(ctx context.Context, id string) (*model.Moderator, error)
	GetModeratorByUserID(ctx context.Context, id string) (*model.Moderator, error)
	GetModeratorByUserIDAndDiscussionID(ctx context.Context, userID, discussionID string) (*model.Moderator, error)
	ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error)
	ListDiscussionsByUserID(ctx context.Context, userID string) (*model.DiscussionsConnection, error)
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
	GetModeratorParticipantsByDiscussionID(ctx context.Context, discussionID string) ([]model.Participant, error)
	UpsertParticipant(ctx context.Context, participant model.Participant) (*model.Participant, error)
	GetPostsByDiscussionIDIter(ctx context.Context, discussionID string) PostIter
	GetPostsByDiscussionIDFromCursorIter(ctx context.Context, discussionID string, cursor string, limit int) PostIter
	GetPostsConnectionByDiscussionID(ctx context.Context, discussionID string, cursor string, limit int) (*model.PostsConnection, error)
	GetLastPostByDiscussionID(ctx context.Context, discussionID string, minutes int) (*model.Post, error)
	GetPostContentByID(ctx context.Context, id string) (*model.PostContent, error)
	PutPost(ctx context.Context, tx *sql2.Tx, post model.Post) (*model.Post, error)
	PutPostContent(ctx context.Context, tx *sql2.Tx, postContent model.PostContent) error
	DeletePostByID(ctx context.Context, postID string, deletedReasonCode model.PostDeletedReason) (*model.Post, error)
	DeleteAllParticipantPosts(ctx context.Context, discussionID string, participantID string, deletedReasonCode model.PostDeletedReason) (int, error)
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
	GetNextShuffleTimeForDiscussionID(ctx context.Context, id string) (*model.DiscussionShuffleTime, error)
	PutNextShuffleTimeForDiscussionID(ctx context.Context, tx *sql2.Tx, id string, shuffleTime *time.Time) (*model.DiscussionShuffleTime, error)

	// Helper functions
	PostIterCollect(ctx context.Context, iter PostIter) ([]*model.Post, error)
	DiscussionAutoPostIterCollect(ctx context.Context, iter AutoPostDiscussionIter) ([]*model.DiscussionAutoPost, error)
	TagIterCollect(ctx context.Context, iter TagIter) ([]*model.Tag, error)
	ContentIterCollect(ctx context.Context, iter ContentIter) ([]*model.ImportedContent, error)
	DiscussionIterCollect(ctx context.Context, iter DiscussionIter) ([]*model.Discussion, error)
	DiscussionInviteIterCollect(ctx context.Context, iter DiscussionInviteIter) ([]*model.DiscussionInvite, error)
	AccessRequestIterCollect(ctx context.Context, iter DiscussionAccessRequestIter) ([]*model.DiscussionAccessRequest, error)

	GetDiscussionsByUserAccess(ctx context.Context, userID string) DiscussionIter
	UpsertDiscussionUserAccess(ctx context.Context, tx *sql2.Tx, discussionID, userID string) (*model.DiscussionUserAccess, error)
	DeleteDiscussionUserAccess(ctx context.Context, tx *sql2.Tx, discussionID, userID string) (*model.DiscussionUserAccess, error)
	GetDiscussionInviteByID(ctx context.Context, id string) (*model.DiscussionInvite, error)
	GetDiscussionRequestAccessByID(ctx context.Context, id string) (*model.DiscussionAccessRequest, error)
	GetDiscussionInvitesByUserIDAndStatus(ctx context.Context, userID string, status model.InviteRequestStatus) DiscussionInviteIter
	GetSentDiscussionInvitesByUserID(ctx context.Context, userID string) DiscussionInviteIter
	GetDiscussionAccessRequestsByDiscussionID(ctx context.Context, discussionID string) DiscussionAccessRequestIter
	GetDiscussionAccessRequestByDiscussionIDUserID(ctx context.Context, discussionID string, userID string) (*model.DiscussionAccessRequest, error)
	GetSentDiscussionAccessRequestsByUserID(ctx context.Context, userID string) DiscussionAccessRequestIter
	GetInvitedTwitterHandlesByDiscussionIDAndInviterID(ctx context.Context, discussionID string, invitingParticipantID string) ([]*string, error)
	PutDiscussionInviteRecord(ctx context.Context, tx *sql2.Tx, invite model.DiscussionInvite) (*model.DiscussionInvite, error)
	PutDiscussionAccessRequestRecord(ctx context.Context, tx *sql2.Tx, request model.DiscussionAccessRequest) (*model.DiscussionAccessRequest, error)
	UpdateDiscussionInviteRecord(ctx context.Context, tx *sql2.Tx, invite model.DiscussionInvite) (*model.DiscussionInvite, error)
	UpdateDiscussionAccessRequestRecord(ctx context.Context, tx *sql2.Tx, request model.DiscussionAccessRequest) (*model.DiscussionAccessRequest, error)
	GetAccessLinkBySlug(ctx context.Context, slug string) (*model.DiscussionAccessLink, error)
	GetAccessLinkByDiscussionID(ctx context.Context, discussionID string) (*model.DiscussionAccessLink, error)
	PutAccessLinkForDiscussion(ctx context.Context, tx *sql.Tx, input model.DiscussionAccessLink) (*model.DiscussionAccessLink, error)

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

type AutoPostDiscussionIter interface {
	Next(discussion *model.DiscussionAutoPost) bool
	Close() error
}

type DiscussionIter interface {
	Next(discussion *model.Discussion) bool
	Close() error
}

type DiscussionInviteIter interface {
	Next(invite *model.DiscussionInvite) bool
	Close() error
}

type DiscussionAccessRequestIter interface {
	Next(request *model.DiscussionAccessRequest) bool
	Close() error
}

func NewDatastore(config config.Config, awsSession *session.Session) Datastore {
	mySession := awsSession
	dbConfig := config.DBConfig
	if dbConfig.Host != "" && dbConfig.Port != 0 {
		mySession = mySession.Copy(awsSession.Config.WithEndpoint(fmt.Sprintf("%s:%d", dbConfig.Host, dbConfig.Port)))
		logrus.Debugf("endpoint: %v", *mySession.Config.Endpoint)
	}

	gormDB, db := NewSQLDatastore(config.SQLDBConfig, awsSession)
	return &delphisDB{
		dbConfig:  dbConfig.TablesConfig,
		sql:       gormDB,
		pg:        db,
		dynamo:    nil,
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
	if d.prepStmts.getPostByIDStmt, err = d.pg.PrepareContext(ctx, getPostByIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare getPostByIDStmt")
		return errors.Wrap(err, "failed to prepare getPostByIDStmt")
	}
	if d.prepStmts.getPostsByDiscussionIDStmt, err = d.pg.PrepareContext(ctx, getPostsByDiscussionIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare getPostsByDiscussionIDStmt")
		return errors.Wrap(err, "failed to prepare getPostsByDiscussionIDStmt")
	}
	if d.prepStmts.getLastPostByDiscussionIDStmt, err = d.pg.PrepareContext(ctx, getLastPostByDiscussionIDStmt); err != nil {
		logrus.WithError(err).Error("failed to prepare getLastPostByDiscussionIDStmt")
		return errors.Wrap(err, "failed to prepare getLastPostByDiscussionIDStmt")
	}
	if d.prepStmts.getPostsByDiscussionIDFromCursorStmt, err = d.pg.PrepareContext(ctx, getPostsByDiscussionIDFromCursorString); err != nil {
		logrus.WithError(err).Error("failed to prepare getPostsByDiscussionIDFromCursorStmt")
		return errors.Wrap(err, "failed to prepare getPostsByDiscussionIDFromCursorStmt")
	}
	if d.prepStmts.putPostStmt, err = d.pg.PrepareContext(ctx, putPostString); err != nil {
		logrus.WithError(err).Error("failed to prepare putPostStmt")
		return errors.Wrap(err, "failed to prepare putPostStmt")
	}
	if d.prepStmts.deletePostByIDStmt, err = d.pg.PrepareContext(ctx, deletePostByIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare deletePostByIDStmt")
		return errors.Wrap(err, "failed to prepare deletePostByIDStmt")
	}
	if d.prepStmts.deletePostByParticipantIDDiscussionIDStmt, err = d.pg.PrepareContext(ctx, deletePostByParticipantIDDiscussionIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare deletePostByParticipantIDDiscussionIDStmt")
		return errors.Wrap(err, "failed to prepare deletePostByParticipantIDDiscussionIDStmt")
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
		logrus.WithError(err).Error("failed to prepare getImportedContentorDiscussionStmt")
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

	// DISCUSSION ACCESS
	if d.prepStmts.getDiscussionsByUserAccessStmt, err = d.pg.PrepareContext(ctx, getDiscussionsByUserAccessString); err != nil {
		logrus.WithError(err).Error("failed to prepare getDiscussionsByUserAccessStmt")
		return errors.Wrap(err, "failed to prepare getDiscussionsByUserAccessStmt")
	}
	if d.prepStmts.upsertDiscussionUserAccessStmt, err = d.pg.PrepareContext(ctx, upsertDiscussionUserAccessString); err != nil {
		logrus.WithError(err).Error("failed to prepare upsertDiscussionUserAccessStmt")
		return errors.Wrap(err, "failed to prepare upsertDiscussionUserAccessStmt")
	}
	if d.prepStmts.deleteDiscussionUserAccessStmt, err = d.pg.PrepareContext(ctx, deleteDiscussionUserAccessString); err != nil {
		logrus.WithError(err).Error("failed to prepare deleteDiscussionUserAccessStmt")
		return errors.Wrap(err, "failed to prepare deleteDiscussionUserAccessStmt")
	}

	// INVITES AND REQUESTS
	if d.prepStmts.getDiscussionInviteByIDStmt, err = d.pg.PrepareContext(ctx, getDiscussionInviteByIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare getDiscussionInviteByIDStmt")
		return errors.Wrap(err, "failed to prepare getDiscussionInviteByIDStmt")
	}
	if d.prepStmts.getDiscussionRequestAccessByIDStmt, err = d.pg.PrepareContext(ctx, getDiscussionRequestAccessByIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare getDiscussionRequestAccessByIDStmt")
		return errors.Wrap(err, "failed to prepare getDiscussionRequestAccessByIDStmt")
	}
	if d.prepStmts.getDiscussionInvitesForUserStmt, err = d.pg.PrepareContext(ctx, getDiscussionInvitesForUserString); err != nil {
		logrus.WithError(err).Error("failed to prepare getDiscussionInvitesForUserStmt")
		return errors.Wrap(err, "failed to prepare getDiscussionInvitesForUserStmt")
	}
	if d.prepStmts.getSentDiscussionInvitesForUserStmt, err = d.pg.PrepareContext(ctx, getSentDiscussionInvitesForUserString); err != nil {
		logrus.WithError(err).Error("failed to prepare getSentDiscussionInvitesForUserStmt")
		return errors.Wrap(err, "failed to prepare getSentDiscussionInvitesForUserStmt")
	}
	if d.prepStmts.getDiscussionAccessRequestsStmt, err = d.pg.PrepareContext(ctx, getDiscussionAccessRequestsString); err != nil {
		logrus.WithError(err).Error("failed to prepare getDiscussionAccessRequestsStmt")
		return errors.Wrap(err, "failed to prepare getDiscussionAccessRequestsStmt")
	}
	if d.prepStmts.getDiscussionAccessRequestByUserIDStmt, err = d.pg.PrepareContext(ctx, getDiscussionAccessRequestByUserIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare getDiscussionAccessRequestByUserIDStmt")
		return errors.Wrap(err, "failed to prepare getDiscussionAccessRequestByUserIDStmt")
	}
	if d.prepStmts.getSentDiscussionAccessRequestsForUserStmt, err = d.pg.PrepareContext(ctx, getSentDiscussionAccessRequestsForUserString); err != nil {
		logrus.WithError(err).Error("failed to prepare getSentDiscussionAccessRequestsForUserStmt")
		return errors.Wrap(err, "failed to prepare getSentDiscussionAccessRequestsForUserStmt")
	}
	if d.prepStmts.putDiscussionInviteRecordStmt, err = d.pg.PrepareContext(ctx, putDiscussionInviteRecordString); err != nil {
		logrus.WithError(err).Error("failed to prepare putDiscussionInviteRecordStmt")
		return errors.Wrap(err, "failed to prepare putDiscussionInviteRecordStmt")
	}
	if d.prepStmts.putDiscussionAccessRequestStmt, err = d.pg.PrepareContext(ctx, putDiscussionAccessRequestString); err != nil {
		logrus.WithError(err).Error("failed to prepare putDiscussionAccessRequestStmt")
		return errors.Wrap(err, "failed to prepare putDiscussionAccessRequestStmt")
	}
	if d.prepStmts.updateDiscussionInviteRecordStmt, err = d.pg.PrepareContext(ctx, updateDiscussionInviteRecordString); err != nil {
		logrus.WithError(err).Error("failed to prepare updateDiscussionInviteRecordStmt")
		return errors.Wrap(err, "failed to prepare updateDiscussionInviteRecordStmt")
	}
	if d.prepStmts.updateDiscussionAccessRequestStmt, err = d.pg.PrepareContext(ctx, updateDiscussionAccessRequestString); err != nil {
		logrus.WithError(err).Error("failed to prepare updateDiscussionAccessRequestStmt")
		return errors.Wrap(err, "failed to prepare updateDiscussionAccessRequestStmt")
	}

	if d.prepStmts.getInvitedTwitterHandlesByDiscussionIDAndInviterIDStmt, err = d.pg.PrepareContext(ctx, getInvitedTwitterHandlesByDiscussionIDAndInviterIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare getInvitedTwitterHandlesByDiscussionIDAndInviterIDStmt")
		return errors.Wrap(err, "failed to prepare getInvitedTwitterHandlesByDiscussionIDAndInviterIDStmt")
	}

	// AccessLinks
	if d.prepStmts.getAccessLinkBySlugStmt, err = d.pg.PrepareContext(ctx, getAccessLinkBySlugString); err != nil {
		logrus.WithError(err).Error("failed to prepare getAccessLinkBySlugStmt")
		return errors.Wrap(err, "failed to prepare getAccessLinkBySlugStmt")
	}
	if d.prepStmts.getAccessLinkByDiscussionIDString, err = d.pg.PrepareContext(ctx, getAccessLinkByDiscussionIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare getAccessLinkByDiscussionIDString")
		return errors.Wrap(err, "failed to prepare getAccessLinkByDiscussionIDString")
	}
	if d.prepStmts.putAccessLinkForDiscussionString, err = d.pg.PrepareContext(ctx, putAccessLinkForDiscussionString); err != nil {
		logrus.WithError(err).Error("failed to prepare putAccessLinkForDiscussionString")
		return errors.Wrap(err, "failed to prepare putAccessLinkForDiscussionString")
	}

	// Discussion Shuffle Time
	if d.prepStmts.getNextShuffleTimeForDiscussionIDString, err = d.pg.PrepareContext(ctx, getNextShuffleTimeForDiscussionIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare getNextShuffleTimeForDiscussionIDString")
		return errors.Wrap(err, "failed to prepare getNextShuffleTimeForDiscussionIDString")
	}

	if d.prepStmts.putNextShuffleTimeForDiscussionIDString, err = d.pg.PrepareContext(ctx, putNextShuffleTimeForDiscussionIDString); err != nil {
		logrus.WithError(err).Error("failed to prepare putNextShuffleTimeForDiscussionIDString")
		return errors.Wrap(err, "failed to prepare putNextShuffleTimeForDiscussionIDString")
	}

	d.ready = true
	return
}

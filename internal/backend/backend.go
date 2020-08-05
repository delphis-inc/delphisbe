package backend

import (
	"context"
	"database/sql"
	"mime/multipart"
	"sync"
	"time"

	"github.com/delphis-inc/delphisbe/internal/mediadb"
	"github.com/delphis-inc/delphisbe/internal/twitter"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/cache"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/delphis-inc/delphisbe/internal/datastore"
	"github.com/delphis-inc/delphisbe/internal/util"
)

type DelphisBackend interface {
	CreateNewDiscussion(ctx context.Context, creatingUser *model.User, anonymityType model.AnonymityType, title string, description string, publicAccess bool, discussionSettings model.DiscussionCreationSettings) (*model.Discussion, error)
	UpdateDiscussion(ctx context.Context, id string, input model.DiscussionInput) (*model.Discussion, error)
	GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error)
	GetDiscussionsByIDs(ctx context.Context, ids []string) (map[string]*model.Discussion, error)
	GetDiscussionsForAutoPost(ctx context.Context) ([]*model.DiscussionAutoPost, error)
	GetDiscussionByModeratorID(ctx context.Context, moderatorID string) (*model.Discussion, error)
	GetDiscussionJoinabilityForUser(ctx context.Context, userObj *model.User, discussionObj *model.Discussion, meParticipant *model.Participant) (*model.CanJoinDiscussionResponse, error)
	SubscribeToDiscussion(ctx context.Context, subscriberUserID string, postChannel chan *model.Post, discussionID string) error
	UnSubscribeFromDiscussion(ctx context.Context, subscriberUserID string, discussionID string) error
	SubscribeToDiscussionEvent(ctx context.Context, subscriberUserID string, eventChannel chan *model.DiscussionSubscriptionEvent, discussionID string) error
	UnSubscribeFromDiscussionEvent(ctx context.Context, subscriberUserID string, discussionID string) error
	ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error)
	ListDiscussionsByUserID(ctx context.Context, userID string, state model.DiscussionUserAccessState) (*model.DiscussionsConnection, error)
	GetModeratorByID(ctx context.Context, id string) (*model.Moderator, error)
	GetModeratorByUserID(ctx context.Context, userID string) (*model.Moderator, error)
	CheckIfModerator(ctx context.Context, userID string) (bool, error)
	CheckIfModeratorForDiscussion(ctx context.Context, userID string, discussionID string) (bool, error)
	AssignFlair(ctx context.Context, participant model.Participant, flairID string) (*model.Participant, error)
	CreateFlair(ctx context.Context, userID string, templateID string) (*model.Flair, error)
	GetFlairByID(ctx context.Context, id string) (*model.Flair, error)
	GetFlairsByUserID(ctx context.Context, userID string) ([]*model.Flair, error)
	RemoveFlair(ctx context.Context, flair model.Flair) (*model.Flair, error)
	UnassignFlair(ctx context.Context, participant model.Participant) (*model.Participant, error)
	ListFlairTemplates(ctx context.Context, query *string) ([]*model.FlairTemplate, error)
	CreateFlairTemplate(ctx context.Context, displayName *string, imageURL *string, source string) (*model.FlairTemplate, error)
	RemoveFlairTemplate(ctx context.Context, flairTemplate model.FlairTemplate) (*model.FlairTemplate, error)
	GetFlairTemplateByID(ctx context.Context, id string) (*model.FlairTemplate, error)
	CreateParticipantForDiscussion(ctx context.Context, discussionID string, userID string, discussionParticipantInput model.AddDiscussionParticipantInput) (*model.Participant, error)
	GetParticipantsByDiscussionIDUserID(ctx context.Context, discussionID string, userID string) (*UserDiscussionParticipants, error)
	GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error)
	GetParticipantByID(ctx context.Context, id string) (*model.Participant, error)
	GetParticipantsByIDs(ctx context.Context, ids []string) (map[string]*model.Participant, error)
	GetModeratorParticipantsByDiscussionID(ctx context.Context, discussionID string) (*UserDiscussionParticipants, error)
	GetTotalParticipantCountByDiscussionID(ctx context.Context, discussionID string) int
	BanParticipant(ctx context.Context, discussionID string, participantID string, requestingUserID string) (*model.Participant, error)
	UpdateParticipant(ctx context.Context, participants UserDiscussionParticipants, currentParticipantID string, input model.UpdateParticipantInput) (*model.Participant, error)
	CreatePost(ctx context.Context, discussionID string, participantID string, input model.PostContentInput) (*model.Post, error)
	CreateAlertPost(ctx context.Context, discussionID string, userObj *model.User, isAnonymous bool) (*model.Post, error)
	PostImportedContent(ctx context.Context, participantID, discussionID, contentID string, postedAt *time.Time, matchingTags []string, dripType model.DripPostType) (*model.Post, error)
	PutImportedContentQueue(ctx context.Context, discussionID, contentID string, postedAt *time.Time, matchingTags []string, dripType model.DripPostType) (*model.ContentQueueRecord, error)
	NotifySubscribersOfCreatedPost(ctx context.Context, post *model.Post, discussionID string) error
	NotifySubscribersOfDeletedPost(ctx context.Context, post *model.Post, discussionID string) error
	NotifySubscribersOfBannedParticipant(ctx context.Context, participant *model.Participant, discussionID string) error
	GetPostsConnectionByDiscussionID(ctx context.Context, discussionID string, cursor string, limit int) (*model.PostsConnection, error)
	GetPostsByDiscussionID(ctx context.Context, userID string, discussionID string) ([]*model.Post, error)
	GetLastPostByDiscussionID(ctx context.Context, discussionID string) (*model.Post, error)
	GetPostContentByID(ctx context.Context, id string) (*model.PostContent, error)
	DeletePostByID(ctx context.Context, discussionID string, postID string, requestingUserID string) (*model.Post, error)
	GetUserProfileByID(ctx context.Context, id string) (*model.UserProfile, error)
	GetUserProfileByUserID(ctx context.Context, userID string) (*model.UserProfile, error)
	CreateUser(ctx context.Context) (*model.User, error)
	GetOrCreateAppleUser(ctx context.Context, input LoginWithAppleInput) (*model.User, error)
	GetOrCreateUser(ctx context.Context, input LoginWithTwitterInput, userObjOverride *model.User) (*model.User, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	UpsertUserDevice(ctx context.Context, deviceID string, userID *string, platform string, token *string) (*model.UserDevice, error)
	GetUserDevicesByUserID(ctx context.Context, userID string) ([]model.UserDevice, error)
	GetUserDeviceByUserIDPlatform(ctx context.Context, userID string, platform string) (*model.UserDevice, error)
	GetViewerByID(ctx context.Context, viewerID string) (*model.Viewer, error)
	GetViewersByIDs(ctx context.Context, viewerIDs []string) (map[string]*model.Viewer, error)
	CreateViewerForDiscussion(ctx context.Context, discussionID string, userID string) (*model.Viewer, error)
	GetViewerForDiscussion(ctx context.Context, discussionID, userID string, createIfNotFound bool) (*model.Viewer, error)
	GetSocialInfosByUserProfileID(ctx context.Context, userProfileID string) ([]model.SocialInfo, error)
	UpsertSocialInfo(ctx context.Context, socialInfo model.SocialInfo) (*model.SocialInfo, error)
	GetMentionedEntities(ctx context.Context, entityIDs []string) (map[string]model.Entity, error)
	NewAccessToken(ctx context.Context, userID string) (*auth.DelphisAccessToken, error)
	ValidateAccessToken(ctx context.Context, token string) (*auth.DelphisAuthedUser, error)
	ValidateRefreshToken(ctx context.Context, token string) (*auth.DelphisRefreshTokenUser, error)
	PutImportedContentAndTags(ctx context.Context, input model.ImportedContentInput) (*model.ImportedContent, error)
	GetDiscussionTags(ctx context.Context, id string) ([]*model.Tag, error)
	PutDiscussionTags(ctx context.Context, discussionID string, tags []string) ([]*model.Tag, error)
	DeleteDiscussionTags(ctx context.Context, discussionID string, tags []string) ([]*model.Tag, error)
	GetImportedContentByID(ctx context.Context, id string) (*model.ImportedContent, error)
	GetUpcomingImportedContentByDiscussionID(ctx context.Context, discussionID string) ([]*model.ImportedContent, error)
	SendNotificationsToSubscribers(ctx context.Context, discussion *model.Discussion, post *model.Post, contentPreview *string) (*SendNotificationResponse, error)
	AutoPostContent()
	GetMediaRecord(ctx context.Context, mediaID string) (*model.Media, error)
	UploadMedia(ctx context.Context, media multipart.File) (string, string, error)
	GetConciergeParticipantID(ctx context.Context, discussionID string) (string, error)
	GetDiscussionAccessByUserID(ctx context.Context, userID string, state model.DiscussionUserAccessState) ([]*model.Discussion, error)
	GetDiscussionInviteByID(ctx context.Context, id string) (*model.DiscussionInvite, error)
	GetDiscussionRequestAccessByID(ctx context.Context, id string) (*model.DiscussionAccessRequest, error)
	GetDiscussionInvitesByUserIDAndStatus(ctx context.Context, userID string, status model.InviteRequestStatus) ([]*model.DiscussionInvite, error)
	GetSentDiscussionInvitesByUserID(ctx context.Context, userID string) ([]*model.DiscussionInvite, error)
	GetDiscussionAccessRequestsByDiscussionID(ctx context.Context, discussionID string) ([]*model.DiscussionAccessRequest, error)
	GetDiscussionAccessRequestByDiscussionIDUserID(ctx context.Context, discussionID string, UserID string) (*model.DiscussionAccessRequest, error)
	GetSentDiscussionAccessRequestsByUserID(ctx context.Context, userID string) ([]*model.DiscussionAccessRequest, error)
	UpsertUserDiscussionAccess(ctx context.Context, userID string, discussionID string, state model.DiscussionUserAccessState) (*model.DiscussionUserAccess, error)
	InviteUserToDiscussion(ctx context.Context, userID, discussionID, invitingParticipantID string) (*model.DiscussionInvite, error)
	InviteTwitterUsersToDiscussion(ctx context.Context, twitterClient twitter.TwitterClient, twitterUserInfos []*model.TwitterUserInput, discussionID, invitingParticipantID string) ([]*model.DiscussionInvite, error)
	GetTwitterUserHandleAutocompletes(ctx context.Context, twitterClient twitter.TwitterClient, query string, discussionID string, invitingParticipantID string) ([]*model.TwitterUserInfo, error)
	GetTwitterAccessToken(ctx context.Context) (string, string, error)
	GetTwitterClientWithUserTokens(ctx context.Context) (twitter.TwitterClient, error)
	GetTwitterClientWithAccessTokens(ctx context.Context, accessToken string, accessTokenSecret string) (twitter.TwitterClient, error)
	DoesTwitterUserFollowUser(ctx context.Context, twitterClient twitter.TwitterClient, firstUser model.SocialInfo, secondUser model.SocialInfo) (bool, error)
	RespondToInvitation(ctx context.Context, inviteID string, response model.InviteRequestStatus, discussionParticipantInput model.AddDiscussionParticipantInput) (*model.DiscussionInvite, error)
	RequestAccessToDiscussion(ctx context.Context, userID, discussionID string) (*model.DiscussionAccessRequest, error)
	RespondToRequestAccess(ctx context.Context, requestID string, response model.InviteRequestStatus, invitingParticipantID string) (*model.DiscussionAccessRequest, error)
	GetAccessLinkBySlug(ctx context.Context, slug string) (*model.DiscussionAccessLink, error)
	GetAccessLinkByDiscussionID(ctx context.Context, discussionID string) (*model.DiscussionAccessLink, error)
	PutAccessLinkForDiscussion(ctx context.Context, discussionID string) (*model.DiscussionAccessLink, error)
	GetNextDiscussionShuffleTime(ctx context.Context, discussionID string) (*model.DiscussionShuffleTime, error)
	PutDiscussionShuffleTime(ctx context.Context, discussionID string, shuffleTime *time.Time) (*model.DiscussionShuffleTime, error)
	ShuffleDiscussionsIfNecessary()
	IncrementDiscussionShuffleCount(ctx context.Context, tx *sql.Tx, id string) (*int, error)
	GetDiscussionIDsToBeShuffledBeforeTime(ctx context.Context, tx *sql.Tx, epoc time.Time) ([]string, error)
}

type delphisBackend struct {
	db              datastore.Datastore
	auth            auth.DelphisAuth
	cache           cache.ChathamCache
	discussionMutex sync.Mutex
	config          config.Config
	timeProvider    util.TimeProvider
	mediadb         mediadb.MediaDB
	twitterBackend  twitter.TwitterBackend
}

func NewDelphisBackend(conf config.Config, awsSession *session.Session) DelphisBackend {
	chathamCache := cache.NewInMemoryCache()
	return &delphisBackend{
		db:              datastore.NewDatastore(conf, awsSession),
		auth:            auth.NewDelphisAuth(&conf.Auth),
		cache:           chathamCache,
		discussionMutex: sync.Mutex{},
		config:          conf,
		timeProvider:    &util.RealTime{},
		mediadb:         mediadb.NewMediaDB(conf, awsSession),
		twitterBackend:  &twitter.TwitterBackendImpl{},
	}
}

func (d *delphisBackend) rollbackTx(ctx context.Context, tx *sql.Tx) error {
	if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
		logrus.WithError(txErr).Error("failed to rollback tx")
		return txErr
	}
	return nil
}
